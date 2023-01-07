package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mozey/httprouter-util/pkg/client"
)

var version string

func main() {
	versionFlag := flag.Bool(
		"version", false, "Print version")
	updateFlag := flag.Bool(
		"update", false, "Request an update from the server")

	flag.Parse()

	if *versionFlag {
		fmt.Println(version)

	} else if *updateFlag {
		latestVersion, err := client.GetLatestVersion()
		fmt.Println(fmt.Sprintf("latest version is %s", latestVersion))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		latestVersion = strings.TrimSpace(
			strings.ReplaceAll(latestVersion, "\n", ""))
		if latestVersion != version {
			fmt.Println("updating...")
			err := client.DoUpdate("http://localhost:8118/client/download?token=123")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			fmt.Println("already on the latest version")
		}

	} else {
		flag.Usage()
	}

	os.Exit(0)
}
