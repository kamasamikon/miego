package libmsa

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
	mscommon "github.com/kamasamikon/miego/ms/common"
)

const (
	msaChaged = "e:/msa/changed"
)

func HostNameGet() string {
	if dat, err := ioutil.ReadFile("/etc/hostname"); err != nil {
		return "N/A"
	} else {
		return strings.TrimSpace(string(dat))
	}
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

func RegisterLoop() {
	waitOK := time.Duration(conf.Int(10, "i:/msb/regWait/ok"))
	waitNG := time.Duration(conf.Int(1, "i:/msb/regWait/ng"))

	MSBAddr := conf.Str("172.17.0.1", "s:/msb/host")
	if ip := os.Getenv("MSBHOST"); ip != "" {
		MSBAddr = ip
	}

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

	j, _ := json.Marshal(&s)
	klog.Dump(s, "MSA: ")

	msRegURL := "http://" + MSBAddr + "/msb/service"
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

func init() {
	// conf.EntryAdd(msaChaged+"=", false)
	// conf.MonitorAdd(msaChaged, func(p string, o, n interface{}) { })
}
