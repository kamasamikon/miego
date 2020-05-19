package otot

import (
	"fmt"
	"strings"
)

type FlowTable struct {
	ID     string
	Class  string
	Column int
	Items  []string
}

func FlowTableNew(ID string, Class string, Column int) *FlowTable {
	return &FlowTable{
		ID:     ID,
		Class:  Class,
		Column: Column,
	}
}

func (ft *FlowTable) SetClass(Class string) *FlowTable {
	ft.Class = Class
	return ft
}
func (ft *FlowTable) SetID(ID string) *FlowTable {
	ft.ID = ID
	return ft
}
func (ft *FlowTable) SetColumn(Column int) *FlowTable {
	ft.Column = Column
	return ft
}
func (ft *FlowTable) SetItem(HTML string) *FlowTable {
	ft.Items = append(ft.Items, HTML)
	return ft
}
func (ft *FlowTable) Gen() string {
	var lines []string

	lines = append(lines, fmt.Sprintf(`<table id="%s" class="%s table is-fullwidth">`, ft.ID, ft.Class))
	lines = append(lines, "<tbody>")
	for i := 0; i < len(ft.Items); i++ {
		if i%ft.Column == 0 {
			if i != 0 {
				lines = append(lines, "</tr>")
			}
			lines = append(lines, "<tr>")
		}
		s := ft.Items[i]
		lines = append(lines, s)
	}

	lines = append(lines, "</tr>")
	lines = append(lines, "</tbody>")
	lines = append(lines, "</table>")

	return strings.Join(lines, "\n")
}
