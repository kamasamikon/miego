package wxcard

import (
	"encoding/json"

	"miego/conf"
	"miego/httpdo"
	"miego/klog"
	"miego/wxcorp"
)

type WXCard struct {
	hashDocURL string

	corpId     string
	corpSecret string
	agentId    int
}

func WxCardNew(hashDocURL string) *WXCard {
	if hashDocURL == "" {
		hashDocURL = conf.Str("", "s:/hashdoc/URL")
	}

	c := WXCard{
		hashDocURL: hashDocURL,
		corpId:     conf.Str("", "s:/wxnotify/corp_id"),
		corpSecret: conf.Str("", "s:/wxnotify/corp_secret"),
		agentId:    int(conf.Int(0, "i:/wxnotify/agent_id")),
	}
	return &c
}

func (c *WXCard) SendStr(title string, desc string, content string) {
	if content != "" {
		ping := map[string]interface{}{
			"Doc":    string(content),
			"Format": "text",
		}

		pong := struct{ UUID string }{}
		httpdo.New(c.hashDocURL + "/doc/add").Ping(&ping).Pong(&pong).Post()

		url := c.hashDocURL + "/doc?uuid=" + pong.UUID
		wxcorp.SendCard(title, desc, url, "@all", c.corpId, c.corpSecret, c.agentId)
	} else {
		text := title + "\r\n\r\n" + desc
		wxcorp.SendText(text, "@all", c.corpId, c.corpSecret, c.agentId)
	}
}

func (c *WXCard) SendObj(title string, desc string, doc interface{}, format string) {
	var content string
	if format == "json" {
		bytes, _ := json.MarshalIndent(doc, "", "    ")
		content = string(bytes)
	} else if format == "markdown" {
		content = doc.(string)
	} else {
		content = doc.(string)
	}

	ping := map[string]interface{}{
		"Format": format,
		"Doc":    string(content),
	}

	pong := struct{ UUID string }{}
	_, err := httpdo.New(c.hashDocURL + "/doc/add").Ping(&ping).Pong(&pong).Post()
	klog.Dump(err)
	klog.Dump(&ping)
	klog.Dump(&pong)

	url := c.hashDocURL + "/doc?uuid=" + pong.UUID
	wxcorp.SendCard(title, desc, url, "@all", c.corpId, c.corpSecret, c.agentId)
}
