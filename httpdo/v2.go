package httpdo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/kamasamikon/miego/klog"
)

type context struct {
	url  string
	ping interface{}
	pong interface{}
	mime string
}

func New(url string) *context {
	return &context{
		url:  url,
		mime: "application/json;charset=utf-8",
	}
}

func (c *context) MIME(mime string) *context {
	c.mime = mime
	return c
}

func (c *context) Ping(ping interface{}) *context {
	c.ping = ping
	return c
}

func (c *context) Pong(pong interface{}) *context {
	c.pong = pong
	return c
}

func (c *context) Post() (n int, err error) {
	var pingString string

	if c.ping == nil {
		pingString = ""
	} else {
		if s, ok := c.ping.(string); ok {
			pingString = s
		} else {
			bytes, err := json.Marshal(c.ping)
			if err != nil {
				return 0, err
			}
			pingString = string(bytes)
		}
	}

	r, err := http.Post(c.url, c.mime, strings.NewReader(pingString))
	if err != nil {
		klog.E(err.Error())
		return 0, err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		klog.E("%d", r.StatusCode)
		return 0, fmt.Errorf("StatusCode == %d", r.StatusCode)
	}

	if c.pong == nil {
		return 0, nil
	}

	if ptr, ok := c.pong.(*string); ok {
		if buf, err := ioutil.ReadAll(r.Body); err != nil {
			return 0, err
		} else {
			*ptr = string(buf)
			return 0, nil
		}
	} else if ptr, ok := c.pong.(*[]byte); ok {
		if buf, err := ioutil.ReadAll(r.Body); err != nil {
			return 0, err
		} else {
			*ptr = buf
			return 0, nil
		}
	} else {
		return 0, json.NewDecoder(r.Body).Decode(c.pong)
	}
}

// Get : HTTPGet convert the response to pongObj structure
func (c *context) Get() (n int, err error) {
	r, err := http.Get(c.url)
	if err != nil {
		klog.E(err.Error())
		return 0, err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		klog.E("URL:%s, Code:%d", c.url, r.StatusCode)
		return 0, fmt.Errorf("StatusCode == %d", r.StatusCode)
	}

	if c.pong == nil {
		return 0, nil
	}

	if ptr, ok := c.pong.(*string); ok {
		if buf, err := ioutil.ReadAll(r.Body); err != nil {
			return 0, err
		} else {
			*ptr = string(buf)
			return 0, nil
		}
	} else if ptr, ok := c.pong.(*[]byte); ok {
		if buf, err := ioutil.ReadAll(r.Body); err != nil {
			return 0, err
		} else {
			*ptr = buf
			return 0, nil
		}
	} else {
		return 0, json.NewDecoder(r.Body).Decode(c.pong)
	}
}
