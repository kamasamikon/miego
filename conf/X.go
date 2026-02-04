package conf

import "github.com/google/uuid"

type xListener func(key string, arg any)

type xItem struct {
	key       string
	listeners map[string]xListener
}

// XXX: no lock
func (cc *ConfCenter) xShout(item *xItem, arg any) {
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

func (cc *ConfCenter) XHas(key string) bool {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	_, ok := cc.xItems[key]
	return ok
}

func (cc *ConfCenter) XAdd(key string) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if _, ok := cc.xItems[key]; !ok {
		cc.xItems[key] = &xItem{
			key: key,
		}
	}
}

func (cc *ConfCenter) XSend(key string, arg any) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if item, ok := cc.xItems[key]; ok {
		cc.xShout(item, arg)
	}
}
func (cc *ConfCenter) Xmit(key string, arg any) {
	cc.XSend(key, arg)
}

func (cc *ConfCenter) XSendf(key string, arg any) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	item, ok := cc.xItems[key]
	if !ok {
		item = &xItem{
			key: key,
		}
		cc.xItems[item.key] = item
	}

	cc.xShout(item, arg)
}
func (cc *ConfCenter) Xmitf(key string, arg any) {
	cc.XSendf(key, arg)
}

func (cc *ConfCenter) XRem(key string) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	delete(cc.xItems, key)
}

// ///////////////////////////////////////////////////////////////////////
// Monitor ///////////////////////////////////////////////////////////////
func (cc *ConfCenter) XListenerAdd(key string, cbName string, cb xListener) string {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if cbName == "" {
		cbName = uuid.New().String()
	}

	if e, ok := cc.xItems[key]; ok {
		if e.listeners == nil {
			e.listeners = make(map[string]xListener)
		}
		if _, ok := e.listeners[cbName]; !ok {
			e.listeners[cbName] = cb
			return cbName
		}
	}
	return ""
}

func (cc *ConfCenter) XListenerRem(key string, cbName string) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if e, ok := cc.xItems[key]; ok {
		if e.listeners != nil {
			delete(e.listeners, cbName)
		}
	}
}
