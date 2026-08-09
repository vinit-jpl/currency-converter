package handlers

import (
	"currency-converter/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ConverterHandler) ConvertCurrency(c *gin.Context) {
	var req dto.ConvertCurrencyRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp, err := h.Service.ConvertCurrency(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ConverterHandler) GetCurrencies(c *gin.Context) {
	resp, err := h.Service.GetCurrencies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
