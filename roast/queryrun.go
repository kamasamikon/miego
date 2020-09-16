package roast

import (
	"github.com/kamasamikon/miego/atox"
	"github.com/kamasamikon/miego/klog"
	"github.com/kamasamikon/miego/xmap"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// Query : dsn:dsn prefix in conf
func Query(dsn string, queryStmt *QueryStatement, mp xmap.Map, FoundRows bool, rowScanner func(*sql.Rows) bool) (int, error) {
	//
	// Go
	//
	rdb, err := sql.Open("mysql", DSN(dsn))
	if err != nil {
		klog.E(err.Error())
		return 0, err
	}
	defer rdb.Close()

	rows, err := rdb.Query(queryStmt.String(mp, FoundRows))
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
	res, err := rdb.Query(`SELECT FOUND_ROWS()`)
	if err != nil {
		klog.E(err.Error())
		return 0, err
	}

	res.Next()

	var lines string
	if err := res.Scan(&lines); err != nil {
		klog.E(err.Error())
		return 0, err
	}

	return atox.Int(lines, 0), nil
}
