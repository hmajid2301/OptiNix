package config

import (
	"fmt"
	"os"

	gap "github.com/muesli/go-app-paths"
)

type Config struct {
	DBFolder string `mapstructure:"db_folder"`
}

func LoadConfig() (Config, error) {
	config := Config{}

	scope := gap.NewScope(gap.User, "optinix")
	dirs, err := scope.DataDirs()
	if err != nil {
		return config, fmt.Errorf("unable to get data directory: %v", err)
	}

	dbFolder, _ := os.UserHomeDir()
	if len(dirs) > 0 {
		dbFolder = dirs[0]
	}

	if folder, ok := os.LookupEnv("OPTINIX_DB_FOLDER"); ok {
		dbFolder = folder
	}

	config.DBFolder = dbFolder
	return config, nil
}
