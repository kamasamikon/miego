package conf

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kamasamikon/miego/klog"
)

// Dump : Print all entries
func Dump(safeMode bool) string {
	keyMaxLength := 0
	var cList []*confEntry

	for k, v := range mapPathEntry {
		cList = append(cList, v)
		if len(k) > keyMaxLength {
			keyMaxLength = len(k)
		}
	}

	sort.Slice(cList, func(i int, j int) bool {
		return strings.Compare(cList[i].path[1:], cList[j].path[1:]) < 0
	})

	fmtstr := fmt.Sprintf(
		"%s%%-%ds%s : %s%%v%s",
		klog.ColorType_D,
		keyMaxLength,
		klog.ColorType_Reset,
		klog.ColorType_W, klog.ColorType_Reset,
	)
	var lines []string
	for _, v := range cList {
		if v.hidden && safeMode {
			continue
		}

		switch v.kind {
		case 'i':
			lines = append(lines, fmt.Sprintf(fmtstr, v.path, v.vInt))

		case 's':
			lines = append(lines, fmt.Sprintf(fmtstr, v.path, v.vStr))

		case 'b':
			lines = append(lines, fmt.Sprintf(fmtstr, v.path, v.vBool))

		case 'o':
			lines = append(lines, fmt.Sprintf(fmtstr, v.path, "..."))

		case 'e':
			lines = append(lines, fmt.Sprintf(fmtstr, v.path, "..."))
		}

	}

	// Add the last \n
	lines = append(lines, "")

	return strings.Join(lines, "\n")
}

// DumpRaw : Dump without Get/Get refs
func DumpRaw(safeMode bool) string {
	var cList []*confEntry
	for _, v := range mapPathEntry {
		cList = append(cList, v)
	}
	sort.Slice(cList, func(i int, j int) bool {
		return strings.Compare(cList[i].path[1:], cList[j].path[1:]) < 0
	})

	var lines []string
	for _, v := range cList {
		if v.hidden && safeMode {
			continue
		}

		switch v.kind {
		case 'i':
			lines = append(lines, fmt.Sprintf("%s=%d", v.path, v.vInt))

		case 's':
			lines = append(lines, fmt.Sprintf("%s=%s", v.path, v.vStr))

		case 'b':
			lines = append(lines, fmt.Sprintf("%s=%t", v.path, v.vBool))

		case 'o':
			lines = append(lines, fmt.Sprintf("%s=%s", v.path, "..."))
		}
	}

	// Add the last \n
	lines = append(lines, "")

	return strings.Join(lines, "\n")
}
