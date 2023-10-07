package repository

import (
	"github.com/ponyo877/dummy_data_generator/internal/model"
	"gorm.io/gorm"
)

type GenerateRepository struct {
	db *gorm.DB
}

func NewGenerateRepository(db *gorm.DB, dbname string) *GenerateRepository {
	return &GenerateRepository{
		db,
	}
}

// Count count current table data
func (r GenerateRepository) Count() (model.Tables, error) {
	uncountedTables, err := r.listTableName()
	if err != nil {
		return nil, err
	}
	var tables model.Tables
	for _, table := range uncountedTables {
		var count int64
		r.db.Table(table.Name).Count(&count)
		tables = append(tables, &model.Table{Name: table.Name, RecordCount: int(count)})
	}
	return tables, nil
}

// listTableName list table name
func (r GenerateRepository) listTableName() (model.Tables, error) {
	var tables model.Tables
	// rows, err := r.db.Raw(`SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'`).Rows()
	rows, err := r.db.Select("tablename").Table("pg_tables").Where(`schemaname = 'public'`).Rows()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, &model.Table{Name: tableName})
	}
	return tables, nil
}

func (r GenerateRepository) Generate() {}
