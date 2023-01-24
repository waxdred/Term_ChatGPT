package main

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/crypto/ssh/terminal"
)

const useHighPerformanceRenderer = false

var (
	fd                 = int(os.Stdout.Fd())
	Wwidth, Wheight, _ = terminal.GetSize(fd)
	WeightSet          = 23
	WeightChat         = Wwidth - WeightSet - 8

	heightPrompt  = 1
	heightSetting = 6
	heightSession = Wheight - heightSetting - 7
	heightChat    = Wheight - heightPrompt - 7
)

var (
	blue       = "#1F6FEB"
	purple     = "#411D37"
	blueSelect = "#58B6C6"
	orange     = "#E7220D"
	grey       = "#8D908B"
	yellow     = "#FF8C00"
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()
	styleborderTop       = lipgloss.NewStyle().Foreground(lipgloss.Color(blue))
	styleborderTopSelect = lipgloss.NewStyle().Foreground(lipgloss.Color(blueSelect))
	styleborderTitle     = lipgloss.NewStyle().Foreground(lipgloss.Color(purple))

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()

	styleBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			PaddingRight(2).PaddingLeft(2).
			Width(WeightChat).MarginLeft(1).Height(heightChat).
			BorderForeground(lipgloss.Color(blue)).
			BorderTop(false)

	styleBorderSelect = styleBorder.Copy().BorderForeground(lipgloss.Color(blueSelect))

	styleBorderSetting = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				PaddingRight(2).PaddingLeft(2).Width(WeightSet).MarginLeft(1).Height(heightSetting).
				BorderForeground(lipgloss.Color(blue)).
				BorderTop(false)

	styleBorderSettingSelect = styleBorderSetting.Copy().BorderForeground(lipgloss.Color(blueSelect))

	styleBorderHistory = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				PaddingRight(2).PaddingLeft(2).
				Width(WeightSet).MarginLeft(1).Height(heightSession).
				BorderForeground(lipgloss.Color(blue)).
				BorderTop(false)

	styleBorderHistorySelect = styleBorderHistory.Copy().BorderForeground(lipgloss.Color(blueSelect))
	styleBorderPrompt        = lipgloss.NewStyle().
					Border(lipgloss.RoundedBorder()).
					Width(WeightChat).MarginLeft(1).Height(heightPrompt).
					BorderForeground(lipgloss.Color(blue)).
					BorderTop(false)

	styleBorderPromptSelect = styleBorderPrompt.Copy().BorderForeground(lipgloss.Color(blueSelect))

	styleUser = lipgloss.NewStyle().MarginLeft(WeightChat/2 + 5).
			Width((WeightChat/2 - 2))

	styleGpt = lipgloss.NewStyle().Width(WeightChat / 2).
			Border(lipgloss.RoundedBorder()).
			PaddingRight(1).PaddingLeft(1)

	TitleGpt = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4"))
	TitleUser = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D66F7"))

	styleSetting = lipgloss.NewStyle().Foreground(lipgloss.Color(orange))

	styleSettingSelectTitle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(grey))
	styleSettingSelectValue = lipgloss.NewStyle().
				Foreground(lipgloss.Color(yellow))

	styleSettingTitle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(orange))
	styleSettingValue = lipgloss.NewStyle().
				Foreground(lipgloss.Color(yellow))
)
