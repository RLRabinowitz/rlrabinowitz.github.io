package update

import (
	"log"
	"rlrabinowitz.github.io/cmd/initialize/module"
	"rlrabinowitz.github.io/cmd/update/provider"
	"strings"
	"sync"
)

type ExecutionResult struct {
	filePath string
	Err      error
}

func Update(args []string, experimental bool) {
	log.Printf("Starting")
	filePaths := getFilePathsToMigrate(args)

	executionCh := make(chan ExecutionResult, len(filePaths))

	var wg sync.WaitGroup
	for _, filePath := range filePaths {
		wg.Add(1)

		go func(filePath string) {
			defer wg.Done()
			var err error
			if isProviderPath(filePath) {
				err = provider.Update(filePath, experimental)
			} else {
				// Right now there's no difference between the update and initialize action of modules
				// as we just go over all tags and create the file from there
				err = module.Initialize(filePath)
			} // TODO Validate path is either provider or module, that amount of parts make sense

			executionCh <- ExecutionResult{
				filePath: filePath,
				Err:      err,
			}

		}(filePath)
	}

	wg.Wait()
	close(executionCh)

	for e := range executionCh {
		if e.Err != nil {
			log.Printf("Failed to update %s: %s", e.filePath, e.Err)
		}
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
