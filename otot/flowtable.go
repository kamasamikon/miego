package otot

import (
	"fmt"
	"strings"
)

type FlowTableItem struct {
	HTML    string
	colspan int
	rowspan int
	class   string
	style   string
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
func (i *FlowTableItem) SetStyle(x string) *FlowTableItem {
	i.style = x
	return i
}

type FlowTable struct {
	ID     string
	Class  string
	Column int
	Items  []FlowTableItem
}

func FlowTableNew(ID string, Class string, Column int) *FlowTable {
	return &FlowTable{
		ID:     ID,
		Class:  Class,
		Column: Column,
	}
}

func (ft *FlowTable) AddOne(HTML string) *FlowTableItem {
	ft.Items = append(ft.Items, FlowTableItem{
		HTML:    HTML,
		colspan: 1,
		rowspan: 1,
	})
	return &ft.Items[len(ft.Items)-1]
}
func (ft *FlowTable) AddSpan(HTML string, colspan int, rowspan int) *FlowTableItem {
	ft.Items = append(ft.Items, FlowTableItem{
		HTML:    HTML,
		colspan: colspan,
		rowspan: rowspan,
	})
	return &ft.Items[len(ft.Items)-1]
}

// AddLabel : shortcut
func (ft *FlowTable) AddLabel(Label string) *FlowTableItem {
	ft.Items = append(ft.Items, FlowTableItem{
		HTML:    fmt.Sprintf(`<label class="label">%s</label>`, Label),
		colspan: 1,
		rowspan: 1,
		style:   "text-align:right; vertical-align:middle;",
	})
	return &ft.Items[len(ft.Items)-1]
}

// AddInput : shortcut
func (ft *FlowTable) AddInput(Model string, colspan int) *FlowTableItem {
	ft.Items = append(ft.Items, FlowTableItem{
		HTML:    fmt.Sprintf(`<input v-model="%s" class="input">`, Model),
		colspan: colspan,
		rowspan: 1,
		style:   "text-align:left; vertical-align:middle;",
	})
	return &ft.Items[len(ft.Items)-1]
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

	ft.Items = append(ft.Items, FlowTableItem{
		HTML:    strings.Join(Lines, "\n"),
		colspan: 1,
		rowspan: 1,
		style:   "text-align:left; vertical-align:middle;",
	})
	return &ft.Items[len(ft.Items)-1]
}

// AddDate : shortcut
func (ft *FlowTable) AddDate(Model string) *FlowTableItem {
	ft.Items = append(ft.Items, FlowTableItem{
		HTML:    fmt.Sprintf(`<input v-model.trim="%s" data-default-date="" class="flatpickr input flatpickr-input active">`, Model),
		colspan: 1,
		rowspan: 1,
		style:   "text-align:left; vertical-align:middle;",
	})
	return &ft.Items[len(ft.Items)-1]
}

// AddAddress : shortcut
func (ft *FlowTable) AddAddress(Province string, City string, District string, Address string, colspan int) *FlowTableItem {
	var Lines []string

	sp := fmt.Sprintf

	Lines = append(Lines, `<div class="distpicker " data-toggle="distpicker">`)

	if Province != "" {
		Lines = append(Lines, `<span class="select">`)
		Lines = append(Lines, sp(`<select v-model.trim="%s" data-province=""></select>`, Province))
		Lines = append(Lines, `</span>`)
	}

	if Province != "" {
		Lines = append(Lines, `<span class="select">`)
		Lines = append(Lines, sp(`<select v-model.trim="%s" data-city=""></select>`, City))
		Lines = append(Lines, `</span>`)
	}

	if Province != "" {
		Lines = append(Lines, `<span class="select">`)
		Lines = append(Lines, sp(`<select v-model.trim="%s" data-district=""></select>`, District))
		Lines = append(Lines, `</span>`)
	}

	if Province != "" {
		Lines = append(Lines, sp(`<input v-model.trim="%s" class="input">`, Address))
	}

	Lines = append(Lines, `</div>`)

	ft.Items = append(ft.Items, FlowTableItem{
		HTML:    strings.Join(Lines, "\n"),
		colspan: colspan,
		rowspan: 1,
		class:   "text-align:right; vertical-align:middle;",
	})
	return &ft.Items[len(ft.Items)-1]
}

func (ft *FlowTable) Gen() string {
	var lines []string

	lines = append(lines, fmt.Sprintf(`<table id="%s" class="%s table is-fullwidth is-narrow">`, ft.ID, ft.Class))
	lines = append(lines, "<tbody>")

	cols := 0
	lines = append(lines, "<tr>")
	for i := 0; i < len(ft.Items); i++ {
		s := ft.Items[i]
		HTML := s.HTML
		colspan := s.colspan
		rowspan := s.rowspan
		class := s.class
		style := s.style

		if cols+colspan > ft.Column {
			cols = 0
			lines = append(lines, "</tr>")
			lines = append(lines, "<tr>")
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
