package events

import "encoding/json"

type ScrapedEvent struct {
	Body     []byte            `json:"body"`
	Headers  map[string]string `json:"headers"`
	Metadata json.RawMessage   `json:"metadata,omitempty"`
}
