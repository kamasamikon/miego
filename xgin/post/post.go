package post

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type PostMap map[string]interface{}

func Map(c *gin.Context) (PostMap, error) {
	if dat, err := ioutil.ReadAll(c.Request.Body); err != nil {
		return nil, err
	} else {
		var m PostMap
		json.Unmarshal(dat, &m)
		return m, nil
	}
}
