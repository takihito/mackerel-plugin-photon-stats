package main

import (
	"fmt"
	"os"

	photonstats "github.com/takihito/mackerel-plugin-photon-stats/lib"
)

const version = "0.5.0"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println("version:", version)
		return
	}
	photonstats.Do()
}
