package roast

import (
	"fmt"
	"strings"

	"github.com/kamasamikon/miego/klog"
	"github.com/kamasamikon/miego/xmap"
)

const (
	// FoundRows: Auto = Yes if mp.Has("PageSize") else No
	FR_Auto = 0
	FR_Yes  = 1
	FR_No   = 2
)

type Column struct {
	TableAlias string
	Field      string
	Default    string
	OutputName string
}

type JoinInfo struct {
	LineAnd string
	LineNon string
	AndList []string
}

type WhereInfo struct {
	Preset  string
	AndList []string
}

type QueryStatement struct {
	ColumnList []*Column

	JoinList  []*JoinInfo
	WhereLine *WhereInfo

	FromTable string
	FromAlias string

	GroupByLine string

	OrderBy string
	Limit   uint
	Offset  uint
}

func OrderLimitOffset2(m xmap.Map) (orderBy string, limit uint, offset uint) {
	pageSize := m.Int("PageSize", 50)
	pageNumber := m.Int("PageNumber", 1)
	limit = uint(pageSize)

	offset = uint(m.Int("Offset", pageSize*(pageNumber-1)))

	orderBy = ""
	cmdOrderBy := m.S("OrderBy")
	if cmdOrderBy != "" {
		segs := strings.Split(cmdOrderBy, "__")

		if cast := m.S("OrderCast"); cast != "" {
			orderBy = fmt.Sprintf("ORDER BY CAST(%s)", cast)
		} else {
			orderBy = fmt.Sprintf("ORDER BY %s", segs[0])
		}

		cmdOrderDir := m.S("OrderDir")
		if cmdOrderDir == "desc" {
			orderBy += " DESC"
		} else if cmdOrderDir == "asc" {
			orderBy += " ASC"
		}
	}

	return orderBy, limit, offset
}

func AND(qList []string) string {
	var andClause string
	if qList != nil {
		andClause = strings.Join(qList, " AND ")
	}
	return andClause
}

func QueryStatementNew() *QueryStatement {
	return &QueryStatement{}
}

func (s *QueryStatement) Column(TableAlias string, Field string, Default string, OutputName string) {
	if OutputName == "" {
		OutputName = Field
	}
	s.ColumnList = append(s.ColumnList, &Column{
		TableAlias: TableAlias,
		Field:      Field,
		Default:    Default,
		OutputName: OutputName,
	})
}

func (s *QueryStatement) ColumnClear() {
	s.ColumnList = nil
}

func (s *QueryStatement) Match(LineAnd string, LineNon string, AndList []string) {
	s.JoinList = append(s.JoinList, &JoinInfo{
		LineAnd: LineAnd,
		LineNon: LineNon,
		AndList: AndList,
	})
}
func (s *QueryStatement) Where(Preset string, AndList []string) {
	s.WhereLine = &WhereInfo{
		Preset:  Preset,
		AndList: AndList,
	}
}

// s.GroupBy("GROUP BY AAA")
func (s *QueryStatement) GroupBy(GroupByLine string) {
	s.GroupByLine = GroupByLine
}

func (s *QueryStatement) From(FromTable string, FromAlias string) {
	s.FromTable = FromTable
	s.FromAlias = FromAlias
}

func (s *QueryStatement) String(mp xmap.Map, FoundRows int) string {
	var lines []string

	// HEAD
	lines = append(lines, "SELECT")

	if FoundRows == FR_Auto {
		if mp.Has("PageNumber") {
			FoundRows = FR_Yes
		}
	}
	if FoundRows == FR_Yes {
		lines = append(lines, "SQL_CALC_FOUND_ROWS")
	}

	// Output fields
	var sublines []string
	for _, c := range s.ColumnList {
		var sss string
		if c.Field[0] == '@' {
			field := c.Field[1:len(c.Field)]
			if c.TableAlias == "" {
				sss = fmt.Sprintf(`    DISTINCT(%s) AS %s`, field, c.OutputName)
			} else {
				sss = fmt.Sprintf(`    DISTINCT(%s.%s) AS %s`, c.TableAlias, field, c.OutputName)
			}
		} else {
			if c.TableAlias == "" {
				sss = fmt.Sprintf(`    IFNULL(%s, "%s") AS %s`, c.Field, c.Default, c.OutputName)
			} else {
				sss = fmt.Sprintf(`    IFNULL(%s.%s, "%s") AS %s`, c.TableAlias, c.Field, c.Default, c.OutputName)
			}
		}
		sublines = append(sublines, sss)
	}
	lines = append(lines, strings.Join(sublines, ",\n"))

	//
	// From: SELECT * __FROM__ TableA
	//
	if s.FromAlias != "" {
		lines = append(lines, "FROM "+s.FromTable+" AS "+s.FromAlias)
	} else {
		lines = append(lines, "FROM "+s.FromTable)
	}

	//
	// Join
	//
	for _, j := range s.JoinList {
		if j.AndList != nil {
			userAdd := AND(j.AndList)

			if j.LineAnd == "" {
				lines = append(lines, "    "+userAdd)
			} else {
				lines = append(lines, "    "+j.LineAnd+" AND "+userAdd)
			}
		} else {
			lines = append(lines, "    "+j.LineNon)
		}
	}

	//
	// Where
	//
	w := s.WhereLine
	if w.AndList != nil {
		userAdd := AND(w.AndList)

		if w.Preset != "" {
			lines = append(lines, "WHERE "+w.Preset+" AND "+userAdd)
		} else {
			lines = append(lines, "WHERE "+userAdd)
		}
	} else {
		if w.Preset != "" {
			lines = append(lines, "WHERE "+w.Preset)
		} else {
			lines = append(lines, "")
		}
	}

	if s.GroupByLine != "" {
		lines = append(lines, s.GroupByLine)
	}

	// OrderBy etc
	s.OrderBy, s.Limit, s.Offset = OrderLimitOffset2(mp)
	lines = append(lines, s.OrderBy)
	lines = append(lines, fmt.Sprintf("LIMIT %d", s.Limit))
	lines = append(lines, fmt.Sprintf("OFFSET %d", s.Offset))

	res := strings.Join(lines, "\n")
	klog.DD(3, "%s", res)
	return res
}
