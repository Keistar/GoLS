package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// スタイルの定義
	titleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#7D56F4")).
			Foreground(lipgloss.Color("#FFFDF5")).
			Bold(true).
			Padding(0, 1)

	dirStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#12B5E5")).
			Bold(true)

	fileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4C94")).
			Bold(true).
			SetString("> ")
)

type model struct {
	cursor int           // 現在選択しているファイルのインデックス
	files  []os.DirEntry // ファイル一覧
}

func initalModel() model {
	files, _ := os.ReadDir(".")
	return model{
		files: files,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.files)-1 {
				m.cursor++
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := titleStyle.Render(" Golphin ") + " (q: Exit / Enter: Open)\n\n"

	for i, file := range m.files {
		cursor := " "
		rowStyle := fileStyle
		if m.cursor == i {
			cursor = selectedStyle.String()
			rowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4C94")).Bold(true)
		}

		icon := "F"
		name := file.Name()
		if file.IsDir() {
			icon = "D"
			if m.cursor != i {
				name = dirStyle.Render(name)
			} else {
				name = rowStyle.Render(name)
			}
		} else {
			name = fileStyle.Render(name)
		}

		s += fmt.Sprintf("%s %s %s\n", cursor, icon, name)
	}
	s += fmt.Sprintf("\n %d items in there.", len(m.files))
	return s
}

func main() {
	p := tea.NewProgram(initalModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("エラーが発生しました: %v", err)
		os.Exit(1)
	}
}
