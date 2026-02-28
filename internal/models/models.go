package models

import "time"

type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // active, paused, completed
	Priority    string    `json:"priority"` // high, medium, low
	Progress    int       `json:"progress"` // 0-100
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Task struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"projectId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // pending, in_progress, completed
	Priority    string    `json:"priority"` // high, medium, low
	DueDate     string    `json:"dueDate"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Weather struct {
	Location   string  `json:"location"`
	Temp       float64 `json:"temp"`
	FeelsLike  float64 `json:"feelsLike"`
	Humidity   int     `json:"humidity"`
	Condition  string  `json:"condition"`
	Icon       string  `json:"icon"`
	UpdatedAt  string  `json:"updatedAt"`
}

type ExchangeRate struct {
	USDToCRC   float64  `json:"usdToCrc"`   // Compra (banco vende USD)
	SellRate   *float64 `json:"sellRate,omitempty"` // Venta (banco compra USD)
	LastUpdate string   `json:"lastUpdate"`
}

type DashboardData struct {
	Projects   []Project     `json:"projects"`
	Tasks      []Task        `json:"tasks"`
	Weather    Weather       `json:"weather"`
	Exchange   ExchangeRate  `json:"exchange"`
	News       []NewsItem    `json:"news"`
	Motivational string      `json:"motivational"`
}

type NewsItem struct {
	Title   string `json:"title"`
	Source  string `json:"source"`
	URL     string `json:"url"`
}
