package roast

import (
	"fmt"
	"strings"

	"github.com/kamasamikon/miego/xmap"
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

	FromLine string

	GroupByLine string

	OrderBy string
	Limit   uint
	Offset  uint
}

func QueryStatementNew() *QueryStatement {
	return &QueryStatement{}
}

func (s *QueryStatement) Column(TableAlias string, Field string, Default string, OutputName string) {
	s.ColumnList = append(s.ColumnList, &Column{
		TableAlias: TableAlias,
		Field:      Field,
		Default:    Default,
		OutputName: OutputName,
	})
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

func (s *QueryStatement) From(FromLine string) {
	s.FromLine = FromLine
}

func (s *QueryStatement) String(mp xmap.Map, FoundRows bool) string {
	var lines []string

	// HEAD
	lines = append(lines, "SELECT")
	if FoundRows {
		lines = append(lines, "SQL_CALC_FOUND_ROWS")
	}

	// Output fields
	var sublines []string
	for _, c := range s.ColumnList {
		var sss string
		if c.Field[0] == '@' {
			field := c.Field[1:len(c.Field)]
			sss = fmt.Sprintf(`    DISTINCT(%s.%s) AS %s`, c.TableAlias, field, c.OutputName)
		} else {
			sss = fmt.Sprintf(`    IFNULL(%s.%s, "%s") AS %s`, c.TableAlias, c.Field, c.Default, c.OutputName)
		}
		sublines = append(sublines, sss)
	}
	lines = append(lines, strings.Join(sublines, ",\n"))

	//
	// From: SELECT * __FROM__ TableA
	//
	lines = append(lines, s.FromLine)

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
	orderBy, limit, offset := OrderLimitOffset2(mp)
	lines = append(lines, orderBy)
	lines = append(lines, fmt.Sprintf("LIMIT %d", limit))
	lines = append(lines, fmt.Sprintf("OFFSET %d", offset))

	// Something else
	s.OrderBy, s.Limit, s.Offset = OrderLimitOffset2(mp)

	return strings.Join(lines, "\n")
}
