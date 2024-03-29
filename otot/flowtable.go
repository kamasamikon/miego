package otot

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

var idIndex = 0

func ElementID() string {
	idIndex++
	return fmt.Sprintf("ft%d", idIndex)
}

type TD struct {
	HTML     string
	colspan  int
	rowspan  int
	class    map[string]int
	styleMap map[string]string
}

func (td *TD) SetHTML(x string) *TD {
	td.HTML = x
	return td
}
func (td *TD) SetColSpan(x int) *TD {
	td.colspan = x
	return td
}
func (td *TD) SetRowSpan(x int) *TD {
	td.rowspan = x
	return td
}

func (td *TD) SetClass(k string) *TD {
	td.class[k] = 1
	return td
}
func (td *TD) RemClass(k string) *TD {
	delete(td.class, k)
	return td
}

func (td *TD) SetStyle(kv ...string) *TD {
	for i := 0; i < len(kv)/2; i++ {
		td.styleMap[kv[2*i]] = kv[2*i+1]
	}
	return td
}
func (td *TD) RemStyle(k ...string) *TD {
	for i := 0; i < len(k); i++ {
		delete(td.styleMap, k[i])
	}
	return td
}

type FlowTable struct {
	ID     string
	Class  string
	Column int
	Items  []*TD

	VueVars []string // 涉及到的Vue的Model变量
}

func FlowTableNew(ID string, Class string, Column int) *FlowTable {
	if ID == "" {
		ID = ElementID()
	}
	return &FlowTable{
		ID:     ID,
		Class:  Class,
		Column: Column,
	}
}

func (ft *FlowTable) Last() *TD {
	return ft.Items[len(ft.Items)-1]
}

func (ft *FlowTable) AddHeader(title string) {
	colspan := ft.Column

	ft.AddRaw(`<label class="label otot-header" style="font-weight: unset; font-size: larger">`+title+`</label>`, colspan, 1)

	ft.AddRaw(`<hr class="otot-header-line" style="margin: unset">`, colspan, 1).SetStyle("padding", "4px 16px")
}

func (ft *FlowTable) AddFoot(HTML string) {
	colspan := ft.Column

	ft.AddRaw(`<hr class="otot-foot-line" style="margin: unset">`, colspan, 1).SetStyle("padding", "4px 16px")

	item := TD{
		HTML:    HTML,
		colspan: colspan,
		rowspan: 1,
	}
	ft.Items = append(ft.Items, &item)
}

func (ft *FlowTable) AddRaw(HTML string, colspan int, rowspan int) *TD {
	item := TD{
		HTML:    HTML,
		colspan: colspan,
		rowspan: rowspan,
		styleMap: map[string]string{
			"vertical-align": "middle",
		},
	}
	ft.Items = append(ft.Items, &item)
	return &item
}

// AddDiv : <div> YOUR_HTML </div>
func (ft *FlowTable) AddDiv(HTML string, colspan int, rowspan int) *TD {
	item := TD{
		HTML:    "<div>" + HTML + "</div>",
		colspan: colspan,
		rowspan: rowspan,
		styleMap: map[string]string{
			"vertical-align": "middle",
		},
	}
	ft.Items = append(ft.Items, &item)
	return &item
}

// AddSpan : <span> YOUR_HTML </span>
func (ft *FlowTable) AddSpan(HTML string, colspan int, rowspan int) *TD {
	item := TD{
		HTML:    "<span>" + HTML + "</span>",
		colspan: colspan,
		rowspan: rowspan,
		styleMap: map[string]string{
			"vertical-align": "middle",
		},
	}
	ft.Items = append(ft.Items, &item)
	return &item
}

// AddTitle : shortcut, AddTitle Bold font
func (ft *FlowTable) AddTitleB(Title string) *TD {
	item := TD{
		HTML:    fmt.Sprintf(`<p class="is-size-5" style="font-weight: bold;">%s</p>`, Title),
		colspan: 1,
		rowspan: 1,
		styleMap: map[string]string{
			"vertical-align": "middle",
		},
	}
	ft.Items = append(ft.Items, &item)
	return &item
}

