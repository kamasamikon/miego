package main

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
	ServiceName string `json:"serviceName"`
	Version     string `json:"version"`
	Desc        string `json:"desc"`
	Upstream    string `json:"upstream"`

	// ipAddress
	IPAddr string `json:"ipAddr"`
	Port   int    `json:"port"`

	// container info
	HostName string `json:"hostName"`

	// Project
	ProjName    string `json:"projName"`
	ProjVersion string `json:"projVersion"`
	ProjTime    string `json:"projTime"`

	// This msa instance
	CreatedAt int64 `json:"createdAt"`

	// Kind: grpc? http?
	Kind string `json:"kind"`

	//
	// Process Information
	//
	Cmd *exec.Cmd
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
	if addr := conf.Str("", "ms/addr"); addr != "" {
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
	exePath := conf.Str("/root/ms/main", "ms/exe")
	workDir := path.Dir(exePath)

	if !exists(exePath) {
		klog.F("BAD SERVICE. %s not found.", exePath)
		return
	}

	waitOK := time.Duration(conf.Int(1, "ms/relaunch/ok"))
	waitNG := time.Duration(conf.Int(1, "ms/relaunch/ng"))

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
	msbHost = conf.Str("172.17.0.1", "msb/host")
	if ip := os.Getenv("MSBHOST"); ip != "" {
		msbHost = ip
	}
}

func (s *KService) regLoop() {
	waitOK := time.Duration(conf.Int(10, "msb/regWait/ok"))
	waitNG := time.Duration(conf.Int(1, "msb/regWait/ng"))

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
		ServiceName: conf.Str("demo", "ms/name"),
		Version:     conf.Str("v1", "ms/version"),
		Desc:        conf.Str("", "ms/desc"),
		Upstream:    conf.Str("", "ms/upstream"),
		Kind:        conf.Str("http", "ms/kind"),

		IPAddr: GetOutboundIP(),
		Port:   int(conf.Int(8888, "ms/port")),

		HostName: hostnameGet(),

		ProjName:    conf.Str("FIXME", "build/dirname"),
		ProjVersion: conf.Str("FIXME", "build/version"),
		ProjTime:    conf.Str("FIXME", "build/time"),

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
