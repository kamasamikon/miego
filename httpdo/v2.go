package httpdo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"miego/klog"
)

type context struct {
	url         string
	ping        interface{}
	pong        interface{}
	contentType string
	header      map[string]string
	cookies     []*http.Cookie
	noRedirect  bool
	timeout     time.Duration
	transport   *http.Transport
	lastErr     error // Last error
}

func New(url string) *context {
	return &context{
		url:         url,
		contentType: "application/json;charset=utf-8",
		timeout:     30 * time.Second,
	}
}

func (c *context) Header(k string, v string) *context {
	if c.lastErr != nil {
		return c
	}
	if c.header == nil {
		c.header = make(map[string]string)
	}
	c.header[k] = v
	return c
}

func (c *context) Headers(items map[string]string) *context {
	if c.lastErr != nil {
		return c
	}
	if c.header == nil {
		c.header = make(map[string]string)
	}
	for k, v := range items {
		c.header[k] = v
	}
	return c
}

func (c *context) Cookie(cookie *http.Cookie) *context {
	if c.lastErr != nil {
		return c
	}
	c.cookies = append(c.cookies, cookie)
	return c
}

func (c *context) ContentType(contentType string) *context {
	if c.lastErr != nil {
		return c
	}
	c.contentType = contentType
	return c
}

func (c *context) Timeout(timeout time.Duration) *context {
	if c.lastErr != nil {
		return c
	}
	c.timeout = timeout
	return c
}

func (c *context) Transport(transport *http.Transport) *context {
	if c.lastErr != nil {
		return c
	}
	c.transport = transport
	return c
}

func (c *context) Redirect(Redirect bool) *context {
	if c.lastErr != nil {
		return c
	}
	c.noRedirect = !Redirect
	return c
}

func (c *context) Ping(ping interface{}) *context {
	if c.lastErr != nil {
		return c
	}
	c.ping = ping
	return c
}

func (c *context) Pong(pong interface{}) *context {
	if c.lastErr != nil {
		return c
	}
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
			// as string
			pingString = s
		} else if s, ok := c.ping.([]byte); ok {
			// as []byte
			pingString = string(s)
		} else {
			// as json object
			bytes, err := json.Marshal(c.ping)
			if err != nil {
				return nil, err
			}
			pingString = string(bytes)
		}
	}

	// New Request
	req, err := http.NewRequest("POST", c.url, strings.NewReader(pingString))
	if err != nil {
		c.lastErr = err
		return nil, err
	}

	// Set Cookie
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
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
		if buf, err := io.ReadAll(r.Body); err != nil {
			return r, err
		} else {
			*ptr = string(buf)
			return r, nil
		}
	} else if ptr, ok := c.pong.(*[]byte); ok {
		if buf, err := io.ReadAll(r.Body); err != nil {
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

	// New Request
	req, err := http.NewRequest("GET", c.url, nil)
	if err != nil {
		c.lastErr = err
		return nil, err
	}

	// Set Cookie
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
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
		if buf, err := io.ReadAll(r.Body); err != nil {
			return r, err
		} else {
			*ptr = string(buf)
			return r, nil
		}
	} else if ptr, ok := c.pong.(*[]byte); ok {
		if buf, err := io.ReadAll(r.Body); err != nil {
			return r, err
		} else {
			*ptr = buf
			return r, nil
		}
	} else {
		return r, json.NewDecoder(r.Body).Decode(c.pong)
	}
}
