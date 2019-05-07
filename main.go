package main

import (
	"fmt"
	"os"

	"github.com/takihito/mackerel-plugin-photon-stats/lib"
)

const version = "0.3.0"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println("version:", version)
		return
	}
	photonstats.Do()
}
