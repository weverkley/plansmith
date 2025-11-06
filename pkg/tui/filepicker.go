package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"plansmith/pkg/logging"
)

type fileSelectedMsg struct {
	path string
}

type FilePicker struct {
	files       []fileItem
	cursor      int
	visible     bool
	currentPath string
	filter      string
	width       int
	height      int
}

type fileItem struct {
	name  string
	isDir bool
	path  string
}

func NewFilePicker() FilePicker {
	return FilePicker{
		files:       []fileItem{},
		cursor:      0,
		visible:     false,
		currentPath: ".",
	}
}

func (f FilePicker) Init() tea.Cmd {
	return nil
}

func (f FilePicker) Update(msg tea.Msg) (FilePicker, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		logging.Debug("FilePicker received key: %s (type: %v)", msg.String(), msg.Type)
		if f.visible {
			logging.Debug("FilePicker is visible, handling navigation keys")
			switch msg.String() {
			case "up", "k":
				if f.cursor > 0 {
					f.cursor--
				}
			case "down", "j":
				if f.cursor < len(f.files)-1 {
					f.cursor++
				}
			case "enter":
				if len(f.files) > 0 {
					selected := f.files[f.cursor]
					if selected.isDir {
						// Navigate into directory
						f.currentPath = selected.path
						f.updateFiles()
						f.cursor = 0
					} else {
						// Select file - send a message
						f.visible = false
						return f, func() tea.Msg {
							return fileSelectedMsg{path: selected.path}
						}
					}
				}
			case "left", "h":
				// Go up one directory
				parent := filepath.Dir(f.currentPath)
				if parent != f.currentPath {
					f.currentPath = parent
					f.updateFiles()
					f.cursor = 0
				}
			case "esc":
				f.visible = false
			}
		}

	case tea.WindowSizeMsg:
		f.width = msg.Width
		f.height = msg.Height
	}

	return f, cmd
}

func (f *FilePicker) updateFiles() {
	// Read directory contents
	entries, err := os.ReadDir(f.currentPath)
	if err != nil {
		logging.Error("Failed to read directory %s: %v", f.currentPath, err)
		return
	}

	f.files = []fileItem{}

	// Add parent directory option if not at root
	if f.currentPath != "." && f.currentPath != "/" {
		f.files = append(f.files, fileItem{
			name:  "..",
			isDir: true,
			path:  filepath.Dir(f.currentPath),
		})
	}

	// Add directory entries first
	var dirs []fileItem
	var files []fileItem

	for _, entry := range entries {
		// Skip hidden files
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		path := filepath.Join(f.currentPath, entry.Name())

		if entry.IsDir() {
			dirs = append(dirs, fileItem{
				name:  entry.Name() + "/",
				isDir: true,
				path:  path,
			})
		} else {
			// Always show all files when filter is set to "all"
			if f.filter == "all" || strings.HasSuffix(entry.Name(), ".md") || strings.HasSuffix(entry.Name(), ".json") {
				files = append(files, fileItem{
					name:  entry.Name(),
					isDir: false,
					path:  path,
				})
			}
		}
	}

	// Sort directories and files
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].name < dirs[j].name
	})

	sort.Slice(files, func(i, j int) bool {
		return files[i].name < files[j].name
	})

	// Combine directories and files
	f.files = append(f.files, dirs...)
	f.files = append(f.files, files...)

	// Show message when no files are found
	if len(f.files) == 0 {
		f.files = append(f.files, fileItem{
			name:  "(No files found)",
			isDir: false,
			path:  "",
		})
	}
}

func (f FilePicker) View() string {
	if !f.visible {
		return ""
	}

	// File picker view
	pickerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1).
		Width(f.width - 10).
		Height(min(20, f.height-10))

	var items []string
	for i, file := range f.files {
		cursor := "  "
		if i == f.cursor {
			cursor = "> "
		}

		icon := "📄 "
		if file.isDir {
			icon = "📁 "
		}

		items = append(items, fmt.Sprintf("%s%s%s", cursor, icon, file.name))
	}

	if len(items) == 0 {
		items = append(items, "No files found")
	}

	fileList := lipgloss.JoinVertical(lipgloss.Left, items...)

	// Show current path
	pathView := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#808080")).
		Render("📂 " + f.currentPath)

	// Show instructions
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#808080")).
		Italic(true).
		Render("↑/↓: Navigate • Enter: Select • ←: Parent • Esc: Close")

	content := lipgloss.JoinVertical(lipgloss.Left,
		pathView,
		"",
		fileList,
		"",
		instructions,
	)

	return pickerStyle.Render(content)
}

func (f *FilePicker) SetFilter(filter string) {
	f.filter = filter
	f.updateFiles()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