// AddLabel : <label> YOUR_HTML </label>
func (ft *FlowTable) AddLabel(Label string) *TD {
	item := TD{
		HTML:    fmt.Sprintf(`<label class="label" style="font-size:unset">%s</label>`, Label),
		colspan: 1,
		rowspan: 1,
		styleMap: map[string]string{
			"vertical-align": "middle",
			"white-space":    "nowrap",
		},
	}
	ft.Items = append(ft.Items, &item)
	return &item
}

func (ft *FlowTable) AddTitle(Label string) *TD {
	item := TD{
		HTML:    fmt.Sprintf(`<label class="label" style="font-size:unset">%s</label>`, Label),
		colspan: 1,
		rowspan: 1,
		styleMap: map[string]string{
			"vertical-align": "middle",
			"white-space":    "nowrap",
		},
	}
	ft.Items = append(ft.Items, &item)
	return &item
}

// AddInput : shortcut
func (ft *FlowTable) AddInput(Model string, colspan int, others ...string) *TD {
	xxx := strings.Join(others, " ")
	item := TD{
		HTML:    fmt.Sprintf(`<input id="%s" name="%s" v-model="%s" class="input is-fullwidth otot-input" %s>`, ElementID(), Model, Model, xxx),
		colspan: colspan,
		rowspan: 1,
		styleMap: map[string]string{
			"vertical-align": "middle",
		},
	}
	ft.Items = append(ft.Items, &item)
	ft.VueVars = append(ft.VueVars, Model)
	return &item
}

// AddText : shortcut
func (ft *FlowTable) AddText(Model string, colspan int, others ...string) *TD {
	xxx := strings.Join(others, " ")
	item := TD{
		HTML:    fmt.Sprintf(`<textarea class="textarea is-fullwidth otot-textarea" id="%s" v-model="%s" %s></textarea>`, ElementID(), Model, xxx),
		colspan: colspan,
		rowspan: 1,
		styleMap: map[string]string{
			"vertical-align": "middle",
		},
	}
	ft.Items = append(ft.Items, &item)
	ft.VueVars = append(ft.VueVars, Model)
	return &item
}

// AddSelect : shortcut
func (ft *FlowTable) AddSelect(Model string, kv ...string) *TD {
	var Lines []string

	sp := fmt.Sprintf

	Lines = append(Lines, `<div class="select is-fullwidth" id="`+ElementID()+`" style="margin:0;padding:0;">`)
	Lines = append(Lines, `  <select v-model="`+Model+`" class="select" id="`+ElementID()+`">`)
	for i := 0; i < len(kv)/2; i++ {
		s := sp(`    <option value="%s">%s</option>`, kv[2*i], kv[2*i+1])
		Lines = append(Lines, s)
	}
	Lines = append(Lines, `  </select>`)
	Lines = append(Lines, `</div>`)

	item := TD{
		HTML:    strings.Join(Lines, "\n"),
		colspan: 1,
		rowspan: 1,
		styleMap: map[string]string{
			"vertical-align": "middle",
		},
	}
	ft.Items = append(ft.Items, &item)
	ft.VueVars = append(ft.VueVars, Model)
	return &item
}

// AddDate : shortcut
func (ft *FlowTable) AddDate(Model string, minDate string, maxDate string) *TD {
	item := TD{
		HTML: fmt.Sprintf(
			`<input id="%s" name="%s" v-model="%s" data-allow-input="true" data-min-date="%s" data-max-date="%s" data-default-date="" class="flatpickr input flatpickr-input active otot-date">`,
			ElementID(),
			Model,
			Model,
			minDate,
			maxDate,
		),
		colspan: 1,
		rowspan: 1,
		styleMap: map[string]string{
			"vertical-align": "middle",
		},
	}
	ft.Items = append(ft.Items, &item)
	ft.VueVars = append(ft.VueVars, Model)
	return &item
}

