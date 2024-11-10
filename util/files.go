package util

import (
	"fmt"
	"os"
	"path/filepath"
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
		fmt.Errorf("Error getting executable path:", err)
		os.Exit(1)
	}

	return filepath.Join(filepath.Dir(executablePath), fileName)
}
