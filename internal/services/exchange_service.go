package services

import (
	"encoding/json"
	"net/http"
	"time"
	"sirius/internal/models"
)

type ExchangeService struct {
	APIKey string
}

func NewExchangeService() *ExchangeService {
	return &ExchangeService{}
}

func (es *ExchangeService) GetExchange() (*models.ExchangeRate, error) {
	url := "https://open.er-api.com/v6/latest/USD"
	
	resp, err := http.Get(url)
	if err != nil {
		return es.getFallbackExchange(), nil
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return es.getFallbackExchange(), nil
	}

	if rates, ok := result["rates"].(map[string]interface{}); ok {
		if crc, ok := rates["CRC"].(float64); ok {
			// Calculate buy/sell with ~2% spread
			buy := crc * 0.98  // Compra (banco vende USD)
			sell := crc * 1.02 // Venta (banco compra USD)
			
			return &models.ExchangeRate{
				USDToCRC:   buy,
				SellRate:   &sell,
				LastUpdate: time.Now().Format("02/01/2006 15:04"),
			}, nil
		}
	}

	return es.getFallbackExchange(), nil
}

func (es *ExchangeService) getFallbackExchange() *models.ExchangeRate {
	sell := 545.00
	buy := 525.00
	return &models.ExchangeRate{
		USDToCRC:   buy,
		SellRate:   &sell,
		LastUpdate: time.Now().Format("02/01/2006 15:04"),
	}
}
