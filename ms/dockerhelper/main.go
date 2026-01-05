package main

import (
	"flag"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"miego/atox"
	"miego/klog"
	"miego/pong"
	"miego/xgin"
	"miego/xmap"

	"github.com/gin-gonic/gin"
)

type ContainerInfo struct {
	ID        string `json:"ID,omitempty"`
	IPAddress string `json:"IPAddress,omitempty"`
	Image     string `json:"Image,omitempty"`
	Names     string `json:"Names,omitempty"`
	Ports     []int  `json:"Ports,omitempty"`
}

func main() {
	Port := flag.Int("port", 11111, "Listen port")
	flag.Parse()

	CombinedOutput := func(cmd *exec.Cmd) string {
		out, err := cmd.CombinedOutput()
		if err != nil {
			klog.E("combined out:\n%s\n", string(out))
			klog.E("cmd.Run() failed with %s\n", err)
			return ""
		}
		return strings.TrimSpace(string(out))
	}

	// 通过容器名获取容器
	GetContainerByName := func(name string) string {
		return CombinedOutput(
			exec.Command(
				"sudo",
				"docker",
				"ps",
				"-aq",
				"--filter",
				fmt.Sprintf(
					"name=^/%s$",
					name,
				),
			),
		)
	}

	// 通过端口获取容器
	GetContainerByPort := func(port string) string {
		return CombinedOutput(
			exec.Command(
				"sudo",
				"docker",
				"ps",
				"-aq",
				"--filter",
				fmt.Sprintf(
					"publish=%s",
					port,
				),
			),
		)
	}

	GetIPAddress := func(cid string) string {
		Output := CombinedOutput(
			exec.Command(
				"sudo",
				"docker",
				"inspect",
				"-f",
				"{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}",
				cid,
			),
		)
		return strings.TrimSpace(Output)
	}

	GetInfo := func(ContainerID string) *ContainerInfo {
		Info := ContainerInfo{}
		Output := CombinedOutput(
			exec.Command(
				"sudo",
				"docker",
				"ps",
				"--filter",
				fmt.Sprintf(
					"id=%s",
					ContainerID,
				),
				"--format",
				"{{.Ports}} @ {{.Names}} @ {{.Image}} @ {{.ID}}",
			),
		)
		if Output == "" {
			return &Info
		}

		segs := strings.Split(Output, "@")

		var Ports []int
		for _, item := range strings.Split(segs[0], ",") {
			{
				re := regexp.MustCompile(`.\d+:(\d+)->\d+\/tcp`)
				arr := re.FindAllSubmatch([]byte(item), -1)
				for _, a := range arr {
					Ports = append(Ports, atox.Int(string(a[1]), 0))
				}
			}
			{
				re := regexp.MustCompile(`[^>\^\d](\d+)\/tcp`)
				arr := re.FindAllSubmatch([]byte(item), -1)
				for _, a := range arr {
					Ports = append(Ports, atox.Int(string(a[1]), 0))
				}
			}
		}
		Info.Ports = Ports
		Info.Names = strings.TrimSpace(segs[1])
		Info.Image = strings.TrimSpace(segs[2])
		Info.ID = strings.TrimSpace(segs[3])
		return &Info
	}

	// return ContainerID
	Gin := xgin.Default()
	Gin.GET("/container", func(c *gin.Context) {
		Port := c.Query("byPort")
		if Port != "" {
			ContainerID := GetContainerByPort(Port)
			pong.OK(c, ContainerID)
			return
		}

		Name := c.Query("byName")
		if Name != "" {
			ContainerID := GetContainerByName(Name)
			pong.OK(c, ContainerID)
			return
		}

		pong.OK(c, "")
	})

	// return ContainerID, ContainerName, Port
	Gin.GET("/info", func(c *gin.Context) {
		var ContainerID string

		Port := c.Query("byPort")
		if Port != "" {
			ContainerID = GetContainerByPort(Port)
		} else {
			Name := c.Query("byName")
			if Name != "" {
				ContainerID = GetContainerByName(Name)
			}
		}

		if ContainerID != "" {
			Info := GetInfo(ContainerID)
			Info.IPAddress = GetIPAddress(ContainerID)
			pong.OK(c, Info)
		} else {
			pong.OK(c, xmap.Make())
		}
	})

	xgin.Go(nil, fmt.Sprintf(":%d", *Port))
}
