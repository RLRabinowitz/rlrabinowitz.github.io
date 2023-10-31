package update

import (
	"log"
	"rlrabinowitz.github.io/cmd/initialize/module"
	"rlrabinowitz.github.io/cmd/update/provider"
	"strings"
)

func Update(args []string) {
	log.Printf("Starting")
	filePaths := getFilePathsToMigrate(args)
	for _, filePath := range filePaths {
		// TODO Range variables
		// TODO Parallelism
		if isProviderPath(filePath) {
			err := provider.Update(filePath)
			if err != nil {
				panic(err)
			}
		} else {
			// Right now there's no difference between the update and initialize action of modules
			// as we just go over all tags and create the file from there
			err := module.Initialize(filePath)
			if err != nil {
				panic(err)
			}
		} // TODO Validate path is either provider or module, that amount of parts make sense
	}
}

func getFilePathsToMigrate(args []string) []string {
	if len(args) != 1 {
		panic("Received wrong number of arguments") // TODO panic
	}
	return strings.Split(args[0], ",")
}

func isProviderPath(filePath string) bool {
	pathParts := strings.Split(filePath, "/")
	return pathParts[0] == "providers" // TODO constant?
}
