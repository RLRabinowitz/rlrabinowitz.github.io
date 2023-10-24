package main

import (
	"fmt"
	"os"
	"rlrabinowitz.github.io/cmd/initialize"
	"rlrabinowitz.github.io/cmd/publish"
)

const (
	Publish    string = "publish"
	Initialize        = "initialize"
	Update            = "update"
)

// TODO use 3rd-party for commands
func getCommandAndArguments() (string, []string) {
	if len(os.Args) < 2 {
		panic("The fuck, missing arguments") // TODO language (and panic)
	}

	return os.Args[1], os.Args[2:]
}

func main() {
	command, args := getCommandAndArguments()
	if command == Publish {
		publish.Publish(args)
	} else if command == Initialize {
		initialize.Run(args)
	} else if command == Update {

	} else {
		panic(fmt.Errorf("unexpected command: %s. Please run one of the following commands: publish, initialize, update", command))
	}
}

// TODO Do not hardcode the "dist" folder
