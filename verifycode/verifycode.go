package verifycode

import (
	"time"

	"miego/klog"
	"miego/xrand"
)

// Verification Code
type VCodeItem struct {
	Code    string
	DueTime int64
}

type VCodeChecker struct {
	Map map[string]*VCodeItem

	// return: result, processed
	CheckHook func(key string, code string) (bool, bool)
}

func VCodeCheckerNew() *VCodeChecker {
	cc := &VCodeChecker{}
	cc.Map = make(map[string]*VCodeItem)
	go func() {
		for {
			if cc == nil {
				klog.E("VCodeChecker is nul")
				break
			}
			cc.Clean()
			time.Sleep(time.Second * 10)
		}
	}()
	return cc
}

func (cc VCodeChecker) Rand(key string, codeLen int, codeKind string, ttl int64) string {
	code := string(xrand.Rand(codeLen, codeKind))
	cc.Add(key, code, ttl)
	return code
}

// Add : ttl: seconds
func (cc VCodeChecker) Add(key string, code string, ttl int64) {
	x := VCodeItem{
		Code:    code,
		DueTime: time.Now().Unix() + ttl,
	}
	cc.Map[key] = &x
}

func (cc VCodeChecker) Rem(key string) {
	delete(cc.Map, key)
}

func (cc VCodeChecker) Check(key string, code string) bool {
	if cc.CheckHook != nil {
		if passed, processed := cc.CheckHook(key, code); processed == true {
			return passed
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
