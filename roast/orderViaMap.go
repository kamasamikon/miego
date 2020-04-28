package roast

import (
	"fmt"
	"strings"

	"github.com/kamasamikon/miego/xmap"
)

// Use in Ping query

func OrderLimitOffset2(m xmap.Map) (orderBy string, limit uint, offset uint) {
	pageSize := m.Int("PageSize", 50)
	pageNumber := m.Int("PageNumber", 1)
	limit = uint(pageSize)

	offset = uint(m.Int("Offset", pageSize*(pageNumber-1)))

	orderBy = ""
	cmdOrderBy := m.Str("OrderBy", "")
	if cmdOrderBy != "" {
		segs := strings.Split(cmdOrderBy, "__")
		orderBy = fmt.Sprintf("ORDER BY %s", segs[0])

		cmdOrderDir := m.Str("OrderDir", "")
		if cmdOrderDir == "desc" {
			orderBy += " DESC"
		} else if cmdOrderDir == "asc" {
			orderBy += " ASC"
		}
	}

	return orderBy, limit, offset
}
