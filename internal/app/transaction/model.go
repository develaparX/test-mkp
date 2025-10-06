package transaction

import (
	"errors"
	"time"
)

type TransactionSummary struct {
	ID            int64   `json:"id"`
	CompanyName   string  `json:"company_name"`
	ProductID     int64   `json:"product_id"`
	ProductName   string  `json:"product_name"`
	Amount        float64 `json:"amount"`
	Count         int64   `json:"count"`
	TaxValue      float64 `json:"tax_value"`
	ServiceFeePct bool    `json:"service_fee_percentage"`
	ServiceFee    float64 `json:"service_fee"`
	LastTrxOn     string  `json:"last_trx_on"`
	IDLastTrx     int64   `json:"id_last_trx"`
	FirstTrxOn    string  `json:"first_trx_on"`
	IDFirstTrx    int64   `json:"id_first_trx"`
}

var CreateTxReq struct {
	ID                     int64  `json:"id" binding:"required,min=1"`
	CustomerID             int64  `json:"customer_id" binding:"required,min=1"`
	TransactionType        string `json:"transaction_type" binding:"required,oneof=purchase refund payment"`
	Amount                 string `json:"amount" binding:"required"`
	TransactionDatetimeStr string `json:"transaction_datetime"`
	TaxAmount              string `json:"tax_amount" binding:"required"`
	TaxType                string `json:"tax_type,omitempty" binding:"omitempty,oneof=VAT GST SALES_TAX"`
	PaymentStatus          string `json:"payment_status" binding:"required,oneof=pending completed failed cancelled"`
	ProductID              int64  `json:"product_id" binding:"required,min=1"`
}

type Transaction struct {
	ID                  int64     `json:"id"`
	CustomerID          int64     `json:"customer_id"`
	TransactionType     string    `json:"transaction_type"`
	Amount              float64   `json:"amount"`
	TransactionDatetime time.Time `json:"transaction_datetime"`
	TaxAmount           float64   `json:"tax_amount"`
	TaxType             string    `json:"tax_type,omitempty"`
	PaymentStatus       string    `json:"payment_status"`
	ProductID           int64     `json:"product_id"`
}

type CustomerActivity struct {
	RowNumber   int64  `json:"row_number"`
	CompanyID   int64  `json:"company_id"`
	CompanyName string `json:"company_name"`
	CustomerID  int64  `json:"customer_id"`
	FullName    string `json:"full_name"`
	CountTrx    int64  `json:"count_trx"`
}

type Pagination struct {
	Page       int64 `json:"page"`
	PageSize   int64 `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int64 `json:"total_pages"`
}

type CustomerActivityResponse struct {
	Data       []CustomerActivity `json:"data"`
	Pagination Pagination         `json:"pagination"`
}

type CustomerActivityFilter struct {
	CompanyID   *int64
	MinTrxCount *int64
	Page        int64
	PageSize    int64
}

type TransactionSummaryFilter struct {
	CompanyID *int64
	ProductID *int64
	StartDate *time.Time
	EndDate   *time.Time
	MinAmount *float64
	MaxAmount *float64
	Page      int64
	PageSize  int64
}

type TransactionSummaryResponse struct {
	Data       []TransactionSummary `json:"data"`
	Pagination Pagination           `json:"pagination"`
}

var (
	ErrInvalidPage            = errors.New("page must be >= 1")
	ErrInvalidPageSize        = errors.New("page_size must be between 1 and 100")
	ErrInvalidAmount          = errors.New("amount must be greater than 0")
	ErrInvalidTaxAmount       = errors.New("tax_amount must be >= 0")
	ErrInvalidTransactionType = errors.New("transaction_type must be one of: purchase, refund, payment")
	ErrInvalidPaymentStatus   = errors.New("payment_status must be one of: pending, completed, failed, cancelled")
	ErrInvalidTaxType         = errors.New("tax_type must be one of: VAT, GST, SALES_TAX, or empty")
	ErrInvalidDateRange       = errors.New("start_date must be before end_date")
	ErrInvalidAmountRange     = errors.New("min_amount must be less than max_amount")
	ErrFutureTransactionDate  = errors.New("transaction_datetime cannot be in the future")
	ErrInvalidCustomerID      = errors.New("customer_id must be greater than 0")
	ErrInvalidProductID       = errors.New("product_id must be greater than 0")

	ErrCustomerNotFound            = errors.New("customer not found")
	ErrProductNotFound             = errors.New("product not found")
	ErrCustomerCompanyMismatch     = errors.New("customer does not belong to the same company as the product")
	ErrRefundExceedsOriginal       = errors.New("refund amount exceeds original purchase amount")
	ErrInsufficientPurchaseHistory = errors.New("insufficient purchase history for refund")
	ErrDuplicateTransactionID      = errors.New("transaction ID already exists")
)

var (
	ValidTransactionTypes = []string{"purchase", "refund", "payment"}
	ValidPaymentStatuses  = []string{"pending", "completed", "failed", "cancelled"}
	ValidTaxTypes         = []string{"VAT", "GST", "SALES_TAX", ""}
)

func (t *Transaction) Validate() error {

	if t.CustomerID <= 0 {
		return ErrInvalidCustomerID
	}
	if t.ProductID <= 0 {
		return ErrInvalidProductID
	}

	if t.Amount <= 0 {
		return ErrInvalidAmount
	}
	if t.TaxAmount < 0 {
		return ErrInvalidTaxAmount
	}

	if !isValidTransactionType(t.TransactionType) {
		return ErrInvalidTransactionType
	}

	if !isValidPaymentStatus(t.PaymentStatus) {
		return ErrInvalidPaymentStatus
	}

	if t.TaxType != "" && !isValidTaxType(t.TaxType) {
		return ErrInvalidTaxType
	}

	if !t.TransactionDatetime.IsZero() && t.TransactionDatetime.After(time.Now()) {
		return ErrFutureTransactionDate
	}

	return nil
}

func (f *TransactionSummaryFilter) Validate() error {

	if f.Page < 1 {
		return ErrInvalidPage
	}
	if f.PageSize < 1 || f.PageSize > 100 {
		return ErrInvalidPageSize
	}

	if f.StartDate != nil && f.EndDate != nil {
		if f.StartDate.After(*f.EndDate) {
			return ErrInvalidDateRange
		}
	}

	if f.MinAmount != nil && f.MaxAmount != nil {
		if *f.MinAmount >= *f.MaxAmount {
			return ErrInvalidAmountRange
		}
	}

	if f.MinAmount != nil && *f.MinAmount < 0 {
		return errors.New("min_amount must be >= 0")
	}
	if f.MaxAmount != nil && *f.MaxAmount < 0 {
		return errors.New("max_amount must be >= 0")
	}

	if f.CompanyID != nil && *f.CompanyID <= 0 {
		return errors.New("company_id must be greater than 0")
	}
	if f.ProductID != nil && *f.ProductID <= 0 {
		return errors.New("product_id must be greater than 0")
	}

	return nil
}

func isValidTransactionType(transactionType string) bool {
	for _, valid := range ValidTransactionTypes {
		if transactionType == valid {
			return true
		}
	}
	return false
}

func isValidPaymentStatus(paymentStatus string) bool {
	for _, valid := range ValidPaymentStatuses {
		if paymentStatus == valid {
			return true
		}
	}
	return false
}

func isValidTaxType(taxType string) bool {
	for _, valid := range ValidTaxTypes {
		if taxType == valid {
			return true
		}
	}
	return false
}
