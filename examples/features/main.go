package main

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyID          = "id"
	columnKeyName        = "name"
	columnKeyDescription = "description"
	columnKeyCount       = "count"
)

type Model struct {
	tableModel table.Model
}

func NewModel() Model {
	headers := []table.Header{
		table.NewHeader(columnKeyID, "ID", 5).WithStyle(lipgloss.NewStyle().Bold(true)),
		table.NewHeader(columnKeyName, "Name", 10),
		table.NewHeader(columnKeyDescription, "Description", 30),
		table.NewHeader(columnKeyCount, "#", 5),
	}

	rows := []table.Row{
		table.NewRow(table.RowData{
			columnKeyID:          "abc",
			columnKeyName:        "Hello",
			columnKeyDescription: "The first table entry, ever",
			columnKeyCount:       4,
		}),
		table.NewRow(table.RowData{
			columnKeyID:          "123",
			columnKeyName:        "Oh no",
			columnKeyDescription: "Super bold!",
			columnKeyCount:       17,
		}).WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)),
		table.NewRow(table.RowData{
			columnKeyID:          "def",
			columnKeyName:        "Yay",
			columnKeyDescription: "This is a really, really, really long description that will get cut off",
			columnKeyCount:       "N/A",
		}),
	}

	return Model{
		tableModel: table.New(headers).
			WithRows(rows).
			HeaderStyle(lipgloss.NewStyle().Bold(true)).
			SelectableRows(true).
			Focused(true),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.tableModel, cmd = m.tableModel.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}

	highlightedRow := m.tableModel.HighlightedRow()
	body.WriteString("Table demo with selectable rows!\nPress space/enter to select a row, q or ctrl+c to quit\n")

	body.WriteString(fmt.Sprintf("Currently looking at ID: %s\n", highlightedRow.Data[columnKeyID]))

	selectedIDs := []string{}

	for _, row := range m.tableModel.SelectedRows() {
		// Slightly dangerous type assumption but fine for demo
		selectedIDs = append(selectedIDs, row.Data[columnKeyID].(string))
	}

	body.WriteString(fmt.Sprintf("SelectedIDs: %s\n", strings.Join(selectedIDs, ", ")))

	body.WriteString(m.tableModel.View())

	return body.String()
}

func main() {
	p := tea.NewProgram(NewModel())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
