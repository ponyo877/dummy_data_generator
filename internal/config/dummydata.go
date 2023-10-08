package config

import (
	"encoding/json"
	"fmt"

	"github.com/ponyo877/dummy_data_generator/internal/model"
	"github.com/spf13/viper"
)

type DummyDataConfig struct {
	Buffer      int          `mapstructure:"buffer"`
	DummyTables []DummyTable `mapstructure:"tables"`
}

type DummyTable struct {
	TableName   string        `mapstructure:"tablename"`
	RecordCount int           `mapstructure:"recordcount"`
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
	Format   string         `mapstructure:"format"`
	Patterns []DummyPattern `mapstructure:"patterns"`
}

type DummyPattern struct {
	Value string `mapstructure:"value"`
	Times int    `mapstructure:"times"`
}

// LoadDummyDataConfig
func LoadDummyDataConfig() (DummyDataConfig, error) {
	if err := viper.ReadInConfig(); err != nil {
		return DummyDataConfig{}, err
	}
	var config DummyDataConfig
	if err := viper.Unmarshal(&config); err != nil {
		return DummyDataConfig{}, err
	}
	str, err := json.Marshal(config)
	if err != nil {
		return DummyDataConfig{}, err
	}
	fmt.Printf("config: %v\n", string(str))
	return config, nil
}

// Tables convert DummyDataConfig to model.Table
func Tables(config DummyDataConfig) (model.Tables, error) {
	var tables model.Tables
	for _, dummyTable := range config.DummyTables {
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
					rule.Max = dummyColumn.Rule.Max
				case "varchar":
					rule.Format = dummyColumn.Rule.Format
				case "timestamp":
					rule.Value = dummyColumn.Rule.Value
				default:
					return nil, fmt.Errorf(errMsg, dummyColumn.Rule.Type, dummyColumn.Type)
				}
			case "pattern":
				switch dummyColumn.Type {
				case "varchar":
					patterns := make([]model.Pattern, len(dummyColumn.Rule.Patterns))
					for i, pattern := range dummyColumn.Rule.Patterns {
						patterns[i] = model.Pattern{Value: pattern.Value, Times: max(pattern.Times, 1)}
					}
					rule.Patterns = patterns
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
			Buffer:      config.Buffer,
			Name:        dummyTable.TableName,
			Columns:     columns,
			RecordCount: dummyTable.RecordCount,
		})
	}
	return tables, nil
}
