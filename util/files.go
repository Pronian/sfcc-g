package util

import (
	"os"
	"path/filepath"
	"sfcc/g/log"
)

// fileName - the name of the file to be found in the same directory as the executable.
//
// Returns the file path in the same directory as the executable
func GetFilePathInExecutableDirectory(fileName string) string {
	if filepath.IsAbs(fileName) {
		return fileName
	}

	executablePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error getting executable path: %v\n", err)
	}

	return filepath.Join(filepath.Dir(executablePath), fileName)
}
