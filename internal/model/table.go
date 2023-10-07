package model

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/text/message"
)

type Table struct {
	Name        string
	Columns     []Column
	RecordCount int
}

type Column struct {
	Name string
	Rule Rule
}

type Rule struct {
	Scope  interface{}
	Format interface{}
}

type Scope[T any] struct {
	Start T
	End   T
}

type Format[T string] struct {
	Zeroppading bool
	StrLength   int
}

type Pattern[T int] struct {
	Value T
	Ratio float64
}

type Constant[T any] struct {
	Value T
}

type Tables []*Table

// StdoutOrg old print tables
func (ts Tables) StdoutOrg() error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.StripEscape|tabwriter.Debug)
	p := message.NewPrinter(message.MatchLanguage("ja"))
	fmt.Fprintln(w, "table\tcount")
	for _, table := range ts {
		fmt.Fprintln(w, table.Name+"\t"+p.Sprint(table.RecordCount))
	}
	return w.Flush()
}

// Stdout print tables
func (ts Tables) Stdout() error {
	p := message.NewPrinter(message.MatchLanguage("ja"))
	var tablesData [][]string
	for _, table := range ts {
		tablesData = append(tablesData, []string{table.Name, p.Sprint(table.RecordCount)})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Table", "Count"})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
	)
	table.AppendBulk(tablesData)
	table.Render()
	return nil
}
