package sd

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/kamasamikon/miego/otot"
)

//
// Speed Dial : 创建一个表格，把各个按钮装入表格的单元格，
// 每个按钮的点击事件就是直接设置 vue.$data["xxx"]
// 配合js.sdPopup, setVueData使用
//

type SD struct {
	b64   bool
	col   int
	items [][]string
}

func New(b64 bool, col int) *SD {
	return &SD{b64: b64, col: col}
}

// args[0] = Title = 弹窗窗口显示的文本
// args[x] = 对应了Vue的Key
// args[x+1] = 对应了Vue的Val
func (sd *SD) Add(kv ...string) {
	sd.items = append(sd.items, kv)
}

// New : col=表格列数 标题，变量名，值 ...
// 点击会调用 setVueData
func (sd *SD) Gen() string {
	buttonQ := `<button class="button is-dark" style="width: 100%;" onclick="app.setVueData(`
	buttonH := `);">%s</button>`

	ft := otot.FlowTableNew("333", "ftwhite", sd.col)
	for i := 0; i < len(sd.items); i++ {
		args := sd.items[i]
		title := args[0]

		if title == "" {
			ft.AddOne("").SetStyle("border", "0")
			continue
		}

		var button = buttonQ
		for j := 0; j < len(args)/2; j++ {
			key := args[2*j+1]
			val := args[2*j+2]
			button += fmt.Sprintf("'%s', '%s', ", key, val)
		}
		button += fmt.Sprintf(buttonH, title)

		ft.AddOne(button).SetStyle("border", "0")
	}

	html := ft.Gen()
	if sd.b64 {
		html = url.QueryEscape(html)
		html = strings.Replace(html, "+", "%20", -1)
		return base64.StdEncoding.EncodeToString([]byte(html))
	} else {
		return html
	}
}
