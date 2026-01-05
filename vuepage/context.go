package vuepage

import (
	"miego/xmap"
)

type Context xmap.Map

// Put : c[Name] = Value
// If Name == "vueData", Value must be a xmap
func (c Context) Put(Name string, Value interface{}) {
	if Name != "vueData" {
		c[Name] = Value
	} else {
		vueData := Value.(xmap.Map)
		for k, v := range vueData {
			c.PutVue(k, v)
		}
	}
}

// PutMap : 把 Name=Value 放到 MapName 中去
func (c Context) PutMap(MapName string, Name string, Value interface{}) {
	//
	// 如果List不存在就创建
	//
	subMap := c[MapName]
	if subMap == nil {
		subMap = xmap.Map{}
		c[MapName] = subMap
	}

	// 这下有了
	if subXMap, ok := subMap.(xmap.Map); ok {
		//
		// 放到MapName指定的List中去
		//
		subXMap[Name] = Value
	}
}

// PutList : 把 Value 放到 ListName 中去
func (c Context) PutList(ListName string, Value interface{}) {
	subObject := c[ListName]
	if subObject == nil {
		subObject = []interface{}{}
	}
	subList := subObject.([]interface{})
	subList = append(subList, Value)
	c[ListName] = subList
}

// PutVue : c["vueData"][Name] = Value
func (c Context) PutVue(Name string, Value interface{}) {
	var vueData xmap.Map
	if tmp := c["vueData"]; tmp == nil {
		vueData = xmap.Make()
	} else {
		vueData = tmp.(xmap.Map)
	}

	vueData.Put(Name, Value)
	c["vueData"] = vueData
}
