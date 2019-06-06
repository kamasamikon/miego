package main

import (
	"encoding/json"
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

	"github.com/davecgh/go-spew/spew"
)

// KService : Micro Service definition
type KService struct {
	// Base info
	ServiceName string `json:"serviceName"`
	Version     string `json:"version"`
	Desc        string `json:"desc"`

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
	exePath := conf.Str("ms/exe", "/root/ms/main")
	workDir := path.Dir(exePath)

	if !exists(exePath) {
		klog.F("BAD SERVICE. %s not found.", exePath)
		return
	}

	waitOK := time.Duration(conf.Int("ms/relaunch/ok", 1))
	waitNG := time.Duration(conf.Int("ms/relaunch/ng", 1))

	go func() {
		//
		// Prepare
		//
		cmd := exec.Command(exePath)
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

		klog.D("cmd: %s", spew.Sdump(cmd))

		for {
			nsBefore := time.Now().UnixNano()

			klog.D("SERVICE RUN: %s, URL: http://%s/ms/%s/%s", exePath, msbHost, s.ServiceName, s.Version)
			err := cmd.Run()
			if err != nil {
				klog.E("SERVICE ERROR EXIT: %s", err.Error())
				time.Sleep(time.Second * waitNG)
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
	msbHost = conf.Str("msb/host", "172.17.0.1")
	if ip := os.Getenv("MSBHOST"); ip != "" {
		msbHost = ip
	}
}

func (s *KService) regLoop() {
	waitOK := time.Duration(conf.Int("msb/regWait/ok", 10))
	waitNG := time.Duration(conf.Int("msb/regWait/ng", 1))

	j, _ := json.Marshal(&s)
	klog.D("KService: %s", spew.Sdump(s))

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
		ServiceName: conf.Str("ms/name", "demo"),
		Version:     conf.Str("ms/version", "v1"),
		Desc:        conf.Str("ms/desc", "TODO: FILL DESC."),

		IPAddr: GetOutboundIP(),
		Port:   int(conf.Int("ms/port", 8888)),

		HostName: hostnameGet(),

		ProjName:    conf.Str("build/dirname", "FIXME"),
		ProjVersion: conf.Str("build/version", "FIXME"),
		ProjTime:    conf.Str("build/time", "FIXME"),

		CreatedAt: time.Now().UnixNano(),
	}

	msbInfoSet()
	service.programRun()
	service.regLoop()
}
