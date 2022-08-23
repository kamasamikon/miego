package conf

import (
	"fmt"
	"sort"
	"strings"
)

// ////////////////////////////////////////////////////////////////////////
// Monitor, callback when configure changed
//
const (
	ColorType_I     = "\x1b[1;32;40m"
	ColorType_Reset = "\x1b[0m"
)

// KConfMonitor is a Callback called when wathed entry modified.
type KConfMonitor func(path string, oVal interface{}, nVal interface{})

// map[s:Path]map[KConfMonitor]int
var mapPathMonitorCallback = make(map[string]map[*KConfMonitor]int)

func MonitorAdd(Path string, Callback KConfMonitor) {
	mapMonitorCallback, ok := mapPathMonitorCallback[Path]
	if !ok {
		mapMonitorCallback = make(map[*KConfMonitor]int)
	}
	mapMonitorCallback[&Callback] = 1
	mapPathMonitorCallback[Path] = mapMonitorCallback
}

func MonitorRem(Path string, Callback KConfMonitor) {
	mapMonitorCallback, ok := mapPathMonitorCallback[Path]
	if ok {
		delete(mapMonitorCallback, &Callback)
	}
}

func MonitorDump() string {
	var lines []string

	pathMaxLength := 0
	var pList []string

	for Path, _ := range mapPathMonitorCallback {
		pList = append(pList, Path)
		if len(Path) > pathMaxLength {
			pathMaxLength = len(Path)
		}
	}

	fmtstr := fmt.Sprintf(
		"%s%%-%ds%s : %%v",
		ColorType_I,
		pathMaxLength,
		ColorType_Reset,
	)

	for _, Path := range pList {
		for Monitor, _ := range mapPathMonitorCallback[Path] {
			lines = append(
				lines,
				fmt.Sprintf(
					fmtstr,
					Path,
					Monitor,
				),
			)
		}
	}

	lines = append(lines, "")
	sort.Slice(lines, func(i int, j int) bool {
		return strings.Compare(lines[i], lines[j]) < 0
	})
	return strings.Join(lines, "\n")
}

func monitorCall(e *confEntry, oVal interface{}, nVal interface{}) {
	if mapMonitorCallback, ok := mapPathMonitorCallback[e.path]; ok {
		for Callback, _ := range mapMonitorCallback {
			if Callback != nil {
				(*Callback)(e.path, oVal, nVal)
			}
		}
	}
}

// OnReady : Called when all configure loaded.
func OnReady(cb func()) {
	MonitorAdd(
		PathReady,
		func(p string, o, n interface{}) {
			cb()
		},
	)
}
