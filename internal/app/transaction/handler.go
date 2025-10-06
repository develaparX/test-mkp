package transaction

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service *TransactionService
}

func NewTransactionHandler(service *TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) GetTransactionSummary(c *gin.Context) {
	summaries, err := h.service.GetTransactionSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summaries)
}
