package config

import (
	"fmt"
	"os"
	"strings"

	gap "github.com/muesli/go-app-paths"
	"github.com/spf13/viper"
)

type Sources struct {
	NixOSURL       string `mapstructure:"nixos_url"`
	HomeManagerURL string `mapstructure:"home_manager_url"`
}

type Config struct {
	DBFolder string  `mapstructure:"db_folder"`
	Sources  Sources `mapstructure:"sources"`
	Retries  int     `mapstructure:"retries"`
	Timeout  int     `mapstructure:"timeout"`
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

	viper.SetDefault("sources.nixos_url", "https://nixos.org/manual/nixos/unstable/options")
	viper.SetDefault("sources.nixos_url", "https://nixos.org/manual/nixos/unstable/options")
	viper.SetDefault("sources.home_manager_url", "https://nix-community.github.io/home-manager/options.xhtml")

	defaultRetries := 3
	defaultTimeout := 30
	viper.SetDefault("retries", defaultRetries)
	viper.SetDefault("timeout", defaultTimeout)
}
