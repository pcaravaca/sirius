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
	Location    string        `json:"location"`
	Temp        float64       `json:"temp"`
	FeelsLike   float64       `json:"feelsLike"`
	Humidity    int           `json:"humidity"`
	Condition   string        `json:"condition"`
	Icon        string        `json:"icon"`
	UpdatedAt   string        `json:"updatedAt"`
	WindSpeed   float64       `json:"windSpeed,omitempty"`
	WindDir     string        `json:"windDir,omitempty"`
	CloudCover  int           `json:"cloudCover,omitempty"`
	Pressure    float64       `json:"pressure,omitempty"`
	Forecast    []ForecastDay `json:"forecast,omitempty"`
}

type ForecastDay struct {
	Date          string  `json:"date"`
	Condition     string  `json:"condition"`
	Icon          string  `json:"icon"`
	TempMax       float64 `json:"tempMax"`
	TempMin       float64 `json:"tempMin"`
	Precipitation float64 `json:"precipitation"`
	PrecipProb    int     `json:"precipProb"`
}

type ExchangeRate struct {
	USDToCRC   float64  `json:"usdToCrc"`   // Compra (banco vende USD)
	SellRate   *float64 `json:"sellRate,omitempty"` // Venta (banco compra USD)
	LastUpdate string   `json:"lastUpdate"`
}

type DashboardData struct {
	Projects    []Project    `json:"projects"`
	Tasks       []Task       `json:"tasks"`
	Weather     Weather      `json:"weather"`
	Exchange    ExchangeRate `json:"exchange"`
	News        []NewsItem   `json:"news"`
	Motivational string      `json:"motivational"`
}

type NewsItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Source      string `json:"source"`
	URL         string `json:"url"`
	Category    string `json:"category"`
	PublishedAt string `json:"publishedAt"`
}

type Agent struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Emoji      string    `json:"emoji"`
	Model      string    `json:"model"`
	Status     string    `json:"status"` // idle, working, completed, error, standby
	Channel    string    `json:"channel"`
	Task       string    `json:"task,omitempty"`
	Output     string    `json:"output,omitempty"`
	StartedAt  time.Time `json:"startedAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	TokensIn   int64     `json:"tokensIn,omitempty"`
	TokensOut  int64     `json:"tokensOut,omitempty"`
	DurationMs int64     `json:"durationMs,omitempty"`
	CostUSD    float64   `json:"costUsd,omitempty"`
	Reasoning  []string  `json:"reasoning,omitempty"`
	LastOutput string    `json:"lastOutput,omitempty"`
}

// SystemInfo holds overall system status
type SystemInfo struct {
	CPU      CPUInfo    `json:"cpu"`
	RAM      RAMInfo    `json:"ram"`
	Disk     DiskInfo   `json:"disk"`
	Uptime   string     `json:"uptime"`
	Ollama   OllamaInfo `json:"ollama"`
	OpenClaw string     `json:"openClaw"`
}

type CPUInfo struct {
	Used  int `json:"used"`
	Cores int `json:"cores"`
}

type RAMInfo struct {
	Used  float64 `json:"used"`
	Total float64 `json:"total"`
}

type DiskInfo struct {
	Used  float64 `json:"used"`
	Total float64 `json:"total"`
}

type OllamaInfo struct {
	Loaded bool           `json:"loaded"`
	Status string         `json:"status"`
	Models []OllamaModel  `json:"models,omitempty"`
	VRAM   VRAMInfo       `json:"vram,omitempty"`
}

type OllamaModel struct {
	Name string `json:"name"`
	Size string `json:"size"`
	VRAM int    `json:"vram"`
}

type VRAMInfo struct {
	Used  float64 `json:"used"`
	Total float64 `json:"total"`
}
