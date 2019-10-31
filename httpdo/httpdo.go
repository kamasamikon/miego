package httpdo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/kamasamikon/miego/klog"
)

const mimeJSON = "application/json;charset=utf-8"

// Post : HTTPPost post json data to peer and convert the response to pongObj structure
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

	r, eb := http.Post(url, mimeJSON, strings.NewReader(pingString))
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

	json.NewDecoder(r.Body).Decode(pongObj)
	return nil
}

// Get : HTTPGet convert the response to pongObj structure
func Get(url string, pongObj interface{}) error {
	r, eb := http.Get(url)
	if eb != nil {
		klog.E(eb.Error())
		return eb
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		klog.E("URL:%s, Code:%d", url, r.StatusCode)
		return errors.New(fmt.Sprintf("StatusCode == %d", r.StatusCode))
	}

	if pongObj == nil {
		return nil
	}
	return json.NewDecoder(r.Body).Decode(pongObj)
}


// Download : Download and save.
func Download(url string, filename string) {
	res, err := http.Get(url)
	if err != nil {
		klog.E("http.Get -> %v", err)
		return
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		klog.E("ioutil.ReadAll -> %s", err.Error())
		return
	}
	defer res.Body.Close()
	if err = ioutil.WriteFile(filename, data, 0777); err != nil {
		klog.E("Error Saving:", filename, err)
	} else {
		klog.D("%s => %s:", url, filename)
	}
}
