package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lithammer/fuzzysearch/fuzzy"

	"github.com/MohitSingh65/file-finder/finder"
)

type item string

func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return string(i) }

type model struct {
	input   textinput.Model
	list    list.Model
	files   []string
	matches []string
}

func initialModel(files []string) model {
	input := textinput.New()
	input.Placeholder = "Search files..."
	input.Focus()
	input.CharLimit = 256
	input.Width = 50

	l := list.New(nil, list.NewDefaultDelegate(), 150, 8)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.DisableQuitKeybindings()

	return model{
		input:   input,
		list:    l,
		files:   files,
		matches: []string{},
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if selected, ok := m.list.SelectedItem().(item); ok {
				exec.Command("xdg-open", selected.Title()).Start()
				return m, tea.Quit
			}
		}
	}

	m.input, cmd = m.input.Update(msg)

	query := m.input.Value()
	if query == "" {
		m.matches = []string{}
		m.list.SetItems(nil)
	} else {
		ranked := fuzzy.RankFind(query, m.files)
		sort.Slice(ranked, func(i, j int) bool {
			return ranked[i].Distance < ranked[j].Distance
		})

		m.matches = make([]string, len(ranked))
		items := make([]list.Item, len(ranked))
		for i, r := range ranked {
			m.matches[i] = r.Target
			items[i] = item(r.Target)
		}
		m.list.SetItems(items)

	}

	m.list, _ = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"\n  %s\n\n%s",
		m.input.View(),
		m.list.View(),
	)
}

func main() {
	fmt.Println("Loading cache files...")

	files, err := finder.LoadFromCache()
	if err != nil {
		fmt.Print("Cache not found, using fallback...")
	}
	go func() {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Println("Failed to get user home:", err)
			return
		}
		updatedFiles, err := finder.IndexFiles(home, []string{".cache", "node_modules"})
		if err != nil {
			log.Println("Background indexing failed:", err)
			return
		}
		if err := finder.SaveToCache(updatedFiles); err != nil {
			log.Println("Failed to write cache:", err)
		}
	}()

	// if no cache, wait for background indexing
	if len(files) == 0 {
		fmt.Println("Indexing files...")
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		files, err = finder.IndexFiles(home, []string{".cache", "node_modules"})
		if err != nil {
			log.Fatal(err)
		}
	}

	// launch the UI as before
	p := tea.NewProgram(initialModel(files))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
