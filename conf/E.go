package conf

import "github.com/google/uuid"

type eListener func(key string, arg any)

type eItem struct {
	key       string
	listeners map[string]eListener
}

// XXX: no lock
func (cc *ConfCenter) eShout(item *eItem, arg any) {
	if item.listeners != nil {
		for _, cb := range item.listeners {
			if cb != nil {
				go cb(item.key, arg)
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

func (cc *ConfCenter) ESend(key string, arg any) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if item, ok := cc.eItems[key]; ok {
		cc.eShout(item, arg)
	}
}

func (cc *ConfCenter) ESendf(key string, arg any) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	item, ok := cc.eItems[key]
	if !ok {
		item = &eItem{
			key: key,
		}
		cc.eItems[item.key] = item
	}

	cc.eShout(item, arg)
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

	if e, ok := cc.eItems[key]; ok {
		if e.listeners == nil {
			e.listeners = make(map[string]eListener)
		}
		if _, ok := e.listeners[cbName]; !ok {
			e.listeners[cbName] = cb
			return cbName
		}
	}
	return ""
}

func (cc *ConfCenter) EListenerRem(key string, cbName string) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if e, ok := cc.eItems[key]; ok {
		if e.listeners != nil {
			delete(e.listeners, cbName)
		}
	}
}
