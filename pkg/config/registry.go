package config

import ()

type RegistryConf struct {
	RegistryURI	string	`yaml:"uri"`
	RequireAuth	bool	`yaml:"require_auth"`
	Username	string	`yaml:"username"`
	Password	string	`yaml:"password"`
}