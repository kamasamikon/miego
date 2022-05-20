package mc

import (
	"sync"
	"time"

	"github.com/twinj/uuid"
)

type Item struct {
	data interface{} // 缓存的数据
	exp  int64       // 过期时间
	tag  string      // 标签，分类用
}

var Items map[string]*Item

var mutex = &sync.Mutex{}

// Set : Add a string to cache and return a UID as a key with exp seconds
func Set(data interface{}, exp int64, uid string, tag string) string {
	mutex.Lock()
	defer mutex.Unlock()

	if uid == "" {
		uid = uuid.NewV4().String()
	}

	// Exp after exp seconds
	Items[uid] = &Item{
		data: data,
		exp:  time.Now().Unix() + exp,
		tag:  tag,
	}

	return uid
}

// Mod : Modify
func Mod(uid string, data interface{}, exp int64, tag string) {
	mutex.Lock()
	defer mutex.Unlock()

	if item, ok := Items[uid]; ok {
		if data != nil {
			item.data = data
		}
		if exp != 0 {
			item.exp = time.Now().Unix() + exp
		}
		if tag != "" {
			item.tag = tag
		}
	}
}

// Find : 通过Tag查找
func Find(tag string) []string {
	mutex.Lock()
	defer mutex.Unlock()

	var uidList []string

	now := time.Now().Unix()
	for uid, item := range Items {
		if item.exp > now {
			uidList = append(uidList, uid)
		}
	}

	return uidList
}

// Rem : 删除一些
func Rem(uidList ...string) {
	mutex.Lock()
	defer mutex.Unlock()

	for i := 0; i < len(uidList); i++ {
		delete(Items, uidList[i])
	}
}

// Get : UID => Original string
func V(uid string) (interface{}, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	now := time.Now().Unix()
	if item, ok := Items[uid]; ok {
		if item.exp > now {
			return item.data, true
		}
	}

	return nil, false
}

// VMust : 获取，如果不存在就添加
func VMust(uid string, New func(uid string) (data interface{}, exp int64, tag string)) interface{} {
	mutex.Lock()
	defer mutex.Unlock()

	now := time.Now().Unix()
	if item, ok := Items[uid]; ok {
		if item.exp > now {
			return item.data
		}
	}

	data, exp, tag := New(uid)
	Items[uid] = &Item{
		data: data,
		exp:  time.Now().Unix() + exp,
		tag:  tag,
	}

	return data
}

func S(uid string) (string, bool) {
	if v, ok := V(uid); !ok {
		return "", false
	} else if s, ok := v.(string); ok {
		return s, true
	}
	return "", false
}

func I(uid string) (int, bool) {
	if v, ok := V(uid); !ok {
		return 0, false
	} else if i, ok := v.(int); ok {
		return i, true
	}
	return 0, false
}

func cleanup() {
	mutex.Lock()
	defer mutex.Unlock()

	nameCache := make(map[string]int)

	now := time.Now().Unix()
	for uid, item := range Items {
		if now < item.exp {
			continue
		}
		nameCache[uid] = 1
	}
	for uid, _ := range nameCache {
		delete(Items, uid)
	}
}

func init() {
	Items = make(map[string]*Item)

	go func() {
		for {
			time.Sleep(time.Second * 30)
			cleanup()
		}
	}()
}
