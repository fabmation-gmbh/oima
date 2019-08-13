package interfaces



// Holds all Informations that are needed to talk with the Registry API
type Registry interface {
	// Initialize Registry (and all required Components (Auth, ...))
	Init()				error

	// List all Repositories found in the Registry
	ListRepositories()  []BaseRepository

	// TODO
	// Test Authentication, API Version (=> Compatibility)
	CheckRegistry()		(bool, error)

	// Fetch _all_ Informations (Repos->Images->Tags) available in the Registry
	FetchAll()			error

	// Stats() returns Statistics of the Registry
	Stats()				Stats

	/// >>>>>>>>>> Getter & Setter <<<<<<<<<<

}

// Holds/ Checks and gets needed Credentials/ Informations
// to communicate with the Registry API
type BaseCredential interface {
	Init(cred *BaseCredential)	error		// Checks and "Initializes" the Credential Struct

	getBearerToken()			error		// getBearerToken() sets/ renews the Token in Credential.Token.BearerToken
}

type Auth interface {
	Init()
}

// A (Docker) Repository is (for example) the 'atlassian-jira' in 'docker.reg.local/atlassian-jira:v1.0.0'
type BaseRepository interface {
	ListImages()		([]BaseImage, error)	// List all available Images
	FetchAllImages()	error					// Fetch _all_ Image Informations (Images->Tags) available in the Repository
}

// An Image represents a **single** Docker Image (with TagName)
type BaseImage interface {
	ListImageTags() 	([]Tag, error)		// List all available Tags of a Image
	FetchAllTags()		error					// Fetch _all_ Tags from the Image

	DeleteSignature(*Tag)					// Delete the Signature of an Tag from the S3-Server


	/// >>>>>>>>>> Getter & Setter <<<<<<<<<<

	GetTags()			[]Tag					// Return copy of Tag slice
	GetTagsPtr()		*[]Tag					// Return Pointer to Tags slice
	SetTags([]Tag)								// Overwrite Field 'Tags' with the new Tag slice

	GetName()			string					// Return Name of the Image

	GetRegistryURI()	string					// Returns the URI of the Registry
}




// Describes a Tag of a Image in a Repository
// Implements the @tag Interface
type Tag struct {
	Name          TagName // Image Tag (eg 'v1.0.0')
	ContentDigest string  // Docker Content Digest
	S3SignFound   bool    // Is a Signature found on the S3 Server
}