package model

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/text/message"
)

type Tables []*Table

type Table struct {
	Name        string
	Columns     []Column
	RecordCount int
	Buffer      int
}

type Column struct {
	Index int
	Name  string
	Type  string
	Rule  Rule
}

type Rule struct {
	Type  string
	Value string
	Min   int
	Max   int
	// Format interface{}
	// Pattern interface{}
}

// hasMax check has max
func (r Rule) hasMax() bool {
	return r.Max != 0
}

// rangeNumber get number in (c.Min, c.Max)
func (r Rule) rangeNumber(idx int) int {
	number := r.Min + idx
	if r.hasMax() && number > r.Max {
		number = r.Min
	}
	return number
}

// queryValue get value for query
func (c Column) queryValue(idx int) string {
	var value string
	switch c.Rule.Type {
	case "const":
		value = c.Rule.Value
		if c.Type == "varchar" {
			value = fmt.Sprintf(`'%s'`, value)
		}
	case "unique":
		value = strconv.Itoa(c.Rule.rangeNumber(idx))
	}
	return value
}

// queryRecord get record for query
func (t Table) queryRecord(idx int) string {
	var record []string
	for _, column := range t.Columns {
		record = append(record, column.queryValue(idx))
	}
	return fmt.Sprintf("(%s)", strings.Join(record, ","))
}

// QueryRecords get records for query
func (t Table) QueryRecords() []string {
	var bufferRecords []string
	var records []string
	for idx := 0; idx < t.RecordCount; idx++ {
		records = append(records, t.queryRecord(idx))
		if idx%t.Buffer == t.Buffer-1 {
			bufferRecords = append(bufferRecords, strings.Join(records, ","))
			records = []string{}
		}
	}
	return bufferRecords
}

// type Format[T string] struct {
// 	Zeroppading bool
// 	StrLength   int
// }

// type Pattern[T int] struct {
// 	Value T
// 	Ratio float64
// }

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
