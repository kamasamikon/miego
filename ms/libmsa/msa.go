package libmsa

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"miego/conf"
	"miego/klog"
	mscommon "miego/ms/common"
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
	if dat, err := os.ReadFile("/etc/hostname"); err != nil {
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
	if addr := conf.S("ms/addr"); addr != "" {
		return addr
	}
	return GetLocalAddr()
}

func RunService() {
	exePath := conf.S("ms/exe")
	if exePath == "" {
		return
	}

	workDir := path.Dir(exePath)

	relaunchOK := time.Duration(conf.I("ms/relaunch/ok", 1))
	relaunchNG := time.Duration(conf.I("ms/relaunch/ng", 1))

	//
	// Prepare
	//
	cmd := exec.Command(exePath)

	cmd.Dir = workDir

	ServiceName := conf.SGet("ms/name", "demo")
	Version := conf.SGet("ms/version", "v1")
	Desc := conf.S("ms/desc")

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
			time.Sleep(time.Second * relaunchNG)
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
			time.Sleep(time.Second * relaunchOK)
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
	ok := resp.StatusCode == 200
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
		x := conf.S(cfgName)
		return x
	}
	return ""
}

// commandLine > env > configure
func RegisterLoop() {
	CreatedAt := time.Now().UnixNano()
	IPAddr := GetOutboundIP()
	HostName := HostNameGet()

	for {
		// Loop
		regWaitOK := conf.I("msb/regWait/ok", 10)
		regWaitNG := conf.I("msb/regWait/ng", 1)

		sleepRegWaitOK := time.Second * time.Duration(regWaitOK)
		sleepRegWaitNG := time.Second * time.Duration(regWaitNG)

		DockerGW := getParam("--dockerGW=", "DOCKER_GATEWAY", "s:/msb/dockerGW")
		if DockerGW == "" {
			time.Sleep(sleepRegWaitNG)
			continue
		}

		MSBName := getParam("--msbName=", "MSBNAME", "s:/msb/name")
		MSBPort := getParam("--msbPort=", "MSBPORT", "s:/msb/port")

		DockerHelperPort := getParam("--dockerHelperPort=", "DOCKERHELPERPORT", "s:/dockerhelper/port")

		// Service
		s := mscommon.KService{
			ServiceName: conf.SGet("ms/name", "demo"),
			Version:     conf.SGet("ms/version", "v1"),
			Desc:        conf.S("ms/desc"),
			Upstream:    conf.S("ms/upstream"),
			Kind:        conf.SGet("ms/kind", "http"),
			IPAddr:      IPAddr,
			Port:        int(conf.I("ms/port", 8888)),
			HostName:    HostName,
			ProjName:    conf.S("build/dirname"),
			ProjVersion: conf.S("build/version"),
			ProjTime:    conf.S("build/time"),
			CreatedAt:   CreatedAt,
			RegInterval: regWaitOK,
		}
		sJson, _ := json.Marshal(&s)
		msDataReader := strings.NewReader(string(sJson))
		klog.Dump(s, "MSA.Srv: ")

		//
		// Save runtime information
		//
		conf.SSetf("msa/serviceName", s.ServiceName)
		conf.SSetf("msa/version", s.Version)
		conf.SSetf("msa/desc", s.Desc)
		conf.SSetf("msa/upstream", s.Upstream)
		conf.SSetf("msa/kind", s.Kind)
		conf.SSetf("msa/IPAddr", s.IPAddr)
		conf.ISetf("msa/port", s.Port)
		conf.SSetf("msa/hostName", s.HostName)
		conf.SSetf("msa/projName", s.ProjName)
		conf.SSetf("msa/projVersion", s.ProjVersion)
		conf.SSetf("msa/projTime", s.ProjTime)
		conf.ISetf("msa/createdAt", s.CreatedAt)

		//
		// MSBPort: DockerGW+MSBPort
		//
		msRegURL = ""
		if MSBPort != "" {
			msRegURL = "http://" + DockerGW + ":" + MSBPort + "/msb/service"
			conf.SSetf("msa/reg/method", "DockerGW+MSBPort")
			conf.SSetf("msa/reg/URL", msRegURL)
			conf.SSetf("msa/reg/when", time.Now().Format("2006/01/02 15:04:05"))
			klog.D("msRegURL: %s", msRegURL)
			for {
				msDataReader.Seek(0, io.SeekStart)
				if !doReg(msRegURL, msDataReader) {
					break
				}
				time.Sleep(sleepRegWaitOK)
			}
		}

		//
		// MSBName: 通过dockerhelper工具（see kamasamikon/hp/dockerhelper）
		//
		msRegURL = ""
		if MSBName != "" {
			if DockerHelperPort == "" {
				DockerHelperPort = "11111"
			}

			dockerHelperURL := "http://" + DockerGW + ":" + DockerHelperPort + "/info?byName=" + MSBName
			client := http.Client{
				Timeout: 5 * time.Second,
			}
			if resp, err := client.Get(dockerHelperURL); err == nil {
				if resp.StatusCode == 200 {
					if payload, err := io.ReadAll(resp.Body); err == nil {
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
				conf.SSetf("msa/reg/method", "dockerhelper")
				conf.SSetf("msa/reg/URL", msRegURL)
				conf.SSetf("msa/reg/when", time.Now().Format("2006/01/02 15:04:05"))
				klog.D("msRegURL: %s", msRegURL)
				for {
					msDataReader.Seek(0, io.SeekStart)
					if !doReg(msRegURL, msDataReader) {
						break
					}
					time.Sleep(sleepRegWaitOK)
				}
			}
		}

		time.Sleep(sleepRegWaitNG)
	}
}
