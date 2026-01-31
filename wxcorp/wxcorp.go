package wxcorp

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"miego/conf"
	"miego/httpdo"
	"miego/klog"
)

type KPong_GetToken struct {
	AccessToken string `json:"access_token"`
}

type KPing_SendText struct {
	ToUser  string `json:"touser"`
	MsgType string `json:"msgtype"`
	AgentID int    `json:"agentid"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	Safe int `json:"safe"`
}

type KPing_SendCard struct {
	ToUser   string `json:"touser"`
	MsgType  string `json:"msgtype"`
	AgentID  int    `json:"agentid"`
	TextCard struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
	} `json:"textcard"`
}

type KPong_SendXyz struct {
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
}

type WxToken struct {
	tokens map[string]string
}

var wxToken *WxToken = nil

var HTTPTransport = http.Transport{
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
}

func token() *WxToken {
	if wxToken == nil {
		wxToken = &WxToken{}
		wxToken.tokens = make(map[string]string)
	}
	return wxToken
}

func (t *WxToken) Get(corpId string, corpSecret string) string {
	seed := fmt.Sprintf("%s@%s", corpId, corpSecret)
	if token, ok := t.tokens[seed]; ok {
		return token
	}

	URLBASE := conf.S("wxnotify/urlbase")
	url := fmt.Sprintf("%sgettoken?corpid=%s&corpsecret=%s", URLBASE, corpId, corpSecret)
	klog.D("%s", url)
	pong := KPong_GetToken{}
	if _, err := httpdo.New(url).Transport(&HTTPTransport).Pong(&pong).Get(); err != nil {
		klog.E("%s", err.Error())
		return ""
	}
	klog.Dump(pong)

	t.tokens[seed] = pong.AccessToken
	return pong.AccessToken
}

func (t *WxToken) Rem(corpId string, corpSecret string) {
	seed := fmt.Sprintf("%s@%s", corpId, corpSecret)
	delete(t.tokens, seed)
}

func SendText(text string, toUser string, corpId string, corpSecret string, agentId int) {
	// Set default
	if corpId == "" {
		corpId = conf.S("wxnotify/corp_id")
	}
	if corpSecret == "" {
		corpSecret = conf.S("wxnotify/corp_secret")
	}
	if agentId == 0 {
		agentId = int(conf.I("wxnotify/agentid", 0))
	}

	// set data
	ping := KPing_SendText{
		ToUser:  toUser,
		MsgType: "text",
		AgentID: agentId,
		Safe:    0,
	}
	ping.Text.Content = text
	pong := &KPong_SendXyz{}

	// Send
	t := token().Get(corpId, corpSecret)
	URLBASE := conf.S("wxnotify/urlbase")
	url := fmt.Sprintf("%smessage/send?access_token=%s", URLBASE, t)
	klog.D(url)
	if _, err := httpdo.New(url).Ping(&ping).Pong(&pong).Post(); err != nil {
		klog.E(err.Error())
		return
	}

	// Retry is necessary
	if pong.ErrCode == 42001 {
		// token expired
		token().Rem(corpId, corpSecret)
		SendText(text, toUser, corpId, corpSecret, agentId)
		return
	}

	klog.Dump(pong)
}

func SendCard(title string, description string, URL string, toUser string, corpId string, corpSecret string, agentId int) {
	// Set default
	if corpId == "" {
		corpId = conf.S("wxnotify/corp_id")
	}
	if corpSecret == "" {
		corpSecret = conf.S("wxnotify/corp_secret")
	}
	if agentId == 0 {
		agentId = int(conf.I("wxnotify/agentid", 0))
	}

	// set data
	ping := KPing_SendCard{
		ToUser:  toUser,
		MsgType: "textcard",
		AgentID: agentId,
	}
	ping.TextCard.Title = title
	ping.TextCard.Description = description
	ping.TextCard.URL = URL
	pong := &KPong_SendXyz{}

	// Send
	t := token().Get(corpId, corpSecret)
	URLBASE := conf.S("wxnotify/urlbase")
	url := fmt.Sprintf("%smessage/send?access_token=%s", URLBASE, t)
	klog.D(url)
	if _, err := httpdo.New(url).Ping(&ping).Pong(&pong).Post(); err != nil {
		klog.E(err.Error())
		return
	}

	// Retry is necessary
	if pong.ErrCode == 42001 || pong.ErrCode == 40014 {
		klog.E("Token Error: %d", pong.ErrCode)
		// token expired
		token().Rem(corpId, corpSecret)
		SendCard(title, description, URL, toUser, corpId, corpSecret, agentId)
		return
	}

	klog.Dump(ping)
	klog.Dump(pong)
}
