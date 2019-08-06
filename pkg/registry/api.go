// Registry Package adds the Possibility to talk with the (Docker) Registry API
package registry

import (
	"encoding/json"
	"fmt"
	"github.com/awnumar/memguard"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
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

	// TODO
	// Test Authentication, API Version (=> Compatibility)
	CheckRegistry()		(bool, error)

	// Fetch _all_ Informations (Repos->Images->Tags) available in the Registry
	FetchAll()			error
}

// Holds/ Checks and gets needed Credentials/ Informations
// to communicate with the Registry API
type credential interface {
	Init(cred *Credential)	error		// Checks and "Initializes" the Credential Struct

	getBearerToken()		error		// getBearerToken() sets/ renews the Token in Credential.Token.BearerToken
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
	ListImageTags() ([]Tag, error)		// List all available Tags of a Image
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
	BearerToken		*memguard.Enclave	`json:"token"`

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
	r.URI = conf.Registry.RegistryURI
	Log.Debugf("Set DockerRegistry URI: %s", r.URI)

	// set parent Pointer back to this struct
	r.Authentication.dockerRegistry = r

	// Initialize Auth Struct
	b, _ := strconv.ParseBool(conf.Registry.RequireAuth)
	r.Authentication.Required = b

	r.Authentication.Init()

	err := r.Authentication.Cred.Init()
	if err != nil {
		Log.PanicF("Could not Initialize Credentials: %s", err.Error())
	}

	// add Repo '/', this is a pseudonym for all Images
	// which are stored at the Root Path (this is possible in e.g. JFrog Artifactory)
	var doBreak bool
	repo := Repository{
		DockerRegistry: r,
		Name:           "/",
		Images:         nil,
	}

	// check if this Pseudonym already exists, because
	// the Application could call Init() twice (or more)
	for _, repoV := range r.Repos { if repoV.Name == "" { doBreak = true; break } }
	if doBreak { doBreak = false }


	// fetch Images from '/'
	Log.Debugf("Fetching Images for Root of Registry '%s", r.URI)

	err = repo.FetchAllImages()
	if err != nil {
		Log.Fatalf("Error while Fetching all Images: %s", err.Error())
		return err
	}
	r.Repos = append(r.Repos, repo)

	return nil
}

//noinspection GoNilness
func (r *DockerRegistry) ListRepositories() []Repository {
	// check if Repos where already fetched
	if len(r.Repos) == 0 || r.Repos == nil {
		// fetch all Informations
		err := r.FetchAll()
		if err != nil {
			Log.Fatalf("Error while Fetching all Informations from Registry '%s': %s", r.URI, err.Error())
			os.Exit(1)
		}
	}

	return r.Repos
}

func (r *DockerRegistry) FetchAll() error {
	authData := authInfo{
		token:   nil,
		authReq: r.Authentication.Required,
	}

	if r.Authentication.Required {
		// check if the Bearer Token is expired, and renew it if needed
		if r.Authentication.Cred.Token.ExpiresOn <= time.Now().Unix() {
			// renew BearerToken
			Log.Debugf("Re-Newing BearerToken because it's expired on %s", time.Unix(r.Authentication.Cred.Token.ExpiresOn, 0))

			err := r.Authentication.Cred.getBearerToken()
			if err != nil {
				Log.Fatalf("Error while re-newing the BearerToken: %s", err.Error())
				return err
			}
		}

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

	var wg sync.WaitGroup
	wg.Add(len(catalog))

	var retErr error
	retErr = nil

	for _, val := range catalog {
		go func(val string) {
			defer wg.Done()

			// if Entry does not contain a '/' it means that it is a Image
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

				// check if a Repo already exists with that Name
				for _, repoV := range r.Repos { if repoV.Name == name { return } }

				// fetch Images from repo
				Log.Debugf("Fetching Images for Repo '%s'", repo.Name)

				err := repo.FetchAllImages()
				if err != nil {
					Log.Fatalf("Error while Fetching all Images: %s", err.Error())
					retErr = err
					wg.Done()
					return
				}

				r.Repos = append(r.Repos, repo)
				Log.Debugf("-- New Repo Entry: %s", repo.Name)
			}
		}(val)
	}

	// wait on for-loop
	wg.Wait()

	// TODO: (1) Test Time with a Lot of Repos/ Images/ Tags
	// TODO: (2) Implement own sort Algorithm to improve runtime?
	// TODO:   => (eg) by calling a DockerRegistry.AddRepo(...) Function
	// TODO:      and this functions adds the Repo at the optimal place
	// sort Repos (this increases the runtime about +6,4281014 %)
	sort.Slice(r.Repos, func(i, j int) bool { return r.Repos[i].Name < r.Repos[j].Name 	})

	if retErr != nil { return retErr }

	return nil
}


