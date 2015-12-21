package publisher

import (
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/publisher"
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
