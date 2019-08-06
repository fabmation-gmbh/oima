package s3

import (
	"github.com/awnumar/memguard"
	"github.com/fabmation-gmbh/oima/pkg/registry"
)

type S3 interface {
	// Initializes the S3 Datatypes
	InitS3() error

	// Checks if a Signature for the Images exists on the S3 Server
	SignatureExists(image *registry.Image)		(bool, error)
}


//noinspection GoNameStartsWithPackageName
type S3Auth interface {
	// Initializes the required Datatypes for Authentication
	InitAuth()								error


	/// >>>>> AccessKeyID & SecretAccessKeyID <<<<<

	// GetAccessKeyID returns the Encrypted AccessKeyID
	GetAccessKeyID()								*memguard.Enclave

	// GetSecretAccessKeyID returns the Encrypted AccessKeyID
	GetSecretAccessKeyID()							*memguard.Enclave

	/// >>>>> Endpoint <<<<<

	// Returns the Endpoint of the S3 Server
	GetEndpoint()	string

	// Returns true if user set UseSSL for S3 to true
	TLSEnabled()	bool
}