package conf

import (
	"fmt"
	"sort"
	"strings"
)

// Dump : Print all entries
func (cc *ConfCenter) Dump(safeMode bool, joinBy string) string {
	keyMaxLength := 0
	var cList []*confEntry

	for p, e := range cc.mapPathEntry {
		cList = append(cList, e)
		if len(p) > keyMaxLength {
			keyMaxLength = len(p)
		}
	}

	sort.Slice(cList, func(i int, j int) bool {
		return strings.Compare(cList[i].path[1:], cList[j].path[1:]) < 0
	})

	fmtstr := fmt.Sprintf(
		"%s%%-%ds%s :(%%04d:%%04d): %s%%v%s",
		ColorTypeD,
		keyMaxLength,
		ColorTypeReset,
		ColorTypeW, ColorTypeReset,
	)
	var lines []string
	for _, e := range cList {
		if e.hidden && safeMode {
			continue
		}

		switch e.kind {
		case 'i':
			vInt := e.vInt
			if e.getter != nil {
				if vv, ok := e.getter(e.path); ok {
					vInt = vv.(int64)
				}
			}
			lines = append(lines, fmt.Sprintf(fmtstr, e.path, e.refGet, e.refSet, vInt))

		case 's':
			vStr := e.vStr
			if e.getter != nil {
				if vv, ok := e.getter(e.path); ok {
					vStr = vv.(string)
				}
			}
			lines = append(lines, fmt.Sprintf(fmtstr, e.path, e.refGet, e.refSet, vStr))

		case 'b':
			vBool := e.vBool
			if e.getter != nil {
				if vv, ok := e.getter(e.path); ok {
					vBool = vv.(bool)
				}
			}
			lines = append(lines, fmt.Sprintf(fmtstr, e.path, e.refGet, e.refSet, vBool))

		case 'o':
			lines = append(lines, fmt.Sprintf(fmtstr, e.path, e.refGet, e.refSet, "..."))

		case 'e':
			lines = append(lines, fmt.Sprintf(fmtstr, e.path, e.refGet, e.refSet, "..."))
		}
	}

	// Add the last \n
	lines = append(lines, "")

	return strings.Join(lines, joinBy)
}

// DumpRaw : Dump without Get/Get refs
func (cc *ConfCenter) DumpRaw(safeMode bool, group bool, joinBy string) string {
	var cList []*confEntry
	for _, e := range cc.mapPathEntry {
		cList = append(cList, e)
	}
	sort.Slice(cList, func(i int, j int) bool {
		return strings.Compare(cList[i].path[1:], cList[j].path[1:]) < 0
	})

	// Group Name
	var groupName string

	var lines []string
	var lastLine string
	for _, e := range cList {
		if e.hidden && safeMode {
			continue
		}

		if group {
			segs := strings.SplitN(e.path, "/", 3)
			if segs[1] != groupName {
				if lastLine != "" {
					lastLine = ""
					lines = append(lines, lastLine)
				}
				groupName = segs[1]
			}
		}

		switch e.kind {
		case 'i':
			vInt := e.vInt
			if e.getter != nil {
				if vv, ok := e.getter(e.path); ok {
					vInt = vv.(int64)
				}
			}
			lastLine = fmt.Sprintf("%s=%d", e.path, vInt)
			lines = append(lines, lastLine)

		case 's':
			vStr := e.vStr
			if e.getter != nil {
				if vv, ok := e.getter(e.path); ok {
					vStr = vv.(string)
				}
			}
			lastLine = fmt.Sprintf("%s=%s", e.path, vStr)
			lines = append(lines, lastLine)

		case 'b':
			vBool := e.vBool
			if e.getter != nil {
				if vv, ok := e.getter(e.path); ok {
					vBool = vv.(bool)
				}
			}
			lastLine = fmt.Sprintf("%s=%t", e.path, vBool)
			lines = append(lines, lastLine)

		case 'o':
			lastLine = fmt.Sprintf("%s=%s", e.path, "...")
			lines = append(lines, lastLine)
		}
	}

	// Add the last \n
	if lastLine != "" {
		lines = append(lines, "")
	}

	return strings.Join(lines, joinBy)
}
