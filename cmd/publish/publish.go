package publish

import (
	"log"
	"rlrabinowitz.github.io/cmd/publish/module"
	"rlrabinowitz.github.io/cmd/publish/provider"
	"strings"
)

func Publish(args []string) {
	log.Printf("Starting")
	filePaths := getFilePathsToMigrate(args)
	for _, filePath := range filePaths {
		// TODO Range variables
		// TODO Parallelism
		if isProviderPath(filePath) {
			err := provider.Publish(filePath)
			if err != nil {
				panic(err)
			}
		} else {
			err := module.Publish(filePath)
			if err != nil {
				panic(err)
			}
		} // TODO Validate path is either provider or module, that amount of parts make sense
	}
}

func getFilePathsToMigrate(args []string) []string {
	if len(args) != 1 {
		panic("The fuck, missing arguments") // TODO language (and panic)
	}
	return strings.Split(args[0], ",")
}

func isProviderPath(filePath string) bool {
	pathParts := strings.Split(filePath, "/")
	return pathParts[0] == "providers" // TODO constant?
}
