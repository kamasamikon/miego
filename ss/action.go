package ss

import (
	"strings"
)

type ActionItem struct {
	Name string
	URL  string
	Act  string
	Hint string
}
type Action struct {
	items []*ActionItem
}

func ActionNew() *Action {
	return &Action{}
}

// Add: (Name, Act|URL, Hint)...
// 2nd: @Xxx => Act=Xxx
func (b *Action) Add(args ...string) *Action {
	for i := 0; i < len(args)/3; i++ {
		a := &ActionItem{}

		Name := args[3*i+0]
		a.Name = Name

		What := args[3*i+1]
		if strings.HasPrefix(What, "@") {
			a.Act = What[1:len(What)]
		} else {
			a.URL = What
		}

		Hint := args[3*i+2]
		a.Hint = Hint

		b.items = append(b.items, a)
	}
	return b
}

func (b *Action) Items() []*ActionItem {
	return b.items
}
