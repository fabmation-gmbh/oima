package errors

/// -=-=-=-=-=-=-=-=-=] CredentialNameAlreadyExists [-=-=-=-=-=-=-=-=-=

// CredNameExistsError occurs when the Application tries to add a new Credential
// to the CredStore but under a already used Name.
type CredNameExistsError struct { message string }

func NewCredNameExistsError() *CredNameExistsError {
	return &CredNameExistsError{
		message: "A Credential already exists with that Name",
	}
}

func NewCredNameExistsErrorMsg(message string) *CredNameExistsError {
	return &CredNameExistsError{
		message: message,
	}
}

func (e *CredNameExistsError) Error() string {
	return e.message
}

/// -=-=-=-=-=-=-=-=-=] CredentialNotExists [-=-=-=-=-=-=-=-=-=

// CredentialNotExists occurs when the returned Enclave from memguard is nil
type CredentialNotExistsError struct { message string }

func NewCredentialNotExistsError() *CredentialNotExistsError {
	return &CredentialNotExistsError{
		message: "The specified Credential was not found in CredStore!",
	}
}

func NewCredentialNotExistsErrorMsg(message string) *CredentialNotExistsError {
	return &CredentialNotExistsError{
		message: message,
	}
}

func (e *CredentialNotExistsError) Error() string {
	return e.message
}

/// -=-=-=-=-=-=-=-=-=] EnclaveEmpty [-=-=-=-=-=-=-=-=-=

// EnclaveEmptyError occurs when the returned Enclave from memguard is nil
type EnclaveEmptyError struct { message string }

func NewEnclaveEmptyError() *EnclaveEmptyError {
	return &EnclaveEmptyError{
		message: "Memguard returned an empty Enclave",
	}
}

func NewEnclaveEmptyErrorMsg(message string) *EnclaveEmptyError {
	return &EnclaveEmptyError{
		message: message,
	}
}

func (e *EnclaveEmptyError) Error() string {
	return e.message
}