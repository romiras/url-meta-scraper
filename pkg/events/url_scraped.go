package events

type UrlScrapedEvent struct {
	URL       string            `json:"url"`
	UpdatedAt int64             `json:"updated_at"`
	Body      string            `json:"body"`
	Headers   map[string]string `json:"headers"`
}
