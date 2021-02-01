package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
	// Additional part
	//
	RefreshTime int64 `json:"refreshTime"`

	//
	// Pretty
	//
	CreatedWhen string `json:"createdWhen"`
	RefreshWhen string `json:"refreshWhen"`
}

// All the service queue here.
var mapServices = make(map[string]*KService)

/////////////////////////////////////////////////////////////////////////
// Services

func msPretty(s *KService, c *gin.Context) {
	a := time.Unix(s.CreatedAt/1e9, 0)
	s.CreatedWhen = a.Format("2006/01/02 15:04:05")

	now := time.Now().UnixNano()
	ago := int(now-s.RefreshTime) / 1e9
	s.RefreshWhen = fmt.Sprintf("%d seconds ago.", ago)
}

func hashKey(serviceName string, version string, ipAddr string, port int) string {
	// return: 'msdemo:v1@127.0.0.1:8765'
	return fmt.Sprintf("%s:%s@%s:%d", serviceName, version, ipAddr, port)
}

func (s *KService) toKey() string {
	return hashKey(s.ServiceName, s.Version, s.IPAddr, s.Port)
}

func msSet(s *KService) bool {
	key := s.toKey()
	if a, ok := mapServices[key]; ok {
		if a.CreatedAt != s.CreatedAt {
			*a = *s
			a.RefreshTime = time.Now().UnixNano()
			return true
		} else {
			a.RefreshTime = time.Now().UnixNano()
			return false
		}
	}

	s.RefreshTime = time.Now().UnixNano()
	mapServices[key] = s
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
func RefreshLoop() {
	for {
		var remKeys []string
		now := time.Now().UnixNano()
		for _, s := range mapServices {
			if diff := now - s.RefreshTime; diff > 10*1000*1000*1000 {
				key := s.toKey()
				remKeys = append(remKeys, key)
			}
		}

		if len(remKeys) > 0 {
			klog.I("Some services should be deleted. %s", remKeys)
			for _, key := range remKeys {
				delete(mapServices, key)
			}

			klog.I("Bad services deleted, reloading nginx")
			nginxConfWrite()
			nginxReload()
		}
		time.Sleep(time.Second * 10)
	}
}

/////////////////////////////////////////////////////////////////////////
// Nginx

func genLocationAndUpstream() (string, string, string) {
	var serviceGroup = make(map[string][]*KService)
	for _, v := range mapServices {
		key := v.ServiceName

		var group []*KService
		group, ok := serviceGroup[key]

		if ok {
			group = append(group, v)
		} else {
			group = []*KService{v}
		}
		serviceGroup[key] = group
	}

	var redirListHttp strings.Builder
	var redirListGrpc strings.Builder
	for key, group := range serviceGroup {
		s := group[0]
		if s.Kind == "http" {
			fmt.Fprintf(&redirListHttp, "        location ^~ /%s/ {\n", key)
			fmt.Fprintf(&redirListHttp, "            proxy_pass http://%s.%s/;\n", s.ServiceName, s.Version)
			fmt.Fprintf(&redirListHttp, "        }\n\n")
		} else if s.Kind == "grpc" {
			// hard code in nginx.tmpl
		}
	}

	var upsList strings.Builder
	for _, group := range serviceGroup {
		s := group[0]
		if s.Upstream != "" {
			fmt.Fprintf(&upsList, "    upstream %s {\n", s.Upstream)
		} else {
			fmt.Fprintf(&upsList, "    upstream %s.%s {\n", s.ServiceName, s.Version)
		}
		for _, a := range group {
			fmt.Fprintf(&upsList, "        server %s:%d;\n", a.IPAddr, a.Port)
		}
		fmt.Fprintf(&upsList, "    }\n\n")
	}

	return redirListHttp.String(), redirListGrpc.String(), upsList.String()
}

func TemplLoad(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		klog.E("TemplLoad NG: %s", err.Error())
		return ""
	}
	return string(data)
}

func nginxConfWrite() error {
	// 1. Load tmpl nginx.conf
	// 2. Insert services information to tmpl
	// 3. Write back to nginx.conf
	redirListHttp, redirListGrpc, upsList := genLocationAndUpstream()

	tmpl := TemplLoad(conf.Str("/etc/nginx/nginx.conf.tmpl", "msb/nginx/tmpl"))

	tmpl = strings.Replace(tmpl, "#@@UPSTREAM_LIST@@", upsList, -1)
	tmpl = strings.Replace(tmpl, "#@@REDIRECT_LIST_HTTP@@", redirListHttp, -1)
	tmpl = strings.Replace(tmpl, "#@@REDIRECT_LIST_GRPC@@", redirListGrpc, -1)

	path := conf.Str("/etc/nginx/nginx.conf", "msb/nginx/conf")
	if err := ioutil.WriteFile(path, []byte(tmpl), os.ModeAppend); err != nil {
		klog.E(err.Error())
		return err
	}

	return nil
}

