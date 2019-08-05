package config

import ()

type Configuration struct {
	Registry RegistryConf `mapstructure:"registry"`
}
