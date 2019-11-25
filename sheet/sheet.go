package sheet

import (
	"fmt"
	"strings"
)

type TH struct {
	Text  string
	Span  int
	Align string
}

type TR struct {
	thList []*TH
}

func (tr *TR) AddTHFull(Text string, Span int, Align string) *TR {
	if Align == "" {
		Align = "left"
	}

	th := TH{
		Text:  Text,
		Span:  Span,
		Align: Align,
	}
	tr.thList = append(tr.thList, &th)
	return tr
}

func (tr *TR) AddTH(Text string) *TR {
	return tr.AddTHFull(Text, 1, "")
}

func (tr *TR) AddTHSpan(Text string, Span int) *TR {
	return tr.AddTHFull(Text, Span, "")
}

func (tr *TR) Print(w *strings.Builder) {
	w.WriteString("<tr>")
	for _, th := range tr.thList {
		w.WriteString(fmt.Sprintf("<th colspan=\"%d\" style=\"text-align:%s;\">%s</th>", th.Span, th.Align, th.Text))
	}
	w.WriteString("</tr>")
}

type TABLE struct {
	trList []*TR
}

func (tab *TABLE) AddTR() *TR {
	tr := &TR{}
	tab.trList = append(tab.trList, tr)
	return tr
}
func (tab *TABLE) Print(w *strings.Builder) {
	w.WriteString("<table width=\"100%\" border=\"0\">")
	for _, tr := range tab.trList {
		tr.Print(w)
	}
	w.WriteString("</table>")
}
