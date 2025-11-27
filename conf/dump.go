package conf

import (
	"fmt"
	"sort"
	"strings"
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
		"%s%%-%ds%s :(%%04d:%%04d): %s%%v%s",
		ColorType_D,
		keyMaxLength,
		ColorType_Reset,
		ColorType_W, ColorType_Reset,
	)
	var lines []string
	for _, v := range cList {
		if v.hidden && safeMode {
			continue
		}

		switch v.kind {
		case 'i':
			lines = append(lines, fmt.Sprintf(fmtstr, v.path, v.refGet, v.refSet, v.vInt))

		case 's':
			lines = append(lines, fmt.Sprintf(fmtstr, v.path, v.refGet, v.refSet, v.vStr))

		case 'b':
			lines = append(lines, fmt.Sprintf(fmtstr, v.path, v.refGet, v.refSet, v.vBool))

		case 'o':
			lines = append(lines, fmt.Sprintf(fmtstr, v.path, v.refGet, v.refSet, "..."))

		case 'e':
			lines = append(lines, fmt.Sprintf(fmtstr, v.path, v.refGet, v.refSet, "..."))
		}

	}

	// Add the last \n
	lines = append(lines, "")

	return strings.Join(lines, "\n")
}

// DumpRaw : Dump without Get/Get refs
func DumpRaw(safeMode bool, group bool) string {
	var cList []*confEntry
	for _, v := range mapPathEntry {
		cList = append(cList, v)
	}
	sort.Slice(cList, func(i int, j int) bool {
		return strings.Compare(cList[i].path[1:], cList[j].path[1:]) < 0
	})

	// Group Name
	var groupName string

	var lines []string
	var lastLine string
	for _, v := range cList {
		if v.hidden && safeMode {
			continue
		}

		if group {
			segs := strings.SplitN(v.path, "/", 3)
			if segs[1] != groupName {
				if lastLine != "" {
					lastLine = ""
					lines = append(lines, lastLine)
				}
				groupName = segs[1]
			}
		}

		switch v.kind {
		case 'i':
			lastLine = fmt.Sprintf("%s=%d", v.path, v.vInt)
			lines = append(lines, lastLine)

		case 's':
			lastLine = fmt.Sprintf("%s=%s", v.path, v.vStr)
			lines = append(lines, lastLine)

		case 'b':
			lastLine = fmt.Sprintf("%s=%t", v.path, v.vBool)
			lines = append(lines, lastLine)

		case 'o':
			lastLine = fmt.Sprintf("%s=%s", v.path, "...")
			lines = append(lines, lastLine)
		}
	}

	// Add the last \n
	if lastLine != "" {
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}
