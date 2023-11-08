package main

import (
	"fmt"
	"os"
	"rlrabinowitz.github.io/cmd/initialize"
	"rlrabinowitz.github.io/cmd/publish"
	"rlrabinowitz.github.io/cmd/update"
)

const (
	Publish            string = "publish"
	Initialize                = "initialize"
	Update                    = "update"
	UpdateExperimental        = "update-experimental"
)

// TODO use 3rd-party for commands
func getCommandAndArguments() (string, []string) {
	if len(os.Args) < 2 {
		panic("Received wrong number of arguments") // TODO panic
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
		update.Update(args, false)
	} else if command == UpdateExperimental {
		update.Update(args, true)
	} else {
		panic(fmt.Errorf("unexpected command: %s. Please run one of the following commands: publish, initialize, update", command))
	}
}

// TODO Do not hardcode the "dist" folder
