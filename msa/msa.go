package main

import (
	"bytes"
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
	ServiceName string `json:"serviceName"`
	Version     string `json:"version"`
	Desc        string `json:"desc"`
	IPAddr      string `json:"ipAddr"`
	Port        int    `json:"port"`
	HostName    string `json:"hostName"`
	ProjName    string `json:"projName"`
	CreatedAt   string `json:"createdAt"`
}

var service *KService

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

func (s *KService) serviceRun() {
	// ms/path: /usr/bin/git config '-l'
	msPath := conf.Str("ms/path", "echo aaa bbb ccc ddd")
	args := strings.Split(msPath, " ")

	cmd := exec.Command(args[0], args[1:]...)

	in := bytes.NewBuffer(nil)
	cmd.Stdin = in

	// var out bytes.Buffer
	// cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		klog.E("Command finished with error: %s", err.Error())
	}
}

func (s *KService) regLoop() {
	msbURL := conf.Str("msb/url", "http://127.0.0.1:7766/")
	waitOK := time.Duration(conf.Int("msb/regWait/ok", 5))
	waitNG := time.Duration(conf.Int("msb/regWait/ng", 1))

	if url := os.Getenv("MSBURL"); url != "" {
		msbURL = url
	}

	j, _ := json.Marshal(&s)
	spew.Dump(s)

	for {
		r, err := http.Post(msbURL, "application/json", strings.NewReader(string(j)))
		if err == nil {
			klog.D("%d", r.StatusCode)
			time.Sleep(time.Second * waitOK)
		} else {
			klog.E("%s", err.Error())
			time.Sleep(time.Second * waitNG)
		}
	}
}

func main() {
	// FIXME: Load from os.Env or os.Args
	conf.Load("./msa.cfg")

	service = &KService{
		ServiceName: conf.Str("ms/name", "FIXME"),
		Version:     conf.Str("ms/version", "FIXME"),
		Desc:        conf.Str("ms/desc", "FIXME"),
		IPAddr:      GetOutboundIP(),
		Port:        int(conf.Int("ms/port", 8888)),
		HostName:    hostnameGet(),
		ProjName:    conf.Str("ms/projName", "FIXME"),
		CreatedAt:   time.Now().Format("FIXME"),
	}

	service.serviceRun()
	service.regLoop()
}
