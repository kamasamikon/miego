package conf

func (cc *ConfCenter) SetSetter(key string, setter setter) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if e, ok := cc.mapPathEntry[key]; ok {
		e.setter = setter
	} else {
		cc.mapPendingSetter[key] = setter
	}
}

func (cc *ConfCenter) SetGetter(key string, getter getter) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if e, ok := cc.mapPathEntry[key]; ok {
		e.getter = getter
	} else {
		cc.mapPendingGetter[key] = getter
	}
}
