package roast

import (
	"database/sql"
	"fmt"
	"strconv"

	"miego/xmap"
	"miego/xtime"
	"miego/xvx/nvn"
)

func Str__Time(mp xmap.Map, field string) {
	mp.Put(field+"__Str", xtime.NumToStr(mp.S(field)))
}
func Str__Nvn(mp xmap.Map, field string, nvnClass string) {
	mp.Put(field+"__Str", nvn.S(mp.S(field), nvnClass))
}

// ///////////////////////////////////////////////////////////////////////
// For Query:
// 1. PageSize: default to 10
// 2. ID: clear RemAt
// 3. UUID: default to last
func PresetQuery(mp xmap.Map, Prefix string) string {
	var preset string

	// If set ID,
	if mp.Has("ID") {
		preset = ""
	} else if !mp.Has("RemAt") {
		preset = Prefix + ".RemAt = 0"
	}
	if mp.Has("UUID") {
		mp.Put(
			"PageSize", "1",
			"OrderBy", "ID",
			"OrderDir", "desc",
		)
	}
	return preset
}

func Int64(a string, def int64) int64 {
	x, e := strconv.ParseInt(a, 0, 64)
	if e != nil {
		return def
	}
	return x
}

func NumParse(numBase int64, vStr string) int64 {
	if len(vStr) > 2 && vStr[1] == '=' {
		switch vStr[0] {
		case '+':
			numBase += Int64(vStr[2:], 0)
		case '-':
			numBase -= Int64(vStr[2:], 0)
		case '*':
			numBase *= Int64(vStr[2:], 0)
		case '/':
			numBase /= Int64(vStr[2:], 0)
		}
	} else {
		numBase = Int64(vStr, 0)
	}
	return numBase
}

// 获取某个表最后更新（包括添加记录和删除记录）的时间
func LastModifyTime(dbSQL *sql.DB, TableName string) uint64 {
	stmtLastOp := fmt.Sprintf("SELECT GREATEST(NewAt, RemAt) AS Last FROM `%s` LIMIT 1", TableName)
	mpList, _ := Raw(dbSQL, stmtLastOp, false)
	if len(mpList) > 0 {
		return mpList[0].AsUint("Last", 0)
	}
	return 0
}
