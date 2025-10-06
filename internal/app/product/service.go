package product

import (
	"errors"
)

var (
	ErrNotFound = errors.New("product not found")
)

type ProductService struct {
	repo *ProductRepo
}

func NewProductService(repo *ProductRepo) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(p *Product) error {

	return s.repo.Create(p)
}

func (s *ProductService) GetByID(id int64) (Product, error) {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return Product{}, err
	}
	if p == nil {
		return Product{}, ErrNotFound
	}
	return *p, nil
}

func (s *ProductService) GetAll() ([]Product, error) {
	products, err := s.repo.GetAll()
	if err != nil {
		return []Product{}, err
	}

	result := make([]Product, len(products))
	for i, p := range products {
		if p != nil {
			result[i] = *p
		}
	}
	return result, nil
}

func (s *ProductService) Update(p *Product) error {
	existing, err := s.repo.GetByID(p.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNotFound
	}
	return s.repo.Update(p)
}

func (s *ProductService) Delete(id int64) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNotFound
	}
	return s.repo.Delete(id)
}
