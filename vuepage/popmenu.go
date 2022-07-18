package vuepage

import (
	"fmt"
	"strings"
)

type PopMenuItem struct {
	Name   string
	URL    string
	Action string
	HTML   string
}
type PopMenu struct {
	Title string
	items []*PopMenuItem
}

func PopMenuNew(title string) *PopMenu {
	return &PopMenu{
		Title: title,
	}
}

func (m *PopMenu) Add(args ...interface{}) *PopMenu {
	for i := 0; i < len(args)/2; i++ {
		var URL string
		var Action string
		var HTML string

		Name := args[2*i+0].(string)

		x := args[2*i+1].(string)
		if x != "" {
			if x[0] == '<' {
				HTML = x
			} else if x[0] == '@' {
				Action = x[1:len(x)]
			} else {
				URL = x
			}
		}

		m.items = append(m.items, &PopMenuItem{
			Name:   Name,
			URL:    URL,
			Action: Action,
			HTML:   HTML,
		})
	}
	return m
}

func (m *PopMenu) Items() []*PopMenuItem {
	return m.items
}

func (m *PopMenu) Gen() string {
	sp := fmt.Sprintf

	var line string
	var lines []string

	for _, i := range m.items {
		if i.URL != "" {
			line = sp(`<a class="dropdown-item" href="%s">%s</a>`, i.URL, i.Name)
		} else if i.Action != "" {
			line = sp(`<a class="dropdown-item" @click="%s">%s</a>`, i.Action, i.Name)
		} else if i.HTML != "" {
			line = i.HTML
		} else {
			line = sp(`<hr class="dropdown-divider">`)
		}
		lines = append(lines, line)
	}

	f := `<div class="dropdown">

  <div class="dropdown-trigger">
    <button class="button is-success" aria-haspopup="true" aria-controls="dropdown-menu">
      <span>%s</span>
      <span class="icon is-small">
        <i class="fas fa-angle-down" aria-hidden="true"></i>
      </span>
    </button>
  </div>

  <div class="dropdown-menu" id="dropdown-menu" role="menu">
    <div class="dropdown-content">
		%s
    </div>
  </div>
</div>
`

	return sp(f, m.Title, strings.Join(lines, ""))

}
