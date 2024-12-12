package menu

import tea "github.com/charmbracelet/bubbletea"

const (
	rooms = iota
	friends
	exit
)

type MenuModel struct {
	items   []string
	cursor  int
	focused int
}

func initialModel() MenuModel {
	items := make([]string, 3)

	items[rooms] = "Chat rooms"
	items[friends] = "Friend list"
	items[exit] = "Exit Clich"

	return MenuModel{
		items:   items,
		cursor:  0,
		focused: 0,
	}

}

func (m MenuModel) Init() tea.Model {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m MenuModel) View() string {
	return ""
}
