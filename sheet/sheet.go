package sheet

import (
	"fmt"
	"strings"
)

type TH struct {
	HTML  string
	Text  string
	Span  int
	Align string
}

type TR struct {
	thList []*TH
}

func (tr *TR) AddTHWithAlign(HTML string, Text string, Span int, Align string) {
	if Align == "" {
		Align = "left"
	}

	th := TH{
		HTML:  HTML,
		Text:  Text,
		Span:  Span,
		Align: Align,
	}
	tr.thList = append(tr.thList, &th)
}

func (tr *TR) AddTH(HTML string, Text string, Span int, Align string) {
	return tr.AddTHWithAlign(HTML, Text, Span, "")
}

func (tr *TR) Print(w *strings.Builder) {
	w.WriteString("<tr>")
	for _, th := range tr.thList {
		w.WriteString(fmt.Sprintf("<th colspan=\"%d\" style=\"text-align:%s;\">", th.Span, th.Align))
		if th.HTML != "" {
			w.WriteString(th.HTML)
		} else {
			w.WriteString(th.Text)
		}
		w.WriteString("</th>")
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
