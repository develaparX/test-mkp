package customer

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	service *CustomerService
}

func NewCustomerHandler(service *CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func parseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (h *CustomerHandler) Create(c *gin.Context) {

	if err := c.ShouldBindJSON(&CreateCustomerReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	birthDate, err := parseDate(CreateCustomerReq.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid birth_date format, expected YYYY-MM-DD"})
		return
	}

	cust := &Customer{
		ID:          CreateCustomerReq.ID,
		FirstName:   CreateCustomerReq.FirstName,
		LastName:    CreateCustomerReq.LastName,
		BirthDate:   birthDate,
		CompanyID:   CreateCustomerReq.CompanyID,
		Email:       CreateCustomerReq.Email,
		PhoneNumber: CreateCustomerReq.PhoneNumber,
		Address:     CreateCustomerReq.Address,
		Gender:      CreateCustomerReq.Gender,
		Photo:       CreateCustomerReq.Photo,
	}

	if err := h.service.Create(cust); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cust)
}

func (h *CustomerHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}

	cust, err := h.service.GetByID(id)
	if err != nil {
		if err == ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cust)
}

func (h *CustomerHandler) GetAll(c *gin.Context) {
	customers, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customers)
}

func (h *CustomerHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}

	if err := c.ShouldBindJSON(&UpdateCustomerReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	birthDate, err := parseDate(UpdateCustomerReq.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid birth_date format, expected YYYY-MM-DD"})
		return
	}

	cust := &Customer{
		ID:          id,
		FirstName:   UpdateCustomerReq.FirstName,
		LastName:    UpdateCustomerReq.LastName,
		BirthDate:   birthDate,
		CompanyID:   UpdateCustomerReq.CompanyID,
		Email:       UpdateCustomerReq.Email,
		PhoneNumber: UpdateCustomerReq.PhoneNumber,
		Address:     UpdateCustomerReq.Address,
		Gender:      UpdateCustomerReq.Gender,
		Photo:       UpdateCustomerReq.Photo,
	}

	if err := h.service.Update(cust); err != nil {
		if err == ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cust)
}

func (h *CustomerHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		if err == ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
