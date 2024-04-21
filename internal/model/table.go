package model

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
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
	Type        string
	Format      string
	Value       string
	Min         int
	Max         int
	MinUnixTime int64
	MaxUnixTime int64
	Index       int
	Patterns    []Pattern
}

type Pattern struct {
	Value string
	Times int
}

// setCurrentIndex set current index
func (t *Table) setCurrentIndex(idx int) {
	t.CurrentIndex = idx
}

// ColumnNames get column names
func (t *Table) ColumnNames() string {
	var columnNames []string
	for _, column := range t.Columns {
		columnNames = append(columnNames, column.Name)
	}
	return strings.Join(columnNames, ",")
}

// setIndex set index
func (r *Rule) setIndex(idx int) {
	r.Index = idx
}

// rangeNumber get number in (c.Min, c.Max)
func (r Rule) rangeNumber() int {
	number := (r.Min+r.Index)%(r.Max-r.Min+1) + r.Min
	return number
}

// patterValue get value for pattern
func (r Rule) patterValue() string {
	var patterns []string
	var num_total int
	for _, pattern := range r.Patterns {
		for i := 0; i < pattern.Times; i++ {
			patterns = append(patterns, pattern.Value)
			num_total++
		}
	}
	return patterns[r.Index%num_total]
}

// hasPattern check if rule has pattern
func (r Rule) hasPattern() bool {
	return len(r.Patterns) > 0
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
			if c.Rule.Format == "NOW" {
				value = fmt.Sprintf(`'%s'`, time.Now().Format(time.DateTime))
			}
		}
	case "pattern":
		switch c.Type {
		case "number":
			value = strconv.Itoa(c.Rule.rangeNumber())
			if c.Rule.hasPattern() {
				value = c.Rule.patterValue()
			}
		case "varchar":
			value = fmt.Sprintf(`'%s'`, c.Rule.patterValue())
		}
	case "random":
		switch c.Type {
		case "timestamp":
			mn := c.Rule.MinUnixTime
			mx := c.Rule.MaxUnixTime
			delta := mx - mn
			sec := rand.Int63n(delta) + mn
			value = fmt.Sprintf(`'%s'`, time.Unix(sec, 0).Format(time.RFC3339))
		}
	}
	return value
}

// queryValues get record for query
func (t Table) queryValues() string {
	var record []string
	for _, column := range t.Columns {
		column.Rule.setIndex(t.CurrentIndex)
		record = append(record, column.queryValue())
	}
	return fmt.Sprintf("(%s)", strings.Join(record, ","))
}

// BufferedValuesList get buffered records for query
func (t Table) BufferedValuesList() []string {
	var bufferValues []string
	var records []string
	for idx := 0; idx < t.RecordCount; idx++ {
		t.setCurrentIndex(idx)
		records = append(records, t.queryValues())
		if idx%t.Buffer == t.Buffer-1 || idx == t.RecordCount-1 {
			bufferValues = append(bufferValues, strings.Join(records, ","))
			records = []string{}
		}
	}
	return bufferValues
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
