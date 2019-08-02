package credential

import (
	"github.com/awnumar/memguard"
	"github.com/fabmation-gmbh/oima/pkg/errors"
)

type Store interface {
	AddCredential(name string, data *[]byte)	error						// Adds a Credential to the CredStore and securely wipes data
	GetCredential(name string)					(*memguard.Enclave, error)	// Returns the Enclave Credential that contains the Credential
	RemoveCredential(name string)				error						// Wipes/ Destroys the Enclave that contains the Credential
}

type CredStore struct {
	credentials		map[string]*memguard.Enclave					// Contains all the Encrypted Credentials
}

func (cred *CredStore) AddCredential(name string, data []byte) error {
	// check if credentials is Initialized
	if cred.credentials == nil {
		cred.credentials = make(map[string]*memguard.Enclave)
	}

	// check if Credential already exist
	if _, err := cred.credentials[name]; err {
		return errors.NewCredNameExistsError()
	}

	// securely move Data to Enclave
	dataBuf := memguard.NewBufferFromBytes(data)
	memguard.ScrambleBytes(data)

	encryptedData := dataBuf.Seal()
	if encryptedData == nil {
		memguard.SafePanic(errors.NewEnclaveEmptyError())
	}

	cred.credentials[name] = encryptedData
	return nil
}

func (cred *CredStore) GetCredential(name string) (*memguard.Enclave, error) {
	// check if Credential exists in CredStore
	if _, err := cred.credentials[name]; !err {
		return nil, errors.NewCredentialNotExistsError()
	}

	// return Enclave
	return cred.credentials[name], nil
}

func (cred *CredStore) RemoveCredential(name string) error {
	// check if Credential exists in CredStore
	if _, err := cred.credentials[name]; !err {
		return errors.NewCredentialNotExistsError()
	}

	delete(cred.credentials, name)
	return nil
}