// Registry Package adds the Possibility to talk with the (Docker) Registry API
package registry

import (
	"encoding/json"
	"fmt"
	"github.com/awnumar/memguard"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty"

	"github.com/fabmation-gmbh/oima/internal"
	. "github.com/fabmation-gmbh/oima/internal/log"
	"github.com/fabmation-gmbh/oima/pkg/config"
	"github.com/fabmation-gmbh/oima/pkg/errors"
)

var conf config.Configuration


type _Tag string             // _Tag of an Image (for example 'v1.0.0' or '0.1.0-beta'
type _RegistryVersion string // Describes the current Version of the Registry API
const (
	VUNK _RegistryVersion	= "UNKNOWN"	// Unknown Registry API Version

	V1	_RegistryVersion	= "v1"		// (Docker) Registry API Version v1
	V2	_RegistryVersion	= "v2"		// (Docker) Registry API Version v1
)

// Holds all Informations that are needed to talk with the Registry API
type registry interface {
	// Initialize Registry (and all required Components (Auth, ...))
	Init()				error

	// List all Repositories found in the Registry
	ListRepositories()  []Repository

	// Test Authentication, API Version (=> Compatibility)
	CheckRegistry()		(bool, error)

	// Fetch _all_ Informations (Repos->Images->Tags) available in the Registry
	FetchAll()			error
}

// Holds/ Checks and gets needed Credentials/ Informations
// to communicate with the Registry API
type credential interface {
	Init(cred *Credential)	error		// Checks and "Initializes" the Credential Struct
}

type auth interface {
	Init()
}

// A (Docker) Repository is (for example) the 'atlassian-jira' in 'docker.reg.local/atlassian-jira:v1.0.0'
type repository interface {
	ListImages()		([]Image, error)	// List all available Images
	FetchAllImages()	error				// Fetch _all_ Image Informations (Images->Tags) available in the Repository
}

// An Image represents a **single** Docker Image (with _Tag)
type image interface {
	ListImageTags() []_Tag 				// List all available Tags of a Image
	FetchAllTags()	error				// Fetch _all_ Tags from the Image
}




// Registry Authentication Information
type Auth struct {
	dockerRegistry	*DockerRegistry		// Pointer to Parent Struct

	Required		bool				// Is Authentication Required
	Cred			Credential			// Needed Credentials
}

// The BearerToken is needed to communicate with the Registry API
type BearerToken struct {
	// The BearerToken is needed to communicate with the Registry API stored securely
	BearerToken		memguard.Enclave	`json:"token"`

	// Date when the token expires as Unix timestamp
	ExpiresOn		int64				`json:"expires_in"`
}

// Contains all Informations and Credentials needed to
// communicate with the Registry API
type Credential struct {
	auth			*Auth				// Pointer to Parent Struct

	Username		string				// Username
	Password		*memguard.Enclave	// Password stored securely
	Token			BearerToken			// The BearerToken is needed to communicate with the Registry API stored securely
}

// Holds all Informations that are needed to talk with the Registry API
// Implements the @registry Interface
type DockerRegistry struct {
	Version			_RegistryVersion	// API Version
	URI				string				// Registry URI
	Authentication	Auth				// Authentication Informations and Credentials
	Repos			[]Repository		// List of all Repos in the Registry
}

// A (Docker) Repository is (for example) the 'atlassian-jira' in 'docker.reg.local/atlassian-jira:v1.0.0'
// or the 'testing/unstable' in 'docker.reg.local/testing/unstable/atlassian-jira:v2.0.0'
// Implements the @repository Interface
type Repository struct {
	DockerRegistry	*DockerRegistry		// Pointer to Parent Struct

	// A Repo can contain unlimited Sub-Repos (and Sub-Sub-Repos, Sub-Sub-Sub-Repos, ...)
	// For Example:
	// docker.io
	// |-- nginx						// nginx Image
	// |   `-- v1.0.0						// Image Version v1.0.0
	// `-- unstable/			// unstable Repo
	//     |-- samba					// samba Image
	//     |   |-- v1.0.0					// Image Version v1.0.0
	//     |   `-- v1.1.0					// Image Version v1.1.0
	//     `-- testing/			// testing Sub-Repo
	//         |-- jira					// jira Image
	//         |   |-- v1.0.0				// Image Version v1.0.0
	//         |   `-- v1.1.0				// Image Version v1.1.0
	//         `-- wiki					// wiki Image
	//             |-- v1.0.0				// Image Version v1.0.0
	//             `-- v2.0.0				// Image Version v2.0.0
	Name			string				// Name of the Repository (eg. 'stable')
	Images			[]Image				// All
}

// An Image represents a **single** Docker Image (with all Tags)
// Implements the @image Interface
type Image struct {
	Repository		*Repository			// Pointer to Parent Struct

	Name			string  			// Image Name (eg. 'nginx')
	Tags			[]Tag  				// List of all available Tags
}

