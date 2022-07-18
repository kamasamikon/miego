package vuepage

import (
	"fmt"

	"github.com/kamasamikon/miego/otot"
	"github.com/kamasamikon/miego/sd"
)

//
// 1. 生成一个表格 otot.FlowTable
// 2. 给这个表格装入一系列的按钮
//

type ButtonChoice struct {
	pc Context
	ft *otot.FlowTable

	divStyleMap  map[string]string
	spanStyleMap map[string]string
	tdStyleMap   map[string]string
}

func ButtonChoiceNew(pc Context, ft *otot.FlowTable) *ButtonChoice {
	return &ButtonChoice{
		pc: pc,
		ft: ft,
	}
}

func (c *ButtonChoice) Style(Kind string, kv ...string) *ButtonChoice {
	var m map[string]string

	switch Kind {
	case "div":
		m = c.divStyleMap
	case "span":
		m = c.spanStyleMap
	case "td":
		m = c.tdStyleMap
	}

	if m != nil {
		for i := 0; i < len(kv)/2; i++ {
			m[kv[2*i]] = kv[2*i+1]
		}
	}

	return c
}

func (c *ButtonChoice) SA(s ...string) []string {
	return s
}

// Add : TitleKey 对应了一个按钮，点击那个按钮后，设置buttons里对应的vue的变量
func (c *ButtonChoice) Add(TitleKey string, buttons ...[]string) {
	ss := sd.New(1, "", "")
	c.pc.PutVue(TitleKey, "")

	// 创建弹出的小窗口，这个小窗口就是一个按钮的表格
	for i := 0; i < len(buttons); i++ {
		var ssArgs []string

		button := buttons[i]
		ButtonText := button[0]

		ssArgs = append(ssArgs, ButtonText)
		ssArgs = append(ssArgs, "setVueData")

		ssArgs = append(ssArgs, TitleKey)
		ssArgs = append(ssArgs, ButtonText)

		Others := button[1:]

		for j := 0; j < len(Others)/2; j++ {
			Key := Others[2*j]
			Val := Others[2*j+1]
			c.pc.PutVue(Key, "")
			ssArgs = append(ssArgs, Key)
			ssArgs = append(ssArgs, Val)
		}

		ss.Add(ssArgs...)
	}

	var divStyle string
	for k, v := range c.divStyleMap {
		divStyle += fmt.Sprintf("%s:%s;", k, v)
	}

	var spanStyle string
	for k, v := range c.spanStyleMap {
		spanStyle += fmt.Sprintf("%s:%s;", k, v)
	}

	var tdStyle []string
	for k, v := range c.spanStyleMap {
		tdStyle = append(tdStyle, k)
		tdStyle = append(tdStyle, v)
	}

	// 创建一个按钮，点击这个按钮就弹出一个小窗口（sd）
	// 这个按钮是追加到FlowTable中的
	bt := fmt.Sprintf(
		`<div class="right button" style="%s" onclick="sdPopupLocal(this)" data-sd="%s"><span style="%s">{{ %s }}</span></div>`,
		divStyle,
		ss.Gen(true),
		spanStyle,
		TitleKey,
	)

	c.ft.AddRaw(bt, 1, 1).SetStyle(tdStyle...)
}
