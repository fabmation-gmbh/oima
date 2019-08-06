package interfaces

type TagName string             // TagName of an Image (for example 'v1.0.0' or '0.1.0-beta')


// Stats contains statistics of a Registry
type Stats struct {
	// Number of Repositories
	Repos				int

	// Number of Images
	Images				int

	// Number of Tags
	Tags				int

	// Number of Signatures found on the S3 Server
	S3Signatures		int
}