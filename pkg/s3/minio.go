package s3

import (
	"fmt"
	"github.com/awnumar/memguard"
	"github.com/minio/minio-go"
	"strings"

	"github.com/fabmation-gmbh/oima/internal"
	. "github.com/fabmation-gmbh/oima/internal/log"
	"github.com/fabmation-gmbh/oima/pkg/config"
	"github.com/fabmation-gmbh/oima/pkg/errors"
	rt "github.com/fabmation-gmbh/oima/pkg/registry/interfaces"
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

// SignatureExists() checks if all Signatures of all
func (s *S3Minio) FetchSignatures(image *rt.BaseImage) error {
	Log.Debugf("Fetching Signatures for Image %s", (*image).GetName())

	//// check if @image is Empty
	//if len((*image).GetTags()) == 0 {
	//	Log.Fatal("Requested Image is Empty (empty Struct)!")
	//	memguard.SafeExit(1)
	//}
	tags := (*image).GetTags()

	// create object Path Prefix
	registryName := PrepareRegPath((*image).GetRegistryURI())
	objPathPrefix := fmt.Sprintf("%s/%s@", registryName, (*image).GetName())

	doneCh := make(chan struct{})
	defer close(doneCh)

	// check Tag Objects at the S3 MinIO Server
	for iTag, _ := range tags {
		objName := fmt.Sprintf("%s%s/signature-1", objPathPrefix, strings.ReplaceAll(tags[iTag].ContentDigest, ":", "="))

		_, err := internal.S3Client.StatObject(s.Auth.BucketName, objName, minio.StatObjectOptions{})
		if err != nil {
			errResponse := minio.ToErrorResponse(err)
			var errMsg string
			var signNotFound = false		// if true than it does not print a Message or Exit the Application

			if errResponse.Code == "AccessDenied" {
				errMsg = fmt.Sprintf("S3 Server returned %s. You have not the Permissions to stat the File!",
										errResponse.Code)
			} else if errResponse.Code == "NoSuchBucket" {
				errMsg = fmt.Sprintf("S3 Server returned %s. A Bucket with the Name '%s' wasn't found!",
					errResponse.Code, s.Auth.BucketName)
			} else if errResponse.Code == "InvalidBucketName" {
				errMsg = fmt.Sprintf("S3 Server returned %s. The Bucket Name (%s) contains invalid chars!",
					errResponse.Code, s.Auth.BucketName)
			} else if errResponse.Code == "NoSuchKey" {
				// Signature File does not exists
				signNotFound = true
				tags[iTag].S3SignFound = false
			} else {
				errMsg = fmt.Sprintf("Unknown Error while getting Object for Image '%s@%s': %s",
										(*image).GetName(), tags[iTag].ContentDigest, err.Error())
			}

			if !signNotFound {
				Log.Critical(errMsg)
				memguard.SafeExit(1)
			}
		} else { tags[iTag].S3SignFound = true }
		//Log.Debugf("++>> Object Path: %s", objName)
	}

	// update Tags from Image
	(*image).SetTags(tags)

	return nil
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