package conf

import "github.com/google/uuid"

type bMonitor func(key string, vnow bool, vnew bool)

type bItem struct {
	monitors map[string]bMonitor
	value    bool
}

// XXX: no lock
func (cc *ConfCenter) bSet(item *bItem, key string, vnew bool) {
	vnow := item.value
	item.value = vnew
	if item.monitors != nil {
		for _, cb := range item.monitors {
			if cb != nil {
				go cb(key, vnow, vnew)
			}
		}
	}
}

//////////////////////////////////////////////////////////////////////////
// BASE: Has, Get, Getb(bool), Set, Setf(force), Rem /////////////////////

func (cc *ConfCenter) BHas(key string) bool {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	_, ok := cc.bItems[key]
	return ok
}

func (cc *ConfCenter) BGet(key string, vdef bool) bool {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if item, ok := cc.bItems[key]; ok {
		return item.value
	}
	return vdef
}
func (cc *ConfCenter) B(key string, vdef bool) bool {
	return cc.BGet(key, vdef)
}
func (cc *ConfCenter) BTrue(key string) bool {
	return cc.BGet(key, true)
}
func (cc *ConfCenter) BFalse(key string) bool {
	return cc.BGet(key, false)
}

func (cc *ConfCenter) BGetb(key string) (bool, bool) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if item, ok := cc.bItems[key]; ok {
		return item.value, true
	}

	return false, false
}

func (cc *ConfCenter) BSet(key string, val bool) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if item, ok := cc.bItems[key]; ok {
		cc.bSet(item, key, val)
	}
}

func (cc *ConfCenter) BSetf(key string, val bool) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	item, ok := cc.bItems[key]
	if !ok {
		item = &bItem{}
		cc.bItems[key] = item
	}

	cc.bSet(item, key, val)
}

func (cc *ConfCenter) BRem(key string) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	delete(cc.bItems, key)
}

// ///////////////////////////////////////////////////////////////////////
// Monitor ///////////////////////////////////////////////////////////////
func (cc *ConfCenter) BMonitorAdd(key string, cbName string, cb bMonitor) string {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if cbName == "" {
		cbName = uuid.New().String()
	}

	if e, ok := cc.bItems[key]; ok {
		if e.monitors == nil {
			e.monitors = make(map[string]bMonitor)
		}
		if _, ok := e.monitors[cbName]; !ok {
			e.monitors[cbName] = cb
			return cbName
		}
	}
	return ""
}

func (cc *ConfCenter) BMonitorRem(key string, cbName string) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if e, ok := cc.bItems[key]; ok {
		if e.monitors != nil {
			delete(e.monitors, cbName)
		}
	}
}

// ///////////////////////////////////////////////////////////////////////
// Others ////////////////////////////////////////////////////////////////
func (cc *ConfCenter) BFlip(key string) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if item, ok := cc.bItems[key]; ok {
		cc.bSet(item, key, !item.value)
	}
}
