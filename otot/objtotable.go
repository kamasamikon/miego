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
