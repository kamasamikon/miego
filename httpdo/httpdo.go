package httpdo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kamasamikon/miego/klog"
)

const mimeJSON = "application/json;charset=utf-8"

// Post : HTTPPost post json data to peer and convert the response to pongObj structure
func Post(URL string, pingObj interface{}, pongObj interface{}) (n int, err error) {
	var pingString string

	if pingObj == nil {
		pingString = ""
	} else {
		if s, ok := pingObj.(string); ok {
			pingString = s
		} else {
			bytes, err := json.Marshal(pingObj)
			if err != nil {
				return 0, err
			}
			pingString = string(bytes)
		}
	}

	r, err := http.Post(URL, mimeJSON, strings.NewReader(pingString))
	if err != nil {
		klog.E(err.Error())
		return 0, err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		klog.E("%d", r.StatusCode)
		return 0, fmt.Errorf("StatusCode == %d", r.StatusCode)
	}

	if pongObj == nil {
		return 0, nil
	}

	if ptr, ok := pongObj.(*string); ok {
		if buf, err := ioutil.ReadAll(r.Body); err != nil {
			return 0, err
		} else {
			*ptr = string(buf)
			return 0, nil
		}
	} else if ptr, ok := pongObj.(*[]byte); ok {
		return r.Body.Read(*ptr)
	} else {
		return 0, json.NewDecoder(r.Body).Decode(pongObj)
	}
}

// Get : HTTPGet convert the response to pongObj structure
func Get(URL string, pongObj interface{}) (n int, err error) {
	r, err := http.Get(URL)
	if err != nil {
		klog.E(err.Error())
		return 0, err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		klog.E("URL:%s, Code:%d", URL, r.StatusCode)
		return 0, fmt.Errorf("StatusCode == %d", r.StatusCode)
	}

	if pongObj == nil {
		return 0, nil
	}

	if ptr, ok := pongObj.(*string); ok {
		if buf, err := ioutil.ReadAll(r.Body); err != nil {
			return 0, err
		} else {
			*ptr = string(buf)
			return 0, nil
		}
	} else if ptr, ok := pongObj.(*[]byte); ok {
		return r.Body.Read(*ptr)
	} else {
		return 0, json.NewDecoder(r.Body).Decode(pongObj)
	}
}

// Download : Download and save.
func Download(URL string, filename string) error {
	res, err := http.Get(URL)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if err = ioutil.WriteFile(filename, data, 0777); err != nil {
		return err
	}
	return nil
}

// Upload : Upload file to remote
func Upload(URL string, params map[string]string, paramName, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", URL, body)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("StatusCode")
	}

	return nil
}
