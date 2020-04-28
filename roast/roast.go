package roast

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/kamasamikon/miego/in"
	"github.com/kamasamikon/miego/xmap"
	"github.com/kamasamikon/miego/xtime"
)

// GE, GT, LE, LT, EQ, NE,
// GE_DATE, GT_DATE, LE_DATE, LT_DATE, EQ_DATE, NE_DATE,
// IN, NI
// IS, LIKE
// NULL, NNULL

// Query: Name__IN=aaa&Name__IN=bbb			// 这个可以重名
// Name: { "IS/LIKE": ["AAA", "BBB"] }

// Map: String and Array, e.g. "LIKE": ["AAA", "BBB"]
type mapSA map[string][]string

//
// 查询JSON: [{Name:Name, Kind:Str}, {Name:Gender, Kind:Choice,Int, Choice:[{Name=男, Var=M},{}]]
//
// 都是 AND, 如果是OR，应该使用 IN
//
// "Name": {"IS": "Auv", "IS": "Hilda", "LIKE": "ngs"} 	=> Where Name = "auv" AND Name = "Hilda" AND Name LIKE "ngs"
// "Gender": {"IN": ["M", "F"]}

// Key: IS, IN, LE, ....
// "Name": {"IS": ["Auv", "Hilda"], "LIKE": ["ngs"]} 	=> Where Name = "auv" AND Name = "Hilda" AND Name LIKE "ngs"
// Key: Name
type QueryMap map[string]mapSA

func QueryMapNew(mpList ...xmap.Map) QueryMap {
	m := make(QueryMap)
	for _, mp := range mpList {
		if mp != nil {
			m.Parse(mp)
		}
	}
	return m
}

func (m QueryMap) Parse(mp xmap.Map) error {
	// k: Name__IS, v:Auv
	// k: Name__IN__0, v:Auv
	for k, v := range mp {
		var s string

		switch v.(type) {
		case bool:
			if t, _ := v.(bool); t {
				s = "true"
			} else {
				s = "false"
			}

		case string:
			if t, _ := v.(string); t != "" {
				s = t
			}
		}
		if s == "" {
			continue
		}

		var Name string
		var Kind string

		// Get :Name and :Kind
		segs := strings.Split(k, "__")
		if len(segs) == 0 {
			continue
		}
		Name = segs[0]
		if len(segs) > 1 {
			Kind = segs[1]
		} else {
			Kind = "GUESS"
		}

		// Save
		sa, ok := m[Name]
		if ok == false {
			sa = make(map[string][]string)
		}

		// if Kind == "IN", Status__IN__Open=true => Status: {"IN", "Open"}
		if Kind == "IN" {
			if in.C(s[0], 'T', 't', 'Y', 'y', '1') {
				Kind = "IN"
				s = segs[2]
			} else {
				continue
			}
		}

		ar, _ := sa[Kind]
		ar = append(ar, s)
		sa[Kind] = ar
		m[Name] = sa
	}

	return nil
}

func (m QueryMap) Has(Name string) bool {
	_, ok := m[Name]
	return ok
}

