package fsutils

import (
	"log"
	"os"
	"path/filepath"
)

func FindProjectRootDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %s", err)
	}

	for {
		// Check if the ".git" folder exists in the current directory
		gitDir := filepath.Join(currentDir, ".git")

		_, err := os.Stat(gitDir)
		if err == nil {
			return currentDir
		}

		// Check if the current directory is the root directory
		if currentDir == "/" {
			log.Fatal("Could not find .git directory")
		}

		currentDir = filepath.Dir(currentDir)
	}
}
