package internal

import(
	"github.com/minio/minio-go/v6"
)

// minioClient will be used to check Signatures from a S3 Server
var minioClient *minio.Client