func (m QueryMap) Use(qList []string, Name string, Table string, NewName string) []string {
	sa, ok := m[Name]
	if ok == false {
		return qList
	}
	if NewName != "" {
		Name = NewName
	}
	if Table != "" {
		Name = Table + "." + Name
	}

	p := fmt.Sprintf

	for kind, arr := range sa {
		if len(arr) == 0 {
			continue
		}

		switch kind {
		case "GE":
			for _, v := range arr {
				qList = append(qList, p(`%s >= "%s"`, Name, v))
			}

		case "GT":
			for _, v := range arr {
				qList = append(qList, p(`%s > "%s"`, Name, v))
			}

		case "LE":
			for _, v := range arr {
				qList = append(qList, p(`%s <= "%s"`, Name, v))
			}

		case "LT":
			for _, v := range arr {
				qList = append(qList, p(`%s < "%s"`, Name, v))
			}

		case "EQ":
			for _, v := range arr {
				qList = append(qList, p(`%s = "%s"`, Name, v))
			}

		case "NE":
			for _, v := range arr {
				qList = append(qList, p(`%s != "%s"`, Name, v))
			}

		case "GE_DATE":
			for _, v := range arr {
				qList = append(qList, p(`%s >= DATE("%s")`, Name, v))
			}

		case "GT_DATE":
			for _, v := range arr {
				qList = append(qList, p(`%s > DATE("%s")`, Name, v))
			}

		case "LE_DATE":
			for _, v := range arr {
				qList = append(qList, p(`%s <= DATE("%s")`, Name, v))
			}

		case "LT_DATE":
			for _, v := range arr {
				qList = append(qList, p(`%s < DATE("%s")`, Name, v))
			}

		case "EQ_DATE":
			for _, v := range arr {
				qList = append(qList, p(`%s = DATE("%s")`, Name, v))
			}

		case "NE_DATE":
			for _, v := range arr {
				qList = append(qList, p(`%s != DATE("%s")`, Name, v))
			}

		case "GE_XDATE":
			for _, v := range arr {
				qList = append(qList, p(`%s >= %d`, Name, xtime.StrToNum(v)))
			}

		case "GT_XDATE":
			for _, v := range arr {
				qList = append(qList, p(`%s > %d`, Name, xtime.StrToNum(v)))
			}

		case "LE_XDATE":
			for _, v := range arr {
				qList = append(qList, p(`%s <= %d`, Name, xtime.StrToNum(v)))
			}

		case "LT_XDATE":
			for _, v := range arr {
				qList = append(qList, p(`%s < %d`, Name, xtime.StrToNum(v)))
			}

		case "EQ_XDATE":
			for _, v := range arr {
				qList = append(qList, p(`%s = %d`, Name, xtime.StrToNum(v)))
			}

		case "NE_XDATE":
			for _, v := range arr {
				qList = append(qList, p(`%s != %d`, Name, xtime.StrToNum(v)))
			}

		case "IN":
			if len(arr) == 1 {
				qList = append(qList, p(`%s = '%s'`, Name, arr[0]))
			} else {
				qList = append(qList, p(`%s IN ("%s")`, Name, strings.Join(arr, `", "`)))
			}

		case "NI":
			if len(arr) == 1 {
				qList = append(qList, p(`%s != "%s"`, Name, arr[0]))
			} else {
				qList = append(qList, p(`%s NOT IN ("%s")`, Name, strings.Join(arr, `", "`)))
			}

		case "IS":
			for _, v := range arr {
				qList = append(qList, p(`%s = "%s"`, Name, v))
			}

		case "LIKE":
			for _, v := range arr {
				qList = append(qList, p(`%s LIKE "%s"`, Name, v))
			}

		case "NULL":
			for _, v := range arr {
				c := v[0]
				if in.C(c, 'T', 't', 'Y', 'y', '1') {
					qList = append(qList, p(`%s IS NULL`, Name))
				} else if in.C(c, 'F', 'f', 'N', 'n', '0') {
					qList = append(qList, p(`%s IS NOT NULL`, Name))
				}
			}

		case "NNULL":
			for _, v := range arr {
				c := v[0]
				if in.C(c, 'T', 't', 'Y', 'y', '1') {
					qList = append(qList, p(`%s IS NOT NULL`, Name))
				} else if in.C(c, 'F', 'f', 'N', 'n', '0') {
					qList = append(qList, p(`%s IS NULL`, Name))
				}
			}

		case "GUESS":
			for _, v := range arr {
				if s, like := likeParse(v, false); like {
					qList = append(qList, fmt.Sprintf(`%s LIKE "%s"`, Name, s))
				} else {
					qList = append(qList, fmt.Sprintf(`%s = "%s"`, Name, s))
				}
			}
		}
	}

	return qList
}

func (m QueryMap) Dump() {
	spew.Dump(m)
}
