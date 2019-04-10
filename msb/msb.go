package msb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/kamasamikon/miego/klog"
)

// KService : Micro service information loaded from configure file.
type KService struct {
	ServiceName string
	Version     string

	IPAddr string
	Port   int

	RefreshTime int
}

// All the service queue here.
var mapServices = make(map[string]*KService)

/////////////////////////////////////////////////////////////////////////
// Services

func hashKey(serviceName string, version string, ipAddr string, port int) string {
	// return: 'msdemo:v1@127.0.0.1:8765'
	return fmt.Sprintf("%s:%s@%s:%d", serviceName, version, ipAddr, port)
}

func (s *KService) toKey() string {
	return hashKey(s.ServiceName, s.Version, s.IPAddr, s.Port)
}

func msSet(postData []byte) bool {
	s := KService{}
	if err := json.Unmarshal(postData, &s); err != nil {
		klog.E("%s", err.Error())
		return false
	}

	key := s.toKey()
	if _, ok := mapServices[key]; ok {
		// Already exist, overwrite?
		klog.E("%s already exists, skip.", key)
		return false
	}
	mapServices[key] = &s

	return true
}

func msGet(serviceName string, version string, ipAddr string, port int) (s *KService) {
	key := hashKey(serviceName, version, ipAddr, port)
	if s, ok := mapServices[key]; ok {
		return s
	}
	klog.E("%s not found.", key)
	return nil
}

func msRem(serviceName string, version string, ipAddr string, port int) bool {
	key := hashKey(serviceName, version, ipAddr, port)
	if _, ok := mapServices[key]; ok {
		delete(mapServices, key)
		return true
	}
	klog.E("%s not found.", key)
	return false
}

/////////////////////////////////////////////////////////////////////////
// Refresh
func timerRefresh() {
	go func() {
		var tobeDel []string
		for {
			now := time.Now().Second()
			for _, s := range mapServices {
				if diff := now - s.RefreshTime; diff > 10 {
					key := s.toKey()
					tobeDel = append(tobeDel, key)
				}
			}

			if len(tobeDel) > 0 {
				klog.I("Some services should be deleted. %s", tobeDel)
				nginxConfWrite()
				nginxReload()
			}
			time.Sleep(time.Second * 10)
		}
	}()
}

/////////////////////////////////////////////////////////////////////////
// Nginx

func genLocationAndUpstream() (string, string) {
	var serviceGroup = make(map[string][]*KService)
	for _, v := range mapServices {
		key := v.ServiceName + "/" + v.Version

		var group []*KService
		group, ok := serviceGroup[key]

		if ok {
			group = append(group, v)
		} else {
			group = []*KService{v}
		}
		serviceGroup[key] = group
	}

	var redir strings.Builder
	for key, group := range serviceGroup {
		s := group[0]
		fmt.Fprintf(&redir, "        location ~ ^/ms/%s/(.+) {\n", key)
		fmt.Fprintf(&redir, "            proxy_pass http://%s.%s/$1;\n", s.ServiceName, s.Version)
		fmt.Fprintf(&redir, "        }\n\n")
	}

	var upstr strings.Builder
	for _, group := range serviceGroup {
		s := group[0]
		fmt.Fprintf(&upstr, "    upstream %s.%s {\n", s.ServiceName, s.Version)
		for _, a := range group {
			fmt.Fprintf(&upstr, "        server %s:%d;\n", a.IPAddr, a.Port)
		}
		fmt.Fprintf(&upstr, "    }\n\n")
	}

	return redir.String(), upstr.String()
}

func templLoad() string {
	data, err := ioutil.ReadFile("./nginx.conf.templ")
	if err != nil {
		klog.E("templLoad NG: %s", err.Error())
		return ""
	}
	return string(data)
}

func nginxConfWrite() error {
	// re-generate nginx.conf file.

	us, lb := genLocationAndUpstream()

	templ := templLoad()

	templ = strings.Replace(templ, "@@UPSTREAM_LIST@@", lb, -1)
	templ = strings.Replace(templ, "@@REDIRECT_LIST@@", us, -1)

	klog.D("%s", templ)

	// path := "/etc/nginx/nginx.conf"
	path := "/tmp/nginx.conf"
	if err := ioutil.WriteFile(path, []byte(templ), os.ModeAppend); err != nil {
		return err
	}

	return nil
}

func nginxReload() {
	exec.Command("/usr/bin/nginx", "-s", "-reload")
}

func init() {
	spew.Config.Indent = "\t"
}
