package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"sirius/internal/models"
)

type WeatherService struct{}

type OpenMeteoResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Current   Current `json:"current"`
	Daily     Daily   `json:"daily"`
}

type Current struct {
	Time                string  `json:"time"`
	Temperature_2m      float64 `json:"temperature_2m"`
	RelativeHumidity_2m int     `json:"relative_humidity_2m"`
	ApparentTemperature float64 `json:"apparent_temperature"`
	WeatherCode         int     `json:"weather_code"`
	WindSpeed_10m       float64 `json:"wind_speed_10m"`
	WindDirection_10m   int     `json:"wind_direction_10m"`
	CloudCover          int     `json:"cloud_cover"`
	SurfacePressure     float64 `json:"surface_pressure"`
}

type Daily struct {
	Time                      []string  `json:"time"`
	WeatherCode               []int     `json:"weather_code"`
	Temperature_2mMax         []float64 `json:"temperature_2m_max"`
	Temperature_2mMin         []float64 `json:"temperature_2m_min"`
	PrecipitationSum          []float64 `json:"precipitation_sum"`
	PrecipitationProbabilityMax []int    `json:"precipitation_probability_max"`
}

func NewWeatherService() *WeatherService {
	return &WeatherService{}
}

func (ws *WeatherService) GetWeather() (*models.Weather, error) {
	url := "https://api.open-meteo.com/v1/forecast?latitude=9.8276&longitude=-84.0922&current=temperature_2m,relative_humidity_2m,apparent_temperature,weather_code,wind_speed_10m,wind_direction_10m,cloud_cover,surface_pressure&daily=weather_code,temperature_2m_max,temperature_2m_min,precipitation_sum,precipitation_probability_max&timezone=auto"
	
	resp, err := http.Get(url)
	if err != nil {
		return ws.getFallbackWeather(), nil
	}
	defer resp.Body.Close()

	var result OpenMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Error decoding: %v\n", err)
		return ws.getFallbackWeather(), nil
	}

	current := result.Current
	
	// Convert daily forecast
	forecast := make([]models.ForecastDay, 0)
	for i := 0; i < len(result.Daily.Time) && i < 5; i++ {
		forecast = append(forecast, models.ForecastDay{
			Date:          result.Daily.Time[i],
			Condition:     ws.getCondition(result.Daily.WeatherCode[i]),
			Icon:          ws.getIcon(result.Daily.WeatherCode[i]),
			TempMax:       result.Daily.Temperature_2mMax[i],
			TempMin:       result.Daily.Temperature_2mMin[i],
			Precipitation: result.Daily.PrecipitationSum[i],
			PrecipProb:    result.Daily.PrecipitationProbabilityMax[i],
		})
	}

	return &models.Weather{
		Location:    "Aserrí, Costa Rica",
		Temp:        current.Temperature_2m,
		FeelsLike:   current.ApparentTemperature,
		Humidity:    current.RelativeHumidity_2m,
		Condition:   ws.getCondition(current.WeatherCode),
		Icon:        ws.getIcon(current.WeatherCode),
		UpdatedAt:   time.Now().Format("15:04"),
		WindSpeed:   current.WindSpeed_10m,
		WindDir:     ws.getWindDirection(current.WindDirection_10m),
		CloudCover:  current.CloudCover,
		Pressure:    current.SurfacePressure,
		Forecast:    forecast,
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

func (ws *WeatherService) getWindDirection(deg int) string {
	dirs := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	idx := (deg + 45) / 45
	if idx >= 8 {
		idx = 0
	}
	return dirs[idx]
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
		WindSpeed:  10.0,
		WindDir:    "NE",
		CloudCover: 20,
		Pressure:   1013.0,
		Forecast:   []models.ForecastDay{},
	}
}
