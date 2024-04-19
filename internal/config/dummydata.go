package config

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/ponyo877/dummy_data_generator/internal/model"
	"github.com/spf13/viper"
)

type DummyTables []DummyTable

type DummyTable struct {
	TableName   string        `mapstructure:"tablename"`
	RecordCount int           `mapstructure:"recordcount"`
	Buffer      int           `mapstructure:"buffer"`
	Columns     []DummyColumn `mapstructure:"columns"`
}

type DummyColumn struct {
	Name string    `mapstructure:"name"`
	Type string    `mapstructure:"type"`
	Rule DummyRule `mapstructure:"rule"`
}

type DummyRule struct {
	Type     string         `mapstructure:"type"`
	Value    string         `mapstructure:"value"`
	Min      int            `mapstructure:"min"`
	Max      int            `mapstructure:"max"`
	MinTime  string         `mapstructure:"min_time"`
	MaxTime  string         `mapstructure:"max_time"`
	Format   string         `mapstructure:"format"`
	Patterns []DummyPattern `mapstructure:"patterns"`
}

type DummyPattern struct {
	Value string `mapstructure:"value"`
	Times int    `mapstructure:"times"`
}

func ParseTimeStringToUnixTime(time_string string, default_unix_time int64) (int64, error) {
	if time_string == "" {
		// if expected parameter is not passed, fallback to default_unix_time
		return default_unix_time, nil
	} else if t, err := time.Parse("2006-01-02 15:04:05", time_string); err != nil {
		return 0, err
	} else {
		return t.Unix(), nil
	}
}

// LoadDummyDataConfig
func LoadDummyDataConfig() (DummyTables, error) {
	dummyTables := DummyTables{}
	configfileRegexs := viper.GetStringSlice("config")
	for _, configfileRegex := range configfileRegexs {
		configfiles, err := filepath.Glob(configfileRegex)
		if err != nil {
			return DummyTables{}, err
		}
		for _, configfile := range configfiles {
			v := viper.New()
			v.SetConfigFile(configfile)
			if err := v.ReadInConfig(); err != nil {
				return DummyTables{}, err
			}
			var dummyTable DummyTable
			if err := v.Unmarshal(&dummyTable); err != nil {
				return DummyTables{}, err
			}
			dummyTables = append(dummyTables, dummyTable)
		}
	}
	// str, err := json.Marshal(dummyTables)
	// if err != nil {
	// 	return DummyTables{}, err
	// }
	// fmt.Printf("config: %v\n", string(str))
	return dummyTables, nil
}

// Tables get tables
func (t DummyTables) Tables() model.Tables {
	var tables model.Tables
	for _, table := range t {
		tables = append(tables, &model.Table{
			Name: table.TableName,
		})
	}
	return tables
}

// ToModels convert to models
func (t DummyTables) ToModels() (model.Tables, error) {
	var tables model.Tables
	for _, dummyTable := range t {
		var columns []model.Column
		for _, dummyColumn := range dummyTable.Columns {
			var rule model.Rule
			errMsg := "Rule:%v is not supported for Type:%v"
			switch dummyColumn.Rule.Type {
			case "const":
				switch dummyColumn.Type {
				case "number", "varchar", "timestamp":
					rule.Value = dummyColumn.Rule.Value
				default:
					return nil, fmt.Errorf(errMsg, dummyColumn.Rule.Type, dummyColumn.Type)
				}
			case "unique":
				switch dummyColumn.Type {
				case "number":
					rule.Min = dummyColumn.Rule.Min
					rule.Max = 2147483647 // intMax
				case "varchar", "timestamp":
					rule.Format = dummyColumn.Rule.Format
				default:
					return nil, fmt.Errorf(errMsg, dummyColumn.Rule.Type, dummyColumn.Type)
				}
			case "pattern":
				switch dummyColumn.Type {
				case "number":
					rule.Min = dummyColumn.Rule.Min
					rule.Max = dummyColumn.Rule.Max
					fallthrough
				case "varchar":
					patterns := make([]model.Pattern, len(dummyColumn.Rule.Patterns))
					for i, pattern := range dummyColumn.Rule.Patterns {
						patterns[i] = model.Pattern{Value: pattern.Value, Times: max(pattern.Times, 1)}
					}
					rule.Patterns = patterns
				default:
					return nil, fmt.Errorf(errMsg, dummyColumn.Rule.Type, dummyColumn.Type)
				}
			case "random":
				switch dummyColumn.Type {
				case "timestamp":
					var err error
					if rule.MinUnixTime, err = ParseTimeStringToUnixTime(
						dummyColumn.Rule.MinTime,
						time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix(),
					); err != nil {
						return nil, err
					}
					if rule.MaxUnixTime, err = ParseTimeStringToUnixTime(
						dummyColumn.Rule.MaxTime,
						time.Now().Unix(),
					); err != nil {
						return nil, err
					}
				default:
					return nil, fmt.Errorf(errMsg, dummyColumn.Rule.Type, dummyColumn.Type)
				}
			default:
				return nil, fmt.Errorf(errMsg, dummyColumn.Rule.Type, "ALL")
			}
			rule.Type = dummyColumn.Rule.Type
			columns = append(columns, model.Column{
				Name: dummyColumn.Name,
				Type: dummyColumn.Type,
				Rule: rule,
			})
		}
		tables = append(tables, &model.Table{
			Buffer:      dummyTable.Buffer,
			Name:        dummyTable.TableName,
			Columns:     columns,
			RecordCount: dummyTable.RecordCount,
		})
	}
	return tables, nil
}
