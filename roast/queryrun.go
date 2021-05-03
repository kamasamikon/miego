package roast

import (
	"github.com/kamasamikon/miego/atox"
	"github.com/kamasamikon/miego/klog"
	"github.com/kamasamikon/miego/xmap"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type preQ func(mp xmap.Map)
type pstQ func(mp xmap.Map, pongs []xmap.Map)

func callPreQ(mp xmap.Map) {
	if fp, ok := mp["__PreQ__"]; ok {
		if f, ok := fp.(preQ); ok {
			f(mp)
		}
	}
}
func callPstQ(mp xmap.Map, pongs []xmap.Map) {
	if fp, ok := mp["__PstQ__"]; ok {
		if f, ok := fp.(pstQ); ok {
			f(mp, pongs)
		}
	}
}

// Raw: Raw(db, "SELECT a, b FROM xxxx") => {"a":"xxx", "b":"xxx"}
func Raw(db *sql.DB, stmt string) ([]xmap.Map, error) {
	klog.D("%s", stmt)
	klog.BT(6, "--------------")
	rows, err := db.Query(stmt)
	if err != nil {
		klog.E(err.Error())
		return nil, err
	}
	defer rows.Close()

	//
	// Prepare
	//
	columns, _ := rows.Columns()
	columnLength := len(columns)
	cache := make([]interface{}, columnLength)
	for index, _ := range cache {
		var a interface{}
		cache[index] = &a
	}

	//
	// Get result
	//
	var pongs []xmap.Map
	for rows.Next() {
		err := rows.Scan(cache...)
		if err != nil {
			klog.E(err.Error())
			continue
		}

		item := make(xmap.Map)
		for i, data := range cache {
			x := *data.(*interface{})
			y := x.([]uint8)
			item[columns[i]] = string(y)
		}
		pongs = append(pongs, item)
	}

	//
	// Done
	//
	klog.D("%d items found", len(pongs))
	return pongs, nil
}

// ViaMap : Return Result and allCount
func ViaMap(db *sql.DB, queryStmt *QueryStatement, mp xmap.Map, FoundRows int) ([]xmap.Map, int, error) {
	callPreQ(mp)

	qStmt, cStmt := queryStmt.String2(mp, FoundRows)

	rows, err := db.Query(qStmt)
	if err != nil {
		klog.E(err.Error())
		return nil, -1, err
	}
	defer rows.Close()

	//
	// Prepare
	//
	columns, _ := rows.Columns()
	columnLength := len(columns)
	cache := make([]interface{}, columnLength)
	for index, _ := range cache {
		var a interface{}
		cache[index] = &a
	}

	//
	// Get result
	//
	var pongs []xmap.Map
	for rows.Next() {
		err := rows.Scan(cache...)
		if err != nil {
			klog.E(err.Error())
			continue
		}

		item := make(xmap.Map)
		for i, data := range cache {
			x := *data.(*interface{})
			y := x.([]uint8)
			item[columns[i]] = string(y)
		}
		pongs = append(pongs, item)
	}

	//
	// Get Rows count
	//
	allCount := 0
	if cStmt != "" {
		res, err := db.Query(cStmt)
		if err == nil {
			defer res.Close()

			res.Next()
			var lines string
			if err := res.Scan(&lines); err == nil {
				allCount = atox.Int(lines, 0)
			} else {
				klog.E(err.Error())
			}
		} else {
			klog.E(err.Error())
		}
		klog.D("AllCount:%d", allCount)
	} else {
		allCount = len(pongs)
		klog.D("allCount:%d", allCount)
	}

	//
	// Done
	//
	callPstQ(mp, pongs)
	return pongs, allCount, nil
}

func ViaScanner(db *sql.DB, queryStmt *QueryStatement, mp xmap.Map, FoundRows int, rowScanner func(*sql.Rows) bool) (int, error) {
	qStmt, cStmt := queryStmt.String2(mp, FoundRows)
	rows, err := db.Query(qStmt)
	if err != nil {
		klog.E(err.Error())
		return 0, err
	}
	defer rows.Close()

	//
	// Next: Scan each rows
	//
	// rowScanner = {
	// 	 var aaa string
	// 	 rows.Scan(&aaa)
	// 	 ...
	// }
	//
	for rows.Next() {
		if ok := rowScanner(rows); !ok {
			break
		}
	}

	allCount := 0
	if cStmt != "" {
		res, err := db.Query(cStmt)
		if err == nil {
			defer res.Close()

			res.Next()
			var lines string
			if err := res.Scan(&lines); err == nil {
				allCount = atox.Int(lines, 0)
			} else {
				klog.E(err.Error())
			}
		} else {
			klog.E(err.Error())
		}
	}

	return allCount, nil
}

func Query(db *sql.DB, dsn string, prefix string, queryStmt *QueryStatement, mp xmap.Map, FoundRows int, rowScanner func(*sql.Rows) bool) (int, error) {
	var err error

	if db == nil {
		if dsn != "" {
			db, err = sql.Open("mysql", dsn)
		} else {
			db, err = sql.Open("mysql", DSN(prefix))
		}
		if err != nil {
			klog.E(err.Error())
			return 0, err
		}
		defer db.Close()
	}

	return ViaScanner(db, queryStmt, mp, FoundRows, rowScanner)
}

func QueryToMap(db *sql.DB, dsn string, prefix string, queryStmt *QueryStatement, mp xmap.Map, FoundRows int) ([]xmap.Map, int, error) {
	var err error

	if db == nil {
		if dsn != "" {
			db, err = sql.Open("mysql", dsn)
		} else {
			db, err = sql.Open("mysql", DSN(prefix))
		}
		if err != nil {
			klog.E(err.Error())
			return nil, 0, err
		}
		defer db.Close()
	}

	return ViaMap(db, queryStmt, mp, FoundRows)
}