// AddAddress : Must set Model and Value
func (ft *FlowTable) AddAddress(mProvince string, mCity string, mDistrict string, Address string, vProvince string, vCity string, vDistrict string) *TD {
	var Lines []string

	sp := fmt.Sprintf

	Lines = append(Lines, `<div class="distpicker " data-toggle="distpicker" id="`+ElementID()+`">`)

	if mProvince != "" {
		Lines = append(Lines, `<span class="select">`)
		Lines = append(Lines, sp(`<select v-model="%s" data-province="%s"></select>`, mProvince, vProvince))
		Lines = append(Lines, `</span>`)
		ft.VueVars = append(ft.VueVars, mProvince)
	}

	if mCity != "" {
		Lines = append(Lines, `<span class="select">`)
		Lines = append(Lines, sp(`<select v-model="%s" data-city="%s"></select>`, mCity, vCity))
		Lines = append(Lines, `</span>`)
		ft.VueVars = append(ft.VueVars, mCity)
	}

	if mDistrict != "" {
		Lines = append(Lines, `<span class="select">`)
		Lines = append(Lines, sp(`<select v-model="%s" data-district="%s"></select>`, mDistrict, vDistrict))
		Lines = append(Lines, `</span>`)
		ft.VueVars = append(ft.VueVars, mDistrict)
	}

	if Address != "" {
		Lines = append(Lines, sp(`<input name="%s" v-model="%s" class="input">`, Address, Address))
		ft.VueVars = append(ft.VueVars, Address)
	}

	Lines = append(Lines, `</div>`)

	item := TD{
		HTML:    strings.Join(Lines, "\n"),
		colspan: 1,
		rowspan: 1,
		styleMap: map[string]string{
			"vertical-align": "middle",
		},
	}
	ft.Items = append(ft.Items, &item)
	return &item
}

func (ft *FlowTable) Gen(b64 bool) string {
	var lines []string

	lines = append(
		lines,
		fmt.Sprintf(`<table id="%s" class="table is-fullwidth is-narrow otot-table otot-table-%d %s">`,
			ft.ID,
			ft.Column,
			ft.Class,
		),
	)
	lines = append(lines, "<tbody>")

	cols := 0
	rows := 0
	lines = append(
		lines,
		fmt.Sprintf(
			`<tr class="otot-tr-nth-%d">`,
			rows,
		),
	)
	for i := 0; i < len(ft.Items); i++ {
		s := ft.Items[i]
		HTML := s.HTML
		colspan := s.colspan
		rowspan := s.rowspan

		if cols+colspan > ft.Column {
			cols = 0
			rows += 1
			lines = append(lines, "</tr>")
			lines = append(
				lines,
				fmt.Sprintf(
					`<tr class="otot-tr-nth-%d">`,
					rows,
				),
			)
		}

		var style string
		for k, v := range s.styleMap {
			style += fmt.Sprintf("%s:%s;", k, v)
		}

		var class string
		for k, _ := range s.class {
			class += k
			class += " "
		}

		line := fmt.Sprintf(
			`<td rowspan="%d" colspan="%d" class="%s %s %s %s %s" style="%s">%s</td>`,
			rowspan,
			colspan,
			fmt.Sprintf("otot-td-rowspan-%d", rowspan),
			fmt.Sprintf("otot-td-colspan-%d", colspan),
			fmt.Sprintf("otot-cell-%d-%d", rows, cols),
			fmt.Sprintf("otot-td-nth-%d", cols),
			class,
			style,
			HTML,
		)
		lines = append(lines, line)
		cols += s.colspan
	}

	lines = append(lines, "</tr>")
	lines = append(lines, "</tbody>")
	lines = append(lines, "</table>")

	html := strings.Join(lines, "\n")

	if b64 {
		html := url.QueryEscape(html)
		html = strings.Replace(html, "+", "%20", -1)
		return base64.StdEncoding.EncodeToString([]byte(html))
	} else {
		return html
	}
}
