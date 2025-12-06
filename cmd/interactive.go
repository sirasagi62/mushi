/*
* Copyright (c) 2020-2025 Charmbracelet, Inc
* Under MIT LICENSE
* See third_party_licenses/github.com/charmbracelet/bubbletea/LICENSE for the full license text
* The following code is based on the list-default and list-sample examples from bubbletea and modified by OKABE Gota.
 */

package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type item string

func (i item) FilterValue() string { return string(i) }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// runInteractiveSelector runs the interactive template selector
func runInteractiveSelector(cacheDir string) (string, error) {
	// キャッシュディレクトリが存在しない場合は、自動的に取得
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		fmt.Println("Cache not found. Cloning github/gitignore repository...")
		if err := cloneCache(cacheDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error cloning cache: %v\n", err)
			os.Exit(1)
		}
	}
	// キャッシュディレクトリ内のすべての .gitignore ファイルを取得
	files, err := os.ReadDir(cacheDir)
	if err != nil {
		return "", err
	}

	// テンプレート名のリストを作成
	items := []list.Item{}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".gitignore") {
			// .gitignore を除去してテンプレート名を取得
			templateName := strings.TrimSuffix(file.Name(), ".gitignore")
			items = append(items, item(templateName))
		}
	}

	// リストを作成
	l := list.New(items, itemDelegate{}, 10, 0)
	l.Title = "Select a gitignore template"
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	// モデルを作成
	m := templateModel{list: l}

	// プログラムを実行
	p := tea.NewProgram(m, tea.WithAltScreen())

	// 結果を取得
	result, err := p.Run()
	if err != nil {
		return "", err
	}

	if m, ok := result.(templateModel); ok {
		return m.choice, nil
	}

	return "", nil
}

// テンプレート選択用のインタラクティブUI
type templateModel struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m templateModel) Init() tea.Cmd {
	return nil
}

func (m templateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		case tea.KeyEnter:
			if len(m.list.VisibleItems()) > 0 {
				m.choice = m.list.SelectedItem().FilterValue()
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m templateModel) View() string {
	if m.quitting {
		return ""
	}
	return "\n" + m.list.View()
}
