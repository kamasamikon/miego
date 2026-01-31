package bf

// Blind Fruits: 芒果DB

import (
	"fmt"

	"miego/conf"
	"miego/klog"

	"gopkg.in/mgo.v2"
)

type XLogger struct {
}

func (c *XLogger) Output(calldepth int, s string) error {
	klog.D(s)
	return nil
}

var xLogger = &XLogger{}

func init() {
	if conf.BFalse("db/mg/debug") {
		mgo.SetLogger(xLogger)
		mgo.SetDebug(true)
	}
}

func Bye(session *mgo.Session) {
	session.Close()
}

func Hey() *mgo.Session {
	addr := conf.S("db/mg/addr")
	port := conf.I("db/mg/port", 3717)

	// Dial
	mgoaddr := fmt.Sprintf("%s:%d", addr, port)
	klog.D("Hey: %s", mgoaddr)
	session, err := mgo.Dial(mgoaddr)
	if err != nil {
		klog.E(err.Error())
		return nil
	}

	// Login
	if user := conf.SGet("db/mg/user", "root"); user != "" {
		pass := conf.S("db/mg/pass")
		auth := conf.SGet("db/mg/auth", "SCRAM-SHA-1")

		err = session.Login(&mgo.Credential{
			Username:  user,
			Password:  pass,
			Mechanism: auth,
			Source:    "admin",
		})

		if err == nil {
			return session
		}
	} else {
		return session
	}

	defer session.Close()
	return nil
}
