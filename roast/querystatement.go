package roast

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"miego/klog"
	"miego/xmap"
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
	Cast       string // ORDER BY cast(xxx.xxx AS yyy)
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

func (s *QueryStatement) OrderLimitOffset2(m xmap.Map) (orderBy string, limit uint, offset uint) {
	pageSize := m.Int("PageSize", 50)
	pageNumber := m.Int("PageNumber", 1)
	limit = uint(pageSize)

	offset = uint(m.Int("Offset", pageSize*(pageNumber-1)))

	orderBy = ""
	// OrderBy: CrtAt__Str or CrtAt or CrtAt_Xxx
	orderByWhat := strings.Split(m.S("OrderBy"), "__")[0]
	if orderByWhat != "" {
		for _, c := range s.ColumnList {
			if c.OutputName == orderByWhat {
				if c.Cast != "" {
					orderBy = fmt.Sprintf("ORDER BY CAST(%s.`%s` AS %s)", c.TableAlias, orderByWhat, c.Cast)
				} else {
					orderBy = fmt.Sprintf("ORDER BY `%s`", orderByWhat)
				}
				cmdOrderDir := m.S("OrderDir")
				if cmdOrderDir == "desc" {
					orderBy += " DESC"
				} else if cmdOrderDir == "asc" {
					orderBy += " ASC"
				}
				break
			}
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

func (s *QueryStatement) Column(TableAlias string, Field string, Default string, OutputName string, Cast string) {
	if OutputName == "" {
		OutputName = Field
	}
	s.ColumnList = append(s.ColumnList, &Column{
		TableAlias: TableAlias,
		Field:      Field,
		Default:    Default,
		OutputName: OutputName,
		Cast:       Cast,
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

// String : qStmt and cStmt
func (s *QueryStatement) String(mp xmap.Map, FoundRows int) (string, string, string) {
	//
	// ColumnLines
	//
	var ColumnLines []string

	// First DISTINCT
	var DistinctItems []string
	for _, c := range s.ColumnList {
		var sss string
		if c.Field[0] == '@' {
			field := c.Field[1:len(c.Field)]
			if c.TableAlias == "" {
				sss = fmt.Sprintf(`%s AS '%s'`, field, c.OutputName)
			} else {
				sss = fmt.Sprintf(`%s.%s AS '%s'`, c.TableAlias, field, c.OutputName)
			}
			DistinctItems = append(DistinctItems, sss)
		}
	}
	var DistinctLine string
	if DistinctItems != nil {
		DistinctLine = "DISTINCT " + strings.Join(DistinctItems, ", ")
		ColumnLines = append(ColumnLines, "    "+DistinctLine)
	}

	// Second IFNULL
	for _, c := range s.ColumnList {
		var sss string
		if c.Field[0] != '@' {
			if c.TableAlias == "" {
				sss = fmt.Sprintf(`    IFNULL(%s, "%s") AS '%s'`, c.Field, c.Default, c.OutputName)
			} else {
				sss = fmt.Sprintf(`    IFNULL(%s.%s, "%s") AS '%s'`, c.TableAlias, c.Field, c.Default, c.OutputName)
			}
			ColumnLines = append(ColumnLines, sss)
		}
	}

	ColumnPart := strings.Join(ColumnLines, ",\n")

	//
	// From: SELECT * __FROM__ TableA
	//
	FromPart := ""

	if s.FromAlias != "" {
		FromPart = "FROM " + s.FromTable + " AS " + s.FromAlias
	} else {
		FromPart = "FROM " + s.FromTable
	}

	//
	// Join
	//
	var JoinLines []string
	for _, j := range s.JoinList {
		if j.AndList != nil {
			userAdd := AND(j.AndList)

			if j.LineAnd == "" {
				JoinLines = append(JoinLines, "    "+userAdd)
			} else {
				JoinLines = append(JoinLines, "    "+j.LineAnd+" AND "+userAdd)
			}
		} else {
			JoinLines = append(JoinLines, "    "+j.LineNon)
		}
	}
	JoinPart := strings.Join(JoinLines, "\n")

	//
	// Where
	//
	var WhereLines []string
	w := s.WhereLine
	if w.AndList != nil {
		userAdd := AND(w.AndList)

		if w.Preset != "" {
			WhereLines = append(WhereLines, "WHERE "+w.Preset+" AND "+userAdd)
		} else {
			WhereLines = append(WhereLines, "WHERE "+userAdd)
		}
	} else {
		if w.Preset != "" {
			WhereLines = append(WhereLines, "WHERE "+w.Preset)
		} else {
			WhereLines = append(WhereLines, "")
		}
	}
	WherePart := strings.Join(WhereLines, "\n")

	GroupPart := s.GroupByLine
	if GroupPart == "" {
		var segs []string
		for _, s := range strings.Split(mp.S("__GroupBy__"), ";") {
			if s != "" {
				segs = append(segs, s)
			}
		}
		if segs != nil {
			GroupPart = "GROUP BY " + strings.Join(segs, ", ")
		}
	}

	// OrderBy etc
	var EtcLines []string
	orderBy, limit, offset := s.OrderLimitOffset2(mp)
	EtcLines = append(EtcLines, orderBy)
	EtcLines = append(EtcLines, fmt.Sprintf("LIMIT %d", limit))
	EtcLines = append(EtcLines, fmt.Sprintf("OFFSET %d", offset))

	EtcPart := strings.Join(EtcLines, "\n")

	// Something else, Will used by post phase
	s.OrderBy, s.Limit, s.Offset = s.OrderLimitOffset2(mp)

	//
	// qStmt: Query Statement
	//
	var qStmt string

	qStmt += "SELECT\n" + ColumnPart + "\n" + FromPart + "\n"
	if JoinPart != "" {
		qStmt += JoinPart + "\n"
	}
	qStmt += WherePart + "\n"
	if GroupPart != "" {
		qStmt += GroupPart + "\n"
	}
	if EtcPart != "" {
		qStmt += EtcPart + "\n"
	}
	if Conf.Noisy {
		klog.D("%s", qStmt)
	}

	//
	// cStmt: Count Statement
	//
	var cStmt string
	var cStmtHash string
	if FoundRows == FR_Auto {
		if mp.Has("PageNumber") {
			FoundRows = FR_Yes
		}
	}
	if FoundRows == FR_Yes {
		var cStmtList []string
		if DistinctLine != "" {
			var DistinctItems []string
			for _, c := range s.ColumnList {
				if c.Field[0] == '@' {
					field := c.Field[1:len(c.Field)]
					if c.TableAlias == "" {
						DistinctItems = append(
							DistinctItems,
							fmt.Sprintf(`%s`, field),
						)
					} else {
						DistinctItems = append(
							DistinctItems,
							fmt.Sprintf(`%s.%s`, c.TableAlias, field),
						)
					}
				}
			}

			DistinctLine = "DISTINCT " + strings.Join(DistinctItems, ", ")
			cStmtList = append(cStmtList, "SELECT COUNT("+DistinctLine+") AS Count")
			cStmtList = append(cStmtList, FromPart)

			if JoinPart != "" {
				cStmtList = append(cStmtList, JoinPart)
			}
			cStmtList = append(cStmtList, WherePart)
		} else {
			if GroupPart != "" {
				cStmtList = append(cStmtList, "SELECT COUNT(*) AS Count FROM (")
				cStmtList = append(cStmtList, "SELECT COUNT(*)")
				cStmtList = append(cStmtList, FromPart)

				if JoinPart != "" {
					cStmtList = append(cStmtList, JoinPart)
				}
				cStmtList = append(cStmtList, WherePart)
				cStmtList = append(cStmtList, GroupPart)

				cStmtList = append(cStmtList, ") a;")
			} else {
				cStmtList = append(cStmtList, "SELECT COUNT(*) AS Count")
				cStmtList = append(cStmtList, FromPart)

				if JoinPart != "" {
					cStmtList = append(cStmtList, JoinPart)
				}
				cStmtList = append(cStmtList, WherePart)
			}
		}

		cStmt = strings.Join(cStmtList, "\n")
		if Conf.Noisy {
			klog.D("%s", cStmt)
		}

		sort.Strings(cStmtList)
		ctx := md5.New()
		for i := range cStmtList {
			ctx.Write([]byte(cStmtList[i]))
		}
		cStmtHash = hex.EncodeToString(ctx.Sum(nil))
	}

	return qStmt, cStmt, cStmtHash
}
