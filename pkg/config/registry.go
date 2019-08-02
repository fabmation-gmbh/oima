package config

import ()

type RegistryConf struct {
	RegistryURI	string	`mapstructure:"uri"`
	RequireAuth	bool	`mapstructure:"require_auth"`
	Username	string	`mapstructure:"username"`
	Password	[]byte	`mapstructure:"password"`
}