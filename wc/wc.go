package wc

import (
	"os"

	"miego/klog"

	"github.com/fsnotify/fsnotify"
)

type HandlerFunc func(path string, event string)

type WatchChanges struct {
	watcher *fsnotify.Watcher
	names   []string

	onEverything HandlerFunc
	onCreate     HandlerFunc
	onWrite      HandlerFunc
	onRemove     HandlerFunc
	onRename     HandlerFunc
	onChmod      HandlerFunc

	done chan bool
}

func New(names ...string) (*WatchChanges, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		klog.E("fsnotify.NewWatcher: %s", err.Error())
		return nil, err
	}

	wc := &WatchChanges{}
	wc.names = names
	wc.done = make(chan bool)
	wc.watcher = watcher

	return wc, nil
}

func (wc *WatchChanges) OnEverything(callback HandlerFunc) *WatchChanges {
	wc.onEverything = callback
	return wc
}
func (wc *WatchChanges) OnCreate(callback HandlerFunc) *WatchChanges {
	wc.onCreate = callback
	return wc
}
func (wc *WatchChanges) OnWrite(callback HandlerFunc) *WatchChanges {
	wc.onWrite = callback
	return wc
}
func (wc *WatchChanges) OnRemove(callback HandlerFunc) *WatchChanges {
	wc.onRemove = callback
	return wc
}
func (wc *WatchChanges) OnRename(callback HandlerFunc) *WatchChanges {
	wc.onRename = callback
	return wc
}
func (wc *WatchChanges) OnChmod(callback HandlerFunc) *WatchChanges {
	wc.onChmod = callback
	return wc
}

func (wc *WatchChanges) Run() error {
	defer wc.watcher.Close()
	go func() {
		for {
			select {
			case event, ok := <-wc.watcher.Events:
				if !ok {
					klog.E("NG")
					return
				}

				if wc.onEverything != nil {
					klog.Dump(event)
					wc.onEverything(event.Name, "ALL")
				}

				if event.Op&fsnotify.Create == fsnotify.Create {
					if wc.onCreate != nil {
						wc.onCreate(event.Name, "CREATE")
					}
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					if wc.onWrite != nil {
						wc.onWrite(event.Name, "WRITE")
					}
				}

				if event.Op&fsnotify.Remove == fsnotify.Remove {
					if wc.onRemove != nil {
						wc.onRemove(event.Name, "REMOVE")
					}
				}

				if event.Op&fsnotify.Rename == fsnotify.Rename {
					if wc.onRename != nil {
						wc.onRename(event.Name, "RENAME")
					}
				}

				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					if wc.onChmod != nil {
						wc.onChmod(event.Name, "CHMOD")
					}
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
