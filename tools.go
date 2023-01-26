package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	TopLeft  = "╭"
	TopRight = "╮"
	Top      = "─"
)

func TopBorderText(size int, title string, margin bool, stl bool) string {
	var ret string
	var style lipgloss.Style
	if stl {
		style = styleborderTopSelect
	} else {
		style = styleborderTop
	}
	tmp := int(size/2) - int(len(title)/2)
	left := strings.Repeat(Top, tmp)
	right := strings.Repeat(Top, tmp)
	if tmp*2+len(title) < size {
		right = strings.Repeat(Top, tmp+1)
	} else if tmp*2+len(title) > size {
		right = strings.Repeat(Top, tmp-1)
	}
	ret = style.Render(
		TopLeft,
	) + style.Render(
		right,
	) + styleborderTitle.Render(title) + style.Render(
		left,
	) + style.Render(
		TopRight,
	)
	if margin {
		style = lipgloss.NewStyle().MarginTop(1).MarginLeft(1)
	} else {
		style = lipgloss.NewStyle().MarginLeft(1)
	}
	return style.Render(ret)
}

func formatSetting(chat *ChatGpt, idx int) string {
	var freq, pres, temp, maxtoken, top, setting string
	mds := strings.TrimPrefix(chat.Model, "text-")
	model := fmt.Sprintf("%s %s", styleSettingTitle.Render("Model:"), styleSettingValue.Render(mds))
	freq = fmt.Sprintf(
		"%s %s",
		styleSettingTitle.Render("Frequency:"),
		styleSettingValue.Render(fmt.Sprintf("%.2f", chat.FrequencyPenalty)),
	)
	pres = fmt.Sprintf(
		"%s %s",
		styleSettingTitle.Render("Presence:"),
		styleSettingValue.Render(fmt.Sprintf("%.2f", chat.PresencePenalty)),
	)
	temp = fmt.Sprintf(
		"%s %s",
		styleSettingTitle.Render("Temperature:"),
		styleSettingValue.Render(fmt.Sprintf("%.2f", chat.Temperature)),
	)
	maxtoken = fmt.Sprintf(
		"%s %s",
		styleSettingTitle.Render("Max token:"),
		styleSettingValue.Render(fmt.Sprintf("%d", chat.MaxTokens)),
	)
	top = fmt.Sprintf(
		"%s %s",
		styleSettingTitle.Render("topP:"),
		styleSettingValue.Render(fmt.Sprintf("%.2f", chat.TopP)),
	)
	switch idx {
	case 1:
		freq = fmt.Sprintf(
			"> %s %s", styleSettingSelectTitle.Render("Frequency:"),
			styleSettingSelectValue.Render(fmt.Sprintf("%.2f", chat.FrequencyPenalty)),
		)
	case 2:
		pres = fmt.Sprintf(
			"> %s %s",
			styleSettingSelectTitle.Render("Presence:"),
			styleSettingSelectValue.Render(fmt.Sprintf("%.2f", chat.PresencePenalty)),
		)
	case 3:
		temp = fmt.Sprintf(
			"> %s %s",
			styleSettingSelectTitle.Render("Temperature:"),
			styleSettingSelectValue.Render(fmt.Sprintf("%.2f", chat.Temperature)),
		)
	case 4:
		maxtoken = fmt.Sprintf(
			"> %s %s",
			styleSettingSelectTitle.Render("Max token:"),
			styleSettingSelectValue.Render(fmt.Sprintf("%d", chat.MaxTokens)),
		)
	case 5:
		top = fmt.Sprintf(
			"> %s %s",
			styleSettingSelectTitle.Render("topP:"),
			styleSettingSelectValue.Render(fmt.Sprintf("%.2f", chat.TopP)),
		)
	}
	setting = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", model, freq, pres, temp, maxtoken, top)
	return setting
}