// Describes a Tag of a Image in a Repository
// Implements the @tag Interface
type Tag struct {
	TagName			_Tag				// Image Tag (eg 'v1.0.0')
	ContentDigest	string				// Docker Content Digest
}


/// >>>>>>>>>> Functions <<<<<<<<<< ///

//noinspection GoNilness
func (r *DockerRegistry) Init() error {
	conf = internal.GetConfig()

	// set Flags
	r.URI = conf.Regitry.RegistryURI
	Log.Debugf("Set DockerRegistry URI: %s", r.URI)

	// set parent Pointer back to this struct
	r.Authentication.dockerRegistry = r

	// Initialize Auth Struct
	b, _ := strconv.ParseBool(conf.Regitry.RequireAuth)
	r.Authentication.Required = b

	r.Authentication.Init()

	err := r.Authentication.Cred.Init()
	if err != nil {
		Log.PanicF("Could not Initialize Credentials: %s", err.Error())
	}
	return nil
}

//noinspection GoNilness
func (r *DockerRegistry) ListRepositories() []Repository {
	// TODO: check if r.Authentication.Cred.Token is expired
	authData := authInfo{
		token:   nil,
		authReq: r.Authentication.Required,
	}

	if r.Authentication.Required {
		token, err := r.Authentication.Cred.Token.BearerToken.Open()
		if err != nil {
			memguard.SafePanic(err)
		}
		defer token.Destroy()

		authData.token = token
	}

	// get Registry Catalog
	catalog, err := getRegistryCatalog(&authData, r.URI, r.Version)
	if err != nil {
		Log.Fatalf("Error while fetching Registry Catalog: %s", err.Error())
		memguard.SafeExit(1)
	}

	// TODO: Move me to FetchAll()
	for _, val := range catalog {
		if strings.Contains(val, "/") {
			repoName := strings.Split(val, "/")
			lenRepos := len(repoName)
			var name string
			for i, v := range repoName {
				if i == (lenRepos - 1) { break }

				// prevent adding a Slash to the last Repo Entry Name
				if i != (lenRepos -2 ) {
					name += fmt.Sprintf("%s/", v)
				} else { name += v }
			}

			repo := Repository{
				DockerRegistry: r,
				Name:           name,
				Images:         nil,
			}
			r.Repos = append(r.Repos, repo)
			Log.Debugf("-- New Repo Entry: %s", repo.Name)
		}
	}

	// TODO: Move me to FetchAll()

	return r.Repos
}
	}

	// fetch all Images (and Image Tags)
	//for _, v := range r.Images { // TODO: implement me }

	return r.Images, nil
}

func (r *Repository) FetchAllImages() error {
	if len(r.Name) == 0 {
		Log.Fatal("[Internal Error] Trying to fetch all Images in a Repo which Name is not set!!")
		return errors.NewRepositoryNameNotDefinedError()
	}

	var authData = authInfo{
		token:   nil,
		authReq: r.DockerRegistry.Authentication.Required,
	}

	if r.DockerRegistry.Authentication.Required {
		token, err := r.DockerRegistry.Authentication.Cred.Token.BearerToken.Open()
		if err != nil {
			memguard.SafePanic(err)
		}
		defer token.Destroy()

		authData.token = token
	}

	// get Registry Catalog
	catalog, err := getRegistryCatalog(&authData, r.DockerRegistry.URI, r.DockerRegistry.Version)
	if err != nil {
		Log.Fatalf("Error while fetching Registry Catalog: %s", err.Error())
		memguard.SafeExit(1)
	}

	for _, v := range catalog {
		// check if Entry is an Image or an Repo
		if strings.HasSuffix(v, "/") {
			newImage := Image{
				Repository: r,
				Name:       v,
				Tags:       nil,
			}

			// fetch Image Tags
			err := newImage.FetchAllTags()
			if err != nil {
				Log.Fatalf("Error while Fetching Tags of Image '%s': %s", newImage.Name, err.Error())
				memguard.SafeExit(1)
			}

			r.Images = append(r.Images, newImage)
			Log.Debugf("--> Add new Image: %s", newImage.Name)
		}
	}

	return nil
}

