package generator

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
	tables, err := s.repository.Count()
	if err != nil {
		return err
	}
	return tables.Stdout()
}

// Generate generate dammy data
func (s Service) Generate() error {
	return nil
}
