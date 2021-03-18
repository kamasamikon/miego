package sd

import (
	"fmt"

	"github.com/kamasamikon/miego/otot"
)

//
// Speed Dial : 创建一个表格，把各个按钮装入表格的单元格，
// 每个按钮的点击事件就是直接设置 vue.$data["xxx"]
// 配合js.sdPopup, setVueData使用
//

/*
   // Speed Dial
   setVueData: function(k, v) {
       this.$data[k] = v;
       swal.close();
   },

   // Speed Dial
   sdPopup: function(URL, Title) {
       $.ajax({
           url: URL,
           type: 'post',
           async: false,

           error: function (res) {
               console.log(res);
           },

           success: function (res) {
               var html = res.Data;
               if (html == undefined) {
                   html = res;
               }

               var template = document.createElement('template');
               template.innerHTML = html;
               var form = template.content.firstElementChild;

               form.id = rands(20);

               var buttons = {}

               swal({
                   closeOnClickOutside: true,
                   closeOnEsc: false,

                   text: Title,
                   content: form,
                   buttons: buttons,
               })
           }
       })
   },
*/

/*
func SDTxt(Kind string, Name string) string {
	URL := fmt.Sprintf("/wx/ig/sd?Kind=%s&Name=%s", Kind, Name)
	return fmt.Sprintf(
		`<button class="button" @click="sdPopup('%s')"> XXX </button>`,
		URL, Name,
	)
}
*/

// New : col=表格列数 标题，变量名，值 ...
func New(col int, args ...string) string {
	button := `<button class="button is-dark" style="width: 100%%;" onclick="app.setVueData('%s', '%s');">%s</button>`

	ft := otot.FlowTableNew("333", "ftwhite", col)
	for i := 0; i < len(args)/3; i++ {
		title := args[3*i+0]
		key := args[3*i+1]
		val := args[3*i+2]

		// hit #title, vueApp.$data[key] = val;
		if title != "" {
			ft.AddOne(fmt.Sprintf(button, key, val, title)).SetStyle("border", "0")
		} else {
			ft.AddOne("").SetStyle("border", "0")
		}
	}

	return ft.Gen()
}
