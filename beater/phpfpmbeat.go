package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"encoding/json"
	"github.com/kozlice/phpfpmbeat/config"
	"net/http"
	"net/url"
)

type Phpfpmbeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
}

func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Phpfpmbeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func (bt *Phpfpmbeat) Run(b *beat.Beat) error {
	logp.Info("phpfpmbeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		for _, u := range bt.config.URLs {
			event, err := bt.collect(b, u)
			if err != nil {
				logp.Err("An error occured: %v", err)
			}
			bt.client.PublishEvent(event)
			logp.Info("Event sent")
		}
	}
}

func (bt *Phpfpmbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

type Phpfpmstatus struct {
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

func (bt *Phpfpmbeat) collect(b *beat.Beat, u *url.URL) (map[string]interface{}, error) {
	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP%s", res.Status)
	}

	s := Phpfpmstatus{}
	err = json.NewDecoder(res.Body).Decode(&s)
	if err != nil {
		return nil, err
	}

	result := common.MapStr{
		"@timestamp":           common.Time(time.Now()),
		"type":                 b.Name,
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
	return result, nil
}
