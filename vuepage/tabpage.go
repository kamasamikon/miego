package vuepage

import (
	"github.com/kamasamikon/miego/xmap"
)

type TabPage struct {
	Name  string
	Title string
	URL   string
}
type TabComponent struct {
	items []*TabPage
}

func TabComponentNew() *TabComponent {
	return &TabComponent{}
}

func (b *TabComponent) Add(Name string, Title string, URL string) *TabComponent {
	if Title == "" {
		Title = Name
	}
	b.items = append(b.items, &TabPage{
		Name:  Name,
		URL:   URL,
		Title: Title,
	})
	return b
}

func (b *TabComponent) SetVars(mp xmap.Map, active int) {
	for _, i := range b.items {
		mp.Put("tabPageIsActive_"+i.Name, false)
	}

	if active < len(b.items) {
		i := b.items[active]
		mp.Put("tabPageIsActive_"+i.Name, true)
		mp.Put("currentTable", i.Name)
	}
}

func (b *TabComponent) Items() []*TabPage {
	return b.items
}