func (r *Repository) ListImages() ([]Image, error) {
	if len(r.Name) == 0 {
		Log.Fatal("[Internal Error] Trying to List Images about a Repo which Name is not set!!")
		return nil, errors.NewRepositoryNameNotDefinedError()
	}

	if r.Images == nil {
		err := r.FetchAllImages()
		if err != nil {
			Log.Fatalf("Error while fetching all Images from Repo '%s': %s", r.Name, err.Error())
			memguard.SafeExit(1)
		}
	}

	// fetch all Images (and Image Tags)
	for _, v := range r.Images {
		// fetch Tags if len of Tags are 0
		if v.Tags == nil || len(v.Tags) == 0 {
			err := v.FetchAllTags()
			if err != nil {
				Log.Fatalf("Error while Fetching Image Tags: %s", err.Error())
				return nil, err
			}
		}
	}

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

	var wg sync.WaitGroup
	wg.Add(len(catalog))

	for _, v := range catalog {
		go func(v string) {
			defer wg.Done()

			if r.Images == nil { r.Images = []Image{} }
			// check if Entry is an Image or an Repo
			if (strings.Contains(v, r.Name) || r.Name == "/") && !strings.HasSuffix(v, "/") {
				// check if Repo is Root of the Registry
				if r.Name != "/"{
					// check if Image is Entry of an Sub-Repo
					// (eg 'nextcloud' is a Entry of the Sub-Repo 'library' in 'docker.io/library/nextcloud')
					if strings.Count(v, "/") > (strings.Count(r.Name, "/") + 1) { return }
				} else { if strings.Count(v, "/") > 0 { return } }

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
		}(v)
	}

	// wait for For-Loop
	wg.Wait()

	// INFO: This is commented out because it increases the runtime about +94 %
	//		 And sorted Images aren't really important
	//// TODO: (1) Test Time with a Lot of Repos/ Images/ Tags
	//// TODO: (2) Implement own sort Algorithm to improve runtime?
	//// TODO:   => (eg) by calling a Repository.AddImage(...) Function
	//// TODO:      and this functions adds the Repo at the optimal place
	//// sort Images
	//sort.Slice(r.Images, func(i, j int) bool { return r.Images[i].Name < r.Images[j].Name 	})

	return nil
}


//noinspection ALL
func (i *Image) FetchAllTags() error {
	Log.Debugf("Fetching Tags for Image %s", i.Name)

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
		Tags []string 	`json:"tags"`	// Array of all Tags of an Image
	}

	var tags _tags

	err = json.Unmarshal(resp.Body(), &tags)
	if err != nil {
		Log.Debugf("Response: %s", resp.Body())
		Log.Fatalf("Error while marshaling Response: %s", err.Error())
		memguard.SafeExit(1)
	}
	var retErr error
	retErr = nil

	var wg sync.WaitGroup
	wg.Add(len(tags.Tags))

	for _, v := range tags.Tags {
		go func(v string) {
			defer wg.Done()

			var imageData = imageInfo{
				name: i.Name,
				tag:  "",
			}

			newTag := Tag{
				TagName:       "",
				ContentDigest: "",
			}

			imageData.tag = v
			newTag.TagName = _Tag(v)

			// get Image-Tag Digest
			newTag.ContentDigest, err = getTagDigest(&authData, imageData,
				i.Repository.DockerRegistry.URI, i.Repository.DockerRegistry.Version)
			if err != nil {
				Log.Fatalf("Error while getting Image-Tag Digest: %s", err.Error())
				retErr = err
				wg.Done()
				return
			}

			// add Tag to the other Tags
			i.Tags = append(i.Tags, newTag)
			Log.Debugf("==> Digest (%s:%s): %s", i.Name, imageData.tag, newTag.ContentDigest)
		}(v)
	}

	// wait for the for loop
	wg.Wait()

	if retErr != nil { return retErr }

	// TODO: (1) Test Time with a Lot of Repos/ Images/ Tags
	// TODO: (2) Implement own sort Algorithm to improve runtime?
	// TODO:   => (eg) by calling a Image.AddTag(...) Function
	// TODO:      and this functions adds the Repo at the optimal place
	// sort Tags (this increases the runtime about +58,581839 %)
	sort.Slice(i.Tags, func(index, j int) bool { return i.Tags[index].TagName < i.Tags[j].TagName })

	return nil
}

