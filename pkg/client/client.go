package client

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/inconshreveable/go-update"
)

// https://github.com/inconshreveable/go-update
func DoUpdate(url string) error {
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

func GetLatestVersion() (latestVersion string, err error) {
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
