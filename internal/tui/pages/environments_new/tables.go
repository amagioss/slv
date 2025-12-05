package environments_new

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.design/x/clipboard"
	"slv.sh/slv/internal/tui/theme"
)

// createResultTable creates a table displaying environment details with copy functionality
func (nep *NewEnvironmentPage) createResultTable() *tview.Table {
	colors := theme.GetCurrentPalette()

	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)
	table.SetBorder(true).
		SetBorderColor(colors.Border).
		SetBackgroundColor(colors.Background)

	row := 0

	// Add Name
	table.SetCell(row, 0, tview.NewTableCell("Name").SetTextColor(colors.TextSecondary).SetAlign(tview.AlignRight))
	table.SetCell(row, 1, tview.NewTableCell(nep.createdEnv.Name).SetTextColor(colors.TextPrimary).SetExpansion(1))
	row++

	// Add Type
	table.SetCell(row, 0, tview.NewTableCell("Type").SetTextColor(colors.TextSecondary).SetAlign(tview.AlignRight))
	table.SetCell(row, 1, tview.NewTableCell(string(nep.createdEnv.EnvType)).SetTextColor(colors.TextPrimary).SetExpansion(1))
	row++

	// Add Email if available
	if nep.createdEnv.Email != "" {
		table.SetCell(row, 0, tview.NewTableCell("Email").SetTextColor(colors.TextSecondary).SetAlign(tview.AlignRight))
		table.SetCell(row, 1, tview.NewTableCell(nep.createdEnv.Email).SetTextColor(colors.TextPrimary).SetExpansion(1))
		row++
	}

	// Add Tags if available
	if len(nep.createdEnv.Tags) > 0 {
		table.SetCell(row, 0, tview.NewTableCell("Tags").SetTextColor(colors.TextSecondary).SetAlign(tview.AlignRight))
		table.SetCell(row, 1, tview.NewTableCell(strings.Join(nep.createdEnv.Tags, ", ")).SetTextColor(colors.TextPrimary).SetExpansion(1))
		row++
	}

	// Add Public Key
	table.SetCell(row, 0, tview.NewTableCell("Public Key").SetTextColor(colors.TextSecondary).SetAlign(tview.AlignRight))
	table.SetCell(row, 1, tview.NewTableCell(nep.createdEnv.PublicKey).SetTextColor(colors.TextPrimary).SetExpansion(1))
	row++

	// Add EDS
	eds, err := nep.createdEnv.ToDefStr(false)
	if err == nil {
		table.SetCell(row, 0, tview.NewTableCell("EDS").SetTextColor(colors.TextSecondary).SetAlign(tview.AlignRight))
		table.SetCell(row, 1, tview.NewTableCell(eds).SetTextColor(colors.Accent).SetExpansion(1))
		row++
	}

	// Add Secret Binding if available
	if nep.createdEnv.SecretBinding != "" {
		table.SetCell(row, 0, tview.NewTableCell("Secret Binding").SetTextColor(colors.TextSecondary).SetAlign(tview.AlignRight))
		table.SetCell(row, 1, tview.NewTableCell(nep.createdEnv.SecretBinding).SetTextColor(colors.TextPrimary).SetExpansion(1))
		row++
	}

	// Add Secret Key if available (for direct service creation)
	if nep.secretKey != "" {
		table.SetCell(row, 0, tview.NewTableCell("Secret Key").SetTextColor(colors.TextSecondary).SetAlign(tview.AlignRight))
		table.SetCell(row, 1, tview.NewTableCell(nep.secretKey).SetTextColor(colors.Warning).SetExpansion(1))
		row++
	}

	// Set up input capture for copying and navigation
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			// Switch to buttons
			nep.GetTUI().GetApplication().SetFocus(nep.resultForm)
			return nil
		}

		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'c', 'C':
				// Copy selected field value
				selectedRow, _ := table.GetSelection()
				if selectedRow >= 0 && selectedRow < table.GetRowCount() {
					cell := table.GetCell(selectedRow, 1)
					if cell != nil {
						value := cell.Text
						if value != "" {
							clipboard.Write(clipboard.FmtText, []byte(value))
							fieldName := table.GetCell(selectedRow, 0).Text
							nep.UpdateStatus(fmt.Sprintf("Copied %s to clipboard", fieldName))
						}
					}
				}
				return nil
			}
		}
		return event
	})

	return table
}
