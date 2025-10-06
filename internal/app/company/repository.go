package company

import "database/sql"

type CompanyRepository struct {
	DB *sql.DB
}

func NewCompanyrepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{DB: db}
}

func (r *CompanyRepository) GetCompany() error {
	return nil
}
