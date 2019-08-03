package registry

import (
	"encoding/json"
	"fmt"
	"github.com/awnumar/memguard"
	"github.com/go-resty/resty"

	. "github.com/fabmation-gmbh/oima/internal/log"
)

func getRegistryCatalog(
	token *memguard.LockedBuffer,
	authRequired bool,
	regURI string,
	version _RegistryVersion) ([]string, error) {

	var uri = fmt.Sprintf("%s/%s/_catalog", regURI, version)

	client := resty.New()
	client.SetHeaders(map[string]string{
		"Docker-Distribution-Api-Version": "registry/2.0",
		"User-Agent":                      "oima-cli",
	})

	if authRequired {
		client.SetHeader("Authorization", fmt.Sprintf("JWT %s", token.String()))
	}

	resp, err := client.R().Get(fmt.Sprintf(uri))
	if err != nil {
		Log.Criticalf("Error while getting Auth. Token: %s", err.Error())
		memguard.SafeExit(1)
	}

	type response struct {
		Entries []string `json:"repositories"`
	}

	var result response
	var entries []string

	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		Log.Debugf("Response: %s", resp.Body())
		Log.Fatalf("Error while marshaling Response: %s", err.Error())
		memguard.SafeExit(1)
	}

	for _, v := range result.Entries {
		entries = append(entries, v)
	}

	return entries, nil
}
