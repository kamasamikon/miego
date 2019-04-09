package conf

import (
	"lib/klog"
	"testing"
	"time"
)

func callbackA(path string, a interface{}, b interface{}) {
	klog.D("Callback A")
	klog.D(path)
	klog.D("%s", a)
	klog.D("%s", b)
}

func callbackC(path string, a interface{}, b interface{}) {
	klog.D("Callback C")
	klog.D(path)
	klog.D("%s", a)
	klog.D("%s", b)
}

func TestMain(t *testing.T) {
	Load("./msa.cfg")
	Monitor("s:/msa/ms/name", callbackA)
	Monitor("s:/msa/ms/name", callbackC)
	Set("s:/msa/ms/name", "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", true)
	klog.D(Str("msa/ms/name", "NoFound"))
	klog.D(Str("msa/ms/namex", "NoFound"))
	Dump()

	time.Sleep(time.Second)
}
