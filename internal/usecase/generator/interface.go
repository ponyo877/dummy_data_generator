package generator

import "github.com/ponyo877/dummy_data_generator/internal/model"

type Reader interface {
	Count() (model.Tables, error)
}

type Writer interface {
	Generate()
}

type Repository interface {
	Reader
	Writer
}

type Usecase interface {
	Count()
	Generate()
}
