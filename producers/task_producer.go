package producers

import (
	"github.com/romiras/url-meta-scraper/options"
)

type TaskProducer interface {
	Produce(t *Task, options ...options.Option) (string, error)
	Close() error
}
