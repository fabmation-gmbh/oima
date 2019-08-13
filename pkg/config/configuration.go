package config

type Configuration struct {
	// (Docker) Registry Configuration
	Registry	RegistryConf `mapstructure:"registry"`

	// S3-Server Configuration
	S3			S3Conf		 `mapstructure:"s3"`
}
