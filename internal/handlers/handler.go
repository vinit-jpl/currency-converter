package handlers

import "currency-converter/internal/service"

type ConverterHandler struct {
	Service *service.ConverterService
}

func NewConverterHandler(service *service.ConverterService) *ConverterHandler {
	return &ConverterHandler{
		Service: service,
	}
}
