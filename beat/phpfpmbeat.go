package beat

import (
	"net/url"
	"time"

	"github.com/elastic/libbeat/beat"
	"github.com/elastic/libbeat/cfgfile"
	"github.com/elastic/libbeat/logp"

    "github.com/kozlice/phpfpmbeat/beat/publisher"
	"github.com/kozlice/phpfpmbeat/beat/collector"
)

const selector = "phpfpmbeat"

type Phpfpmbeat struct {
	PfbConfig ConfigSettings

	period    time.Duration
	urls      []*url.URL

	done      chan struct{}
}

func New() *Phpfpmbeat {
	return &Phpfpmbeat{}
}

func (pfb *Phpfpmbeat) Config(b *beat.Beat) error {
	err := cfgfile.Read(&pfb.PfbConfig, "")
	if err != nil {
		logp.Err("Error reading configuration file: %v", err)
		return err
	}

	// URLs
	var urlConfig []string
	if pfb.PfbConfig.Input.URLs != nil {
		urlConfig = pfb.PfbConfig.Input.URLs
	} else {
		urlConfig = []string{"http://127.0.0.1/status"}
	}

	pfb.urls = make([]*url.URL, len(urlConfig))
	for i := 0; i < len(urlConfig); i++ {
		u, err := url.Parse(urlConfig[i])

		if err != nil {
			logp.Err("Invalid PHP-FPM status page: %v", err)
			return err
		}

		q := u.Query()
		q.Set("json", "")
		u.RawQuery = q.Encode()

		pfb.urls[i] = u
	}

	// Polling interval
	if pfb.PfbConfig.Input.Period != nil {
		pfb.period = time.Duration(*pfb.PfbConfig.Input.Period) * time.Second
	} else {
		pfb.period = 10 * time.Second
	}

	logp.Debug(selector, "Watch %v", pfb.urls)
	logp.Debug(selector, "Period %v", pfb.period)

	return nil
}

func (pfb *Phpfpmbeat) Setup(b *beat.Beat) error {
	// pfb.events = b.Events
	pfb.done = make(chan struct{})
	return nil
}

func (pfb *Phpfpmbeat) Run(b *beat.Beat) error {
	logp.Debug(selector, "Run phpfpmbeat")

	var err error

	ticker := time.NewTicker(pfb.period)
	defer ticker.Stop()

    c := collector.New()
    p := publisher.New(b.Events)

	for {
		select {
		case <-pfb.done:
			return nil
		case <-ticker.C:
		}

		timerStart := time.Now()

		for _, u := range pfb.urls {
			s, err := c.Collect(*u)
			if (err != nil) {
				logp.Err("Failed to read PHP-FPM status: %v", err)
			} else {
                p.Publish(s)
            }
		}

		timerEnd := time.Now()
		duration := timerEnd.Sub(timerStart)
		if duration.Nanoseconds() > pfb.period.Nanoseconds() {
			logp.Warn("Ignoring tick(s) due to processing taking longer than one period")
		}
	}

	return err
}

func (pfb *Phpfpmbeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (pfb *Phpfpmbeat) Stop() {
	logp.Debug(selector, "Stop phpfpmbeat")
	close(pfb.done)
}
