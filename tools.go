package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	TopLeft  = "╭"
	TopRight = "╮"
	Top      = "─"
)

func TopBorderText(size int, title string, margin bool) string {
	var ret string
	var style lipgloss.Style
	ret = strings.Repeat(Top, int((size-len(title))/2))
	ret = TopLeft + ret + title + ret + TopRight
	if margin {
		style = lipgloss.NewStyle().MarginTop(1).MarginLeft(1)
	} else {
		style = lipgloss.NewStyle().MarginLeft(1)
	}
	return style.Render(ret)
}
