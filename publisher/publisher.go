package publisher

import (
	"time"

	"github.com/elastic/libbeat/common"
	"github.com/elastic/libbeat/publisher"
)

type PhpfpmPublisher struct {
	client publisher.Client
}

func NewFpmPublisher(c publisher.Client) *PhpfpmPublisher {
	return &PhpfpmPublisher{client: c}
}

func (fp *PhpfpmPublisher) Publish(data map[string]interface{}) {
	fp.client.PublishEvent(common.MapStr{
		"@timestamp": common.Time(time.Now()),
		"type":       "phpfpm",
		"phpfpm":     data,
	})
}
