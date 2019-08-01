package internal

import "fmt"

var Commit="commit unknown"
var BuildDate="build date unknown"
var Version="version unknown"

func GetVersion() string {
	return fmt.Sprintf("%s\nGit Commit: %s\nBuild Date: %s", Version, Commit, BuildDate)
}