package config

import (
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
	Type    string       `mapstructure:"type"`
	Value   string       `mapstructure:"value"`
	Min     int          `mapstructure:"min"`
	Max     int          `mapstructure:"max"`
	Format  string       `mapstructure:"format"`
	Pattern DummyPattern `mapstructure:"pattern"`
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
	// str, err := json.Marshal(config)
	// if err != nil {
	// 	return DummyDataConfig{}, err
	// }
	// fmt.Printf("config: %v\n", string(str))
	return config, nil
}

// Tables convert DummyDataConfig to model.Table
func Tables(config DummyDataConfig) (model.Tables, error) {
	var tables model.Tables
	for _, dummyTable := range config.DummyTables {
		var columns []model.Column
		for _, dummyColumn := range dummyTable.Columns {
			var rule model.Rule
			switch dummyColumn.Rule.Type {
			case "const":
				if dummyColumn.Type != "number" && dummyColumn.Type != "varchar" {
					return nil, fmt.Errorf("Rule:%vはType:%vはサポートしていません", dummyColumn.Rule.Type, dummyColumn.Type)
				}
				rule.Value = dummyColumn.Rule.Value
			case "unique":
				if dummyColumn.Type != "number" {
					return nil, fmt.Errorf("Rule:%vはType:%vはサポートしていません", dummyColumn.Rule.Type, dummyColumn.Type)
				}
				rule = model.Rule{Min: dummyColumn.Rule.Min, Max: dummyColumn.Rule.Max}
			default:
				return nil, fmt.Errorf("Rule:%vはサポートしていません", dummyColumn.Rule.Type)
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
