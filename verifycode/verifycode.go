package verifycode

import (
	"time"

	"github.com/kamasamikon/miego/klog"
	"github.com/kamasamikon/miego/xrand"
)

//
// Verification Code
//
type VCodeItem struct {
	Code    string
	DueTime int64
}

type VCodeChecker struct {
	Map       map[string]*VCodeItem
	CheckHook func(key string, code string) (bool, error)
}

func VCodeCheckerNew() *VCodeChecker {
	cc := &VCodeChecker{}
	cc.Map = make(map[string]*VCodeItem)
	go func() {
		for {
			if cc != nil {
				klog.E("VCodeChecker is nul")
				break
			}
			cc.Clean()
			time.Sleep(time.Second * 10)
		}
	}()
	return cc
}

func (cc VCodeChecker) Rand(key string, codeLen int, codeKind string, ttl int64) {
	code := string(xrand.Rand(codeLen, codeKind))
	cc.Add(key, code, ttl)
}

// Add : ttl: seconds
func (cc VCodeChecker) Add(key string, code string, ttl int64) {
	x := VCodeItem{
		Code:    code,
		DueTime: time.Now().Unix() + ttl,
	}
	cc.Map[key] = &x
}

func (cc VCodeChecker) Check(key string, code string) bool {
	if cc.CheckHook != nil {
		if ok, err := cc.CheckHook(key, code); err != nil {
			return ok
		}
	}

	now := time.Now().Unix()
	if v, ok := cc.Map[key]; ok {
		if v.Code != code {
			klog.D("Bad: %s: %s vs %s", key, v.Code, code)
			return false
		}
		if v.DueTime < now {
			delete(cc.Map, key)
			klog.D("Exp: %s", key)
			return false
		}
		klog.D("Valid: %s", key)
		return true
	}

	klog.D("Miss: %s", key)
	return false
}

func (cc VCodeChecker) Clean() {
	now := time.Now().Unix()
	for k, v := range cc.Map {
		if v.DueTime < now {
			klog.D("Clean: %s", k)
			delete(cc.Map, k)
		}
	}
}
