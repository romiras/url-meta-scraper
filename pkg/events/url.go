package events

type UrlEvent struct {
	URL      string `json:"url"`
	Attempts uint   `json:"attempts,omitempty"`
}
