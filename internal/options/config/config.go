package config

import (
	"fmt"
	"os"
	"path/filepath"

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

	setDefaults()
	err := viper.BindEnv("sources.nixos_url", "NIXOS_URL")
	if err != nil {
		return config, err
	}

	err = viper.BindEnv("sources.home_manager_url", "HOME_MANAGER_URL")
	if err != nil {
		return config, err
	}

	err = viper.BindEnv("db_folder", "DB_FOLDER)")
	if err != nil {
		return config, err
	}
	viper.SetConfigType("env")
	viper.AutomaticEnv()

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

func setDefaults() {
	// TODO: what if it not set
	state := os.Getenv("XDG_STATE_HOME")
	configPath := filepath.Join(state, "optinix")
	viper.SetDefault("db_path", configPath)

	viper.SetDefault("sources.nixos_url", "https://nixos.org/manual/nixos/unstable/options")
	viper.SetDefault("sources.nixos_url", "https://nixos.org/manual/nixos/unstable/options")
	viper.SetDefault("sources.home_manager_url", "https://nix-community.github.io/home-manager/options.xhtml")

	defaultRetries := 3
	defaultTimeout := 30
	viper.SetDefault("retries", defaultRetries)
	viper.SetDefault("timeout", defaultTimeout)
}
