package main

import (
    "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type model struct {
    width       int
    height      int
    sidebar     sidebar
    mainContent mainContent
}

type sidebar struct {
    items   []string
    cursor  int
}

type mainContent struct {
    text string
}

func initialModel() model {
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
}

func (m model) View() string {
}

func main() {
}
