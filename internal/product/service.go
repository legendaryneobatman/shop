package product

type Service struct {
	repository         *Repository
	categoryRepository ICategoryRepository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}
