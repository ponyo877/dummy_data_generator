package repository

import (
	"fmt"

	"github.com/ponyo877/dummy_data_generator/internal/model"
	"golang.org/x/sync/errgroup"
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
func (r GenerateRepository) Count(targetTables model.Tables) (model.Tables, error) {
	var tables model.Tables
	for _, table := range targetTables {
		var count int64
		r.db.Table(table.Name).Count(&count)
		tables = append(tables, &model.Table{Name: table.Name, RecordCount: int(count)})
	}
	return tables, nil
}

// Generate generate dummy data
func (r GenerateRepository) Generate(tables model.Tables) error {
	pbmap := make(map[string]model.Bar)
	p, wait := model.NewProgress()
	defer wait()
	for _, table := range tables {
		pbmap[table.Name] = p.AddBar(int64(table.RecordCount), table.Name)
	}
	var eg errgroup.Group
	queryTemplate := "INSERT INTO %s (%s) VALUES %s"
	for _, table := range tables {
		// parallel insert per table
		func(table *model.Table) {
			eg.Go(func() error {
				bufferedValuesList := table.BufferedValuesList()
				rest := table.RecordCount
				columns := table.ColumnNames()
				for _, bufferedValues := range bufferedValuesList {
					query := fmt.Sprintf(queryTemplate, table.Name, columns, bufferedValues)
					if err := r.db.Exec(query).Error; err != nil {
						return err
					}
					pbmap[table.Name].IncrInt64(int64(min(rest, table.Buffer)))
					rest -= table.Buffer
				}
				return nil
			})
		}(table)
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
