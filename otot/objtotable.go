package otot

import (
	"fmt"
	"strings"

	"miego/atox"

	"github.com/fatih/structs"
)

func Convert(key string, val string, obj interface{}) string {
	m := structs.Map(obj)

	var lines []string
	lines = append(lines, fmt.Sprintf("|%s|%s|", key, val))
	lines = append(lines, "|---|---|")

	for k, v := range m {
		lines = append(lines, fmt.Sprintf("|%s|%s|", k, fmt.Sprint(v)))
	}

	return strings.Join(lines, "\n")
}

func ViaKV(class string, id string, column int, v ...string) string {
	var lines []string

	table := fmt.Sprintf(`<table id="%s" class="table is-fullwidth %s">`, id, class)
	lines = append(lines, table)
	lines = append(lines, "<tbody>")

	for colidx, i := 0, 0; i < len(v)/2; i++ {
		if colidx%column == 0 {
			if colidx != 0 {
				lines = append(lines, "</tr>")
			}
			lines = append(lines, "<tr>")
		}

		// XXX: Tricky, If the title part starts with @N@, N is a number
		// this will be set to colspan = N
		colspan := "1"
		rowspan := "1"
		k := v[2*i]
		if len(k) > 2 && k[0] == '@' && k[2] == '@' {
			colspan = k[1:2]
			k = k[3:]
		}
		v := v[2*i+1]
		if len(v) > 2 && v[0] == '@' && v[2] == '@' {
			rowspan = v[1:2]
			v = v[3:]
		}
		td := fmt.Sprintf(
			`<td colspan="%s" rowspan="%s"><span class="ototKey">%s</span><span class="ototVal">%s</span></td>`,
			colspan, rowspan,
			k, v,
		)
		lines = append(lines, td)
		colidx += atox.Int(colspan, 1)
	}

	lines = append(lines, "</tr>")
	lines = append(lines, "</tbody>")
	lines = append(lines, "</table>")

	return strings.Join(lines, "\n")
}

func ViaS(class string, id string, column int, v ...string) string {
	var lines []string

	lines = append(lines, fmt.Sprintf(`<table id="%s" class="table is-fullwidth %s">`, id, class))
	lines = append(lines, "<tbody>")
	for i := 0; i < len(v); i++ {
		if i%column == 0 {
			if i != 0 {
				lines = append(lines, "</tr>")
			}
			lines = append(lines, "<tr>")
		}
		s := v[i]
		lines = append(lines, s)
	}

	lines = append(lines, "</tr>")
	lines = append(lines, "</tbody>")
	lines = append(lines, "</table>")

	return strings.Join(lines, "\n")
}
