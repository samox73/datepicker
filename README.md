# Datepicker

A terminal-based interactive date picker for shell scripts, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Installation

```bash
go install github.com/samox73/datepicker@latest
```

## Usage

### Direct usage
```bash
datepicker
```

### In shell scripts
```bash
# Capture selected date in a variable
selected_date=$(./datepicker)
echo "You selected: $selected_date"

# With custom format
iso_date=$(./datepicker --format "2006-01-02")
echo "ISO format: $iso_date"
```

## Format Options

The `--format` flag accepts Go time layout strings:

```bash
./datepicker --format "02.01.2006"        # DD.MM.YYYY (default)
./datepicker --format "2006-01-02"        # YYYY-MM-DD
./datepicker --format "01/02/2006"        # MM/DD/YYYY
./datepicker --format "January 2, 2006"   # Month Day, Year
./datepicker --format "Mon, 02 Jan 06"    # Weekday, DD Mon YY
```

**Go time format reference:** `Mon Jan 2 15:04:05 MST 2006`

## Keybindings

| Key | Action |
|-----|--------|
| `↑/k` | Move up |
| `↓/j` | Move down |
| `←/h` | Move left |
| `→/l` | Move right |
| `n` | Next month |
| `m` | Previous month |
| `Enter` | Select date |
| `q` / `Ctrl+C` | Quit without selecting |

## Features

- 📅 Interactive calendar navigation
- 🎨 Highlighted cursor position
- 📤 Outputs to stdout (perfect for shell scripts)
- 🎯 TUI displays on stderr (doesn't interfere with output)
- ⌨️  Vim-style keybindings
- 🔧 Customizable date format

## License

MIT
