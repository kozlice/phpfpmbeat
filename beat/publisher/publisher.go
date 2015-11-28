package publisher

import (
	"time"

	"github.com/elastic/libbeat/common"
	"github.com/elastic/libbeat/publisher"
)

// StubPublisher is a Publisher that publishes Nginx Stub status.
type PhpfpmPublisher struct {
	client publisher.Client
}

// NewPhpfpmPublisher constructs a new PhpfpmPublisher.
func New(c publisher.Client) *PhpfpmPublisher {
	return &PhpfpmPublisher{client: c}
}

// Publish Phpfpm status.
func (pfb *PhpfpmPublisher) Publish(s map[string]interface{}) {
	pfb.client.PublishEvent(common.MapStr{
		"@timestamp": common.Time(time.Now()),
		"type":       "phpfpm",
		"phpfpm":     s,
	})
}
