package optionstest

import (
	"os"
)

func GetHost(path string) string {
	fullPath := "http://localhost:8080" + path
	if os.Getenv("CI") == "true" {
		fullPath = "http://docker:8080" + path
	}

	return fullPath
}
