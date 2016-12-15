package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/kozlice/phpfpmbeat/beater"
)

func main() {
	err := beat.Run("phpfpmbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
