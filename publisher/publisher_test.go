package publisher

import (
	"testing"

	"encoding/json"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/publisher"
	"github.com/stretchr/testify/assert"
)

type PhpfpmTestClient struct {
	Channel chan common.MapStr
}

func (c PhpfpmTestClient) PublishEvent(event common.MapStr, opts ...publisher.ClientOption) bool {
	c.Channel <- event
	return true
}

func (c PhpfpmTestClient) PublishEvents(events []common.MapStr, opts ...publisher.ClientOption) bool {
	for _, event := range events {
		c.Channel <- event
	}
	return true
}

func TestPhpfpmPublisher(t *testing.T) {
	c := make(chan common.MapStr, 16)
	p := NewFpmPublisher(&PhpfpmTestClient{Channel: c})

	s := map[string]interface{}{}

	p.Publish(s)
	assert.Equal(t, 1, len(c))

	se := <-c
	var sm map[string]interface{}
	if err := json.Unmarshal([]byte(se.String()), &sm); assert.NoError(t, err) {
		assert.Equal(t, "phpfpm", sm["type"])
	}
}
