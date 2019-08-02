// Registry Package adds the Possibility to talk with the (Docker) Registry API
package registry

import (
	"time"

	"github.com/awnumar/memguard"
)

type _Tag string             // _Tag of an Image (for example 'v1.0.0' or '0.1.0-beta'
type _RegistryVersion string // Describes the current Version of the Registry API
const (
	V1	_RegistryVersion	= "v1"		// (Docker) Registry API Version v1
	V2	_RegistryVersion	= "v2"		// (Docker) Registry API Version v1
)

// Registry Authentication Information
type Auth struct {
	Required		bool				// Is Authentication Required
	Cred			Credential			// Needed Credentials
}

// Holds/ Checks and gets needed Credentials/ Informations
// to communicate with the Registry API
type credential interface {
	Init(cred *Credential)	error		// Checks and "Initializes" the Credential Struct
	check(cred *Credential)	error		// Check if the Credentials works
}

// The BearerToken is needed to communicate with the Registry API
type BearerToken struct {
	BearerToken		memguard.Enclave	// The BearerToken is needed to communicate with the Registry API stored securely
	ExpiresOn		int64				// Date when the token expires as Unix timestamp
}

// Contains all Informations and Credentials needed to
// communicate with the Registry API
type Credential struct {
	Username		string				// Username
	Password		memguard.Enclave	// Password stored securely
	Token			BearerToken			// The BearerToken is needed to communicate with the Registry API stored securely
}

// Holds all Informations that are needed to talk with the Registry API
type registry interface {
	ListRepositories()  []Repository	// List all Repositories found in the Registry
	CheckRegistry()		(bool, error)	// Test Authentication, API Version (=> Compatibility)
	FetchAll()			error			// Fetch _all_ Informations (Repos->Images->Tags) available in the Registry
}

// Holds all Informations that are needed to talk with the Registry API
// Implements the @registry Interface
type DockerRegistry struct {
	Version			_RegistryVersion	// API Version
	URI				string				// Registry URI
	Authentication	Auth				// Authentication Informations and Credentials
}

// A (Docker) Repository is (for example) the 'atlassian-jira' in 'docker.reg.local/atlassian-jira:v1.0.0'
type repository interface {
	ListImages()		[]Image			// List all available Images
	FetchAllImages()	error			// Fetch _all_ Image Informations (Images->Tags) available in the Repository
}

// A (Docker) Repository is (for example) the 'atlassian-jira' in 'docker.reg.local/atlassian-jira:v1.0.0'
// or the 'testing/unstable' in 'docker.reg.local/testing/unstable/atlassian-jira:v2.0.0'
// Implements the @repository Interface
type Repository struct {
	Name			string				// Name of the Repository (eg. 'atlassian-jira' or 'testing/unstable')
	Images			[]Image				// All
}

// An Image represents a **single** Docker Image (with _Tag)
type image interface {
	ListImageTags() []_Tag 				// List all available Tags of a Image
	FetchAllTags()	error				// Fetch _all_ Tags from the Image
}

// An Image represents a **single** Docker Image (with all Tags)
// Implements the @image Interface
type Image struct {
	Name			string  			// Image Name (eg. 'nginx')
	Tags			[]Tag  				// List of all available Tags
}

// Describes a Tag of a Image in a Repository
// Implements the @tag Interface
type Tag struct {
	TagName			_Tag				// Image Tag (eg 'v1.0.0')
	ContentDigest	string				// Docker Content Digest
}
