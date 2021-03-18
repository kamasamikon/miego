package roast

import (
	"fmt"
	"strconv"

	"github.com/twinj/uuid"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func IDNew() string {
	return uuid.NewV4().String()
}

func IDNxt(db *sql.DB, TableName string, where string) uint64 {
	stmtfmt := "SELECT UUID FROM `%s` %s ORDER BY `UUID` DESC LIMIT 1"
	if where != "" {
		where = " WHERE " + where
	}
	stmt := fmt.Sprintf(stmtfmt, TableName, where)
	rows, err := db.Query(stmt)
	if err != nil {
		return 1
	}
	defer rows.Close()
	for rows.Next() {
		var UUID string
		err := rows.Scan(
			&UUID,
		)
		if err != nil {
			return 1
		}
		if num, err := strconv.ParseUint(UUID, 10, 64); err != nil {
			return 1
		} else {
			return num + 1
		}
	}
	return 1
}
