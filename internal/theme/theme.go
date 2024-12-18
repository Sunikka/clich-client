package theme

import (
	"fmt"
	"os"
	"os/user"

	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

type Theme struct {
	PrimaryColor   lipgloss.Color `yaml:"primary-color"`
	SecondaryColor lipgloss.Color `yaml:"secondary-color"`
	HighlightColor lipgloss.Color `yaml:"highlight-color"`
}

func Init() (*Theme, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("Failed to get current user: %w", err)
	}

	confpath := usr.HomeDir + "/.config/clich/default_theme.yml"
	file, err := os.Open(confpath)
	if err != nil {
		return nil, fmt.Errorf("Error opening file; %w", err)
	}

	defer file.Close()

	var theme Theme
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&theme)
	if err != nil {
		fmt.Println("Error decoding config file: ", err)
		return nil, nil
	}

	return &theme, nil
}

// TODO: Implement theme switching
func Switch() {

}
