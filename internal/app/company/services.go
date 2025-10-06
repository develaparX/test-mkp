package company

import (
	"errors"
)

var (
	ErrNotFound = errors.New("company not found")
)

type CompanyService struct {
	repo *CompanyRepo
}

func NewCompanyService(repo *CompanyRepo) *CompanyService {
	return &CompanyService{repo: repo}
}

func (s *CompanyService) Create(company *Company) error {
	return s.repo.Create(company)
}

func (s *CompanyService) GetByID(id int64) (Company, error) {
	company, err := s.repo.GetByID(id)
	if err != nil {
		return Company{}, err
	}
	if company == nil {
		return Company{}, ErrNotFound
	}
	return *company, nil
}

func (s *CompanyService) GetAll() ([]Company, error) {
	companies, err := s.repo.GetAll()
	if err != nil {
		return []Company{}, err
	}
	
	result := make([]Company, len(companies))
	for i, c := range companies {
		if c != nil {
			result[i] = *c
		}
	}
	return result, nil
}

func (s *CompanyService) Update(company *Company) error {
	existing, err := s.repo.GetByID(company.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNotFound
	}
	return s.repo.Update(company)
}

func (s *CompanyService) Delete(id int64) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNotFound
	}
	return s.repo.Delete(id)
}
