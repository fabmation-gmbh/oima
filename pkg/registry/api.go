// Registry Package adds the Possibility to talk with the (Docker) Registry API
package registry

import (
	"encoding/json"
	"fmt"
	"github.com/awnumar/memguard"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/fabmation-gmbh/oima/pkg/config"
	"github.com/fabmation-gmbh/oima/internal"
	. "github.com/fabmation-gmbh/oima/internal/log"
	_http "github.com/fabmation-gmbh/oima/pkg/http"
)

var conf config.Configuration


type _Tag string             // _Tag of an Image (for example 'v1.0.0' or '0.1.0-beta'
type _RegistryVersion string // Describes the current Version of the Registry API
const (
	V1	_RegistryVersion	= "v1"		// (Docker) Registry API Version v1
	V2	_RegistryVersion	= "v2"		// (Docker) Registry API Version v1
)

type auth interface {
	Init()
}

// Registry Authentication Information
type Auth struct {
	dockerRegistry	*DockerRegistry		// Pointer to Parent Struct

	Required		bool				// Is Authentication Required
	Cred			Credential			// Needed Credentials
}

// Holds/ Checks and gets needed Credentials/ Informations
// to communicate with the Registry API
type credential interface {
	Init(cred *Credential)	error		// Checks and "Initializes" the Credential Struct
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
type registry interface {
	Init()				error			// Initialize Registry (and all required Components (Auth, ...))
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
	DockerRegistry	*DockerRegistry		// Pointer to Parent Struct

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


/// >>>>>>>>>> Functions <<<<<<<<<< ///
func (r *DockerRegistry) Init() error {
	conf = internal.GetConfig()

	// set Flags
	r.URI = conf.Regitry.RegistryURI
	Log.Debugf("Set DockerRegistry URI: %s", r.URI)

	// set parent Pointer back to this struct
	r.Authentication.dockerRegistry = r
	if b, _ := strconv.ParseBool(conf.Regitry.RequireAuth); b {
		r.Authentication.Required = true
	} else {
		r.Authentication.Required = false
	}
	r.Authentication.Init()

	err := r.Authentication.Cred.Init()
	if err != nil {
		Log.PanicF("Could not Initialize Credentials: %s", err.Error())
	}

	return nil
}

func (a *Auth) Init() { a.Cred.auth = a }

//noinspection GoNilness
func (c *Credential) Init()	error {
	var uri = fmt.Sprintf("%s/api/docker/docker/v2/token", c.auth.dockerRegistry.URI)
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

	// get Bearer Token
	httpClient := _http.NewClient()
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		Log.Fatal(err.Error())
		memguard.SafeExit(1)
	}

	// set Request Parameters
	req.Header.Set("User-Agent", "oima-client")
	if b, _ := strconv.ParseBool(conf.Regitry.RequireAuth); b {
		req.SetBasicAuth(conf.Regitry.Username, password.String())
	}

	response, err := httpClient.Do(req)
	if err != nil {
		Log.Fatalf("Error while making Request: %s", err.Error())
		memguard.SafeExit(1)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		Log.Fatalf("Error while Reading Response: %s", err.Error())
		memguard.SafeExit(1)
	}

	err = json.Unmarshal(body, &c.Token)
	if err != nil {
		Log.Fatalf("Error while marshaling Response: %s", err.Error())
		memguard.SafeExit(1)
	}

	// convert Seconds in BearerToken.ExpiresOn into Unix Timestamp
	c.Token.ExpiresOn = time.Now().Unix() + c.Token.ExpiresOn
	Log.Debugf("Bearer Token: '%s'", c.Token.BearerToken)
	Log.Debugf("Bearer Token Expires On %d (%s)", c.Token.ExpiresOn, time.Unix(c.Token.ExpiresOn, 0))

	return nil
}