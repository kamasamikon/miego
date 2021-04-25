package mc

import (
	"sync"
	"time"

	"github.com/kamasamikon/miego/roast"
)

// ID <> ExpTime, ID <> RealData
var map_UUID_Exp map[string]int64
var map_UUID_Data map[string]interface{}

var mutex = &sync.Mutex{}

// Set : Add a string to cache and return a UUID as a key
func Set(data interface{}, exp int64, UUID string) string {
	mutex.Lock()
	defer mutex.Unlock()

	if UUID == "" {
		UUID = roast.IDNew()
	}

	// Exp after exp seconds
	map_UUID_Exp[UUID] = time.Now().Unix() + exp
	map_UUID_Data[UUID] = data

	return UUID
}

func Rem(uuid string) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, ok := map_UUID_Exp[uuid]; ok {
		map_UUID_Exp[uuid] = 0
	}
}

// Get : UUID => Original string
func V(uuid string) interface{} {
	mutex.Lock()
	defer mutex.Unlock()

	nowUnix := time.Now().Unix()
	if Data, ok := map_UUID_Data[uuid]; ok {
		if Exp, ok := map_UUID_Exp[uuid]; ok {
			if Exp < nowUnix {
				return Data
			}
		}
	}
	return nil
}

func VMust(uuid string, New func(uuid string) (data interface{}, exp int64)) interface{} {
	mutex.Lock()
	defer mutex.Unlock()

	nowUnix := time.Now().Unix()
	if Data, ok := map_UUID_Data[uuid]; ok {
		if Exp, ok := map_UUID_Exp[uuid]; ok {
			if Exp < nowUnix {
				return Data
			}
		}
	}

	data, exp := New(uuid)
	map_UUID_Exp[uuid] = time.Now().Unix() + exp
	map_UUID_Data[uuid] = data

	return data
}

func S(uuid string) string {
	v := V(uuid)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func I(uuid string) int {
	v := V(uuid)
	if i, ok := v.(int); ok {
		return i
	}
	return 0
}

func init() {
	map_UUID_Exp = make(map[string]int64)
	map_UUID_Data = make(map[string]interface{})

	nameCache := make(map[string]int)

	var mutex = &sync.Mutex{}

	go func() {
		time.Sleep(time.Second * 30)

		mutex.Lock()
		nowUnix := time.Now().Unix()
		for UUID, _ := range map_UUID_Data {
			if Exp, ok := map_UUID_Exp[UUID]; ok {
				if nowUnix < Exp {
					continue
				}
			}
			nameCache[UUID] = 1
		}
		for UUID, _ := range nameCache {
			delete(map_UUID_Exp, UUID)
			delete(map_UUID_Data, UUID)
		}
		mutex.Unlock()
	}()
}