package roast

import (
	"fmt"
	"sort"
	"strings"

	"crypto/md5"
	"encoding/hex"

	"github.com/kamasamikon/miego/klog"
	"github.com/kamasamikon/miego/xmap"
)

// String : qStmt and cStmt
func (s *QueryStatement) String2(mp xmap.Map, FoundRows int) (string, string, string) {
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
				sss = fmt.Sprintf(`%s AS %s`, field, c.OutputName)
			} else {
				sss = fmt.Sprintf(`%s.%s AS %s`, c.TableAlias, field, c.OutputName)
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
				sss = fmt.Sprintf(`    IFNULL(%s, "%s") AS %s`, c.Field, c.Default, c.OutputName)
			} else {
				sss = fmt.Sprintf(`    IFNULL(%s.%s, "%s") AS %s`, c.TableAlias, c.Field, c.Default, c.OutputName)
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
	orderBy, limit, offset := OrderLimitOffset2(mp)
	EtcLines = append(EtcLines, orderBy)
	EtcLines = append(EtcLines, fmt.Sprintf("LIMIT %d", limit))
	EtcLines = append(EtcLines, fmt.Sprintf("OFFSET %d", offset))

	EtcPart := strings.Join(EtcLines, "\n")

	// Something else, Will used by post phase
	s.OrderBy, s.Limit, s.Offset = OrderLimitOffset2(mp)

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
	klog.D("%s", qStmt)

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
		klog.D("%s", cStmt)

		sort.Strings(cStmtList)
		ctx := md5.New()
		for i := range cStmtList {
			ctx.Write([]byte(cStmtList[i]))
		}
		cStmtHash = hex.EncodeToString(ctx.Sum(nil))
	}

	return qStmt, cStmt, cStmtHash
}
