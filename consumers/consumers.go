package consumers

import "time"

type (
	IConsumer interface {
		Consume() (<-chan interface{}, error)
		Close() error
	}

	Delivery struct {
		ContentType     string // MIME content type
		ContentEncoding string // MIME content encoding
		Body            []byte
		MessageId       string    // application use - message identifier
		Timestamp       time.Time // application use - message timestamp
		Type            string    // application use - message type name
	}
)
