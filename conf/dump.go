package conf

import (
	"fmt"
	"sort"
	"strings"
)

// Dump : Print all entries
func (cc *ConfCenter) Dump(joinBy string) string {
	keyMaxLength := 0

	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	for p := range cc.iItems {
		if len(p) > keyMaxLength {
			keyMaxLength = len(p)
		}
	}
	for p := range cc.sItems {
		if len(p) > keyMaxLength {
			keyMaxLength = len(p)
		}
	}
	for p := range cc.bItems {
		if len(p) > keyMaxLength {
			keyMaxLength = len(p)
		}
	}
	for p := range cc.eItems {
		if len(p) > keyMaxLength {
			keyMaxLength = len(p)
		}
	}

	// FIXME: keyMaxLength应该+3
	fmtstr := fmt.Sprintf(
		"%s%%-%ds%s %s%%v%s",
		ColorTypeD,
		keyMaxLength,
		ColorTypeReset,
		ColorTypeW, ColorTypeReset,
	)

	var lines []string

	for _, e := range cc.iItems {
		lines = append(lines, fmt.Sprintf(fmtstr, "i:/"+e.key, e.value))
	}
	for _, e := range cc.sItems {
		lines = append(lines, fmt.Sprintf(fmtstr, "s:/"+e.key, e.value))
	}
	for _, e := range cc.bItems {
		lines = append(lines, fmt.Sprintf(fmtstr, "b:/"+e.key, e.value))
	}
	for _, e := range cc.eItems {
		lines = append(lines, fmt.Sprintf(fmtstr, "e:/"+e.key, "..."))
	}

	sort.Slice(lines, func(i int, j int) bool {
		return strings.Compare(lines[i][1:], lines[j][1:]) < 0
	})

	// Add the last \n
	lines = append(lines, "")

	return strings.Join(lines, joinBy)
}

// Dump : Print all entries
func (cc *ConfCenter) DumpMap() map[string]string {
	keyMaxLength := 0
	var cList []*confEntry

	cc.mutex.Lock()
	for p, e := range cc.mapPathEntry {
		cList = append(cList, e)
		if len(p) > keyMaxLength {
			keyMaxLength = len(p)
		}
	}
	cc.mutex.Unlock()

	var j map[string]string = make(map[string]string)

	for _, e := range cList {
		switch e.kind {
		case 'i':
			vInt := e.vInt
			j[e.path] = fmt.Sprintf("%v", vInt)

		case 's':
			vStr := e.vStr
			j[e.path] = vStr

		case 'b':
			vBool := e.vBool
			j[e.path] = fmt.Sprintf("%v", vBool)

		case 'o':
			j[e.path] = "..."

		case 'e':
			j[e.path] = "..."
		}
	}

	return j
}

// DumpRaw : Dump without Get/Get refs
func (cc *ConfCenter) DumpRaw(group bool, joinBy string) string {
	var cList []*confEntry

	cc.mutex.Lock()
	for _, e := range cc.mapPathEntry {
		cList = append(cList, e)
	}
	cc.mutex.Unlock()

	sort.Slice(cList, func(i int, j int) bool {
		return strings.Compare(cList[i].path[1:], cList[j].path[1:]) < 0
	})

	// Group Name
	var groupName string

	var lines []string
	var lastLine string
	for _, e := range cList {
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
			lastLine = fmt.Sprintf("%s=%d", e.path, vInt)
			lines = append(lines, lastLine)

		case 's':
			vStr := e.vStr
			lastLine = fmt.Sprintf("%s=%s", e.path, vStr)
			lines = append(lines, lastLine)

		case 'b':
			vBool := e.vBool
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
