package roast

import (
	"github.com/kamasamikon/miego/atox"
	"github.com/kamasamikon/miego/klog"
	"github.com/kamasamikon/miego/xmap"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// ViaMap : Return Result and allCount
func ViaMap(db *sql.DB, queryStmt *QueryStatement, mp xmap.Map, FoundRows bool) ([]xmap.Map, int, error) {
	rows, err := db.Query(queryStmt.String(mp, FoundRows))
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

	if !FoundRows {
		return pongs, -1, nil
	}

	//
	// Get Rows count
	//
	res, err := db.Query(`SELECT FOUND_ROWS()`)
	if err != nil {
		klog.E(err.Error())
		return pongs, -1, nil
	}
	defer res.Close()

	res.Next()

	var lines string
	if err := res.Scan(&lines); err != nil {
		klog.E(err.Error())
		return pongs, -1, nil
	}

	//
	// Done
	//
	return pongs, atox.Int(lines, 0), nil
}

func ViaScanner(db *sql.DB, queryStmt *QueryStatement, mp xmap.Map, FoundRows bool, rowScanner func(*sql.Rows) bool) (int, error) {
	rows, err := db.Query(queryStmt.String(mp, FoundRows))
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

	if !FoundRows {
		return 0, nil
	}

	//
	// Get Rows count
	//
	res, err := db.Query(`SELECT FOUND_ROWS()`)
	if err != nil {
		klog.E(err.Error())
		return 0, err
	}
	defer res.Close()

	res.Next()

	var lines string
	if err := res.Scan(&lines); err != nil {
		klog.E(err.Error())
		return 0, err
	}

	return atox.Int(lines, 0), nil
}

func Query(dsn string, prefix string, queryStmt *QueryStatement, mp xmap.Map, FoundRows bool, rowScanner func(*sql.Rows) bool) (int, error) {
	var db *sql.DB
	var err error

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

	return ViaScanner(db, queryStmt, mp, FoundRows, rowScanner)
}

func QueryToMap(dsn string, prefix string, queryStmt *QueryStatement, mp xmap.Map, FoundRows bool) ([]xmap.Map, int, error) {
	var db *sql.DB
	var err error

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

	return ViaMap(db, queryStmt, mp, FoundRows)
}
