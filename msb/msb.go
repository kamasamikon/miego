package msb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
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
		return false
	}

	key := s.toKey()
	if _, ok := mapServices[key]; ok {
		// Already exist, overwrite?
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
	return nil
}

func msRem(serviceName string, version string, ipAddr string, port int) bool {
	key := hashKey(serviceName, version, ipAddr, port)
	if _, ok := mapServices[key]; ok {
		delete(mapServices, key)
		return true
	}
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
				nginxConfWrite()
				nginxReload()
			}
			time.Sleep(time.Second * 10)
		}
	}()
}

/////////////////////////////////////////////////////////////////////////
// Nginx

func lbGen() string {
	/*
			servAddr == "msdemo.v2"

			location ~ ^/ms/msdemo/v2/(.+) {
				proxy_pass http://msb/msdemo.v2/$1;
			}

		var lines strings.Builder

		for _, v := range mapServices {
			// location /ms/
			servAddr := v.ipAddr
			orgAddr := v.ipAddr.Replace(".", "/", 1)
			lines.WriteString(fmt.Sprintf("\t\tlocation /ms/%s/ {", orgAddr))
			lines.WriteString(fmt.Sprintf("\t\t\tproxy_pass http://%s/;", servAddr))
			lines.WriteString(fmt.Sprintf("\t\t}"))
		}

		return lines.String()
	*/
	return ""
}

func upstreamGen() string {
	/*
			server msdemo:7788

			upstream msdemo {
				1.1.1.1:7788
				2.2.2.2:7788
			}

		var lines strings.Builder

		for _, v := range mapServices {
		}

	*/
	return "TODO"
}

func templLoad() string {
	data, err := ioutil.ReadFile("./nginx.conf.templ")
	if err != nil {
		return ""
	}
	return string(data)
}

func nginxConfWrite() error {
	// re-generate nginx.conf file.

	us := upstreamGen()
	lb := lbGen()

	templ := templLoad()

	templ = strings.Replace(templ, "XXXXXXXXX", us, -1)
	templ = strings.Replace(templ, "YYYYYYYYY", lb, -1)

	path := "/etc/nginx/nginx.conf"
	if err := ioutil.WriteFile(path, []byte(templ), os.ModeAppend); err != nil {
		return err
	}

	return nil
}

func nginxReload() {
	exec.Command("/usr/bin/nginx", "-s", "-reload")
}

func init() {
}
