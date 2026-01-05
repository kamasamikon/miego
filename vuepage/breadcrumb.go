package vuepage

type BreadcrumbItem struct {
	Name   string
	URL    string
	Active bool
}
type Breadcrumb struct {
	items []*BreadcrumbItem
}

func BreadcrumbNew() *Breadcrumb {
	return &Breadcrumb{}
}

func (b *Breadcrumb) Add(args ...interface{}) *Breadcrumb {
	for i := 0; i < len(args)/2; i++ {
		b.items = append(b.items, &BreadcrumbItem{
			Name: args[2*i+0].(string),
			URL:  args[2*i+1].(string),
		})
	}
	return b
}

func (b *Breadcrumb) Items() []*BreadcrumbItem {
	for i := range b.items {
		b.items[i].Active = false
	}
	if count := len(b.items); count > 0 {
		b.items[count-1].Active = true
	}
	return b.items
}
