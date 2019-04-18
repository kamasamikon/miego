package main

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
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
var msbURL string

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
	var exePath string
	var workDir string

	if exists("/root/ms/main") {
		exePath = "/root/ms/main"
		workDir = "/root/ms"
	} else {
		klog.F("BAD SERVICE. /root/ms/main not found.")
		return
	}

	go func() {
		for {
			nsBefore := time.Now().UnixNano()

			cmd := exec.Command(exePath)
			cmd.Dir = workDir
			cmd.Env = append(cmd.Env, "MSBURL="+msbURL)

			// in := bytes.NewBuffer(nil)
			// cmd.Stdin = in

			// XXX: Not save to buffer, because the output maybe too long.
			// var out bytes.Buffer
			// cmd.Stdout = &out

			klog.D("RUN /ms/main")
			err := cmd.Run()
			if err != nil {
				klog.E("Command finished with error: %s", err.Error())
			}
			klog.D("Normal exit, but it will be restarted now.")

			nsAfter := time.Now().UnixNano()
			if nsAfter-nsBefore < 1*1000*1000*1000 {
				klog.C("Service quit too frequenty.")
			}
		}
	}()
}

func msbInfoSet() {
	msbURL = conf.Str("msa/msb/url", "http://172.17.0.1/msb")
	if url := os.Getenv("MSBURL"); url != "" {
		msbURL = url
	}
}

func (s *KService) regLoop() {
	waitOK := time.Duration(conf.Int("msa/msb/regWait/ok", 5))
	waitNG := time.Duration(conf.Int("msa/msb/regWait/ng", 1))

	j, _ := json.Marshal(&s)
	spew.Dump(s)

	msRegURL := msbURL + "/service"
	klog.D("msRegURL: %s", msRegURL)
	for {
		_, err := http.Post(msRegURL, "application/json", strings.NewReader(string(j)))
		if err == nil {
			time.Sleep(time.Second * waitOK)
		} else {
			klog.E("%s @%s", err.Error(), msRegURL)
			time.Sleep(time.Second * waitNG)
		}
	}
}

func main() {
	conf.Load("/root/msa.cfg")
	conf.Load("/root/usr.cfg")
	conf.Load("/root/ms/msa.cfg")
	conf.Load("/root/ms/usr.cfg")

	service = &KService{
		ServiceName: conf.Str("msa/ms/name", "FIXME"),
		Version:     conf.Str("msa/ms/version", "FIXME"),
		Desc:        conf.Str("msa/ms/desc", "FIXME"),

		IPAddr: GetOutboundIP(),
		Port:   int(conf.Int("msa/ms/port", 8888)),

		HostName: hostnameGet(),

		ProjName:    conf.Str("msa/build/dirname", "FIXME"),
		ProjVersion: conf.Str("msa/build/version", "FIXME"),
		ProjTime:    conf.Str("msa/build/time", "FIXME"),

		CreatedAt: time.Now().UnixNano(),
	}

	msbInfoSet()
	service.programRun()
	service.regLoop()
}
