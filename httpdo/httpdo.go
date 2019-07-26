package httpdo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/kamasamikon/miego/klog"
)

const MIMEJSON = "application/json;charset=utf-8"

// HTTPPost post json data to peer and convert the response to pongObj structure
func Post(url string, pingObj interface{}, pongObj interface{}) error {
	var pingString string

	if pingObj == nil {
		pingString = ""
	} else {
		if s, ok := pingObj.(string); ok {
			pingString = s
		} else {
			bytes, ea := json.Marshal(pingObj)
			if ea != nil {
				return ea
			}
			pingString = string(bytes)
		}
	}
	klog.D(pingString)

	r, eb := http.Post(url, MIMEJSON, strings.NewReader(pingString))
	if eb != nil {
		klog.E(eb.Error())
		return eb
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		klog.E("%d", r.StatusCode)
		return fmt.Errorf("StatusCode == %d", r.StatusCode)
	}

	if pongObj == nil {
		return nil
	}
	return json.NewDecoder(r.Body).Decode(pongObj)
}

// HTTPGet convert the response to pongObj structure
func Get(url string, pongObj interface{}) error {
	r, eb := http.Get(url)
	if eb != nil {
		klog.E(eb.Error())
		return eb
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		klog.E("%d", r.StatusCode)
		return errors.New(fmt.Sprintf("StatusCode == %d", r.StatusCode))
	}

	if pongObj == nil {
		return nil
	}
	return json.NewDecoder(r.Body).Decode(pongObj)
}
