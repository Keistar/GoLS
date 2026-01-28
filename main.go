package main

import (
	"fmt"
	"os"
	"path/filepath"

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
	path   string        // Current path
	info   string
}

func initalModel() model {
	path, _ := os.Getwd()
	if len(os.Args) > 1 {
		argPath := os.Args[1]
		fmt.Println(argPath)
		absPath, err := filepath.Abs(argPath)
		if err == nil {
			path = absPath
		}
	}
	files, _ := os.ReadDir(path)
	return model{
		path:  path,
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
		case "enter", "right", "l":
			if len(m.files) == 0 {
				return m, nil
			}

			selected := m.files[m.cursor]
			if selected.IsDir() {
				newPath := filepath.Join(m.path, selected.Name())
				newFiles, err := os.ReadDir(newPath)
				if err == nil {
					m.path = newPath
					m.files = newFiles
					m.cursor = 0
				}
			}
		case "backspace", "left", "h":
			parentPath := filepath.Dir(m.path)

			if parentPath == m.path {
				return m, nil
			}

			newFiles, err := os.ReadDir(parentPath)
			if err == nil {
				m.path = parentPath
				m.files = newFiles
				m.cursor = 0
			}
		}
	}
	if len(m.files) > 0 {
		fileInfo, err := m.files[m.cursor].Info()
		if err == nil {
			size := fileInfo.Size()
			modTime := fileInfo.ModTime().Format("2006-01-02 15:04")
			m.info = fmt.Sprintf("Size: %d bytes | Mod: %s", size, modTime)
		}
	}
	return m, nil
}

func (m model) View() string {
	s := titleStyle.Render(" GoLS ") + " (q: Exit / Enter: Open / Backspace: Return)\n\n"
	s += lipgloss.NewStyle().Foreground(lipgloss.Color("#777777")).Render(" 📍 "+m.path) + "\n\n"

	for i, file := range m.files {
		cursor := " "
		rowStyle := fileStyle
		if m.cursor == i {
			cursor = selectedStyle.String()
			rowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4C94")).Bold(true)
		}

		icon := "📄"
		name := file.Name()
		if file.IsDir() {
			icon = "📁"
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
	infoBar := lipgloss.NewStyle().
		Background(lipgloss.Color("#353535")).
		Foreground(lipgloss.Color("#AAAAAA")).
		Width(50). // 表示幅を固定
		Padding(0, 1)

	s += "\n" + infoBar.Render(m.info)
	return s
}

func main() {
	p := tea.NewProgram(initalModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("エラーが発生しました: %v", err)
		os.Exit(1)
	}
}
