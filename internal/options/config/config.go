package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Sources struct {
	NixOSURL       string `mapstructure:"NIXOS_URL"`
	HomeManagerURL string `mapstructure:"HOME_MANAGER_URL"`
}

type Config struct {
	Sources Sources `mapstructure:"sources"`
	Retries int     `mapstructure:"retries"`
	Timeout int     `mapstructure:"timeout"`
}

func LoadConfig() (*Config, error) {
	config := &Config{}

	viper.SetDefault("sources.nixos_url", "https://nixos.org/manual/nixos/unstable/options")
	viper.SetDefault("sources.home_manager_url", "https://nix-community.github.io/home-manager/options.xhtml")
	defaultRetries := 3
	defaultTimeout := 30
	viper.SetDefault("retries", defaultRetries)
	viper.SetDefault("timeout", defaultTimeout)
	viper.SetConfigName("optinix)")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, err
		}
	}

	err := viper.Unmarshal(config)
	if err != nil {
		return config, fmt.Errorf("unable to decode into config struct, %v", err)
	}

	return config, nil
}
