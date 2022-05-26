package main

import (
	"flag"

	config "github.com/rhmdnd/compserv/pkg/config"
)

func main() {
	var configDir = flag.String("config-dir", "configs/", "Path to YAML configuration directory containing a config.yaml file.")
	flag.Parse()
	config.ParseConfig(*configDir)
}
