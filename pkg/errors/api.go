package errors

/// -=-=-=-=-=-=-=-=-=] RepositoryNameNotDefinedError [-=-=-=-=-=-=-=-=-=

// RepositoryNameNotDefinedError occurs when the Application tries to get Information
// about a Repository but the Repository Name is not defined/ set.
type RepositoryNameNotDefinedError struct { message string }

func NewRepositoryNameNotDefinedError() *RepositoryNameNotDefinedError {
	return &RepositoryNameNotDefinedError{
		message: "Internal Error: Repository Name not defined",
	}
}

func NewRepositoryNameNotDefinedErrorMsg(message string) *RepositoryNameNotDefinedError {
	return &RepositoryNameNotDefinedError{
		message: message,
	}
}

func (e *RepositoryNameNotDefinedError) Error() string {
	return e.message
}