package roast

// Get First or Last one

import (
	"database/sql"
	"fmt"

	"miego/xmap"
)

// Last : Get the latest row
func Last(db *sql.DB, mp xmap.Map, Query func(db *sql.DB, mp xmap.Map) []xmap.Map) xmap.Map {
	mp.SafePut(
		"PageSize", "1",
		"OrderBy", "ID",
		"OrderDir", "desc",
	)
	arr := Query(db, mp)
	if arr != nil {
		return arr[0]
	} else {
		return nil
	}
}

// Last : Get the oldest row
func First(db *sql.DB, mp xmap.Map, Query func(db *sql.DB, mp xmap.Map) []xmap.Map) xmap.Map {
	mp.SafePut(
		"PageSize", "1",
		"OrderBy", "ID",
		"OrderDir", "asc",
	)
	arr := Query(db, mp)
	if arr != nil {
		return arr[0]
	} else {
		return nil
	}
}

func LastOne(db *sql.DB, TableName string, FieldName string) string {
	stmtfmt := "SELECT %s FROM `%s` ORDER BY `%s` DESC LIMIT 1"
	stmt := fmt.Sprintf(stmtfmt, FieldName, TableName, FieldName)
	rows, err := db.Query(stmt)
	if err != nil {
		return ""
	}
	defer rows.Close()

	rows.Next()
	var UUID string
	err = rows.Scan(&UUID)
	if err != nil {
		return ""
	}
	return UUID
}

func FirstOne(db *sql.DB, TableName string, FieldName string) string {
	stmtfmt := "SELECT %s FROM `%s` ORDER BY `%s` ASC LIMIT 1"
	stmt := fmt.Sprintf(stmtfmt, FieldName, TableName, FieldName)
	rows, err := db.Query(stmt)
	if err != nil {
		return ""
	}
	defer rows.Close()

	rows.Next()
	var UUID string
	err = rows.Scan(&UUID)
	if err != nil {
		return ""
	}
	return UUID
}
