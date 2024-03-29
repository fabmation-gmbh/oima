package config


type S3Conf struct {
	// Enable or Disable the S3 Component
	Enabled			bool		`mapstructure:"enabled"`

	// S3 Endpoint of the S3 Server. For Example "play.min.io"
	Endpoint		string		`mapstructure:"endpoint"`

	// Access key is the user ID that uniquely identifies the S3 account
	AccessKeyID		string		`mapstructure:"accessKeyID"`

	// Secret key is the password to the S3 account
	SecretAccessKey	string		`mapstructure:"secretAccessKey"`

	// Use SSL to connect to the S3 Server
	UseSSL			bool		`mapstructure:"useSSL"`

	// BucketName defines the Bucket where the Signatures are stored
	BucketName		string		`mapstructure:"bucketName"`
}