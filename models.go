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
	gogpt "github.com/sashabaranov/go-openai"
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
	chat.user = os.Getenv("USER") + ":"
	chat.Model = "text-davinci-003"
	chat.c = gogpt.NewClient(chat.api)
	chat.ctx = context.Background()
	chat.routine = false
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
	textinput       textarea.Model
	rename          bool
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
	api             bool
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
		reCmd tea.Cmd
	)
	if !m.ready {
		m.viewport = viewport.New(WeightChat, Wheight-8)
		m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
		text := wrap.String(m.content, WeightChat)
		m.viewport.SetContent(text)
		m.ready = true
		m.viewport.GotoBottom()
	}
	if m.chatGpt.rep != "" {
		m.chatGpt.Lock()
		defer m.chatGpt.Unlock()
		render := fmt.Sprintf("%s\n%s", TitleGpt.Render("Chat_GPT:"), m.chatGpt.rep)
		m.chatGpt.routine = false
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
			m.curr_session.Msg = m.curr_session.Msg + "chatGpt: " + m.chatGpt.rep
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
			m.curr_session.Content = m.curr_session.Content + m.content
			m.curr_session.Msg = m.curr_session.Msg + " " + m.chatGpt.rep
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
		m.chatGpt.rep = ""
		m.curr_session.save()
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
		if c := msg.String(); c == "k" || c == "u" || c == "b" || c == "d" || c == "j" {
			if m.prompt {
				m.textarea, tiCmd = m.textarea.Update(msg)
			} else if m.rename {
				m.textinput, reCmd = m.textinput.Update(msg)
			}
			return m, tea.Batch(tiCmd, vpCmd, spCmd, reCmd)
		}
		switch msg.String() {
		case "ctrl+v":
			if m.prompt {
				clip, err := clipboard.ReadAll()
				if err != nil {
					return m, nil
				}
				value := m.textarea.Value() + clip
				m.textarea.SetValue(value)
			}
		case "ctrl+y":
			if m.last_answer != "" {
				clipboard.WriteAll(m.last_answer)
			}
		case "+":
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
		case "-":
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
			if m.rename && m.textinput.Value() != "" && m.session {
				input := strings.Replace(m.textinput.Value(), " ", "_", -1)
				err := m.sessions.rename(input, m.selectorSession)
				if err != nil {
					return m, nil
				}
				m.rename = false
				m.textinput.Reset()
				m.sessions = m.sessions.init()
			} else if m.prompt && m.textarea.Value() != "" {
				m.typing = false
				m.curr_session.Msg = m.curr_session.Msg + " " + m.textarea.Value()
				render := wrap.String(m.textarea.Value(), WeightChat/3)
				render = fmt.Sprintf("%s\n%s", TitleUser.Render(m.chatGpt.user), render)
				m.content = m.content + "\n" + styleUser.Render(render) + "\n"
				go requetOpenAI(m.chatGpt, m.curr_session, m.textarea.Value())
				if m.chatGpt.routine {
					m.textarea.Reset()
					return m, tea.Batch(tiCmd, vpCmd, spCmd, reCmd)
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
				m.selectorSession = 1
				m.session = false
				m.prompt = true
			}
		case tea.KeyTab:
			m.rename = false
			m.textinput.Reset()
			if m.prompt {
				m.prompt = false
				m.setting = true
			} else if m.setting {
				m.setting = false
				m.session = true
			} else if m.session {
				m.session = false
				m.prompt = true
			}
		case tea.KeyCtrlK:
			if m.setting {
				m.selectorSetting--
				if m.selectorSetting <= 0 {
					m.selectorSetting = 5
				}
			} else if m.session && m.selectorSession > 1 {
				m.selectorSession--
			}
		case tea.KeyCtrlN:
			if m.typing {
				m.curr_session = Session{}
				m.content = ""
				m.viewport.SetContent("")
			}
		case tea.KeyCtrlJ:
			if m.setting {
				m.selectorSetting++
			} else if m.session && m.selectorSession < int64(m.sessions.Len()) {
				m.selectorSession++
			}
		case tea.KeyCtrlD:
			if m.session {
				m.sessions = m.sessions.deleteFile(int(m.selectorSession))
			}
		case tea.KeyCtrlR:
			if m.session {
				m.rename = true
				m.textinput.Focus()
			}
		}
	}

	if m.prompt {
		m.textarea, tiCmd = m.textarea.Update(msg)
	}
	m.viewport, vpCmd = m.viewport.Update(msg)
	if m.rename {
		m.textinput, reCmd = m.textinput.Update(msg)
	}
	m.spinner, spCmd = m.spinner.Update(msg)
	return m, tea.Batch(tiCmd, vpCmd, spCmd, reCmd)
}

func (m model) View() string {
	var prompt lipgloss.Style
	var session lipgloss.Style
	var setting lipgloss.Style
	var rename string
	if !m.api {
		m.viewport.SetContent(NoApi)
	}
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
		TopBorderText(WeightChat, " Chat GPT: "+m.curr_session.Title, true, false),
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
	status := ""
	if m.chatGpt.routine {
		status = statusLoading
	} else {
		status = statusStart
	}
	if m.rename && m.session {
		rename = m.textinput.View()
	} else {
		rename = StyleUTF + StyleCreate
		rename = StylehelperLoader.Render(rename)
	}
	if m.session {
		ret = lipgloss.JoinVertical(
			lipgloss.Top,
			ret,
			status+StylehelperValue.Render(helperSession)+StylehelperLoader.Render(rename))
	} else if m.prompt {
		ret = lipgloss.JoinVertical(lipgloss.Top,
			ret,
			status+StylehelperValue.Render(helperInput)+StylehelperLoader.Render(rename))
	} else if m.setting {
		ret = lipgloss.JoinVertical(lipgloss.Top,
			ret,
			status+StylehelperValue.Render(helperSetting)+StylehelperLoader.Render(rename))
	} else {
		ret = lipgloss.JoinVertical(lipgloss.Top,
			ret,
			status+StylehelperValue.Render("")+rename)
	}
	return ret
}
