package roast

import (
	"fmt"
	"strings"
)

// AND : "AND a = 3 AND b = 4
func AND(qList []string) string {
	var andClause string
	if qList != nil {
		andClause = fmt.Sprintf("AND %s", strings.Join(qList, " AND "))
	}
	return andClause
}

// WHERE : WHERE a = 3 AND b = 4
func WHERE(qList []string) string {
	var whereClause string
	if qList != nil {
		whereClause = fmt.Sprintf("WHERE %s", strings.Join(qList, " AND "))
	}
	return whereClause
}
