package tui

import (
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ChatMessage struct {
	Role      string // "user", "assistant", "system"
	Content   string
	Timestamp time.Time
}

type Chat struct {
	messages  []ChatMessage
	viewport  viewport.Model
	spinner   spinner.Model
	isLoading bool
	width     int
	height    int
}

func NewChat() Chat {
	vp := viewport.New(80, 20)
	sp := spinner.New()
	sp.Spinner = spinner.Dot

	return Chat{
		messages: []ChatMessage{},
		viewport: vp,
		spinner:  sp,
	}
}

func (c Chat) Init() tea.Cmd {
	return nil
}

func (c *Chat) Update(msg tea.Msg) (*Chat, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.width = msg.Width
		c.height = msg.Height
		c.viewport.Width = msg.Width

	case spinner.TickMsg:
		if c.isLoading {
			c.spinner, cmd = c.spinner.Update(msg)
		}
	}

	c.viewport, cmd = c.viewport.Update(msg)
	return c, cmd
}

func (c Chat) View() string {
	var content []string

	// Add messages
	for _, msg := range c.messages {
		content = append(content, c.renderMessage(msg))
	}

	// Add loading indicator if needed
	// This is rendered in the main view, not here

	// Set content in viewport
	contentStr := lipgloss.JoinVertical(lipgloss.Left, content...)
	if len(contentStr) == 0 {
		// Show a placeholder when there are no messages
		contentStr = helpTextStyle.Render("No messages yet. Type 'help' to get started.")
	}

	c.viewport.SetContent(contentStr)

	return c.viewport.View()
}

func (c Chat) renderMessage(msg ChatMessage) string {
	var roleStyle lipgloss.Style
	var roleIcon string

	switch msg.Role {
	case "user":
		roleStyle = userMessageStyle
		roleIcon = userIcon
	case "assistant":
		roleStyle = assistantMessageStyle
		roleIcon = assistantIcon
	case "system":
		roleStyle = systemMessageStyle
		roleIcon = systemIcon
	case "error":
		roleStyle = errorMessageStyle
		roleIcon = errorIcon
	default:
		roleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#9E9E9E"))
		roleIcon = ""
	}

	// Format the message content with path/file highlighting
	formattedContent := formatMessageContent(msg.Content)

	// Combine icon, role, and content with proper spacing
	return lipgloss.JoinVertical(lipgloss.Left,
		roleStyle.Render(roleIcon+" "+strings.ToTitle(msg.Role)),
		formattedContent,
		"", // Add spacing between messages
	)
}

func (c Chat) renderMessages() []string {
	var content []string

	// Add messages
	for _, msg := range c.messages {
		content = append(content, c.renderMessage(msg))
	}

	return content
}

var ( // Regex for path and file detection
	// pathRegex matches common path patterns like ./foo/bar, ../baz.txt, /usr/local/bin
	pathRegex = regexp.MustCompile(`(?:\./|\.\./|/)(?:[\w-]+\/)*[\w.-]+`)
	// fileRegex matches standalone filenames like file.txt, archive.tar.gz
	fileRegex = regexp.MustCompile(`\b[\w-]+\.[\w-.]+\b`)
)

func formatMessageContent(content string) string {
	// Apply path highlighting first
	content = pathRegex.ReplaceAllStringFunc(content, func(s string) string {
		return pathStyle.Render(s)
	})

	// Apply file highlighting. This will highlight filenames that are not part of a path
	// already highlighted. If a filename is part of a path, the path style will take precedence
	// because it was applied first. This is a simpler heuristic.
	content = fileRegex.ReplaceAllStringFunc(content, func(s string) string {
		return fileStyle.Render(s)
	})

	// Simple word wrapping
	lines := strings.Split(content, "\n")
	var wrappedLines []string

	for _, line := range lines {
		if lipgloss.Width(line) > 80 {
			// Wrap long lines, this is still imperfect with ANSI codes
			// but should prevent the "0mist" issue by having more precise regexes.
			// A proper ANSI-aware word wrapper is needed for a perfect solution.
			wrappedLines = append(wrappedLines, lipgloss.NewStyle().MaxWidth(80).Render(line))
		} else {
			wrappedLines = append(wrappedLines, line)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, wrappedLines...)
}

func (c *Chat) AddMessage(role, content string) {
	c.messages = append(c.messages, ChatMessage{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	})

	// Update the viewport with the new content and scroll to bottom
	c.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, c.renderMessages()...))
	c.viewport.GotoBottom()
}

func (c *Chat) SetLoading(loading bool) {
	c.isLoading = loading
	if c.isLoading {
		c.spinner.Tick()
	}
}

func (c Chat) IsLoading() bool {
	return c.isLoading
}
