package msa

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"lib/conf"
	"lib/klog"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"

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
	msPath := conf.Str("ms/path", "/usr/bin/git")

	// wait till OK
	cmd := exec.Command(msPath, "help")
	in := bytes.NewBuffer(nil)
	cmd.Stdin = in

	var out bytes.Buffer
	cmd.Stdout = &out

	go func() {
		in.WriteString("node E:/design/test.js\n")
	}()

	err := cmd.Start()
	if err != nil {
		klog.E("Command finished with error: %v", err)
	}
	klog.D("%s", cmd.Args)

	err = cmd.Wait()
	if err != nil {
		klog.E("Command finished with error: %v", err)
	}
	klog.D(out.String())
}

func (s *KService) regLoop() {
	msbURL := conf.Str("msb/url", "http://127.0.0.1:7766/")
	waitOK := time.Duration(conf.Int("msb/regWait/ok", 5))
	waitNG := time.Duration(conf.Int("msb/regWait/ng", 1))

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

func init() {
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
