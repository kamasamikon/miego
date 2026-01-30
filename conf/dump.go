package conf

import (
	"fmt"
	"sort"
	"strings"
)

// Dump : Print all entries
func (cc *ConfCenter) Dump(joinBy string) string {
	keyMaxLength := 0
	var lines []string

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

	fmtstr := fmt.Sprintf(
		"\033[%sm%%-%ds\033[0m\033[%sm%%v\033[0m",
		"31", keyMaxLength+3, "34",
	)

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
	var dict map[string]string = make(map[string]string)

	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	for _, e := range cc.iItems {
		dict["i:/"+e.key] = fmt.Sprintf("%v", e.value)
	}
	for _, e := range cc.sItems {
		dict["s:/"+e.key] = e.value
	}
	for _, e := range cc.bItems {
		dict["b:/"+e.key] = fmt.Sprintf("%v", e.value)
	}
	for _, e := range cc.eItems {
		dict["e:/"+e.key] = "..."
	}

	return dict
}

// DumpRaw : Dump without Get/Get refs
func (cc *ConfCenter) DumpRaw(joinBy string) string {
	var lines []string

	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	fmtstr := "%s=%v"

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
