package libmsa

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/kamasamikon/miego/conf"
	"github.com/kamasamikon/miego/klog"
	mscommon "github.com/kamasamikon/miego/ms/common"
)

const (
	msaChaged = "e:/msa/changed"
)

var HTTPTransport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second, // 连接超时时间
		KeepAlive: 60 * time.Second, // 保持长连接的时间
	}).DialContext, // 设置连接的参数
	MaxIdleConns:          500,              // 最大空闲连接
	IdleConnTimeout:       60 * time.Second, // 空闲连接的超时时间
	ExpectContinueTimeout: 30 * time.Second, // 等待服务第一个响应的超时时间
	MaxIdleConnsPerHost:   100,              // 每个host保持的空闲连接数
}

var msRegURL string

func GetRegURL() string {
	return msRegURL
}

func HostNameGet() string {
	if dat, err := ioutil.ReadFile("/etc/hostname"); err != nil {
		return "N/A"
	} else {
		return strings.TrimSpace(string(dat))
	}
}

// GetLocalAddr : Get preferred outbound ip of this machine
func GetLocalAddr() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// GetOutboundIP : Get preferred outbound ip of this machine
func GetOutboundIP() string {
	if addr := conf.Str("", "s:/ms/addr"); addr != "" {
		return addr
	}
	return GetLocalAddr()
}

func RunService() {
	exePath := conf.Str("", "s:/ms/exe")
	if exePath == "" {
		return
	}

	workDir := path.Dir(exePath)

	waitOK := time.Duration(conf.Int(1, "i:/ms/relaunch/ok"))
	waitNG := time.Duration(conf.Int(1, "i:/ms/relaunch/ng"))

	//
	// Prepare
	//
	cmd := exec.Command(exePath)

	cmd.Dir = workDir

	ServiceName := conf.Str("demo", "s:/ms/name")
	Version := conf.Str("v1", "s:/ms/version")
	Desc := conf.Str("", "s:/ms/desc")

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "MS_NAME="+ServiceName)
	cmd.Env = append(cmd.Env, "MS_VERSION="+Version)
	cmd.Env = append(cmd.Env, "MS_DESC="+Desc)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	//
	// Go ...
	//
	for {
		nsBefore := time.Now().UnixNano()

		err := cmd.Run()
		if err != nil {
			klog.E("cmd.Run ERROR: %s", err.Error())
			time.Sleep(time.Second * waitNG)
			if cmd.Process != nil {
				if err := cmd.Process.Kill(); err != nil {
					klog.E("cmd.Kill ERROR: %s", err.Error())
				}
			}
		} else {
			nsAfter := time.Now().UnixNano()

			klog.D("SERVICE NORMAL EXIT: RESTART NOW.")
			if nsAfter-nsBefore < 1*1000*1000*1000 {
				klog.C("Service quit too frequenty.")
			}
			time.Sleep(time.Second * waitOK)
		}
	}
}

func doReg(msRegURL string, msDataReader io.Reader) bool {
	client := http.Client{
		Timeout:   5 * time.Second,
		Transport: HTTPTransport,
	}
	resp, err := client.Post(msRegURL, "application/json", msDataReader)
	if err != nil {
		klog.E("%s, %s", msRegURL, err.Error())
		return false
	}
	defer resp.Body.Close()
	ok := resp != nil && resp.StatusCode == 200
	if !ok {
		klog.E("%s, StatusCode: %d", msRegURL, resp.StatusCode)
	}
	return ok
}

func getParam(argPrefix string, envName string, cfgName string) string {
	if argPrefix != "" {
		for _, argv := range os.Args {
			if strings.HasPrefix(argv, argPrefix) {
				x := argv[len(argPrefix):]
				return x
			}
		}
	}
	if envName != "" {
		if x := os.Getenv(envName); x != "" {
			return x
		}
	}
	if cfgName != "" {
		x := conf.Str("", cfgName)
		return x
	}
	return ""
}

// commandLine > env > configure
func RegisterLoop() {
	for {
		DockerGW := getParam("--dockerGW=", "DOCKER_GATEWAY", "s:/msb/dockerGW")
		MSBName := getParam("--msbName=", "MSBNAME", "s:/msb/name")
		MSBPort := getParam("--msbPort=", "MSBPORT", "s:/msb/port")

		DockerHelperPort := getParam("--dockerHelperPort=", "DOCKERHELPERPORT", "s:/dockerhelper/port")

		// Loop
		waitOK := time.Second * time.Duration(conf.Int(10, "i:/msb/regWait/ok"))
		waitNG := time.Second * time.Duration(conf.Int(1, "i:/msb/regWait/ng"))

		// Service
		s := mscommon.KService{
			ServiceName: conf.Str("demo", "s:/ms/name"),
			Version:     conf.Str("v1", "s:/ms/version"),
			Desc:        conf.Str("", "s:/ms/desc"),
			Upstream:    conf.Str("", "s:/ms/upstream"),
			Kind:        conf.Str("http", "s:/ms/kind"),
			IPAddr:      GetOutboundIP(),
			Port:        int(conf.Int(8888, "i:/ms/port")),
			HostName:    HostNameGet(),
			ProjName:    conf.Str("FIXME", "s:/build/dirname"),
			ProjVersion: conf.Str("FIXME", "s:/build/version"),
			ProjTime:    conf.Str("FIXME", "s:/build/time"),
			CreatedAt:   time.Now().UnixNano(),
		}
		sJson, _ := json.Marshal(&s)
		msDataReader := strings.NewReader(string(sJson))
		klog.Dump(os.Environ(), "MSA.Env: ")
		klog.Dump(s, "MSA.Srv: ")

		//
		// MSBPort: DockerGW+MSBPort
		//
		msRegURL = ""
		if DockerGW != "" && MSBPort != "" {
			msRegURL = "http://" + DockerGW + ":" + MSBPort + "/msb/service"
			klog.D("msRegURL: %s", msRegURL)
			for {
				msDataReader.Seek(io.SeekStart, 0)
				if !doReg(msRegURL, msDataReader) {
					break
				}
				time.Sleep(waitOK)
			}
		}

		//
		// MSBName: 通过dockerhelper
		//
		msRegURL = ""
		if DockerGW != "" && MSBName != "" {
			if DockerHelperPort == "" {
				DockerHelperPort = "11111"
			}

			dockerHelperURL := "http://" + DockerGW + ":" + DockerHelperPort + "/info?byName=" + MSBName
			client := http.Client{
				Timeout: 5 * time.Second,
			}
			if resp, err := client.Get(dockerHelperURL); err == nil {
				if resp.StatusCode == 200 {
					if payload, err := ioutil.ReadAll(resp.Body); err == nil {
						var dict map[string]interface{}
						if err := json.Unmarshal(payload, &dict); err == nil {
							klog.Dump(dict)
							IPAddress := dict["Data"].(map[string]interface{})["IPAddress"].(string)
							msRegURL = "http://" + IPAddress + "/msb/service"
						}
					}
				}
				resp.Body.Close()
			}

			if msRegURL != "" {
				klog.D("msRegURL: %s", msRegURL)
				for {
					msDataReader.Seek(io.SeekStart, 0)
					if !doReg(msRegURL, msDataReader) {
						break
					}
					time.Sleep(waitOK)
				}
			}
		}

		time.Sleep(waitNG)
	}
}

func init() {
	// conf.EntryAdd(msaChaged+"=", false)
	// conf.MonitorAdd(msaChaged, func(p string, o, n interface{}) { })
}
