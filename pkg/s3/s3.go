package s3

import (
	"regexp"
	"strings"
)

// This File contains base Functions which could be used by
// all S3-Provider Implementations

// PrepareRegPath() checks and manipulates the Registry URI
// so that it's usable for S3 Clients
// Example:
//		Input:		https://docker.internal.int/
//		Output:		docker.internal.int
func PrepareRegPath(uri string) string {
	r, _ := regexp.Compile("http(s)?://")					// remove 'https://' or 'http://'
	return strings.ReplaceAll(r.ReplaceAllString(uri, ""), "/", "") // remove '/' Postfix
}