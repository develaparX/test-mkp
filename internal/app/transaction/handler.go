package transaction

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service *TransactionService
}

func NewTransactionHandler(service *TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) Create(c *gin.Context) {

	if err := c.ShouldBindJSON(&CreateTxReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":                   "validation failed",
			"details":                 err.Error(),
			"valid_transaction_types": ValidTransactionTypes,
			"valid_payment_statuses":  ValidPaymentStatuses,
			"valid_tax_types":         ValidTaxTypes,
		})
		return
	}

	amount, err := strconv.ParseFloat(CreateTxReq.Amount, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount format, must be a valid number"})
		return
	}
	if amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be greater than 0"})
		return
	}

	taxAmount, err := strconv.ParseFloat(CreateTxReq.TaxAmount, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tax_amount format, must be a valid number"})
		return
	}
	if taxAmount < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tax_amount must be >= 0"})
		return
	}

	var trxTime time.Time
	if CreateTxReq.TransactionDatetimeStr != "" {
		trxTime, err = time.Parse(time.RFC3339, CreateTxReq.TransactionDatetimeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid transaction_datetime format, use RFC3339",
				"example": "2023-12-25T10:30:00Z",
			})
			return
		}

		if trxTime.After(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "transaction_datetime cannot be in the future"})
			return
		}
	}

	if CreateTxReq.TransactionType == "refund" && taxAmount > amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tax_amount cannot be greater than refund amount"})
		return
	}

	t := &Transaction{
		ID:                  CreateTxReq.ID,
		CustomerID:          CreateTxReq.CustomerID,
		TransactionType:     CreateTxReq.TransactionType,
		Amount:              amount,
		TaxAmount:           taxAmount,
		PaymentStatus:       CreateTxReq.PaymentStatus,
		ProductID:           CreateTxReq.ProductID,
		TransactionDatetime: trxTime,
		TaxType:             CreateTxReq.TaxType,
	}

	if err := h.service.Create(t); err != nil {
		switch {
		case err == ErrInvalidAmount || err == ErrInvalidTaxAmount ||
			err == ErrInvalidTransactionType || err == ErrInvalidPaymentStatus ||
			err == ErrInvalidTaxType || err == ErrFutureTransactionDate ||
			err == ErrInvalidCustomerID || err == ErrInvalidProductID:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, t)
}

func (h *TransactionHandler) GetAll(c *gin.Context) {
	txs, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, txs)
}

func (h *TransactionHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	tx, err := h.service.GetByID(id)
	if err != nil {
		if err.Error() == "Transaction not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tx)
}

func (h *TransactionHandler) GetTransactionSummary(c *gin.Context) {
	summaries, err := h.service.GetTransactionSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summaries)
}

func (h *TransactionHandler) GetTransactionSummaryFiltered(c *gin.Context) {
	var filter TransactionSummaryFilter

	page := int64(1)
	if pStr := c.DefaultQuery("page", "1"); pStr != "" {
		p, err := strconv.ParseInt(pStr, 10, 64)
		if err != nil || p < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid page parameter",
				"details": "page must be a positive integer >= 1",
			})
			return
		}
		page = p
	}

	pageSize := int64(10)
	if psStr := c.DefaultQuery("page_size", "10"); psStr != "" {
		ps, err := strconv.ParseInt(psStr, 10, 64)
		if err != nil || ps < 1 || ps > 100 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid page_size parameter",
				"details": "page_size must be between 1 and 100",
			})
			return
		}
		pageSize = ps
	}

	filter.Page = page
	filter.PageSize = pageSize

	if cidStr := c.Query("company_id"); cidStr != "" {
		cid, err := strconv.ParseInt(cidStr, 10, 64)
		if err != nil || cid <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid company_id parameter",
				"details": "company_id must be a positive integer",
			})
			return
		}
		filter.CompanyID = &cid
	}

	if pidStr := c.Query("product_id"); pidStr != "" {
		pid, err := strconv.ParseInt(pidStr, 10, 64)
		if err != nil || pid <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid product_id parameter",
				"details": "product_id must be a positive integer",
			})
			return
		}
		filter.ProductID = &pid
	}

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid start_date format",
				"details": "start_date must be in YYYY-MM-DD format",
				"example": "2023-01-01",
			})
			return
		}
		filter.StartDate = &startDate
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid end_date format",
				"details": "end_date must be in YYYY-MM-DD format",
				"example": "2023-12-31",
			})
			return
		}
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		filter.EndDate = &endDate
	}

	if minAmountStr := c.Query("min_amount"); minAmountStr != "" {
		minAmount, err := strconv.ParseFloat(minAmountStr, 64)
		if err != nil || minAmount < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid min_amount parameter",
				"details": "min_amount must be a non-negative number",
			})
			return
		}
		filter.MinAmount = &minAmount
	}

	if maxAmountStr := c.Query("max_amount"); maxAmountStr != "" {
		maxAmount, err := strconv.ParseFloat(maxAmountStr, 64)
		if err != nil || maxAmount < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid max_amount parameter",
				"details": "max_amount must be a non-negative number",
			})
			return
		}
		filter.MaxAmount = &maxAmount
	}

	if filter.StartDate != nil && filter.EndDate != nil {
		if filter.StartDate.After(*filter.EndDate) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid date range",
				"details": "start_date must be before or equal to end_date",
			})
			return
		}
	}

	if filter.MinAmount != nil && filter.MaxAmount != nil {
		if *filter.MinAmount >= *filter.MaxAmount {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid amount range",
				"details": "min_amount must be less than max_amount",
			})
			return
		}
	}

	resp, err := h.service.GetTransactionSummaryWithFilter(filter)
	if err != nil {
		switch {
		case err == ErrInvalidPage || err == ErrInvalidPageSize ||
			err == ErrInvalidDateRange || err == ErrInvalidAmountRange:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	response := gin.H{
		"data":       resp.Data,
		"pagination": resp.Pagination,
		"applied_filters": gin.H{
			"company_id": filter.CompanyID,
			"product_id": filter.ProductID,
			"start_date": filter.StartDate,
			"end_date":   filter.EndDate,
			"min_amount": filter.MinAmount,
			"max_amount": filter.MaxAmount,
		},
	}

	c.JSON(http.StatusOK, response)
}

func (h *TransactionHandler) GetCustomerActivity(c *gin.Context) {
	var companyID *int64
	if cidStr := c.Query("company_id"); cidStr != "" {
		cid, err := strconv.ParseInt(cidStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company_id"})
			return
		}
		companyID = &cid
	}

	var minTrxCount *int64
	if minStr := c.Query("min_trx"); minStr != "" {
		min, err := strconv.ParseInt(minStr, 10, 64)
		if err != nil || min < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid min_trx (must be non-negative integer)"})
			return
		}
		minTrxCount = &min
	}

	page := int64(1)
	if pStr := c.DefaultQuery("page", "1"); pStr != "" {
		p, err := strconv.ParseInt(pStr, 10, 64)
		if err != nil || p < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
			return
		}
		page = p
	}

	pageSize := int64(10)
	if psStr := c.DefaultQuery("page_size", "10"); psStr != "" {
		ps, err := strconv.ParseInt(psStr, 10, 64)
		if err != nil || ps < 1 || ps > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "page_size must be between 1 and 100"})
			return
		}
		pageSize = ps
	}

	resp, err := h.service.GetCustomerActivity(companyID, minTrxCount, page, pageSize)
	if err != nil {
		switch {
		case err == ErrInvalidPage || err == ErrInvalidPageSize:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}
