package vuepage

type ToolbarItem struct {
	Name   string
	URL    string
	Action string
	HTML   string
}
type Toolbar struct {
	items []*ToolbarItem
}

func ToolbarNew() *Toolbar {
	return &Toolbar{}
}

func (b *Toolbar) Add(args ...interface{}) *Toolbar {
	for i := 0; i < len(args)/2; i++ {
		var URL string
		var Action string
		var HTML string

		x := args[2*i+1].(string)
		if x != "" {
			if x[0] == '<' {
				HTML = x
			} else if x[0] == '@' {
				Action = x[1:]
			} else {
				URL = x
			}
		}

		b.items = append(b.items, &ToolbarItem{
			Name:   args[2*i+0].(string),
			URL:    URL,
			Action: Action,
			HTML:   HTML,
		})
	}
	return b
}

func (b *Toolbar) Items() []*ToolbarItem {
	return b.items
}
