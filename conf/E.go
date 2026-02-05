package conf

import "github.com/google/uuid"

type eListener func(key string, arg any)

type eItem struct {
	listeners map[string]eListener
}

// XXX: no lock
func (cc *ConfCenter) eShout(item *eItem, key string, arg any) {
	if item.listeners != nil {
		for _, cb := range item.listeners {
			if cb != nil {
				go cb(key, arg)
			}
		}
	}
}

//////////////////////////////////////////////////////////////////////////
// BASE: Has, Get, Getb(bool), Set, Setf(force), Rem /////////////////////

func (cc *ConfCenter) EHas(key string) bool {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	_, ok := cc.eItems[key]
	return ok
}

func (cc *ConfCenter) EAdd(key string) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if _, ok := cc.eItems[key]; !ok {
		cc.eItems[key] = &eItem{}
	}
}

func (cc *ConfCenter) ESend(key string, arg any) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if item, ok := cc.eItems[key]; ok {
		cc.eShout(item, key, arg)
	}
}
func (cc *ConfCenter) Emit(key string, arg any) {
	cc.ESend(key, arg)
}

func (cc *ConfCenter) ESendf(key string, arg any) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	item, ok := cc.eItems[key]
	if !ok {
		item = &eItem{}
		cc.eItems[key] = item
	}

	cc.eShout(item, key, arg)
}
func (cc *ConfCenter) Emitf(key string, arg any) {
	cc.ESendf(key, arg)
}

func (cc *ConfCenter) ERem(key string) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	delete(cc.eItems, key)
}

// ///////////////////////////////////////////////////////////////////////
// Monitor ///////////////////////////////////////////////////////////////
func (cc *ConfCenter) EListenerAdd(key string, cbName string, cb eListener) string {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if cbName == "" {
		cbName = uuid.New().String()
	}

	if item, ok := cc.eItems[key]; ok {
		if item.listeners == nil {
			item.listeners = make(map[string]eListener)
		}
		if _, ok := item.listeners[cbName]; !ok {
			item.listeners[cbName] = cb
			return cbName
		}
	}
	return ""
}

func (cc *ConfCenter) EListenerRem(key string, cbName string) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if item, ok := cc.eItems[key]; ok {
		if item.listeners != nil {
			delete(item.listeners, cbName)
		}
	}
}
