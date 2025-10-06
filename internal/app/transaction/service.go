package transaction

import (
	"database/sql"
	"errors"
	"sinibeli/internal/app/customer"
	"sinibeli/internal/app/product"
	"time"
)

type TransactionService struct {
	Repo         *TransactionRepo
	CustomerRepo *customer.CustomerRepo
	ProductRepo  *product.ProductRepo
}

func NewTransactionService(repo *TransactionRepo, customerRepo *customer.CustomerRepo, productRepo *product.ProductRepo) *TransactionService {
	return &TransactionService{
		Repo:         repo,
		CustomerRepo: customerRepo,
		ProductRepo:  productRepo,
	}
}

func (s *TransactionService) Create(t *Transaction) error {

	if err := t.Validate(); err != nil {
		return err
	}

	if t.TransactionDatetime.IsZero() {
		t.TransactionDatetime = time.Now()
	}

	customer, err := s.CustomerRepo.GetByID(t.CustomerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrCustomerNotFound
		}
		return err
	}
	if customer == nil {
		return ErrCustomerNotFound
	}

	product, err := s.ProductRepo.GetByID(t.ProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrProductNotFound
		}
		return err
	}
	if product == nil {
		return ErrProductNotFound
	}

	existing, err := s.Repo.GetByID(t.ID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if existing != nil {
		return ErrDuplicateTransactionID
	}

	switch t.TransactionType {
	case "refund":

		purchaseHistory, err := s.Repo.GetTransactionsByCustomerAndProduct(t.CustomerID, t.ProductID)
		if err != nil {
			return err
		}

		if len(purchaseHistory) == 0 {
			return ErrInsufficientPurchaseHistory
		}

		totalPurchased := 0.0
		totalRefunded := 0.0
		var latestPurchase *Transaction

		for _, tx := range purchaseHistory {
			if tx.TransactionType == "purchase" {
				totalPurchased += tx.Amount
				if latestPurchase == nil || tx.TransactionDatetime.After(latestPurchase.TransactionDatetime) {
					latestPurchase = &tx
				}
			} else if tx.TransactionType == "refund" {
				totalRefunded += tx.Amount
			}
		}

		availableRefund := totalPurchased - totalRefunded
		if t.Amount > availableRefund {
			return ErrRefundExceedsOriginal
		}

		if t.Amount <= 0 {
			return errors.New("refund amount must be greater than 0")
		}

		if latestPurchase == nil {
			return errors.New("no purchase found for refund")
		}

		refundDeadline := latestPurchase.TransactionDatetime.AddDate(0, 0, 30)
		if t.TransactionDatetime.After(refundDeadline) {
			return errors.New("refund period has expired (30 days limit)")
		}

	case "purchase":

		if t.Amount < 1.0 {
			return errors.New("minimum purchase amount is 1.00")
		}

		maxPurchaseAmount := 1000000.0
		if t.Amount > maxPurchaseAmount {
			return errors.New("purchase amount exceeds maximum allowed limit")
		}

		if t.TaxAmount > 0 {
			taxPercentage := (t.TaxAmount / t.Amount) * 100

			if taxPercentage > 50.0 {
				return errors.New("tax amount seems unreasonably high (>50% of purchase amount)")
			}

			switch t.TaxType {
			case "VAT":
				if taxPercentage > 25.0 {
					return errors.New("VAT rate exceeds typical maximum (25%)")
				}
			case "GST":
				if taxPercentage > 15.0 {
					return errors.New("GST rate exceeds typical maximum (15%)")
				}
			case "SALES_TAX":
				if taxPercentage > 12.0 {
					return errors.New("sales tax rate exceeds typical maximum (12%)")
				}
			}
		}

	case "payment":
		if t.Amount <= 0 {
			return errors.New("payment amount must be greater than 0")
		}
	}

	return s.Repo.Create(t)
}

func (s *TransactionService) GetByID(id int64) (Transaction, error) {
	t, err := s.Repo.GetByID(id)
	if err != nil {
		return Transaction{}, err
	}
	if t == nil {
		return Transaction{}, errors.New("transaction not found")
	}
	return *t, nil
}

func (s *TransactionService) GetAll() ([]Transaction, error) {
	transactions, err := s.Repo.GetAll()
	if err != nil {
		return []Transaction{}, err
	}

	result := make([]Transaction, len(transactions))
	for i, t := range transactions {
		if t != nil {
			result[i] = *t
		}
	}
	return result, nil
}

func (s *TransactionService) GetTransactionSummary() ([]TransactionSummary, error) {
	return s.Repo.GetTransactionSummary()
}

func (s *TransactionService) GetTransactionSummaryWithFilter(filter TransactionSummaryFilter) (TransactionSummaryResponse, error) {

	if err := filter.Validate(); err != nil {
		return TransactionSummaryResponse{}, err
	}

	data, total, err := s.Repo.GetTransactionSummaryWithFilter(filter)
	if err != nil {
		return TransactionSummaryResponse{}, err
	}

	totalPages := total / filter.PageSize
	if total%filter.PageSize > 0 {
		totalPages++
	}

	return TransactionSummaryResponse{
		Data: data,
		Pagination: Pagination{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *TransactionService) GetCustomerActivity(
	companyID *int64,
	minTrxCount *int64,
	page int64,
	pageSize int64,
) (CustomerActivityResponse, error) {

	if page < 1 {
		return CustomerActivityResponse{}, ErrInvalidPage
	}
	if pageSize < 1 || pageSize > 100 {
		return CustomerActivityResponse{}, ErrInvalidPageSize
	}

	filter := CustomerActivityFilter{
		CompanyID:   companyID,
		MinTrxCount: minTrxCount,
		Page:        page,
		PageSize:    pageSize,
	}

	data, total, err := s.Repo.GetCustomerActivity(filter)
	if err != nil {
		return CustomerActivityResponse{}, err
	}

	totalPages := total / pageSize
	if total%pageSize > 0 {
		totalPages++
	}

	return CustomerActivityResponse{
		Data: data,
		Pagination: Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}
