package main

import (
	"log"

	run "github.com/ayupov-ayaz/anti-brute-force/cli/cmd"
)

func main() {
	if err := run.Run(); err != nil {
		log.Fatal(err)
	}
}
