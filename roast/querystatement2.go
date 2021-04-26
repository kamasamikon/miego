package roast

import (
	"fmt"
	"strings"

	"github.com/kamasamikon/miego/klog"
	"github.com/kamasamikon/miego/xmap"
)

// String : qStmt and cStmt
func (s *QueryStatement) String2(mp xmap.Map, FoundRows int) (string, string) {
	//
	// ColumnLines
	//
	var ColumnLines []string
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
		ColumnLines = append(ColumnLines, sss)
	}
	ColumnPart := strings.Join(ColumnLines, ",\n")

	//
	// From: SELECT * __FROM__ TableA
	//
	FromPart := s.FromLine

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
	if FoundRows == FR_Auto {
		if mp.Has("PageSize") {
			FoundRows = FR_Yes
		}
	}
	if FoundRows == FR_Yes {
		if GroupPart != "" {
			cStmt += "SELECT COUNT(*) AS Count FROM (\n"
			cStmt += "SELECT COUNT(*) \n"
			cStmt += FromPart + "\n"

			if JoinPart != "" {
				cStmt += JoinPart + "\n"
			}
			cStmt += WherePart + "\n"
			cStmt += GroupPart + "\n"

			cStmt += ") a;"
		} else {
			cStmt += "SELECT COUNT(*) AS Count \n"
			cStmt += FromPart + "\n"

			if JoinPart != "" {
				cStmt += JoinPart + "\n"
			}
			cStmt += WherePart + "\n"
		}
		klog.D("%s", cStmt)
	}

	return qStmt, cStmt
}
