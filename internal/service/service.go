package service

import "currency-converter/internal/config"

type ConverterService struct {
	Cfg *config.Config
}

func NewConverterService(cfg *config.Config) *ConverterService {
	return &ConverterService{
		Cfg: cfg,
	}
}
