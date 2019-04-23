package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/kamasamikon/miego/klog"
	"github.com/kamasamikon/miego/xgin"
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

	//
	// Additional part
	//
	RefreshTime int64 `json:"refreshTime"`

	//
	// Pretty
	//
	CreatedWhen string
	RefreshWhen string
}

// All the service queue here.
var mapServices = make(map[string]*KService)

/////////////////////////////////////////////////////////////////////////
// Services

func msPretty(s *KService) {
	// TODO: Fill CreatedWhen and RefreshWhen
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
			klog.W("U: '%s'.", key)
			*a = *s
			a.RefreshTime = time.Now().UnixNano()
			return true
		} else {
			klog.W("S: '%s'", key)
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
		klog.E("Waiting 10 seconds before next loop.")
		time.Sleep(time.Second * 10)
	}
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
		fmt.Fprintf(&redir, "        location ^~ /ms/%s/ {\n", key)
		fmt.Fprintf(&redir, "            proxy_pass http://%s.%s/;\n", s.ServiceName, s.Version)
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

func TemplLoad(path string) string {
	if path == "" {
		path = "/etc/nginx/nginx.conf.templ"
	}

	klog.D("use templ '%s'", path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		klog.E("TemplLoad NG: %s", err.Error())
		return ""
	}
	return string(data)
}

func nginxConfWrite() error {
	us, lb := genLocationAndUpstream()

	templ := TemplLoad("")

	templ = strings.Replace(templ, "#@@UPSTREAM_LIST@@", lb, -1)
	templ = strings.Replace(templ, "#@@REDIRECT_LIST@@", us, -1)

	klog.D("%s", templ)

	path := "/etc/nginx/nginx.conf"
	if err := ioutil.WriteFile(path, []byte(templ), os.ModeAppend); err != nil {
		klog.E(err.Error())
		return err
	}

	return nil
}

func nginxReload() {
	klog.D("Reload nginx")
	cmd := exec.Command("/usr/sbin/nginx", "-s", "reload")
	err := cmd.Run()
	if err != nil {
		klog.E(err.Error())
	}
}

func serverSet(c *gin.Context) {
	var s KService
	err := c.BindJSON(&s)
	if err != nil {
		spew.Dump(s)
	}

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
	spew.Config.Indent = "\t"

	Gin := xgin.New(false)

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

	go RefreshLoop()

	// XXX: Must be 9100, it is defined in /etc/nginx/nginx.conf
	Gin.Run(":9100")
}
