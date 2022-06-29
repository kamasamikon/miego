package wxcorp

import (
	"fmt"

	"github.com/kamasamikon/miego/httpdo"
	"github.com/kamasamikon/miego/klog"
)

var (
	CORP_ID     string
	CORP_SECRET string
	AGENTID     string
	URLBASE     string
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

func init() {
	CORP_ID = conf.String("", "wxnotify/corp_id")
	CORP_SECRET = conf.String("", "wxnotify/corp_secret")
	AGENTID = conf.String("", "wxnotify/agentid")
	URLBASE = conf.String("", "wxnotify/urlbase")
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

	url := fmt.Sprintf("%sgettoken?corpid=%s&corpsecret=%s", URLBASE, corpId, corpSecret)
	pong := KPong_GetToken{}
	if _, err := httpdo.Get(url, &pong); err != nil {
		fmt.Println(err.Error())
		return ""
	}

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
		corpId = CORP_ID
	}
	if corpSecret == "" {
		corpSecret = CORP_SECRET
	}
	if agentId == 0 {
		agentId = AGENTID
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
	url := fmt.Sprintf("%smessage/send?access_token=%s", URLBASE, t)
	klog.D(url)
	if _, err := httpdo.Post(url, &ping, &pong); err != nil {
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
		corpId = CORP_ID
	}
	if corpSecret == "" {
		corpSecret = CORP_SECRET
	}
	if agentId == 0 {
		agentId = AGENTID
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
	url := fmt.Sprintf("%smessage/send?access_token=%s", URLBASE, t)
	klog.D(url)
	if _, err := httpdo.Post(url, &ping, &pong); err != nil {
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
