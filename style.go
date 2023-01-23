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
	heightPrompt       = 1
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()

	styleBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			PaddingRight(2).PaddingLeft(2).Width(Wwidth - int(Wwidth/3)).MarginLeft(1).Height(Wheight - 8).
			BorderTop(false)

	styleBorderPrompt = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Width(Wwidth - int(Wwidth/3)).MarginLeft(1).Height(heightPrompt).
				BorderTop(false)

	styleUser = lipgloss.NewStyle().MarginLeft((Wwidth / 3) - 5).
			Width((Wwidth - int(Wwidth/3)) / 2).
			Border(lipgloss.RoundedBorder())

	styleGpt = lipgloss.NewStyle().Width((Wwidth - int(Wwidth/3)) / 2).
			Border(lipgloss.RoundedBorder()).
			PaddingRight(1).PaddingLeft(1)

	TitleGpt = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4"))
	TitleUser = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D66F7"))
)
