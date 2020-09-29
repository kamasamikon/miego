package otot

import (
	"fmt"
	"strings"
)

type FlowTableItem struct {
	HTML     string
	colspan  int
	rowspan  int
	class    string
	styleMap map[string]string
}

func (i *FlowTableItem) SetHTML(x string) *FlowTableItem {
	i.HTML = x
	return i
}
func (i *FlowTableItem) SetColSpan(x int) *FlowTableItem {
	i.colspan = x
	return i
}
func (i *FlowTableItem) SetRowSpan(x int) *FlowTableItem {
	i.rowspan = x
	return i
}
func (i *FlowTableItem) SetClass(x string) *FlowTableItem {
	i.class = x
	return i
}

func (i *FlowTableItem) SetStyle(k string, v string) *FlowTableItem {
	i.styleMap[k] = v
	return i
}
func (i *FlowTableItem) RemStyle(k string) *FlowTableItem {
	delete(i.styleMap, k)
	return i
}

type FlowTable struct {
	ID     string
	Class  string
	Column int
	Items  []*FlowTableItem
}

func FlowTableNew(ID string, Class string, Column int) *FlowTable {
	return &FlowTable{
		ID:     ID,
		Class:  Class,
		Column: Column,
	}
}

func (ft *FlowTable) AddOne(HTML string) *FlowTableItem {
	item := FlowTableItem{
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
func (ft *FlowTable) AddSpan(HTML string, colspan int, rowspan int) *FlowTableItem {
	item := FlowTableItem{
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
func (ft *FlowTable) AddTitleB(Title string) *FlowTableItem {
	item := FlowTableItem{
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
func (ft *FlowTable) AddLabel(Label string) *FlowTableItem {
	item := FlowTableItem{
		HTML:    fmt.Sprintf(`<label class="label">%s</label>`, Label),
		colspan: 1,
		rowspan: 1,
		styleMap: map[string]string{
			"text-align":     "right",
			"vertical-align": "middle",
		},
	}
	ft.Items = append(ft.Items, &item)
	return &item
}

func (ft *FlowTable) AddTitle(Label string) *FlowTableItem {
	item := FlowTableItem{
		HTML:    fmt.Sprintf(`<label class="label">%s</label>`, Label),
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

// AddInput : shortcut
func (ft *FlowTable) AddInput(Model string, colspan int) *FlowTableItem {
	item := FlowTableItem{
		HTML:    fmt.Sprintf(`<input v-model="%s" class="input">`, Model),
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
func (ft *FlowTable) AddSelect(Model string, kv ...string) *FlowTableItem {
	var Lines []string

	sp := fmt.Sprintf

	Lines = append(Lines, `<div class="select">`)
	Lines = append(Lines, `  <select v-model="`+Model+`" class="select">`)
	for i := 0; i < len(kv)/2; i++ {
		s := sp(`    <option value="%s">%s</option>`, kv[2*i], kv[2*i+1])
		Lines = append(Lines, s)
	}
	Lines = append(Lines, `  </select>`)
	Lines = append(Lines, `</div>`)

	item := FlowTableItem{
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
func (ft *FlowTable) AddDate(Model string, minDate string, maxDate string) *FlowTableItem {
	item := FlowTableItem{
		HTML:    fmt.Sprintf(`<input v-model.trim="%s" data-min-date="%s" data-max-date="%s" data-default-date="" class="flatpickr input flatpickr-input active">`, Model, minDate, maxDate),
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
func (ft *FlowTable) AddAddress(mProvince string, mCity string, mDistrict string, Address string, vProvince string, vCity string, vDistrict string) *FlowTableItem {
	var Lines []string

	sp := fmt.Sprintf

	Lines = append(Lines, `<div class="distpicker " data-toggle="distpicker">`)

	if mProvince != "" {
		Lines = append(Lines, `<span class="select">`)
		Lines = append(Lines, sp(`<select v-model.trim="%s" data-province="%s"></select>`, mProvince, vProvince))
		Lines = append(Lines, `</span>`)
	}

	if mCity != "" {
		Lines = append(Lines, `<span class="select">`)
		Lines = append(Lines, sp(`<select v-model.trim="%s" data-city="%s"></select>`, mCity, vCity))
		Lines = append(Lines, `</span>`)
	}

	if mDistrict != "" {
		Lines = append(Lines, `<span class="select">`)
		Lines = append(Lines, sp(`<select v-model.trim="%s" data-district="%s"></select>`, mDistrict, vDistrict))
		Lines = append(Lines, `</span>`)
	}

	if Address != "" {
		Lines = append(Lines, sp(`<input v-model.trim="%s" class="input">`, Address))
	}

	Lines = append(Lines, `</div>`)

	item := FlowTableItem{
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

func (ft *FlowTable) Gen() string {
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
		class := s.class
		styleMap := s.styleMap

		if cols+colspan > ft.Column {
			cols = 0
			lines = append(lines, "</tr>")
			lines = append(lines, "<tr>")
		}

		var style string
		for k, v := range styleMap {
			style += fmt.Sprintf("%s:%s;", k, v)
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

	return strings.Join(lines, "\n")
}
