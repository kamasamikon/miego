package httpdo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/kamasamikon/miego/klog"
)

type context struct {
	url         string
	ping        interface{}
	pong        interface{}
	contentType string
	header      map[string]string
	cookie      map[string]string
	noRedirect  bool
	timeout     time.Duration
	transport   *http.Transport
}

func New(url string) *context {
	return &context{
		url:         url,
		contentType: "application/json;charset=utf-8",
		timeout:     30 * time.Second,
	}
}

func (c *context) Header(k string, v string) *context {
	if c.header == nil {
		c.header = make(map[string]string)
	}
	c.header[k] = v
	return c
}

func (c *context) Cookie(k string, v string) *context {
	if c.cookie == nil {
		c.cookie = make(map[string]string)
	}
	c.cookie[k] = v
	return c
}

func (c *context) ContentType(contentType string) *context {
	c.contentType = contentType
	return c
}

func (c *context) Timeout(timeout time.Duration) *context {
	c.timeout = timeout
	return c
}

func (c *context) Transport(transport *http.Transport) *context {
	c.transport = transport
	return c
}

func (c *context) Redirect(Redirect bool) *context {
	c.noRedirect = !Redirect
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

func (c *context) Post() (resp *http.Response, err error) {
	client := &http.Client{
		Timeout: c.timeout,
	}
	if c.transport != nil {
		client.Transport = c.transport
	}

	if c.noRedirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	// PostData
	var pingString string
	if c.ping == nil {
		pingString = ""
	} else {
		if s, ok := c.ping.(string); ok {
			pingString = s
		} else {
			bytes, err := json.Marshal(c.ping)
			if err != nil {
				return nil, err
			}
			pingString = string(bytes)
		}
	}

	// New Request
	req, err := http.NewRequest("POST", c.url, strings.NewReader(pingString))

	// Set Cookie
	for k, v := range c.cookie {
		req.AddCookie(
			&http.Cookie{
				Name:     k,
				Value:    v,
				HttpOnly: true,
			},
		)
	}

	// Set Header, include contentType
	for k, v := range c.header {
		req.Header.Add(k, v)
	}
	req.Header.Add("Content-Type", c.contentType)

	//
	// Go
	//
	r, err := client.Do(req)
	if err != nil {
		klog.E(err.Error())
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		klog.E("%d", r.StatusCode)
		return r, fmt.Errorf("StatusCode == %d", r.StatusCode)
	}

	if c.pong == nil {
		return r, nil
	}

	if ptr, ok := c.pong.(*string); ok {
		if buf, err := ioutil.ReadAll(r.Body); err != nil {
			return r, err
		} else {
			*ptr = string(buf)
			return r, nil
		}
	} else if ptr, ok := c.pong.(*[]byte); ok {
		if buf, err := ioutil.ReadAll(r.Body); err != nil {
			return r, err
		} else {
			*ptr = buf
			return r, nil
		}
	} else {
		return r, json.NewDecoder(r.Body).Decode(c.pong)
	}
}

// Get : HTTPGet convert the response to pongObj structure
func (c *context) Get() (resp *http.Response, err error) {
	client := &http.Client{}

	if c.noRedirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	// New Request
	req, err := http.NewRequest("GET", c.url, nil)

	// Set Cookie
	for k, v := range c.cookie {
		req.AddCookie(
			&http.Cookie{
				Name:     k,
				Value:    v,
				HttpOnly: true,
			},
		)
	}

	// Set Header, include contentType
	for k, v := range c.cookie {
		req.Header.Add(k, v)
	}
	req.Header.Add("Content-Type", c.contentType)

	//
	// Go
	//
	r, err := client.Do(req)
	if err != nil {
		klog.E(err.Error())
		return r, err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		klog.E("URL:%s, Code:%d", c.url, r.StatusCode)
		return r, fmt.Errorf("StatusCode == %d", r.StatusCode)
	}

	if c.pong == nil {
		return r, nil
	}

	if ptr, ok := c.pong.(*string); ok {
		if buf, err := ioutil.ReadAll(r.Body); err != nil {
			return r, err
		} else {
			*ptr = string(buf)
			return r, nil
		}
	} else if ptr, ok := c.pong.(*[]byte); ok {
		if buf, err := ioutil.ReadAll(r.Body); err != nil {
			return r, err
		} else {
			*ptr = buf
			return r, nil
		}
	} else {
		return r, json.NewDecoder(r.Body).Decode(c.pong)
	}
}
