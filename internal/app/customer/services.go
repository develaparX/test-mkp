package customer

import (
	"errors"
)

var (
	ErrNotFound = errors.New("customer not found")
)

type CustomerService struct {
	repo *CustomerRepo
}

func NewCustomerService(repo *CustomerRepo) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) Create(c *Customer) error {
	return s.repo.Create(c)
}

func (s *CustomerService) GetByID(id int64) (Customer, error) {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return Customer{}, err
	}
	if c == nil {
		return Customer{}, ErrNotFound
	}
	return *c, nil
}

func (s *CustomerService) GetAll() ([]Customer, error) {
	customers, err := s.repo.GetAll()
	if err != nil {
		return []Customer{}, err
	}
	
	result := make([]Customer, len(customers))
	for i, c := range customers {
		if c != nil {
			result[i] = *c
		}
	}
	return result, nil
}

func (s *CustomerService) Update(c *Customer) error {
	existing, err := s.repo.GetByID(c.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNotFound
	}
	return s.repo.Update(c)
}

func (s *CustomerService) Delete(id int64) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNotFound
	}
	return s.repo.Delete(id)
}
