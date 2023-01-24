package main

// An example program demonstrating the pager component from the Bubbles
// component library.

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	var chatGpt ChatGpt
	chatGpt.Init()
	// Load some text for our viewport
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Prompt = "> "
	ta.Focus()
	ta.CharLimit = 1000
	ta.SetWidth(WeightChat)
	ta.SetHeight(heightPrompt)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)
	p := tea.NewProgram(
		model{
			chatGpt:         &chatGpt,
			content:         "",
			textarea:        ta,
			prompt:          true,
			chat:            false,
			session:         false,
			setting:         false,
			selectorSetting: 1,
			selectorSession: 0,
			typing:          true,
		},
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
