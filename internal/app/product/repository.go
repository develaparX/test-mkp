package product

import (
	"database/sql"
	"fmt"
)

type ProductRepo struct {
	DB *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{DB: db}
}

func (r *ProductRepo) Create(p *Product) error {
	query := `
		INSERT INTO product (id, product_name, service_fee, service_fee_percentage)
		VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, p.ID, p.ProductName, p.ServiceFee, p.ServiceFeePercentage)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (r *ProductRepo) GetByID(id int64) (*Product, error) {
	query := `
		SELECT id, product_name, service_fee, service_fee_percentage
		FROM product WHERE id = $1`
	row := r.DB.QueryRow(query, id)

	var p Product
	var serviceFee float64
	var isPercentage bool

	err := row.Scan(&p.ID, &p.ProductName, &serviceFee, &isPercentage)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan product: %w", err)
	}
	p.ServiceFee = serviceFee
	p.ServiceFeePercentage = isPercentage

	return &p, nil
}

func (r *ProductRepo) GetAll() ([]*Product, error) {
	query := `
		SELECT id, product_name, service_fee, service_fee_percentage
		FROM product`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	products := make([]*Product, 0)
	for rows.Next() {
		var p Product
		var serviceFee float64
		var isPercentage bool

		err := rows.Scan(&p.ID, &p.ProductName, &serviceFee, &isPercentage)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product row: %w", err)
		}
		p.ServiceFee = serviceFee
		p.ServiceFeePercentage = isPercentage

		products = append(products, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return products, nil
}

func (r *ProductRepo) Update(p *Product) error {
	query := `
		UPDATE product
		SET product_name = $1, service_fee = $2, service_fee_percentage = $3
		WHERE id = $4`
	res, err := r.DB.Exec(query, p.ProductName, p.ServiceFee, p.ServiceFeePercentage, p.ID)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no product found with id %d", p.ID)
	}

	return nil
}

func (r *ProductRepo) Delete(id int64) error {
	query := `DELETE FROM product WHERE id = $1`
	res, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no product found with id %d", id)
	}

	return nil
}
