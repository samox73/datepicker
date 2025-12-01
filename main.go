package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	datepicker "github.com/ethanefung/bubble-datepicker"
	flag "github.com/spf13/pflag"
)

type model struct {
	datepicker datepicker.Model
	help       help.Model
	selected   time.Time
	quitting   bool
}

// keyMap defines our custom keybindings for help display
type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Enter key.Binding
	Quit  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Left, k.Right, k.Enter, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Enter, k.Quit},
	}
}

func initialModel() model {
	dp := datepicker.New(time.Now())
	dp.Styles = datepicker.DefaultStyles()
	dp.SelectDate() // Highlight the current date

	return model{
		datepicker: dp,
		help:       help.New(),
	}
}

func getKeyMap(dp datepicker.Model) keyMap {
	return keyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if m.datepicker.Selected {
				m.selected = m.datepicker.Time
				m.quitting = true
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.datepicker, cmd = m.datepicker.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	keys := getKeyMap(m.datepicker)
	helpView := m.help.View(keys)

	return fmt.Sprintf(
		"\n%s\n\n%s\n",
		m.datepicker.View(),
		helpView,
	)
}

func validateFormat(format string) error {
	// Check for common mistakes
	if format == "" {
		return fmt.Errorf("format string cannot be empty")
	}

	// Test the format with a known date
	testDate := time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)
	formatted := testDate.Format(format)

	// Check if the formatted output looks reasonable
	// If the format is invalid, it might produce unexpected results
	if len(formatted) == 0 {
		return fmt.Errorf("format string produces empty output")
	}

	// Check for the common mistake of using 2001 or other years
	if strings.Contains(format, "2001") || strings.Contains(format, "2002") ||
	   strings.Contains(format, "2003") || strings.Contains(format, "2004") ||
	   strings.Contains(format, "2005") || strings.Contains(format, "2007") {
		return fmt.Errorf("invalid year token in format string. Use '2006' for 4-digit year or '06' for 2-digit year")
	}

	// Warn about common mistakes (missing year token)
	if (format != "02.01.2006") &&
	   (!strings.Contains(format, "2006") && !strings.Contains(format, "06")) {
		fmt.Fprintf(os.Stderr, "Warning: Format string doesn't contain year token (2006 or 06)\n")
	}

	return nil
}

func main() {
	// Create a new flag set with just the binary name
	flags := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)

	// Define the format flag with default DD.MM.YYYY format
	formatFlag := flags.String("format", "02.01.2006", "Output date format using Go time layout")

	// Customize usage to include examples
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", filepath.Base(os.Args[0]))
		flags.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Go time format reference: Mon Jan 2 15:04:05 MST 2006
  01 or 1  = Month
  02 or 2  = Day
  2006     = Year (4-digit)
  06       = Year (2-digit)

Examples:
  --format "02.01.2006"        DD.MM.YYYY (default)
  --format "2006-01-02"        YYYY-MM-DD
  --format "01/02/2006"        MM/DD/YYYY
  --format "January 2, 2006"   Month Day, Year
  --format "Mon, 02 Jan 06"    Weekday, DD Mon YY
`)
	}

	flags.Parse(os.Args[1:])

	// Validate the format string
	if err := validateFormat(*formatFlag); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Run '%s --help' for format examples\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(), tea.WithOutput(os.Stderr))
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	m := finalModel.(model)
	if !m.selected.IsZero() {
		// Output in the specified format
		fmt.Println(m.selected.Format(*formatFlag))
	}
}
