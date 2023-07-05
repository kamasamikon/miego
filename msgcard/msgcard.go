package msgcard

import (
	"encoding/json"

	"github.com/kamasamikon/miego/conf"
	"github.com/kamasamikon/miego/httpdo"
	"github.com/kamasamikon/miego/klog"
)

type MsgCard struct {
	hashDocURL string

	title   string // 标题
	desc    string // 子标题 subtitle
	content string // 正文
	format  string // 正文的格式

	corpSecret string
	agentId    int
}

func MsgCardNew(hashDocURL string) *MsgCard {
	if hashDocURL == "" {
		hashDocURL = conf.Str("", "s:/hashdoc/URL")
	}

	c := MsgCard{
		hashDocURL: hashDocURL,
		corpSecret: conf.Str("", "s:/wxnotify/corp_secret"),
		agentId:    int(conf.Int(0, "i:/wxnotify/agent_id")),
	}
	return &c
}

func (c *MsgCard) SendStr(title string, desc string, content string) {
	c.SendObj(title, desc, content, "text")
}

func (c *MsgCard) SendObj(title string, desc string, content interface{}, format string) {
	var s string
	if format == "json" {
		bytes, _ := json.MarshalIndent(content, "", "    ")
		s = string(bytes)
	} else if format == "markdown" {
		s = content.(string)
	} else {
		s = content.(string)
	}

	ping := map[string]interface{}{
		"Format": format,
		"Doc":    string(s),
	}

	pong := struct{ UUID string }{}
	_, err := httpdo.New(c.hashDocURL + "/doc/add").Ping(&ping).Pong(&pong).Post()
	klog.Dump(err)
	klog.Dump(&ping)
	klog.Dump(&pong)

	// url := c.hashDocURL + "/doc?uuid=" + pong.UUID
	// wxcorp.SendCard(title, desc, url, "@all", c.corpId, c.corpSecret, c.agentId)
}
