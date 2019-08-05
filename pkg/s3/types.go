package s3

import (
	"github.com/awnumar/memguard"
)

//noinspection GoNameStartsWithPackageName
type S3Minio struct {
	Auth		*S3AuthMinio
}


//noinspection GoNameStartsWithPackageName
type S3AuthMinio struct {
	Endpoint			string      // S3 Server Endpoint
	BucketName			string		// S3 BucketName
	UseSSL				bool        // Communicate with S3 Server over SSL


	accessKeyID       *memguard.Enclave // Contains the Access Key ID securely
	secretAccessKeyID *memguard.Enclave // Contains the Secret Access Key ID securely
}