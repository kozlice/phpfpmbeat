package publisher

import (
	"time"

	"github.com/elastic/libbeat/common"
	"github.com/elastic/libbeat/publisher"
)

type PhpfpmPublisher struct {
	client publisher.Client
}

func New(c publisher.Client) *PhpfpmPublisher {
	return &PhpfpmPublisher{client: c}
}

func (pfb *PhpfpmPublisher) Publish(s map[string]interface{}) {
	pfb.client.PublishEvent(common.MapStr{
		"@timestamp": common.Time(time.Now()),
		"type":       "phpfpm",
		"phpfpm":     s,
	})
}
