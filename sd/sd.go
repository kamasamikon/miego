package sd

import (
	"fmt"

	"miego/otot"
)

//
// Speed Dial : 创建一个表格，把各个按钮装入表格的单元格，
// 每个按钮的点击事件就是直接设置 vue.$data["xxx"]
// 配合js.sdPopup, setVueData使用
//

type SD struct {
	col   int
	id    string
	class string
	items [][]string
}

// New : 创建一个SD
// b64 - Gen()生成的HTML编码成Base64，否则直接输出HTML。
//
// ss = sd.New(true, 1)
//
//	for _, s := range arr {
//			ss.Add(OrgnName, "HospitalName", OrgnName, "HospitalID", OrgnID)
//	}
//
// bt = fmt.Sprintf(`<div onclick="sdPopupLocal(this)" data-sd="%s">{{ HospitalName }}</div>`, ss.Gen())
// ft.AddOne(bt)
func New(col int, id string, class string) *SD {
	return &SD{col: col, id: id, class: class}
}

// Add : "辅仁大学", "setVueData", "SchoolName", "FuRen"
// args[0] = Title = 弹窗窗口显示的文本
// args[1] = Function = 点击按钮时调用的函数，默认是setVueData
// args[x] = 对应了Vue的Key
// args[x+1] = 对应了Vue的Val
func (sd *SD) Add(kv ...string) {
	sd.items = append(sd.items, kv)
}

// New : col=表格列数 标题，变量名，值 ...
// 点击会调用 setVueData
func (sd *SD) Gen(b64 bool) string {
	buttonQ := `<button class="button sd-button" style="width: 100%%;" onclick="%s(`
	buttonH := `);">%s</button>`

	ft := otot.FlowTableNew(sd.id, sd.class, sd.col)
	for i := 0; i < len(sd.items); i++ {
		args := sd.items[i]
		title := args[0]
		function := args[1]
		if function == "" {
			function = "setVueData" // 默认是这个
		}

		if title == "" {
			// 占位
			ft.AddRaw("", 1, 1).SetStyle("border", "0")
			continue
		}

		var button = fmt.Sprintf(buttonQ, function)
		for j := 1; j < len(args)/2; j++ {
			key := args[2*j+0]
			val := args[2*j+1]
			button += fmt.Sprintf("'%s', '%s', ", key, val)
		}
		button += fmt.Sprintf(buttonH, title)

		ft.AddRaw(button, 1, 1).SetStyle("border", "0", "padding", "1px")
	}

	return ft.Gen(b64)
}
