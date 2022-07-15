package conf

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/kamasamikon/miego/klog"
)

var mutex = &sync.Mutex{}
var nextMonitorID int

// ////////////////////////////////////////////////////////////////////////
// Monitor, callback when configure changed
//

// KConfMonitor is a Callback called when wathed entry modified.
type KConfMonitor func(path string, monitorID string, oVal interface{}, nVal interface{})

// map[s:Path]map[s:MonitorID]KConfMonitor
var mapPathMonitorCallback = make(map[string]map[string]KConfMonitor)

func MonitorAdd(Path string, Callback KConfMonitor) string {
	mutex.Lock()
	defer mutex.Unlock()
	MonitorID := fmt.Sprintf("%d", nextMonitorID)
	nextMonitorID++

	mapMonitorCallback, ok := mapPathMonitorCallback[Path]
	if !ok {
		mapMonitorCallback = make(map[string]KConfMonitor)
	}
	mapMonitorCallback[MonitorID] = Callback
	mapPathMonitorCallback[Path] = mapMonitorCallback

	return MonitorID
}

func MonitorRem(Path string, MonitorID string) {
	mapMonitorCallback, ok := mapPathMonitorCallback[Path]
	if ok {
		delete(mapMonitorCallback, MonitorID)
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

	sort.Slice(pList, func(i int, j int) bool {
		return strings.Compare(pList[i], pList[j]) < 0
	})

	fmtstr := fmt.Sprintf(
		"%s%%-%ds%s : %%v",
		klog.ColorType_I,
		pathMaxLength,
		klog.ColorType_Reset,
	)

	for _, Path := range pList {
		for MonitorID, _ := range mapPathMonitorCallback[Path] {
			lines = append(
				lines,
				fmt.Sprintf(
					fmtstr,
					Path,
					MonitorID,
				),
			)
		}
	}

	lines = append(lines, "")
	return strings.Join(lines, "\n")
}

// 同步调用，处理函数应该把真正的逻辑放到goroutine中去
func monitorCall(e *confEntry, oVal interface{}, nVal interface{}) {
	if mapMonitorCallback, ok := mapPathMonitorCallback[e.path]; ok {
		for MonitorID, Callback := range mapMonitorCallback {
			if Callback != nil {
				Callback(e.path, MonitorID, oVal, nVal)
			}
		}
	}
}
