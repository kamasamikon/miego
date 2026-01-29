package conf

// ////////////////////////////////////////////////////////////////////////
// Monitor, callback when configure changed
//

func (cc *ConfCenter) MonitorAdd(key string, cb Monitor) int {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if e, ok := cc.mapPathEntry[key]; ok {
		for idx := range e.monitors {
			if e.monitors[idx] == nil {
				e.monitors[idx] = cb
				return idx
			}
		}
		e.monitors = append(e.monitors, cb)
		return len(e.monitors) - 1
	}
	return -1
}

func (cc *ConfCenter) MonitorRem(key string, idx int) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if e, ok := cc.mapPathEntry[key]; ok {
		if idx >= 0 && idx < len(e.monitors) {
			e.monitors[idx] = nil
		}
	}
}

func (cc *ConfCenter) monitorCall(e *confEntry, oVal any, nVal any) {
	for _, cb := range e.monitors {
		if cb != nil {
			go cb(e.path, oVal, nVal)
		}
	}
}
