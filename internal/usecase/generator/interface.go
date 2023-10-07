package generator

import "github.com/ponyo877/dummy_data_generator/internal/model"

type Reader interface {
	ListTableName() (model.Tables, error)
	Count(model.Tables) (model.Tables, error)
}

type Writer interface {
	Generate(model.Tables) error
}

type Repository interface {
	Reader
	Writer
}

type Usecase interface {
	Count()
	Generate()
}
