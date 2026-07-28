// Command localbox is the LocalBox CLI entrypoint.
package main

import (
	"fmt"
	"os"
)

var version = "0.0.0-dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "version") {
		fmt.Println("localbox " + version)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "status" {
		line, err := runStatus()
		if err != nil {
			fmt.Fprintf(os.Stderr, "localbox: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(line)
		return
	}

	fmt.Fprintln(os.Stderr, "localbox: not yet implemented — see README.md and CLAUDE.md")
	os.Exit(1)
}
