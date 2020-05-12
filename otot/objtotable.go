package otot

import (
	"fmt"
	"strings"

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

	lines = append(lines, fmt.Sprintf(`<table id="%s" class="%s table is-fullwidth">`, id, class))
	lines = append(lines, "<tbody>")
	for i := 0; i < len(v)/2; i++ {
		if i%column == 0 {
			if i != 0 {
				lines = append(lines, "</tr>")
			}
			lines = append(lines, "<tr>")
		}
		k := v[2*i]
		v := v[2*i+1]
		lines = append(lines, fmt.Sprintf(`<td><span class="ototKey">%s</span><span class="ototVal">%s</span></td>`, k, v))
	}

	lines = append(lines, "</tr>")
	lines = append(lines, "</tbody>")
	lines = append(lines, "</table>")

	return strings.Join(lines, "\n")
}

func ViaS(class string, id string, column int, v ...string) string {
	var lines []string

	lines = append(lines, fmt.Sprintf(`<table id="%s" class="%s table is-fullwidth">`, id, class))
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
