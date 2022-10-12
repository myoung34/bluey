package cli

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/myoung34/bluey/bluey"
	"os"
)

func ParseArgs() blue.Config {
	parser := argparse.NewParser("bluey", "A pluggable system to receive and transmit bluetooth events as a proxy")

	configFile := parser.String("c", "config", &argparse.Options{
		Required: true,
		Help:     "Configuration file location",
	})

	if err := parser.Parse(os.Args); err != nil {
		fmt.Println(parser.Usage(err))
		os.Exit(1)
	}

	return blue.ParseConfig(*configFile)
}
