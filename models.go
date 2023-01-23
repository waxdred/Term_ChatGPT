package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wrap"
	gogpt "github.com/sashabaranov/go-gpt3"
)

type ChatGpt struct {
	api              string
	history          []string
	ctx              context.Context
	c                *gogpt.Client
	Temperature      float64
	TopP             float64
	FrequencyPenalty float64
	PresencePenalty  float64
	MaxTokens        int64
	user             string
}

func (chat *ChatGpt) Init() {
	chat.api = os.Getenv("OPENAI_API_KEY")
	if chat.api == "" {
		os.Exit(-1)
	}
	chat.user = os.Getenv("USER") + ":"
	chat.c = gogpt.NewClient(chat.api)
	chat.ctx = context.Background()
	chat.Temperature = 0.9
	chat.TopP = 1
	chat.FrequencyPenalty = 0
	chat.PresencePenalty = 0.6
	chat.MaxTokens = 400
}

type model struct {
	chatGpt  *ChatGpt
	apikey   string
	content  string
	ready    bool
	message  []string
	viewport viewport.Model
	textarea textarea.Model
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "esc" {
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		}
		switch msg.Type {
		case tea.KeyEnter:
			render := fmt.Sprintf("%s\n\n%s", TitleUser.Render(m.chatGpt.user), m.textarea.Value())
			m.content = m.content + "\n" + styleUser.Render(render) + "\n"
			chatGpt, err := requetOpenAI(m.chatGpt, m.textarea.Value())
			if err != nil {
				m.viewport.SetContent(err.Error())
				return m, nil
			}
			if chatGpt != "" {
				render := fmt.Sprintf("%s\n%s", TitleGpt.Render("Chat_GPT:"), chatGpt)
				m.content = m.content + styleGpt.Render(render)
				m.viewport.SetContent(m.content)
				m.textarea.Reset()
				m.viewport.GotoBottom()
			}
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(Wwidth-int(Wwidth/3), Wheight-8)
			m.viewport.YPosition = 4
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			text := wrap.String(m.content, Wwidth-int(Wwidth/3))
			m.viewport.SetContent(text)
			m.ready = true
			m.viewport.YPosition = Wheight + 1
		} else {
			m.viewport.Width = Wwidth - int(Wwidth/3)
			m.viewport.Height = Wheight - 8
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, vpCmd = m.viewport.Update(msg)
	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	msg := fmt.Sprintf("%s", m.viewport.View())
	msg = lipgloss.JoinVertical(lipgloss.Top, TopBorderText(Wwidth-int(Wwidth/3), " Chat GPT ", true), styleBorder.Render(msg))
	msg = lipgloss.JoinVertical(
		lipgloss.Top,
		msg,
		TopBorderText(Wwidth-int(Wwidth/3), " Prompt ", false),
		styleBorderPrompt.Render(m.textarea.View()),
	)
	return msg
}
