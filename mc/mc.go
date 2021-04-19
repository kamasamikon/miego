package mc

import (
	"sync"
	"time"

	"github.com/kamasamikon/miego/roast"
)

// ID <> ExpTime, ID <> RealData
var map_UUID_Exp map[string]int64
var map_UUID_Data map[string]string

var mutex = &sync.Mutex{}

// Set : Add a string to cache and return a UUID as a key
func Set(data string, exp int64) string {
	mutex.Lock()
	defer mutex.Unlock()

	UUID := roast.IDNew()

	// Exp after exp seconds
	map_UUID_Exp[UUID] = time.Now().Unix() + exp
	map_UUID_Data[UUID] = data

	return UUID
}

// Get : UUID => Original string
func Get(data string) string {
	mutex.Lock()
	defer mutex.Unlock()

	nowUnix := time.Now().Unix()
	if Data, ok := map_UUID_Data[data]; ok {
		if Exp, ok := map_UUID_Exp[data]; ok {
			if Exp < nowUnix {
				return Data
			}
		}
	}
	return ""
}

func init() {
	map_UUID_Exp = make(map[string]int64)
	map_UUID_Data = make(map[string]string)

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
