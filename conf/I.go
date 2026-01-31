package conf

import "github.com/google/uuid"

type iMonitor func(key string, vnow int64, vnew int64)

type iItem struct {
	key      string
	monitors map[string]iMonitor
	value    int64
}

// XXX: no lock
func (cc *ConfCenter) iSet(item *iItem, val any) {
	var vnew int64
	switch v := val.(type) {
	case int64:
		vnew = int64(v)
	case int32:
		vnew = int64(v)
	case int:
		vnew = int64(v)
	case int16:
		vnew = int64(v)
	case int8:
		vnew = int64(v)
	case uint64:
		vnew = int64(v)
	case uint32:
		vnew = int64(v)
	case uint:
		vnew = int64(v)
	case uint16:
		vnew = int64(v)
	case uint8:
		vnew = int64(v)
	}

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

func (cc *ConfCenter) IHas(key string) bool {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	_, ok := cc.iItems[key]
	return ok
}

func (cc *ConfCenter) IGet(key string, vdef int64) int64 {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if item, ok := cc.iItems[key]; ok {
		return item.value
	}
	return vdef
}
func (cc *ConfCenter) I(key string, vdef int64) int64 {
	return cc.IGet(key, vdef)
}

func (cc *ConfCenter) IGetb(key string) (int64, bool) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if item, ok := cc.iItems[key]; ok {
		return item.value, true
	}

	return 0, false
}

func (cc *ConfCenter) ISet(key string, val any) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if item, ok := cc.iItems[key]; ok {
		cc.iSet(item, val)
	}
}

func (cc *ConfCenter) ISetf(key string, val any) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	item, ok := cc.iItems[key]
	if !ok {
		item = &iItem{
			key: key,
		}
		cc.iItems[item.key] = item
	}

	cc.iSet(item, val)
}

func (cc *ConfCenter) IRem(key string) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	delete(cc.iItems, key)
}

// ///////////////////////////////////////////////////////////////////////
// Monitor ///////////////////////////////////////////////////////////////
func (cc *ConfCenter) IMonitorAdd(key string, cbName string, cb iMonitor) string {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if cbName == "" {
		cbName = uuid.New().String()
	}

	if e, ok := cc.iItems[key]; ok {
		if e.monitors == nil {
			e.monitors = make(map[string]iMonitor)
		}

		if _, ok := e.monitors[cbName]; !ok {
			e.monitors[cbName] = cb
			return cbName
		}
	}
	return ""
}

func (cc *ConfCenter) IMonitorRem(key string, cbName string) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if e, ok := cc.iItems[key]; ok {
		if e.monitors != nil {
			delete(e.monitors, cbName)
		}
	}
}

// ///////////////////////////////////////////////////////////////////////
// Others ////////////////////////////////////////////////////////////////
func (cc *ConfCenter) IInc(key string, inc int64) int64 {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if item, ok := cc.iItems[key]; ok {
		vNew := item.value + inc
		cc.iSet(item, vNew)
		return item.value
	}
	return -1
}
