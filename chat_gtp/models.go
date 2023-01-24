package main

import (
	"context"
	"fmt"
	"os"

	"github.com/atotto/clipboard"
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
	Model            string
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
	chat.Model = "text-davinci-003"
	chat.c = gogpt.NewClient(chat.api)
	chat.ctx = context.Background()
	chat.Temperature = 0.9
	chat.TopP = 1
	chat.FrequencyPenalty = 0
	chat.PresencePenalty = 0.6
	chat.MaxTokens = 400
}

type model struct {
	chatGpt         *ChatGpt
	apikey          string
	content         string
	ready           bool
	message         []string
	viewport        viewport.Model
	textarea        textarea.Model
	prompt          bool
	chat            bool
	session         bool
	setting         bool
	selectorSetting int8
	selectorSession int64
	typing          bool
	last_answer     string
}

func (m model) Init() tea.Cmd {
	m.typing = true
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)
	if m.typing {
		m.textarea, tiCmd = m.textarea.Update(msg)
	}
	m.viewport, vpCmd = m.viewport.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "esc" || k == "ctrl+c" {
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		}
		switch msg.String() {
		case "ctrl+y":
			clipboard.WriteAll(m.last_answer)
		case "ctrl+k":
			if m.setting {
				if m.selectorSetting == 1 && m.chatGpt.FrequencyPenalty < 2.0 {
					m.chatGpt.FrequencyPenalty = m.chatGpt.FrequencyPenalty + 0.1
				} else if m.selectorSetting == 2 && m.chatGpt.PresencePenalty < 2.0 {
					m.chatGpt.PresencePenalty = m.chatGpt.PresencePenalty + 0.1
				} else if m.selectorSetting == 3 && m.chatGpt.Temperature < 2 {
					m.chatGpt.Temperature = m.chatGpt.Temperature + 0.1
				} else if m.selectorSetting == 4 && m.chatGpt.MaxTokens < 4000 {
					m.chatGpt.MaxTokens += 10
				} else if m.selectorSetting == 5 && m.chatGpt.TopP < 1 {
					m.chatGpt.TopP = m.chatGpt.TopP + 0.1
				}
			}
		case "ctrl+j":
			if m.setting {
				if m.selectorSetting == 1 && m.chatGpt.FrequencyPenalty > -2.0 {
					m.chatGpt.FrequencyPenalty = m.chatGpt.FrequencyPenalty - 0.1
				} else if m.selectorSetting == 2 && m.chatGpt.PresencePenalty > -2.0 {
					m.chatGpt.PresencePenalty = m.chatGpt.PresencePenalty - 0.1
				} else if m.selectorSetting == 3 && m.chatGpt.Temperature > 0.10 {
					m.chatGpt.Temperature = m.chatGpt.Temperature - 0.1
				} else if m.selectorSetting == 4 && m.chatGpt.MaxTokens > 0 {
					m.chatGpt.MaxTokens -= 10
				} else if m.selectorSetting == 5 && m.chatGpt.TopP > 0.10 {
					m.chatGpt.TopP = m.chatGpt.TopP - 0.1
				}
			}
		}
		switch msg.Type {
		case tea.KeyEnter:
			m.typing = false
			if m.prompt && m.textarea.Value() != "" {
				render := fmt.Sprintf("%s\n%s", TitleUser.Render(m.chatGpt.user), m.textarea.Value())
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
					m.last_answer = chatGpt
					m.textarea.Reset()
					m.viewport.GotoBottom()
				}
				m.typing = true
			}
		case tea.KeyTab:
			m.typing = false
			if m.prompt {
				m.prompt = false
				m.setting = true
			} else if m.setting {
				m.setting = false
				m.session = true
			} else if m.session {
				m.session = false
				m.prompt = true
				m.typing = true
			}
		case tea.KeyUp:
			if m.setting {
				m.selectorSetting--
				if m.selectorSetting <= 0 {
					m.selectorSetting = 5
				}
			}
		case tea.KeyDown:
			if m.setting {
				m.selectorSetting++
				if m.selectorSetting >= 6 {
					m.selectorSetting = 1
				}
			}
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(WeightChat, Wheight-8)
			m.viewport.YPosition = 4
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			text := wrap.String(m.content, WeightChat)
			m.viewport.SetContent(text)
			m.ready = true
			m.viewport.YPosition = WeightChat + 1
		} else {
			m.viewport.Width = WeightChat
			m.viewport.Height = Wheight - 8
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, vpCmd = m.viewport.Update(msg)
	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	var prompt lipgloss.Style
	var session lipgloss.Style
	var setting lipgloss.Style
	if m.prompt {
		prompt = styleBorderPromptSelect
	} else {
		prompt = styleBorderPrompt
	}
	if m.setting {
		setting = styleBorderSettingSelect
	} else {
		setting = styleBorderSetting
	}
	if m.session {
		session = styleBorderHistorySelect
	} else {
		session = styleBorderHistory
	}
	if !m.ready {
		return "\n  Initializing..."
	}
	msg := fmt.Sprintf("%s", m.viewport.View())
	ChatGpt := lipgloss.JoinVertical(
		lipgloss.Top,
		TopBorderText(WeightChat, " Chat GPT ", true, false),
		styleBorder.Render(msg),
	)
	ChatGptPrompt := lipgloss.JoinVertical(
		lipgloss.Top,
		ChatGpt,
		TopBorderText(WeightChat, " Prompt ", false, m.prompt),
		prompt.Render(m.textarea.View()),
	)
	Setting := lipgloss.JoinVertical(
		lipgloss.Top,
		TopBorderText(WeightSet, " Setting ", true, m.setting),
		setting.Render(formatSetting(m.chatGpt, int(m.selectorSetting))),
	)
	History := lipgloss.JoinVertical(
		lipgloss.Top,
		Setting,
		TopBorderText(WeightSet, " Session ", false, m.session),
		session.Render(""),
	)
	ret := lipgloss.JoinHorizontal(lipgloss.Top, ChatGptPrompt, History)
	return ret
}
