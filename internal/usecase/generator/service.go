package generator

import (
	"github.com/ponyo877/dummy_data_generator/internal/config"
)

// Service Generator usecase
type Service struct {
	repository Repository
}

// NewService create new service
func NewService(r Repository) *Service {
	return &Service{
		repository: r,
	}
}

// Count count dammy data
func (s Service) Count() error {
	dummyDataConfig, err := config.LoadDummyDataConfig()
	if err != nil {
		return err
	}
	tables, err := s.repository.Count(dummyDataConfig.Tables())
	if err != nil {
		return err
	}
	return tables.Stdout()
}

// Generate generate dammy data
func (s Service) Generate() error {
	dummyDataConfig, err := config.LoadDummyDataConfig()
	if err != nil {
		return err
	}
	tables, err := dummyDataConfig.ToModels()
	if err != nil {
		return err
	}
	// str, err2 := json.Marshal(tables)
	// if err2 != nil {
	// 	return err
	// }
	// fmt.Printf("model: %v\n", string(str))
	if err := s.repository.Generate(tables); err != nil {
		return err
	}
	return nil
}
