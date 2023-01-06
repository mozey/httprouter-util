package main

import (
	"flag"
	"fmt"
	"github.com/inconshreveable/go-update"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var version string

// https://github.com/inconshreveable/go-update
func doUpdate(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer (func() {
		_ = resp.Body.Close()
	})()
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		// error handling
	}
	return err
}

func getLatestVersion() (latestVersion string, err error) {
	client := &http.Client{}

	req, err := http.NewRequest(
		"GET", "http://localhost:8118/client/version?token=123", nil)
	if err != nil {
		return latestVersion, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return latestVersion, err
	}
	defer (func() {
		_ = resp.Body.Close()
	})()

	if resp.StatusCode != http.StatusOK {
		return latestVersion, fmt.Errorf(
			"%v %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return latestVersion, err
	}

	return string(b), nil
}

func main() {
	versionFlag := flag.Bool(
		"version", false, "Print version")
	updateFlag := flag.Bool(
		"update", false, "Request an update from the server")

	flag.Parse()

	if *versionFlag {
		fmt.Println(version)

	} else if *updateFlag {
		latestVersion, err := getLatestVersion()
		fmt.Println(fmt.Sprintf("latest version is %s", latestVersion))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		latestVersion = strings.TrimSpace(
			strings.ReplaceAll(latestVersion, "\n", ""))
		if latestVersion != version {
			fmt.Println("updating...")
			err := doUpdate("http://localhost:8118/client/download?token=123")
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
