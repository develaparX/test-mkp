package customer

import (
	"database/sql"
	"fmt"
)

type CustomerRepo struct {
	DB *sql.DB
}

func NewCustomerRepo(db *sql.DB) *CustomerRepo {
	return &CustomerRepo{DB: db}
}

func (r *CustomerRepo) Create(c *Customer) error {
	query := `
		INSERT INTO customer (
			id, first_name, last_name, birth_date, email, 
			phone_number, address, gender, company, photo
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	var birthDate interface{}
	if !c.BirthDate.IsZero() {
		birthDate = c.BirthDate.Format("2006-01-02")
	} else {
		birthDate = nil
	}

	var email interface{}
	if c.Email != "" {
		email = c.Email
	} else {
		email = nil
	}

	var phoneNumber interface{}
	if c.PhoneNumber != "" {
		phoneNumber = c.PhoneNumber
	} else {
		phoneNumber = nil
	}

	var address interface{}
	if c.Address != "" {
		address = c.Address
	} else {
		address = nil
	}

	var gender interface{}
	if c.Gender != "" {
		gender = c.Gender
	} else {
		gender = nil
	}

	var photo interface{}
	if c.Photo != "" {
		photo = c.Photo
	} else {
		photo = nil
	}

	_, err := r.DB.Exec(
		query,
		c.ID,
		c.FirstName,
		c.LastName,
		birthDate,
		email,
		phoneNumber,
		address,
		gender,
		c.CompanyID,
		photo,
	)
	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}
	return nil
}

func (r *CustomerRepo) GetByID(id int64) (*Customer, error) {
	query := `
		SELECT id, first_name, last_name, birth_date, email,
		       phone_number, address, gender, company, photo
		FROM customer WHERE id = $1`

	row := r.DB.QueryRow(query, id)

	var c Customer
	var birthDate sql.NullTime
	var email, phone, addr, gender, photo sql.NullString

	err := row.Scan(
		&c.ID,
		&c.FirstName,
		&c.LastName,
		&birthDate,
		&email,
		&phone,
		&addr,
		&gender,
		&c.CompanyID,
		&photo,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan customer: %w", err)
	}

	if birthDate.Valid {
		c.BirthDate = birthDate.Time
	}
	if email.Valid {
		c.Email = email.String
	}
	if phone.Valid {
		c.PhoneNumber = phone.String
	}
	if addr.Valid {
		c.Address = addr.String
	}
	if gender.Valid {
		c.Gender = gender.String
	}
	if photo.Valid {
		c.Photo = photo.String
	}

	return &c, nil
}

func (r *CustomerRepo) GetAll() ([]*Customer, error) {
	query := `
		SELECT id, first_name, last_name, birth_date, email,
		       phone_number, address, gender, company, photo
		FROM customer`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query customers: %w", err)
	}
	defer rows.Close()

	customers := make([]*Customer, 0)
	for rows.Next() {
		var c Customer
		var birthDate sql.NullTime
		var email, phone, addr, gender, photo sql.NullString

		err := rows.Scan(
			&c.ID,
			&c.FirstName,
			&c.LastName,
			&birthDate,
			&email,
			&phone,
			&addr,
			&gender,
			&c.CompanyID,
			&photo,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan customer row: %w", err)
		}

		if birthDate.Valid {
			c.BirthDate = birthDate.Time
		}
		if email.Valid {
			c.Email = email.String
		}
		if phone.Valid {
			c.PhoneNumber = phone.String
		}
		if addr.Valid {
			c.Address = addr.String
		}
		if gender.Valid {
			c.Gender = gender.String
		}
		if photo.Valid {
			c.Photo = photo.String
		}

		customers = append(customers, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return customers, nil
}

func (r *CustomerRepo) Update(c *Customer) error {
	query := `
		UPDATE customer
		SET first_name = $1, last_name = $2, birth_date = $3, email = $4,
		    phone_number = $5, address = $6, gender = $7, company = $8, photo = $9
		WHERE id = $10`

	var birthDate interface{}
	if !c.BirthDate.IsZero() {
		birthDate = c.BirthDate.Format("2006-01-02")
	} else {
		birthDate = nil
	}

	var email interface{}
	if c.Email != "" {
		email = c.Email
	} else {
		email = nil
	}

	var phoneNumber interface{}
	if c.PhoneNumber != "" {
		phoneNumber = c.PhoneNumber
	} else {
		phoneNumber = nil
	}

	var address interface{}
	if c.Address != "" {
		address = c.Address
	} else {
		address = nil
	}

	var gender interface{}
	if c.Gender != "" {
		gender = c.Gender
	} else {
		gender = nil
	}

	var photo interface{}
	if c.Photo != "" {
		photo = c.Photo
	} else {
		photo = nil
	}

	res, err := r.DB.Exec(
		query,
		c.FirstName,
		c.LastName,
		birthDate,
		email,
		phoneNumber,
		address,
		gender,
		c.CompanyID,
		photo,
		c.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no customer found with id %d", c.ID)
	}

	return nil
}

func (r *CustomerRepo) Delete(id int64) error {
	query := `DELETE FROM customer WHERE id = $1`
	res, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no customer found with id %d", id)
	}

	return nil
}
