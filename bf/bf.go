package bf

// Blind Fruits: 芒果DB

import (
	"fmt"

	"github.com/kamasamikon/miego/conf"
	"github.com/kamasamikon/miego/klog"
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
	if conf.Bool(false, "b:/db/mg/debug") {
		mgo.SetLogger(xLogger)
		mgo.SetDebug(true)
	}
}

func Bye(session *mgo.Session) {
	session.Close()
}

func Hey() *mgo.Session {
	addr := conf.Str("", "s:/db/mg/addr")
	port := conf.Int(3717, "i:/db/mg/port")

	// Dial
	mgoaddr := fmt.Sprintf("%s:%d", addr, port)
	klog.D("Hey: %s", mgoaddr)
	session, err := mgo.Dial(mgoaddr)
	if err != nil {
		klog.E(err.Error())
		return nil
	}

	// Login
	if user := conf.Str("root", "s:/db/mg/user"); user != "" {
		pass := conf.Str("", "s:/db/mg/pass")
		auth := conf.Str("SCRAM-SHA-1", "s:/db/mg/auth")

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
