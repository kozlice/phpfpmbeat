package collector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Status struct {
	Pool               string `json:"pool"`
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

func New() *PhpfpmCollector {
	return &PhpfpmCollector{}
}

func (c *PhpfpmCollector) Collect(u url.URL) (map[string]interface{}, error) {
	var err error
	var st Status
	var v map[string]interface{}

	res, err := http.Get(u.String())
	if err != nil {
		return v, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return v, fmt.Errorf("HTTP%s", res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(&st)
	if err != nil {
		return v, err
	}

	v = map[string]interface{}{
		"pool":                 st.Pool,
		"start_since":          st.StartSince,
		"accepted_conn":        st.AcceptedConn,
		"total_processes":      st.TotalProcesses,
		"idle_processes":       st.IdleProcesses,
		"active_processes":     st.ActiveProcesses,
		"max_active_processes": st.MaxActiveProcesses,
		"max_children_reached": st.MaxChildrenReached,
		"listen_queue":         st.ListenQueue,
		"listen_queue_len":     st.ListenQueueLen,
		"max_listen_queue":     st.MaxListenQueue,
		"slow_requests":        st.SlowRequests,
	}
	return v, nil
}
