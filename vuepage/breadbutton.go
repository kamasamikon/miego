package vuepage

type BreadButtonItem struct {
	Name   string
	URL    string
	Action string
}
type BreadButton struct {
	items []*BreadButtonItem
}

func BreadButtonNew() *BreadButton {
	return &BreadButton{}
}

func (b *BreadButton) Add(args ...interface{}) *BreadButton {
	for i := 0; i < len(args)/2; i++ {
		var URL string
		var Action string

		x := args[2*i+1].(string)
		if x[0] == '@' {
			Action = x[1:len(x)]
		} else {
			URL = x
		}

		b.items = append(b.items, &BreadButtonItem{
			Name:   args[2*i+0].(string),
			URL:    URL,
			Action: Action,
		})
	}
	return b
}

func (b *BreadButton) Items() []*BreadButtonItem {
	return b.items
}