//noinspection ALL
func (i *Image) FetchAllTags() error {
	// check Image Name
	if len(i.Name) == 0 {
		Log.Fatal("[Internal Error] Trying to fetch all Tags of an Image which Name is not set!!")
		return errors.NewImageNameNotDefinedError()
	}

	var uri = fmt.Sprintf("%s/%s/%s/tags/list",
			i.Repository.DockerRegistry.URI,
			i.Repository.DockerRegistry.Version,
			i.Name)
	var authData = authInfo{
		token:   nil,
		authReq: i.Repository.DockerRegistry.Authentication.Required,
	}

	client := resty.New()
	client.SetHeaders(map[string]string{
		"Docker-Distribution-Api-Version": "registry/2.0",
		"User-Agent":                      "oima-cli",
	})

	if i.Repository.DockerRegistry.Authentication.Required {
		token, err := i.Repository.DockerRegistry.Authentication.Cred.Token.BearerToken.Open()
		if err != nil {
			memguard.SafePanic(err)
		}
		defer token.Destroy()

		client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token.String()))
		authData.token = token
	}

	resp, err := client.R().Get(uri)
	if err != nil {
		Log.Criticalf("Error while getting Auth. Token: %s", err.Error())
		memguard.SafeExit(1)
	}

	type _tags struct {
		tags	[]string		// Array of all Tags of an Image
	}

	var tags = _tags{}

	err = json.Unmarshal(resp.Body(), &tags)
	if err != nil {
		Log.Debugf("Response: %s", resp.Body())
		Log.Fatalf("Error while marshaling Response: %s", err.Error())
		memguard.SafeExit(1)
	}

	var imageData = imageInfo{
		name: i.Name,
		tag:  nil,
	}

	newTag := Tag{
		TagName:       v,
		ContentDigest: nil,
	}
	for _, v := range tags {
		imageData.tag = v

		// get Image-Tag Digest
		newTag.ContentDigest, err = getTagDigest(&authData, imageData,
				i.Repository.DockerRegistry.URI, i.Repository.DockerRegistry.Version)
		if err != nil {
			Log.Fatalf("Error while getting Image-Tag Digest: %s", err.Error())
			return err
		}

		// add Tag to the other Tags
		i.Tags = append(i.Tags, newTag)
		Log.Debugf("Digest (%s:%s): %s", i.Name, imageData.tag, newTag.ContentDigest)
	}

	return nil
}

func (a *Auth) Init() { a.Cred.auth = a }

//noinspection GoNilness
func (c *Credential) Init()	error {
	var password *memguard.LockedBuffer

	// get Password
	pwdEnclave, err := internal.Cred.GetCredential("password")
	if err != nil {
		Log.PanicF("Error while getting Credential from CredStore: %s", err.Error())
	}
	c.Password = pwdEnclave

	password, err = pwdEnclave.Open()
	if err != nil {
		memguard.SafePanic(err)
	}
	defer password.Destroy()

	// get Registry Version
	c.auth.dockerRegistry.Version, err = getRegistryVersion(c)
	if err != nil {
		Log.Errorf("Error while getting Registry API Version: %s", err.Error())
		memguard.SafeExit(1)
	}

	if c.auth.dockerRegistry.Version == V1 {
		Log.Errorf("Registry API Version is v1. Version 1 isn't supported yet!")
		memguard.SafeExit(1)
	}

	if c.auth.Required {
		// get Bearer Token
		client := resty.New()
		client.SetHeaders(map[string]string{
			"Docker-Distribution-Api-Version": "registry/2.0",
			"User-Agent": "oima-cli",
		})

		var password *memguard.LockedBuffer

		// get Password
		password, err := c.Password.Open()
		if err != nil { memguard.SafePanic(err) }
		defer password.Destroy()

		client.SetBasicAuth(c.Username, password.String())

		var uri = fmt.Sprintf("%s/api/docker/docker/%s/token",
			c.auth.dockerRegistry.URI, c.auth.dockerRegistry.Version)
		resp, err := client.R().Get(uri)
		if err != nil {
			Log.Criticalf("Error while getting Auth. Token: %s", err.Error())
			memguard.SafeExit(1)
		}

		err = json.Unmarshal(resp.Body(), &c.Token)
		if err != nil {
			Log.Debugf("Response: %s", resp.Body())
			Log.Fatalf("Error while marshaling Response: %s", err.Error())
			memguard.SafeExit(1)
		}

		Log.Debugf("Response: %s", c.Token)

		// convert Seconds in BearerToken.ExpiresOn into Unix Timestamp
		c.Token.ExpiresOn = time.Now().Unix() + c.Token.ExpiresOn
		Log.Debugf("Bearer Token: '%s'", c.Token.BearerToken)
		Log.Debugf("Bearer Token Expires On %d (%s)", c.Token.ExpiresOn, time.Unix(c.Token.ExpiresOn, 0))
	} else { Log.Notice("Authentication not required, so no need to get a Bearer Token") }

	return nil
}

/// ------------- Internal Functions

//noinspection GoNilness
func getRegistryVersion(c *Credential) (_RegistryVersion, error) {
	var version _RegistryVersion
	client := resty.New()

	client.SetHeader("User-Agent", "oima-cli")

	if c.auth.Required {
		var password *memguard.LockedBuffer

		// get Password
		password, err := c.Password.Open()
		if err != nil { memguard.SafePanic(err) }
		defer password.Destroy()

		client.SetBasicAuth(c.Username, password.String())
	}

	resp, err := client.R().Get(fmt.Sprintf("%s/v2/", c.auth.dockerRegistry.URI))
	if err != nil { return VUNK, err }

	if resp.StatusCode() == 404 {
		version = V1
	} else {
		version = V2
	}

	return version, nil
}