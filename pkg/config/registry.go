package config

import ()

type RegistryConf struct {
	RegistryURI	string	`mapstructure:"uri"`
	RequireAuth	string	`mapstructure:"require_auth"`
	Username	string	`mapstructure:"username"`
	Password	string	`mapstructure:"password"`
}