func (i* Image) ListImageTags() ([]Tag, error) {
	if len(i.Name) == 0 {
		Log.Fatal("[Internal Error] Trying to List Image Tags from an Image which Name is not set!!")
		return nil, errors.NewImageNameNotDefinedError()
	}

	if i.Tags == nil || len(i.Tags) == 0 {
		err := i.FetchAllTags()
		if err != nil {
			Log.Fatalf("Error while fetching all Tags of Image '%s': %s", i.Name, err.Error())
			return nil, err
		}
	}

	return i.Tags, nil
}


func (a *Auth) Init() { a.Cred.auth = a }


//noinspection GoNilnesss
func (c *Credential) Init()	error {
	// get Password
	pwdEnclave, err := internal.Cred.GetCredential("password")
	if err != nil { Log.PanicF("Error while getting Credential from CredStore: %s", err.Error()) }
	c.Password = pwdEnclave

	// get and set Username
	c.Username = conf.Registry.Username

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
		err = c.getBearerToken()
		if err != nil {
			Log.Fatalf("Error while getting Bearer Token: %s", err.Error())
			memguard.SafeExit(1)
		}
	} else { Log.Notice("Authentication not required, so no need to get a Bearer Token") }

	return nil
}

func (c *Credential) getBearerToken() error {
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
		return err
	}

	type _token struct {
		// The BearerToken is needed to communicate with the Registry API
		BearerToken		string				`json:"token"`

		// Date when the token expires as Unix timestamp
		ExpiresOn		int64				`json:"expires_in"`
	}
	tokenData := _token{}

	err = json.Unmarshal(resp.Body(), &tokenData)
	if err != nil {
		Log.Debugf("Response: %s", resp.Body())
		Log.Fatalf("Error while marshaling Response: %s", err.Error())
		return err
	}

	// prepare Token
	if c.Token.BearerToken == nil { c.Token.BearerToken = memguard.NewEnclaveRandom(len(tokenData.BearerToken)) }

	bearerToken, err := c.Token.BearerToken.Open()
	if err != nil { memguard.SafePanic(err) }
	defer bearerToken.Destroy()

	// make bearerToken immutable
	bearerToken.Melt()

	c.Token.ExpiresOn = tokenData.ExpiresOn
	bearerToken.Copy([]byte(tokenData.BearerToken))


	Log.Debugf("Response: %s", resp.Body())

	// convert Seconds in BearerToken.ExpiresOn into Unix Timestamp
	c.Token.ExpiresOn = time.Now().Unix() + c.Token.ExpiresOn
	Log.Debugf("Bearer Token: '%s'", bearerToken.String())
	Log.Debugf("Bearer Token Expires On %d (%s)", c.Token.ExpiresOn, time.Unix(c.Token.ExpiresOn, 0))

	// return Encrypted Data back to Token Struct
	c.Token.BearerToken = bearerToken.Seal()
	memguard.ScrambleBytes([]byte(tokenData.BearerToken))

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