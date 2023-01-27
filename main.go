package main

// An example program demonstrating the pager component from the Bubbles
// component library.

import (
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	temperature = flag.Float64("temp", 0, "chatGpt temperature")
	top         = flag.Float64("top", 0.5, "chatGpt topP")
	frequency   = flag.Float64("freg", 0.5, "chatGpt frequency")
	presence    = flag.Float64("pres", 0.5, "chatGpt presence")
	token       = flag.Int64("token", 100, "chatGpt presence")
)

func main() {
	flag.Parse()
	var chatGpt ChatGpt
	var sessions Sessions
	var api bool
	sessions.init()
	chatGpt.Init()
	if chatGpt.api == "" {
		api = false
	} else {
		api = true
	}
	s := spinner.New()
	s.Spinner = spinner.Dot
	// Load some text for our viewport
	ta := textarea.New()
	ti := textarea.New()
	ta.Placeholder = "Send a message..."
	ti.Placeholder = "Rename"
	ta.Prompt = " "
	ti.Prompt = "> "
	ta.Focus()
	ta.CharLimit = 1000
	ti.CharLimit = WeightSet - 4
	ta.SetWidth(WeightChat)
	ta.SetHeight(heightPrompt)
	ti.SetWidth(WeightSet + 4)
	ti.SetHeight(heightPrompt)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ti.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)
	p := tea.NewProgram(
		model{
			textinput:       ti,
			rename:          false,
			spinner:         s,
			sessions:        sessions.init(),
			chatGpt:         &chatGpt,
			content:         "",
			textarea:        ta,
			prompt:          true,
			chat:            false,
			session:         false,
			setting:         false,
			selectorSetting: 1,
			selectorSession: 1,
			typing:          true,
			api:             api,
		},
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
