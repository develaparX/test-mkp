package customer

import "time"

type Customer struct {
	ID          int64     `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	BirthDate   time.Time `json:"birth_date,omitempty"`
	Email       string    `json:"email,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Address     string    `json:"address,omitempty"`
	Gender      string    `json:"gender,omitempty"`
	CompanyID   int64     `json:"company_id"`
	Photo       string    `json:"photo,omitempty"`
}

var CreateCustomerReq struct {
	ID          int64  `json:"id" binding:"required"`
	FirstName   string `json:"first_name" binding:"required,max=50"`
	LastName    string `json:"last_name" binding:"required,max=50"`
	BirthDate   string `json:"birth_date"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Gender      string `json:"gender"`
	CompanyID   int64  `json:"company_id" binding:"required"`
	Photo       string `json:"photo"`
}

var UpdateCustomerReq struct {
	FirstName   string `json:"first_name" binding:"required,max=50"`
	LastName    string `json:"last_name" binding:"required,max=50"`
	BirthDate   string `json:"birth_date"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Gender      string `json:"gender"`
	CompanyID   int64  `json:"company_id" binding:"required"`
	Photo       string `json:"photo"`
}
