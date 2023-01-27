package main

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wrap"
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
	purple     = "#6D00E8"
	blueSelect = "#8D908B"
	orange     = "#E7220D"
	grey       = "#8D908B"
	yellow     = "#FF8C00"
	pink       = "#FF4D86"
	pinkDark   = "#B341E7"
	darkGrey   = "#343433"
	greyHelper = "#515150"

	styleborderTop       = lipgloss.NewStyle().Foreground(lipgloss.Color(blue))
	styleborderTopSelect = lipgloss.NewStyle().Foreground(lipgloss.Color(blueSelect))
	styleborderTitle     = lipgloss.NewStyle().Foreground(lipgloss.Color(purple))

	styleSpinner = lipgloss.NewStyle().Foreground(lipgloss.Color(purple))

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

	StyleUTF = lipgloss.NewStyle().Background(lipgloss.Color(pinkDark)).
			Height(1).Width(7).
			Render(" UTF-8 ")
	StyleCreate = lipgloss.NewStyle().Background(lipgloss.Color(purple)).
			Height(1).Width(WeightSet - 7).
			Render("  By waxdred ")

	StylehelperTitle = lipgloss.NewStyle().
				MarginLeft(2).
				Height(1).Width(8)

	statusStart      = StylehelperTitle.Background(lipgloss.Color(pink)).Render(" STATUS")
	statusLoading    = StylehelperTitle.Background(lipgloss.Color(purple)).Render(" STATUS")
	StylehelperValue = lipgloss.NewStyle().Background(lipgloss.Color(darkGrey)).
				Height(1).Width(WeightChat - 5)
	StylehelperLoader = lipgloss.NewStyle().Background(lipgloss.Color(purple)).
				Height(1).Width(WeightSet)
	colorHelper   = lipgloss.NewStyle().Foreground(lipgloss.Color(greyHelper))
	Stylehelper   = lipgloss.NewStyle().Foreground(lipgloss.Color(grey)).MarginLeft(2)
	helperInput   = colorHelper.Render(" <C-n>: new Session <C-y>: copy <Tab>: Cycle over windows")
	helperSetting = colorHelper.Render(" <C-k>: up <C-j>: down (+/-)")
	helperSession = colorHelper.Render(" <C-k>: up <C-j>: down <C-r>: rename <C-d>: delete")
	errorApi      = "Error OPENAI_API_KEY env missing:\nadd OPENAI_API_KEY=<api> to your env\nfor get your api:\nhttps://beta.openai.com/account/api-keys"
	errorApiWrap  = wrap.String(errorApi, WeightChat/2)
	NoApi         = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF2727")).
			Border(lipgloss.RoundedBorder()).
			Render(errorApiWrap)
)
