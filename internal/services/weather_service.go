package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"sirius/internal/models"
)

type WeatherService struct {
	APIKey string
}

func NewWeatherService() *WeatherService {
	return &WeatherService{
		APIKey: "", // Set via OPENWEATHER_API_KEY env
	}
}

func (ws *WeatherService) GetWeather() (*models.Weather, error) {
	// Coordinates: 9°49'39.5"N 84°05'32.0"W (Barrio Altavista, Aserri)
	lat := "9.8276"
	lon := "-84.0922"
	
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current=temperature_2m,relative_humidity_2m,apparent_temperature,weather_code", lat, lon)
	
	resp, err := http.Get(url)
	if err != nil {
		return ws.getFallbackWeather(), nil
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ws.getFallbackWeather(), nil
	}

	current := result["current"].(map[string]interface{})
	
	weatherCode := int(current["weather_code"].(float64))
	condition := ws.getCondition(weatherCode)

	return &models.Weather{
		Location:   "Aserrí, Costa Rica",
		Temp:       current["temperature_2m"].(float64),
		FeelsLike:  current["apparent_temperature"].(float64),
		Humidity:   int(current["relative_humidity_2m"].(float64)),
		Condition:  condition,
		Icon:       ws.getIcon(weatherCode),
		UpdatedAt:  time.Now().Format("15:04"),
	}, nil
}

func (ws *WeatherService) getCondition(code int) string {
	switch {
	case code == 0: return "Despejado"
	case code <= 3: return "Parcialmente nublado"
	case code <= 49: return "Niebla"
	case code <= 59: return "Llovizna"
	case code <= 69: return "Lluvia"
	case code <= 79: return "Nieve"
	case code <= 99: return "Tormenta"
	default: return "Desconocido"
	}
}

func (ws *WeatherService) getIcon(code int) string {
	switch {
	case code == 0: return "☀️"
	case code <= 3: return "⛅"
	case code <= 49: return "🌫️"
	case code <= 59: return "🌧️"
	case code <= 69: return "🌧️"
	case code <= 79: return "❄️"
	case code <= 99: return "⛈️"
	default: return "🌡️"
	}
}

func (ws *WeatherService) getFallbackWeather() *models.Weather {
	return &models.Weather{
		Location:   "Aserrí, Costa Rica",
		Temp:       24.0,
		FeelsLike:  25.0,
		Humidity:   70,
		Condition:  "Despejado",
		Icon:       "☀️",
		UpdatedAt:  time.Now().Format("15:04"),
	}
}
