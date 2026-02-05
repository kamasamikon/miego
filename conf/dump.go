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
		"94", keyMaxLength+4, "7;93",
	)

	for key, item := range cc.iItems {
		lines = append(lines, fmt.Sprintf(fmtstr, "i:/"+key, item.value))
	}
	for key, item := range cc.sItems {
		lines = append(lines, fmt.Sprintf(fmtstr, "s:/"+key, item.value))
	}
	for key, item := range cc.bItems {
		lines = append(lines, fmt.Sprintf(fmtstr, "b:/"+key, item.value))
	}
	for key := range cc.eItems {
		lines = append(lines, fmt.Sprintf(fmtstr, "e:/"+key, "..."))
	}
	for key, item := range cc.xItems {
		if item.getter != nil {
			lines = append(lines, key, fmt.Sprintf("%v", item.getter(key)))
		} else {
			lines = append(lines, key, "...")
		}
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

	for key, item := range cc.iItems {
		dict["i:/"+key] = fmt.Sprintf("%v", item.value)
	}
	for key, item := range cc.sItems {
		dict["s:/"+key] = strings.TrimSpace(item.value)
	}
	for key, item := range cc.bItems {
		dict["b:/"+key] = fmt.Sprintf("%v", item.value)
	}
	for key := range cc.eItems {
		dict["e:/"+key] = "..."
	}
	for key, item := range cc.xItems {
		if item.getter != nil {
			dict["e:/"+key] = fmt.Sprintf("%v", item.getter(key))
		} else {
			dict["e:/"+key] = "..."
		}
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
			lines = append(lines, fmt.Sprintf("%s=<<EOF\n%s\nEOF", key, val))
		} else {
			lines = append(lines, fmt.Sprintf("%s=%s", key, val))
		}
	}

	// Add the last \n
	lines = append(lines, "")

	return strings.Join(lines, joinBy)
}
