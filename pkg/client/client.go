package client

import (
	"io/ioutil"
	"net/http"

	"github.com/inconshreveable/go-update"
	"github.com/mozey/httprouter-util/pkg/config"
	"github.com/pkg/errors"
)

type Client struct {
	Config *config.Config
}

func NewHandler(conf *config.Config) (c *Client) {
	c = &Client{}
	c.Config = conf
	return c
}

// https://github.com/inconshreveable/go-update
func (c *Client) DoUpdate(token string) error {
	resp, err := http.Get(c.Config.ExecTemplateClientDownloadUrl(token))
	if err != nil {
		return err
	}
	defer (func() {
		_ = resp.Body.Close()
	})()
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *Client) GetLatestVersion(token string) (
	latestVersion string, err error) {

	client := &http.Client{}

	req, err := http.NewRequest(
		"GET", c.Config.ExecTemplateClientVersionUrl(token), nil)
	if err != nil {
		return latestVersion, errors.WithStack(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return latestVersion, errors.WithStack(err)
	}
	defer (func() {
		_ = resp.Body.Close()
	})()

	if resp.StatusCode != http.StatusOK {
		return latestVersion, errors.Errorf(
			"%v %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return latestVersion, errors.WithStack(err)
	}

	return string(b), nil
}
