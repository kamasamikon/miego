package wc

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/kamasamikon/miego/klog"
)

type WatchChanges struct {
	watcher  *fsnotify.Watcher
	names    []string
	callback func(name string)
	done     chan bool
}

func WCNew(callback func(name string), names ...string) *WatchChanges {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		klog.E("error:%s", err.Error())
		return nil
	}

	wc := &WatchChanges{}
	wc.names = names
	wc.callback = callback
	wc.done = make(chan bool)
	wc.watcher = watcher

	return wc
}

func (wc *WatchChanges) Run() error {
	defer wc.watcher.Close()
	go func() {
		for {
			select {
			case event, ok := <-wc.watcher.Events:
				if !ok {
					klog.E("!OK")
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					wc.callback(event.Name)
				}
			case err, ok := <-wc.watcher.Errors:
				if !ok {
					klog.E("!OK")
					return
				}
				klog.E("error:%s", err.Error())
			}
		}
	}()

	for _, name := range wc.names {
		if _, err := os.Stat(name); err == nil {
			err := wc.watcher.Add(name)
			if err != nil {
				klog.E("Name:%s, Error:%s", name, err.Error())
			}
		} else {
			klog.E("Name:%s, Error:%s", name, err.Error())
		}
	}

	<-wc.done
	return nil
}

func (wc *WatchChanges) Bye() {
	wc.done <- true
}
