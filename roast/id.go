package roast

import (
	"fmt"
	"strconv"

	"github.com/twinj/uuid"

	"database/sql"
)

func IDNew() string {
	return uuid.NewV4().String()
}

func IDNxt(db *sql.DB, TableName string) uint64 {
	stmtfmt := "SELECT MAX(id) AS UUID FROM `%s`"
	stmt := fmt.Sprintf(stmtfmt, TableName)
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
