package internal

import (
	"github.com/minio/minio-go"
)

// TODO: think about an Concept to integrate it with our Interfaces
// TODO: to provide a possibility to extend any other S3 Provider easily
// minioClient will be used to check Signatures from a S3 Server
var S3Client *minio.Client