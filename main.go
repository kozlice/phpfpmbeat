package main

import (
	phpfpmbeat "github.com/kozlice/phpfpmbeat/beat"

	"github.com/elastic/libbeat/beat"
)

var Version = "1.0.0-beta1"
var Name = "phpfpmbeat"

func main() {
	beat.Run(Name, Version, phpfpmbeat.New())
}
