package config

import (
	"fmt"
	"os"
	"strings"

	gap "github.com/muesli/go-app-paths"
	"github.com/spf13/viper"
)

type Sources struct {
	NixOS       string `mapstructure:"nixos_path"`
	HomeManager string `mapstructure:"home_manager_path"`
	Darwin      string `mapstructure:"darwin_path"`
}

type Config struct {
	DBFolder string  `mapstructure:"db_folder"`
	Sources  Sources `mapstructure:"sources"`
}

func LoadConfig() (*Config, error) {
	config := &Config{}

	viper.SetEnvPrefix("optinix")
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	viper.AutomaticEnv()

	scope := gap.NewScope(gap.User, "optinix")
	dirs, err := scope.DataDirs()
	if err != nil {
		return config, fmt.Errorf("unable to get data directory: %v", err)
	}

	setDefaults(dirs)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, err
		}
	}

	err = viper.Unmarshal(config)
	if err != nil {
		return config, fmt.Errorf("unable to decode into config struct, %v", err)
	}

	return config, nil
}

func setDefaults(dirs []string) {
	dbFolder, _ := os.UserHomeDir()
	if len(dirs) > 0 {
		dbFolder = dirs[0]
	}

	viper.SetDefault("db_folder", dbFolder)
	viper.SetDefault("sources.nixos_path", "nix/nixos-options.nix")
	viper.SetDefault("sources.home_manager_path", "nix/hm-options.nix")
	viper.SetDefault("sources.darwin_path", "nix/darwin-options.nix")
}
