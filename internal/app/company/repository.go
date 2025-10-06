package company

import (
	"database/sql"
	"fmt"
)

type CompanyRepo struct {
	DB *sql.DB
}

func NewCompanyRepo(db *sql.DB) *CompanyRepo {
	return &CompanyRepo{DB: db}
}

func (r *CompanyRepo) Create(company *Company) error {
	query := `INSERT INTO company (id, name, type, address, city) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.DB.Exec(query, company.ID, company.Name, company.Type, company.Address, company.City)
	if err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}
	return nil
}

func (r *CompanyRepo) GetByID(id int64) (*Company, error) {
	query := `SELECT id, name, type, address, city FROM company WHERE id = $1`
	row := r.DB.QueryRow(query, id)

	var c Company
	err := row.Scan(&c.ID, &c.Name, &c.Type, &c.Address, &c.City)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get company by id: %w", err)
	}
	return &c, nil
}

func (r *CompanyRepo) GetAll() ([]*Company, error) {
	query := `SELECT id, name, type, address, city FROM company`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all companies: %w", err)
	}
	defer rows.Close()

	companies := make([]*Company, 0)
	for rows.Next() {
		var c Company
		err := rows.Scan(&c.ID, &c.Name, &c.Type, &c.Address, &c.City)
		if err != nil {
			return nil, fmt.Errorf("failed to scan company row: %w", err)
		}
		companies = append(companies, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return companies, nil
}

func (r *CompanyRepo) Update(company *Company) error {
	query := `UPDATE company SET name = $1, type = $2, address = $3, city = $4 WHERE id = $5`
	res, err := r.DB.Exec(query, company.Name, company.Type, company.Address, company.City, company.ID)
	if err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no company found with id %d", company.ID)
	}

	return nil
}

func (r *CompanyRepo) Delete(id int64) error {
	query := `DELETE FROM company WHERE id = $1`
	res, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no company found with id %d", id)
	}

	return nil
}
