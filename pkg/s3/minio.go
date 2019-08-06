package s3

import (
	"github.com/awnumar/memguard"
	"github.com/fabmation-gmbh/oima/internal"
	. "github.com/fabmation-gmbh/oima/internal/log"
	"github.com/fabmation-gmbh/oima/pkg/config"
	"github.com/fabmation-gmbh/oima/pkg/errors"
	"github.com/fabmation-gmbh/oima/pkg/registry"
	"github.com/minio/minio-go"
)

var conf config.Configuration


func (s *S3Minio) InitS3() error {
	// initialize Config
	conf = internal.GetConfig()

	// check if S3 is enabled
	if !conf.S3.Enabled {
		Log.Fatalf("[Internal Error] Calling initS3() but S3 Component is disabled!")
		memguard.SafeExit(1)
	}

	// TODO: implement provider

	// initialize Auth
	s.Auth = &S3AuthMinio{}
	err := s.Auth.InitAuth()
	if err != nil {
		Log.Fatalf("Error while initializing MinIO S3 Authentication: %s", err.Error())
		return err
	}

	return nil
}

func (s *S3Minio) SignatureExists(image *registry.Image) (bool, error) {


	return true, nil
}


func (auth *S3AuthMinio) InitAuth() error {

	// check if S3 is enabled
	if !conf.S3.Enabled {
		Log.Fatalf("[Internal Error] Calling initS3() but S3 Component is disabled!")
		memguard.SafeExit(1)
	}

	_accessKeyID, err := internal.Cred.GetCredential("s3_accessKeyID")
	if err != nil {
		if _, ok := err.(*errors.CredentialNotExistsError); ok {
			Log.Fatalf("Demanded Credential (key: %s) does not exists in CredStore: %s",
						"s3_accessKeyID", err.Error())
			memguard.SafeExit(1)
		}

		Log.Criticalf("Error while getting encrypted Credential")
		memguard.SafeExit(1)
	}
	auth.accessKeyID = _accessKeyID

	_secretAccessKeyID, err := internal.Cred.GetCredential("s3_secretAccessKeyID")
	if err != nil {
		if _, ok := err.(*errors.CredentialNotExistsError); ok {
			Log.Fatalf("Demanded Credential (key: %s) does not exists in CredStore: %s",
				"s3_secretAccessKeyID", err.Error())
			memguard.SafeExit(1)
		}

		Log.Criticalf("Error while getting encrypted Credential")
		memguard.SafeExit(1)
	}
	auth.secretAccessKeyID = _secretAccessKeyID

	// open credentials
	accessKeyID, err := auth.accessKeyID.Open()
	if err != nil { memguard.SafePanic(err) }
	defer accessKeyID.Destroy()

	secretAccessKeyID, err := auth.secretAccessKeyID.Open()
	if err != nil { memguard.SafePanic(err) }
	defer secretAccessKeyID.Destroy()


	// initialize objects of Struct
	auth.Endpoint = conf.S3.Endpoint
	auth.UseSSL = conf.S3.UseSSL
	auth.BucketName = conf.S3.BucketName

	// initialize Minio Client object
	internal.S3Client, err = minio.New(auth.Endpoint, conf.S3.AccessKeyID, conf.S3.SecretAccessKey, auth.UseSSL)
	if err != nil {
		Log.Fatalf("Error while initializing MinIO Client: %s", err.Error())
		memguard.SafeExit(1)
	}

	Log.Debugf("MinIO S3 Client initialization finished")

	return nil
}

func (auth *S3AuthMinio) GetAccessKeyID() *memguard.Enclave { return auth.accessKeyID }

func (auth *S3AuthMinio) GetSecretAccessKeyID() *memguard.Enclave { return auth.secretAccessKeyID }

func (auth *S3AuthMinio) GetEndpoint() string { return auth.Endpoint }

func (auth *S3AuthMinio) TLSEnabled() bool { return auth.UseSSL }