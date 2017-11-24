package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/takihito/mackerel-plugin-photon-stats/lib"
)

const version = "0.0.4"

var (
	showVersion = flag.Bool("version", false, "show version")
)

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("version: %s \n", version)
		os.Exit(0)
	}

	photonstats.Do()
}
