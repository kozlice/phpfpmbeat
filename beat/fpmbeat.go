package beat

import (
	"net/url"
	"time"

	"github.com/elastic/libbeat/beat"
	"github.com/elastic/libbeat/cfgfile"
	"github.com/elastic/libbeat/logp"

	"github.com/kozlice/phpfpmbeat/collector"
	"github.com/kozlice/phpfpmbeat/publisher"
)

const selector = "phpfpmbeat"

type PhpfpmBeat struct {
	PfbConfig ConfigSettings
	period    time.Duration
	urls      []*url.URL
	done      chan struct{}
}

func NewPhpfpmBeat() *PhpfpmBeat {
	return &PhpfpmBeat{}
}

func (fb *PhpfpmBeat) Config(b *beat.Beat) error {
	err := cfgfile.Read(&fb.PfbConfig, "")
	if err != nil {
		logp.Err("Error reading configuration file: %v", err)
		return err
	}

	// URLs
	var urlConfig []string
	if fb.PfbConfig.Input.URLs != nil {
		urlConfig = fb.PfbConfig.Input.URLs
	} else {
		urlConfig = []string{"http://127.0.0.1/status"}
	}

	fb.urls = make([]*url.URL, len(urlConfig))
	for i := 0; i < len(urlConfig); i++ {
		u, err := url.Parse(urlConfig[i])

		if err != nil {
			logp.Err("Invalid PHP-FPM status page: %v", err)
			return err
		}

		q := u.Query()
		q.Set("json", "")
		u.RawQuery = q.Encode()

		fb.urls[i] = u
	}

	// Polling interval
	if fb.PfbConfig.Input.Period != nil {
		fb.period = time.Duration(*fb.PfbConfig.Input.Period) * time.Second
	} else {
		fb.period = 10 * time.Second
	}

	logp.Debug(selector, "Watch %v", fb.urls)
	logp.Debug(selector, "Period %v", fb.period)

	return nil
}

func (fb *PhpfpmBeat) Setup(b *beat.Beat) error {
	fb.done = make(chan struct{})
	return nil
}

func (fb *PhpfpmBeat) Run(b *beat.Beat) error {
	logp.Debug(selector, "Run phpfpmbeat")

	var err error

	ticker := time.NewTicker(fb.period)
	defer ticker.Stop()

	c := collector.NewPhpfpmCollector()
	p := publisher.NewFpmPublisher(b.Events)

	// TODO: Different scheme
	for {
		select {
		case <-fb.done:
			return nil
		case <-ticker.C:
		}

		timerStart := time.Now()

		for _, u := range fb.urls {
			s, err := c.Collect(*u)
			if err != nil {
				logp.Err("Failed to read PHP-FPM status: %v", err)
			} else {
				p.Publish(s)
			}
		}

		timerEnd := time.Now()
		duration := timerEnd.Sub(timerStart)
		if duration.Nanoseconds() > fb.period.Nanoseconds() {
			logp.Warn("Ignoring tick(s) due to processing taking longer than one period")
		}
	}

	return err
}

func (fb *PhpfpmBeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (fb *PhpfpmBeat) Stop() {
	logp.Debug(selector, "Stop phpfpmbeat")
	close(fb.done)
}
