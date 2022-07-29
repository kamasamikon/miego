package msa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/kamasamikon/miego/conf"
	"github.com/kamasamikon/miego/klog"
)

// KService : Micro Service definition
type KService struct {
	// Base info
	ServiceName string `json:"serviceName"` // 服务的名字
	Version     string `json:"version"`     // 服务的版本号，一般是v1，但基本没有用
	Desc        string `json:"desc"`        // 服务的一句话说明
	Upstream    string `json:"upstream"`    // 对应了nginx中，proxy_pass对应的upstream

	// ipAddress
	IPAddr string `json:"ipAddr"` // 本机的地址
	Port   int    `json:"port"`   // 服务监听的地址

	// container info
	HostName string `json:"hostName"` // 服务所在容器的主机名

	// Project
	ProjName    string `json:"projName"`    // 源码的目录名
	ProjVersion string `json:"projVersion"` // 源码的Git版本
	ProjTime    string `json:"projTime"`    // 工程编译的时间

	// This msa instance
	CreatedAt int64 `json:"createdAt"` // 服务启动的时间

	// Kind: grpc? http?
	Kind string `json:"kind"` // 是GRPC还是HTTP？

	// Process Information
	Cmd *exec.Cmd // 如果服务的单独的进程，通过这个来运行
}

var service *KService
var msbHost string

func hostnameGet() string {
	b, e := ioutil.ReadFile("/etc/hostname")
	if e != nil {
		return "N/A"
	}
	s := string(b)
	return strings.TrimSpace(s)
}

// GetOutboundIP : Get preferred outbound ip of this machine
func GetOutboundIP() string {
	if addr := conf.Str("", "s:/ms/addr"); addr != "" {
		return addr
	}

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return os.IsExist(err)
	}
	return true
}

func (s *KService) programRun() {
	exePath := conf.Str("/root/ms/main", "s:/ms/exe")
	workDir := path.Dir(exePath)

	if !exists(exePath) {
		klog.F("BAD SERVICE. %s not found.", exePath)
		return
	}

	waitOK := time.Duration(conf.Int(1, "i:/ms/relaunch/ok"))
	waitNG := time.Duration(conf.Int(1, "i:/ms/relaunch/ng"))

	go func() {
		//
		// Prepare
		//
		cmd := exec.Command(exePath)
		s.Cmd = cmd

		cmd.Dir = workDir
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "MS_NAME="+s.ServiceName)
		cmd.Env = append(cmd.Env, "MS_VERSION="+s.Version)
		cmd.Env = append(cmd.Env, "MS_DESC="+s.Desc)

		// in := bytes.NewBuffer(nil)
		// cmd.Stdin = in
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout

		for {
			nsBefore := time.Now().UnixNano()

			klog.D("SERVICE RUN: %s, URL: http://%s/ms/%s/%s", exePath, msbHost, s.ServiceName, s.Version)
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
	}()
}

func msbInfoSet() {
	msbHost = conf.Str("172.17.0.1", "s:/msb/host")
	if ip := os.Getenv("MSBHOST"); ip != "" {
		msbHost = ip
	}
}

func (s *KService) regLoop() {
	waitOK := time.Duration(conf.Int(10, "i:/msb/regWait/ok"))
	waitNG := time.Duration(conf.Int(1, "i:/msb/regWait/ng"))

	j, _ := json.Marshal(&s)
	klog.Dump(s, "KService: ")

	msRegURL := "http://" + msbHost + "/msb/service"
	for {
		r := strings.NewReader(string(j))
		resp, err := http.Post(msRegURL, "application/json", r)
		if err == nil {
			resp.Body.Close()
			time.Sleep(time.Second * waitOK)
		} else {
			klog.E("%s @%s", err.Error(), msRegURL)
			time.Sleep(time.Second * waitNG)
		}
	}
}

func main() {
	conf.Load("./msa.cfg")
	conf.Load("./usr.cfg")
	conf.Load("./ms/msa.cfg")
	conf.Load("./ms/usr.cfg")

	service = &KService{
		ServiceName: conf.Str("demo", "s:/ms/name"),
		Version:     conf.Str("v1", "s:/ms/version"),
		Desc:        conf.Str("", "s:/ms/desc"),
		Upstream:    conf.Str("", "s:/ms/upstream"),
		Kind:        conf.Str("http", "s:/ms/kind"),

		IPAddr: GetOutboundIP(),
		Port:   int(conf.Int(8888, "i:/ms/port")),

		HostName: hostnameGet(),

		ProjName:    conf.Str("FIXME", "s:/build/dirname"),
		ProjVersion: conf.Str("FIXME", "s:/build/version"),
		ProjTime:    conf.Str("FIXME", "s:/build/time"),

		CreatedAt: time.Now().UnixNano(),
	}

	msbInfoSet()
	service.programRun()
	go service.regLoop()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("MSA exiting ...")

	if service.Cmd != nil && service.Cmd.Process != nil {
		service.Cmd.Process.Signal(syscall.SIGINT)
	}
}
