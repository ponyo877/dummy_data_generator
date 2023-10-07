package repository

import (
	"fmt"

	"github.com/ponyo877/dummy_data_generator/internal/model"
	"gorm.io/gorm"
)

type GenerateRepository struct {
	db *gorm.DB
}

func NewGenerateRepository(db *gorm.DB) *GenerateRepository {
	return &GenerateRepository{
		db,
	}
}

// Count count current table data
func (r GenerateRepository) Count(uncountedTables model.Tables) (model.Tables, error) {
	var tables model.Tables
	for _, table := range uncountedTables {
		var count int64
		r.db.Table(table.Name).Count(&count)
		tables = append(tables, &model.Table{Name: table.Name, RecordCount: int(count)})
	}
	return tables, nil
}

// ListTableName list table name
func (r GenerateRepository) ListTableName() (model.Tables, error) {
	var tables model.Tables
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

// Generate generate dummy data
func (r GenerateRepository) Generate(tables model.Tables) error {
	for _, table := range tables {
		queryHeader := fmt.Sprintf("INSERT INTO %s VALUES ", table.Name)
		bufferedQuerys := table.QueryRecords()
		for _, bufferedQuery := range bufferedQuerys {
			query := queryHeader + bufferedQuery
			if err := r.db.Exec(query).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
