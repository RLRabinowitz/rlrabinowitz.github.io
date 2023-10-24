package files

import (
	"os"
	"path"
)

func WriteToFile(filePath string, data []byte) error {
	err := os.MkdirAll(path.Dir(filePath), 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
