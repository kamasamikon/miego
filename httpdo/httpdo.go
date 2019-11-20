package httpdo

import (
	"bytes"
	"encoding/json"
	"errors"
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
func Post(URL string, pingObj interface{}, pongObj interface{}) error {
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

	r, eb := http.Post(URL, mimeJSON, strings.NewReader(pingString))
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
func Get(URL string, pongObj interface{}) error {
	r, eb := http.Get(URL)
	if eb != nil {
		klog.E(eb.Error())
		return eb
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		klog.E("URL:%s, Code:%d", URL, r.StatusCode)
		return errors.New(fmt.Sprintf("StatusCode == %d", r.StatusCode))
	}

	if pongObj == nil {
		return nil
	}
	return json.NewDecoder(r.Body).Decode(pongObj)
}

// Download : Download and save.
func Download(URL string, filename string) {
	res, err := http.Get(URL)
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
		klog.D("%s => %s:", URL, filename)
	}
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
