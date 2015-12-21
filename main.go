package main

import (
	phpfpmbeat "github.com/kozlice/phpfpmbeat/beat"

	"github.com/elastic/beats/libbeat/beat"
)

var Version = "1.0.0-beta2"
var Name = "phpfpmbeat"

func main() {
	beat.Run(Name, Version, phpfpmbeat.NewPhpfpmBeat())
}
