package conf

type xGetter func(key string) any
type xSetter func(key string, arg any)

type xItem struct {
	getter xGetter
	setter xSetter
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
		cc.xItems[key] = &xItem{}
	}
}

func (cc *ConfCenter) XSet(key string, arg any) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if item, ok := cc.xItems[key]; ok {
		if item.setter != nil {
			item.setter(key, arg)
		}
	}
}
func (cc *ConfCenter) XGet(key string) any {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if item, ok := cc.xItems[key]; ok {
		if item.setter != nil {
			return item.getter(key)
		}
	}
	return nil
}

func (cc *ConfCenter) XRem(key string) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	delete(cc.xItems, key)
}

// ///////////////////////////////////////////////////////////////////////
// Monitor ///////////////////////////////////////////////////////////////
func (cc *ConfCenter) XSetSetter(key string, setter xSetter) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if item, ok := cc.xItems[key]; ok {
		item.setter = setter
	}
}

func (cc *ConfCenter) XSetGetter(key string, getter xGetter) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if item, ok := cc.xItems[key]; ok {
		item.getter = getter
	}
}
