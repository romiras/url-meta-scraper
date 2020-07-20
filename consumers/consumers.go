package consumers

import (
	"time"

	"github.com/romiras/url-meta-scraper/log"
	"github.com/streadway/amqp"
)

type (
	HandleFunc func(<-chan amqp.Delivery /*chan error,*/, log.Logger)

	IConsumer interface {
		Consume(handleFunc HandleFunc, logger log.Logger) error
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
