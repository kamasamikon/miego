package conf

import (
	"github.com/google/uuid"
)

type sMonitor func(key string, vnow string, vnew string)

type sItem struct {
	key      string
	monitors map[string]sMonitor
	value    string
}

// XXX: no lock
func (cc *ConfCenter) sSet(item *sItem, vnew string) {
	vnow := item.value
	item.value = vnew
	if item.monitors != nil {
		for _, cb := range item.monitors {
			if cb != nil {
				go cb(item.key, vnow, vnew)
			}
		}
	}
}

//////////////////////////////////////////////////////////////////////////
// BASE: Has, Get, Getb(bool), Set, Setf(force), Rem /////////////////////

func (cc *ConfCenter) SHas(key string) bool {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	_, ok := cc.sItems[key]
	return ok
}

func (cc *ConfCenter) SGet(key string, vdef string) string {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if item, ok := cc.sItems[key]; ok {
		return item.value
	}
	return vdef
}

func (cc *ConfCenter) SGetb(key string) (string, bool) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if item, ok := cc.sItems[key]; ok {
		return item.value, true
	}

	return "", false
}

func (cc *ConfCenter) SSet(key string, val string) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if item, ok := cc.sItems[key]; ok {
		cc.sSet(item, val)
	}
}

func (cc *ConfCenter) SSetf(key string, val string) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	item, ok := cc.sItems[key]
	if !ok {
		item = &sItem{
			key: key,
		}
		cc.sItems[item.key] = item
	}

	cc.sSet(item, val)
}

func (cc *ConfCenter) SRem(key string) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	delete(cc.sItems, key)
}

// ///////////////////////////////////////////////////////////////////////
// Monitor ///////////////////////////////////////////////////////////////
func (cc *ConfCenter) SMonitorAdd(key string, cbName string, cb sMonitor) string {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if cbName == "" {
		cbName = uuid.New().String()
	}

	if e, ok := cc.sItems[key]; ok {
		if e.monitors == nil {
			e.monitors = make(map[string]sMonitor)
		}

		if _, ok := e.monitors[cbName]; !ok {
			e.monitors[cbName] = cb
			return cbName
		}
	}
	return ""
}

func (cc *ConfCenter) SMonitorRem(key string, cbName string) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if e, ok := cc.sItems[key]; ok {
		if e.monitors != nil {
			delete(e.monitors, cbName)
		}
	}
}

// ///////////////////////////////////////////////////////////////////////
// Others ////////////////////////////////////////////////////////////////
