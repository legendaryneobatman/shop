package list

type IService interface {
	Create(userID int, list List) (int, error)
	GetAll(userID int) ([]List, error)
	GetWithPagination(userID int, limit int, offset int) ([]List, error)
	GetByID(listID string) (List, error)
	Update(listID string, input List) (List, error)
}

type Service struct {
	repo *Repository
}

func NewListService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(userID int, list List) (int, error) {
	return s.repo.Create(userID, list)
}
func (s *Service) GetAll(userID int) ([]List, error) {
	return s.repo.GetAll(userID)
}
func (s *Service) GetByID(listID string) (List, error) { return s.repo.GetByID(listID) }
func (s *Service) Update(listID string, input List) (List, error) {
	return s.repo.Update(listID, input)
}
func (s *Service) GetWithPagination(userID int, limit int, offset int) ([]List, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	return s.repo.GetWithPagination(userID, limit, offset)
}
