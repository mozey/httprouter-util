package client

import (
	"crypto"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/inconshreveable/go-update"
	"github.com/mozey/httprouter-util/pkg/config"
	"github.com/mozey/httprouter-util/pkg/share"
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
func (c *Client) DoUpdate(token, checksumHex string) error {
	checksumHex = strings.ReplaceAll(checksumHex, "\n", "")
	checksumBytes, err := hex.DecodeString(checksumHex)
	if err != nil {
		return errors.WithStack(err)
	}

	resp, err := http.Get(c.Config.ExecTemplateClientDownloadUrl(token))
	if err != nil {
		return err
	}
	defer (func() {
		_ = resp.Body.Close()
	})()
	err = update.Apply(resp.Body, update.Options{
		Hash:       crypto.SHA256,
		Checksum:   checksumBytes,
		TargetMode: os.FileMode(0755),
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *Client) GetLatestVersion(token string) (
	clientVersion share.ClientVersion, err error) {

	client := &http.Client{}

	req, err := http.NewRequest(
		"GET", c.Config.ExecTemplateClientVersionUrl(token), nil)
	if err != nil {
		return clientVersion, errors.WithStack(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return clientVersion, errors.WithStack(err)
	}
	defer (func() {
		_ = resp.Body.Close()
	})()

	if resp.StatusCode != http.StatusOK {
		return clientVersion, errors.Errorf(
			"%v %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return clientVersion, errors.WithStack(err)
	}

	err = json.Unmarshal(b, &clientVersion)
	if err != nil {
		return clientVersion, errors.WithStack(err)
	}

	return clientVersion, nil
}
