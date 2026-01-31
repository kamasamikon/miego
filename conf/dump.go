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

	for _, item := range cc.iItems {
		lines = append(lines, fmt.Sprintf(fmtstr, "i:/"+item.key, item.value))
	}
	for _, item := range cc.sItems {
		lines = append(lines, fmt.Sprintf(fmtstr, "s:/"+item.key, item.value))
	}
	for _, item := range cc.bItems {
		lines = append(lines, fmt.Sprintf(fmtstr, "b:/"+item.key, item.value))
	}
	for _, item := range cc.eItems {
		lines = append(lines, fmt.Sprintf(fmtstr, "e:/"+item.key, "..."))
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
	dict := make(map[string]string)

	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	for _, item := range cc.iItems {
		dict["i:/"+item.key] = fmt.Sprintf("%v", item.value)
	}
	for _, item := range cc.sItems {
		dict["s:/"+item.key] = strings.TrimSpace(item.value)
	}
	for _, item := range cc.bItems {
		dict["b:/"+item.key] = fmt.Sprintf("%v", item.value)
	}
	for _, item := range cc.eItems {
		dict["e:/"+item.key] = "..."
	}

	return dict
}

// DumpRaw : Dump without Get/Get refs
func (cc *ConfCenter) DumpRaw(joinBy string) string {
	dict := cc.DumpMap()

	var keys []string
	for key := range dict {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i int, j int) bool {
		return strings.Compare(keys[i][1:], keys[j][1:]) < 0
	})

	var lines []string
	for _, key := range keys {
		val := dict[key]

		if strings.IndexByte(val, '\n') >= 0 {
			lines = append(lines, fmt.Sprintf("s:/%s=<<EOF", key))
			lines = append(lines, fmt.Sprintf("%s", val))
			lines = append(lines, fmt.Sprintf("%s", "EOF"))
		} else {
			lines = append(lines, fmt.Sprintf("s:/%s=%s", key, val))
		}
	}

	// Add the last \n
	lines = append(lines, "")

	return strings.Join(lines, joinBy)
}
