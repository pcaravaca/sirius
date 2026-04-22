package services

import (
	"encoding/json"
	"net/http"
	"time"
	"sirius/internal/models"
)

type NewsService struct{}

func NewNewsService() *NewsService {
	return &NewsService{}
}

func (ns *NewsService) GetNews() ([]models.NewsItem, error) {
	news := []models.NewsItem{}
	
	// Fetch Costa Rica news
	crNews, err := ns.fetchCostaRicaNews()
	if err == nil {
		news = append(news, crNews...)
	}
	
	// Fetch Bitcoin/Crypto news
	cryptoNews, err := ns.fetchCryptoNews()
	if err == nil {
		news = append(news, cryptoNews...)
	}
	
	// Fetch Tech news
	techNews, err := ns.fetchTechNews()
	if err == nil {
		news = append(news, techNews...)
	}
	
	// If no news, return some default items
	if len(news) == 0 {
		news = ns.getDefaultNews()
	}
	
	// Sort by date (newest first)
	for i := 0; i < len(news)-1; i++ {
		for j := i + 1; j < len(news); j++ {
			if news[i].PublishedAt > news[j].PublishedAt {
				news[i], news[j] = news[j], news[i]
			}
		}
	}
	
	return news, nil
}

func (ns *NewsService) fetchCostaRicaNews() ([]models.NewsItem, error) {
	url := "https://newsdata.io/api/1/news?apikey=pub_demo&q=Costa%20Rica&language=es&category=politics,business"
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	news := []models.NewsItem{}
	if results, ok := result["results"].([]interface{}); ok {
		for i, item := range results {
			if i >= 3 {
				break
			}
			if r, ok := item.(map[string]interface{}); ok {
				news = append(news, models.NewsItem{
					Title:       ns.getString(r, "title"),
					Description: ns.getString(r, "description"),
					Source:      ns.getString(r, "source_id"),
					URL:         ns.getString(r, "link"),
					Category:    "cr",
					PublishedAt: ns.getString(r, "pubDate"),
				})
			}
		}
	}
	
	return news, nil
}

func (ns *NewsService) fetchCryptoNews() ([]models.NewsItem, error) {
	url := "https://newsdata.io/api/1/news?apikey=pub_demo&q=Bitcoin&language=en&category=business"
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	news := []models.NewsItem{}
	if results, ok := result["results"].([]interface{}); ok {
		for i, item := range results {
			if i >= 2 {
				break
			}
			if r, ok := item.(map[string]interface{}); ok {
				news = append(news, models.NewsItem{
					Title:       ns.getString(r, "title"),
					Description: ns.getString(r, "description"),
					Source:      ns.getString(r, "source_id"),
					URL:         ns.getString(r, "link"),
					Category:    "crypto",
					PublishedAt: ns.getString(r, "pubDate"),
				})
			}
		}
	}
	
	return news, nil
}

func (ns *NewsService) fetchTechNews() ([]models.NewsItem, error) {
	url := "https://newsdata.io/api/1/news?apikey=pub_demo&q=technology&language=en&category=technology"
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	news := []models.NewsItem{}
	if results, ok := result["results"].([]interface{}); ok {
		for i, item := range results {
			if i >= 2 {
				break
			}
			if r, ok := item.(map[string]interface{}); ok {
				news = append(news, models.NewsItem{
					Title:       ns.getString(r, "title"),
					Description: ns.getString(r, "description"),
					Source:      ns.getString(r, "source_id"),
					URL:         ns.getString(r, "link"),
					Category:    "tech",
					PublishedAt: ns.getString(r, "pubDate"),
				})
			}
		}
	}
	
	return news, nil
}

func (ns *NewsService) getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func (ns *NewsService) getDefaultNews() []models.NewsItem {
	return []models.NewsItem{
		{
			Title:       "💰 Tipo de cambio: Dólar se mantiene estable",
			Description: "El dólar americano se mantiene respecto al Colón costarricense",
			Source:      "Sirius",
			Category:    "cr",
			PublishedAt: time.Now().Format("2006-01-02 15:04"),
		},
		{
			Title:       "₿ Bitcoin continúa su tendencia al alza",
			Description: "El mercado de criptomonedas muestra señales positivas",
			Source:      "Sirius",
			PublishedAt: time.Now().Format("2006-01-02 15:04"),
		},
		{
			Title:       "💻 Nuevas tecnologías en desarrollo",
			Description: "El mundo tech avanza rápidamente con nuevas herramientas",
			Source:      "Sirius",
			Category:    "tech",
			PublishedAt: time.Now().Format("2006-01-02 15:04"),
		},
	}
}
