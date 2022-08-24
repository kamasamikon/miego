package klogin

import "github.com/kamasamikon/miego/xmap"

func LoginData(keys ...string) []string {
	var arr []string

	if c.ContentType() == "application/json" {
		mp := xmap.MapBody(c)
		for key := range keys {
			arr = append(arr, mp.S(key))
		}
	} else {
		for key := range keys {
			arr = append(arr, c.PostForm(key))
		}
	}

	return arr
}
