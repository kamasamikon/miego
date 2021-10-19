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

func (ft *FlowTable) AddOne(HTML string) *TD {
	item := TD{
		HTML:    HTML,
		colspan: 1,
		rowspan: 1,
		styleMap: map[string]string{
			"text-align":     "left",
			"vertical-align": "middle",
		},
	}
	ft.Items = append(ft.Items, &item)
	return &item
}
func (ft *FlowTable) AddSpan(HTML string, colspan int, rowspan int) *TD {
	item := TD{
		HTML:    HTML,
		colspan: colspan,
		rowspan: rowspan,
		styleMap: map[string]string{
			"text-align":     "left",
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
		colspan: ft.Column,
		rowspan: 1,
		styleMap: map[string]string{
			"text-align":     "left",
			"vertical-align": "middle",
		},
	}
	ft.Items = append(ft.Items, &item)
	return &item
}

// AddLabel : shortcut
func (ft *FlowTable) AddLabel(Label string) *TD {
	item := TD{
		HTML:    fmt.Sprintf(`<label class="label" style="font-size:unset">%s</label>`, Label),
		colspan: 1,
		rowspan: 1,
		styleMap: map[string]string{
			"text-align":     "right",
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
			"text-align":     "left",
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
		HTML:    fmt.Sprintf(`<input id="%s" v-model="%s" class="input is-fullwidth" %s>`, ElementID(), Model, xxx),
		colspan: colspan,
		rowspan: 1,
		styleMap: map[string]string{
			"text-align":     "left",
			"vertical-align": "middle",
		},
	}
	ft.Items = append(ft.Items, &item)
	return &item
}

// AddText : shortcut
func (ft *FlowTable) AddText(Model string, colspan int, others ...string) *TD {
	xxx := strings.Join(others, " ")
	item := TD{
		HTML:    fmt.Sprintf(`<textarea class="textarea is-fullwidth" id="%s" v-model="%s" %s></textarea>`, ElementID(), Model, xxx),
		colspan: colspan,
		rowspan: 1,
		styleMap: map[string]string{
			"text-align":     "left",
			"vertical-align": "middle",
		},
	}
	ft.Items = append(ft.Items, &item)
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
			"text-align":     "left",
			"vertical-align": "middle",
		},
	}
	ft.Items = append(ft.Items, &item)
	return &item
}

// AddDate : shortcut
func (ft *FlowTable) AddDate(Model string, minDate string, maxDate string) *TD {
	item := TD{
		HTML: fmt.Sprintf(
			`<input id="%s" v-model="%s" data-allow-input="true" data-min-date="%s" data-max-date="%s" data-default-date="" class="flatpickr input flatpickr-input active">`,
			ElementID(),
			Model,
			minDate,
			maxDate,
		),
		colspan: 1,
		rowspan: 1,
		styleMap: map[string]string{
			"text-align":     "left",
			"vertical-align": "middle",
		},
	}
	ft.Items = append(ft.Items, &item)
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
	}

	if mCity != "" {
		Lines = append(Lines, `<span class="select">`)
		Lines = append(Lines, sp(`<select v-model="%s" data-city="%s"></select>`, mCity, vCity))
		Lines = append(Lines, `</span>`)
	}

	if mDistrict != "" {
		Lines = append(Lines, `<span class="select">`)
		Lines = append(Lines, sp(`<select v-model="%s" data-district="%s"></select>`, mDistrict, vDistrict))
		Lines = append(Lines, `</span>`)
	}

	if Address != "" {
		Lines = append(Lines, sp(`<input v-model="%s" class="input">`, Address))
	}

	Lines = append(Lines, `</div>`)

	item := TD{
		HTML:    strings.Join(Lines, "\n"),
		colspan: 1,
		rowspan: 1,
		styleMap: map[string]string{
			"text-align":     "left",
			"vertical-align": "middle",
		},
	}
	ft.Items = append(ft.Items, &item)
	return &item
}

func (ft *FlowTable) Gen(b64 bool) string {
	var lines []string

	lines = append(lines, fmt.Sprintf(`<table id="%s" class="table is-fullwidth is-narrow %s">`, ft.ID, ft.Class))
	lines = append(lines, "<tbody>")

	cols := 0
	lines = append(lines, "<tr>")
	for i := 0; i < len(ft.Items); i++ {
		s := ft.Items[i]
		HTML := s.HTML
		colspan := s.colspan
		rowspan := s.rowspan

		if cols+colspan > ft.Column {
			cols = 0
			lines = append(lines, "</tr>")
			lines = append(lines, "<tr>")
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
			`<td rowspan="%d" colspan="%d" class="%s" style="%s">%s</td>`,
			rowspan, colspan, class, style, HTML,
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
