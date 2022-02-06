package main

import (
	"fmt"

	"github.com/carlware/promtail-go"
	"github.com/carlware/promtail-go/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	host     = ""
	username = ""
	password = ""
	labels   = "level,my_label"
)

func main() {

	// create a promtail client
	promtail, pErr := client.NewSimpleClient(host, username, password,
		client.WithStaticLabels(map[string]interface{}{
			"app": "awesome-service",
		}),
		client.WithStreamConverter(promtail.NewRawStreamConv(labels, "=")),
	)
	if pErr != nil {
		panic(pErr)
	}

	// now can be used as a Writer for any logger

	// setup any logger and use as the output, can be combined with many outputs using io.MultiWriter
	output := zerolog.ConsoleWriter{Out: promtail}
	output.FormatMessage = func(i interface{}) string {
		_, ok := i.(string)
		if ok {
			return fmt.Sprintf("%-50s", i)
		} else {
			return ""
		}
	}
	output.FormatLevel = func(i interface{}) string {
		_, ok := i.(string)
		if ok {
			return fmt.Sprintf("level=%-7s", i)
		} else {
			return "level=info"
		}
	}
	log.Logger = log.Output(output)

	log.Info().
		Str("my_label", "some_value").
		Msg("the message")

}
