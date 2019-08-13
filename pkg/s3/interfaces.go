package s3

import (
	"github.com/awnumar/memguard"
	rt "github.com/fabmation-gmbh/oima/pkg/registry/interfaces"
)

type S3 interface {
	// Initializes the S3 Datatypes
	InitS3() error

	// Fetches Signatures of all Tags of an Image and sets the Bool
	// in image.Tags[n].S3SignFound = "Signature Found?"
	FetchSignatures(image *rt.BaseImage)		error

	// Deletes a Signature of the specific Tag from the S3-Server
	// @ctxPath	describes the Full converted Path of registry+image+tag_info
	//			do **NOT** add any annotations like '/signature-1'!!
	DeleteSignature(ctxPath string, tag *rt.Tag)
}


//noinspection GoNameStartsWithPackageName
type S3Auth interface {
	// Initializes the required Datatypes for Authentication
	InitAuth()					error


	/// >>>>> AccessKeyID & SecretAccessKeyID <<<<<

	// GetAccessKeyID returns the Encrypted AccessKeyID
	GetAccessKeyID()			*memguard.Enclave

	// GetSecretAccessKeyID returns the Encrypted AccessKeyID
	GetSecretAccessKeyID()		*memguard.Enclave

	/// >>>>> Endpoint <<<<<

	// Returns the Endpoint of the S3 Server
	GetEndpoint()	string

	// Returns true if user set UseSSL for S3 to true
	TLSEnabled()	bool
}
