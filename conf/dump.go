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

	cc.mutex.Lock()
	for p, e := range cc.mapPathEntry {
		cList = append(cList, e)
		if len(p) > keyMaxLength {
			keyMaxLength = len(p)
		}
	}
	cc.mutex.Unlock()

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
			lines = append(lines, fmt.Sprintf(fmtstr, e.path, e.refGet, e.refSet, vInt))

		case 's':
			vStr := e.vStr
			lines = append(lines, fmt.Sprintf(fmtstr, e.path, e.refGet, e.refSet, vStr))

		case 'b':
			vBool := e.vBool
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

// Dump : Print all entries
func (cc *ConfCenter) DumpJson(safeMode bool) map[string]string {
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
		if e.hidden && safeMode {
			continue
		}

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
func (cc *ConfCenter) DumpRaw(safeMode bool, group bool, joinBy string) string {
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
