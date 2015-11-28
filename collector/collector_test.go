package collector

import (
	"testing"

	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/stretchr/testify/assert"
)

func TestPhpfpmCollector(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := []byte(`{
    "pool": "www",
    "process manager": "dynamic",
    "start since": 16527,
    "accepted conn": 941,
    "total processes": 4,
    "idle processes": 2,
    "active processes": 1,
    "max active processes": 3,
    "max children reached": 4,
    "listen queue": 0,
    "listen queue len": 128,
    "max listen queue": 13,
    "slow requests": 4
}`)
		w.Write(data)
	}))
	defer ts.Close()

	c := &PhpfpmCollector{}
	u, _ := url.Parse(ts.URL)
	s, _ := c.Collect(*u)

	assert.Equal(t, s["pool"], "www")
	assert.Equal(t, s["process_manager"], "dynamic")
	assert.Equal(t, s["start_since"], 16527)
	assert.Equal(t, s["accepted_conn"], 941)
	assert.Equal(t, s["total_processes"], 4)
	assert.Equal(t, s["idle_processes"], 2)
	assert.Equal(t, s["active_processes"], 1)
	assert.Equal(t, s["max_active_processes"], 3)
	assert.Equal(t, s["max_children_reached"], 4)
	assert.Equal(t, s["listen_queue"], 0)
	assert.Equal(t, s["listen_queue_len"], 128)
	assert.Equal(t, s["max_listen_queue"], 13)
	assert.Equal(t, s["slow_requests"], 4)
}
