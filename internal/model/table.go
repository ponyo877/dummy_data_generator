package model

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/text/message"
)

type Tables []*Table

type Table struct {
	Name         string
	Columns      []Column
	RecordCount  int
	Buffer       int
	CurrentIndex int
}

type Column struct {
	Name string
	Type string
	Rule Rule
}

type Rule struct {
	Type     string
	Format   string
	Value    string
	Min      int
	Max      int
	Index    int
	Patterns []Pattern
}

type Pattern struct {
	Value string
	Times int
}

// setCurrentIndex set current index
func (t *Table) setCurrentIndex(idx int) {
	t.CurrentIndex = idx
}

// setIndex set index
func (r *Rule) setIndex(idx int) {
	r.Index = idx
}

// hasMax check has max
func (r Rule) hasMax() bool {
	return r.Max != 0
}

// rangeNumber get number in (c.Min, c.Max)
func (r Rule) rangeNumber() int {
	number := r.Min + r.Index
	if r.hasMax() && number > r.Max {
		number = r.Min
	}
	return number
}

// patterValue get value for pattern
func (r Rule) patterValue() string {
	var petterns []string
	var num_total int
	for _, pattern := range r.Patterns {
		for i := 0; i < pattern.Times; i++ {
			petterns = append(petterns, pattern.Value)
			num_total++
		}
	}
	return petterns[r.Index%num_total]
}

// genUUID generate UUID
func genUUID() string {
	uuid, _ := uuid.NewRandom()
	return uuid.String()
}

// genULID generate ULID
func genULID() string {
	return ulid.Make().String()
}

// queryValue get value for query
func (c Column) queryValue() string {
	var value string
	switch c.Rule.Type {
	case "const":
		switch c.Type {
		case "number":
			value = c.Rule.Value
		case "varchar", "timestamp":
			value = fmt.Sprintf(`'%s'`, c.Rule.Value)
		}
	case "unique":
		switch c.Type {
		case "number":
			value = strconv.Itoa(c.Rule.rangeNumber())
		case "varchar":
			switch c.Rule.Format {
			case "UUID":
				value = genUUID()
			case "ULID":
				value = genULID()
			}
			value = fmt.Sprintf(`'%s'`, value)
		case "timestamp":
			if c.Rule.Value == "now" {
				value = fmt.Sprintf(`'%s'`, time.Now().Format(time.RFC3339))
			}
			// log.Panicf("now_timestamp: %v", value)
		}
	case "pattern":
		if c.Type == "varchar" {
			value = fmt.Sprintf(`'%s'`, c.Rule.patterValue())
		}
	}
	return value
}

// queryRecord get record for query
func (t Table) queryRecord() string {
	var record []string
	for _, column := range t.Columns {
		column.Rule.setIndex(t.CurrentIndex)
		record = append(record, column.queryValue())
	}
	return fmt.Sprintf("(%s)", strings.Join(record, ","))
}

// QueryRecords get records for query
func (t Table) QueryRecords() []string {
	var bufferRecords []string
	var records []string
	for idx := 0; idx < t.RecordCount; idx++ {
		t.setCurrentIndex(idx)
		records = append(records, t.queryRecord())
		if idx%t.Buffer == t.Buffer-1 {
			bufferRecords = append(bufferRecords, strings.Join(records, ","))
			records = []string{}
		}
	}
	return bufferRecords
}

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
