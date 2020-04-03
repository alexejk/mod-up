package main

import (
	"fmt"
	"os"

	"alexejk.io/mod-up/cmd"
	"alexejk.io/mod-up/internal/version"
)

var (
	appVersion    = "x.y.z"
	appName       = "~unknown~"
	appCommit     = "~unknown~"
	appCommitDate = "~unknown~"
	appCommitTime = "~unknown~"
)


func main() {

	v := version.NewVersionInfo(appName, appVersion, appCommit, appCommitDate, appCommitTime)

	if err := cmd.RootCmd(v).Execute(); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
		return
	}
}
