package events

type UrlFailedEvent struct {
	URL      string `json:"url"`
	FailedAt int64  `json:"failed_at"`
	Attempts uint   `json:"attempts"`
}
