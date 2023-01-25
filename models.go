package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wrap"
	gogpt "github.com/sashabaranov/go-gpt3"
)

type ChatGpt struct {
	sync.Mutex
	routine          bool
	rep              string
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
	if *temperature >= 0 && *temperature <= 2.0 {
		chat.Temperature = *temperature
	} else {
		chat.Temperature = 0
	}
	if *top >= 0 && *top <= 2.0 {
		chat.TopP = *top
	} else {
		chat.TopP = 0
	}
	if *frequency >= -2.0 && *frequency <= 2.0 {
		chat.FrequencyPenalty = *frequency
	} else {
		chat.FrequencyPenalty = 0
	}
	if *presence >= -2.0 && *presence <= 2.0 {
		chat.PresencePenalty = *presence
	} else {
		chat.PresencePenalty = 0
	}
	if *token >= 0 && *token <= 4000 {
		chat.MaxTokens = *token
	} else {
		chat.MaxTokens = 100
	}
}

type model struct {
	spinner         spinner.Model
	sessions        Sessions
	curr_session    Session
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
	return tea.Batch(textarea.Blink, spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		spCmd tea.Cmd
	)
	if m.typing {
		m.textarea, tiCmd = m.textarea.Update(msg)
	}
	if m.chatGpt.rep != "" {
		render := fmt.Sprintf("%s\n%s", TitleGpt.Render("Chat_GPT:"), m.chatGpt.rep)
		m.content = m.content + styleGpt.Render(render)
		m.viewport.SetContent(m.content)
		m.last_answer = m.chatGpt.rep
		m.viewport.GotoBottom()
		if m.curr_session.Title == "" {
			t := time.Now()
			timestamp := t.Format("2006-01-02 15:04:05")
			name := strings.Replace(timestamp, " ", "", -1)
			m.curr_session.Id = int64(len(m.sessions) + 1)
			m.curr_session.Title = name
			m.curr_session.Content = m.content
			m.curr_session.Created_at = timestamp
			m.curr_session.Setting = setting{
				Temperature:      m.chatGpt.Temperature,
				TopP:             m.chatGpt.TopP,
				FrequencyPenalty: m.chatGpt.FrequencyPenalty,
				PresencePenalty:  m.chatGpt.PresencePenalty,
				MaxTokens:        m.chatGpt.MaxTokens,
			}
			m.sessions = append(m.sessions, m.curr_session)
		} else {
			m.curr_session.Content = m.content
			m.curr_session.Setting = setting{
				Temperature:      m.chatGpt.Temperature,
				TopP:             m.chatGpt.TopP,
				FrequencyPenalty: m.chatGpt.FrequencyPenalty,
				PresencePenalty:  m.chatGpt.PresencePenalty,
				MaxTokens:        m.chatGpt.MaxTokens,
			}
		}
		m.typing = true
		m.textarea.Placeholder = "Send a message..."
		m.textarea.Focus()
		m.curr_session.save()
		m.chatGpt.rep = ""
		m.chatGpt.routine = false
	}
	m.spinner, spCmd = m.spinner.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "esc" || k == "ctrl+c" {
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		}
		switch msg.String() {
		case "ctrl+y":
			if m.last_answer != "" {
				clipboard.WriteAll(m.last_answer)
			}
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
			if m.prompt && m.typing && m.textarea.Value() != "" {
				m.typing = false
				m.chatGpt.routine = false
				render := fmt.Sprintf("%s\n%s", TitleUser.Render(m.chatGpt.user), m.textarea.Value())
				m.content = m.content + "\n" + styleUser.Render(render) + "\n"
				go requetOpenAI(m.chatGpt, m.textarea.Value())
				if m.chatGpt.routine {
					return m, tea.Batch(tiCmd, vpCmd, spCmd)
				}
				m.textarea.Reset()
				m.textarea.Placeholder = "Loading..."
				m.viewport.SetContent(m.content)
				m.viewport.GotoBottom()
			}
			if m.session {
				m.curr_session = m.sessions[m.selectorSession-1]
				m.content = m.curr_session.Content
				m.viewport.SetContent(m.content)
				m.viewport.GotoBottom()
				m.chatGpt.MaxTokens = m.curr_session.Setting.MaxTokens
				m.chatGpt.PresencePenalty = m.curr_session.Setting.PresencePenalty
				m.chatGpt.FrequencyPenalty = m.curr_session.Setting.FrequencyPenalty
				m.chatGpt.TopP = m.curr_session.Setting.TopP
				m.chatGpt.Temperature = m.curr_session.Setting.Temperature
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
			} else if m.session && m.selectorSession > 1 {
				m.selectorSession--
			}
		case tea.KeyDown:
			if m.setting {
				m.selectorSetting++
			} else if m.session && m.selectorSession < int64(m.sessions.Len()) {
				m.selectorSession++
			}
		case tea.KeyCtrlD:
			if m.session {
				m.sessions = m.sessions.deleteFile(int(m.selectorSession))
			}
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(WeightChat, Wheight-8)
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			text := wrap.String(m.content, WeightChat)
			m.viewport.SetContent(text)
			m.ready = true
			m.viewport.GotoBottom()
		} else {
			m.viewport.Width = WeightChat
			m.viewport.Height = Wheight - 8
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, vpCmd = m.viewport.Update(msg)
	return m, tea.Batch(tiCmd, vpCmd, spCmd)
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
		prompt.Render(" "+styleSpinner.Render(m.spinner.View())+m.textarea.View()),
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
		session.Render(m.sessions.getList(m.selectorSession)),
	)
	ret := lipgloss.JoinHorizontal(lipgloss.Top, ChatGptPrompt, History)
	if m.session {
		ret = lipgloss.JoinVertical(
			lipgloss.Top,
			ret,
			StylehelperTitle+StylehelperValue.Render(helperSession)+StylehelperLoader.Render(""))
	} else if m.prompt {
		ret = lipgloss.JoinVertical(lipgloss.Top,
			ret,
			StylehelperTitle+StylehelperValue.Render(helperInput)+StylehelperLoader.Render(""))
	} else if m.setting {
		ret = lipgloss.JoinVertical(lipgloss.Top,
			ret,
			StylehelperTitle+StylehelperValue.Render(helperSetting)+StylehelperLoader.Render(""))
	} else {
		ret = lipgloss.JoinVertical(lipgloss.Top,
			ret,
			StylehelperTitle+StylehelperValue.Render("")+StylehelperLoader.Render(""))
	}
	return ret
}