func nginxReload() {
	nginx := conf.Str("/usr/sbin/nginx", "msb/nginx/exec")
	cmd := exec.Command(nginx, "-s", "reload")
	err := cmd.Run()
	if err != nil {
		klog.E(err.Error())
	}
}

func serverSet(c *gin.Context) {
	var s KService
	c.BindJSON(&s)

	klog.Dump(s)
	if ok := msSet(&s); ok {
		nginxConfWrite()
		nginxReload()
	}

	pong := gin.H{}
	c.JSON(200, &pong)
}

func serverGet(c *gin.Context) {
	serviceName := c.Param("name")
	if serviceName == "" {
		serviceName = "ALL"
	}

	version := c.Param("version")
	if version == "" {
		version = "ALL"
	}

	ipAddr := c.Param("ipaddr")
	if ipAddr == "" {
		ipAddr = "ALL"
	}

	var iport int64
	port := c.Param("port")
	if port == "" || port == "ALL" {
		iport = -1
	} else {
		if x, err := strconv.ParseInt(port, 10, 64); err != nil {
			klog.E(err.Error())
			c.JSON(404, nil)
			return
		} else {
			iport = x
		}
	}

	var services []*KService
	for _, v := range mapServices {
		if serviceName != "ALL" && v.ServiceName != serviceName {
			continue
		}
		if version != "ALL" && v.Version != version {
			continue
		}
		if ipAddr != "ALL" && v.IPAddr != ipAddr {
			continue
		}
		if iport != -1 && int64(v.Port) != iport {
			continue
		}

		msPretty(v, c)
		services = append(services, v)
	}

	if services == nil {
		c.JSON(404, nil)
	} else {
		c.JSON(200, services)
	}
}

func serverRem(c *gin.Context) {
	serviceName := c.Param("name")
	if serviceName == "" {
		serviceName = "ALL"
	}

	version := c.Param("version")
	if version == "" {
		version = "ALL"
	}

	ipAddr := c.Param("ipaddr")
	if ipAddr == "" {
		ipAddr = "ALL"
	}

	var iport int64
	port := c.Param("port")
	if port == "" || port == "ALL" {
		iport = -1
	} else {
		if x, err := strconv.ParseInt(port, 10, 64); err != nil {
			klog.E(err.Error())
			c.JSON(404, nil)
			return
		} else {
			iport = x
		}
	}

	var services []*KService
	for _, v := range mapServices {
		if serviceName != "ALL" && v.ServiceName != serviceName {
			continue
		}
		if version != "ALL" && v.Version != version {
			continue
		}
		if ipAddr != "ALL" && v.IPAddr != ipAddr {
			continue
		}
		if iport != -1 && int64(v.Port) != iport {
			continue
		}

		services = append(services, v)
	}

	if services == nil {
		c.JSON(404, nil)
	} else {
		for _, s := range services {
			klog.D(s.ServiceName)
			key := s.toKey()
			delete(mapServices, key)
		}

		nginxConfWrite()
		nginxReload()

		c.JSON(200, services)
	}
}

func main() {
	// s:/msb/nginx/conf=/etc/nginx/nginx.conf
	// s:/msb/nginx/tmpl=/etc/nginx/nginx.conf
	// s:/msb/nginx/exec=/usr/sbin/nginx
	conf.Load("./etc/msb.cfg")
	conf.Load("./msb.cfg")

	gin.SetMode(gin.ReleaseMode)
	Gin := gin.New()

	Gin.POST("/service", serverSet)

	Gin.GET("/service", serverGet)
	Gin.GET("/service/:name", serverGet)
	Gin.GET("/service/:name/:version", serverGet)
	Gin.GET("/service/:name/:version/:ipaddr", serverGet)
	Gin.GET("/service/:name/:version/:ipaddr/:port", serverGet)

	Gin.DELETE("/service", serverRem)
	Gin.DELETE("/service/:name", serverRem)
	Gin.DELETE("/service/:name/:version", serverRem)
	Gin.DELETE("/service/:name/:version/:ipaddr", serverRem)
	Gin.DELETE("/service/:name/:version/:ipaddr/:port", serverRem)

	Gin.GET("/nginx", func(c *gin.Context) {
		tmpl := TemplLoad("/etc/nginx/nginx.conf")
		c.String(200, tmpl)
	})

	Gin.GET("/reload", func(c *gin.Context) {
		nginxConfWrite()
		nginxReload()
		c.String(200, "DONE")
	})

	go RefreshLoop()

	// XXX: Must be 9100, it is defined in /etc/nginx/nginx.conf
	Gin.Run(":9100")
}
