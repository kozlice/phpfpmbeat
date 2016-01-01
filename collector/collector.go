package collector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type PhpfpmStatus struct {
	Pool               string `json:"pool"`
	ProcessManager     string `json:"process manager"`
	StartSince         int    `json:"start since"`
	AcceptedConn       int    `json:"accepted conn"`
	TotalProcesses     int    `json:"total processes"`
	IdleProcesses      int    `json:"idle processes"`
	ActiveProcesses    int    `json:"active processes"`
	MaxActiveProcesses int    `json:"max active processes"`
	MaxChildrenReached int    `json:"max children reached"`
	ListenQueue        int    `json:"listen queue"`
	ListenQueueLen     int    `json:"listen queue len"`
	MaxListenQueue     int    `json:"max listen queue"`
	SlowRequests       int    `json:"slow requests"`
}

type PhpfpmCollector struct{}

func NewPhpfpmCollector() *PhpfpmCollector {
	return &PhpfpmCollector{}
}

func (c *PhpfpmCollector) Collect(u url.URL) (map[string]interface{}, error) {
	var (
		err error
		s PhpfpmStatus
		v map[string]interface{}
	)

	res, err := http.Get(u.String())
	if err != nil {
		return v, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return v, fmt.Errorf("HTTP%s", res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(&s)
	if err != nil {
		return v, err
	}

	v = map[string]interface{}{
		"pool":                 s.Pool,
		"process_manager":      s.ProcessManager,
		"start_since":          s.StartSince,
		"accepted_conn":        s.AcceptedConn,
		"total_processes":      s.TotalProcesses,
		"idle_processes":       s.IdleProcesses,
		"active_processes":     s.ActiveProcesses,
		"max_active_processes": s.MaxActiveProcesses,
		"max_children_reached": s.MaxChildrenReached,
		"listen_queue":         s.ListenQueue,
		"listen_queue_len":     s.ListenQueueLen,
		"max_listen_queue":     s.MaxListenQueue,
		"slow_requests":        s.SlowRequests,
	}
	return v, nil
}
