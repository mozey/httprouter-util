package main

import (
	"crypto"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mozey/httprouter-util/pkg/client"
	"github.com/mozey/httprouter-util/pkg/config"
	"github.com/mozey/logutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// configBase64 must be compiled into executable with ldflags,
// if it is not set then read config from env
var configBase64 string

func main() {
	// Make console logs human readable
	logutil.SetupLogger(true)

	// Config
	if configBase64 != "" {
		err := config.SetEnvBase64(configBase64)
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}
	}
	conf := config.New()

	// Flags
	versionFlag := flag.Bool(
		"version", false, "Print version")
	updateFlag := flag.Bool(
		"update", false, "Request an update from the server")
	tokenFlag := flag.String(
		"token", "", "Auth token")
	checksumFlag := flag.String(
		"checksum", "", "Print checksum for specified path")
	flag.Parse()

	c := client.NewHandler(conf)

	if *versionFlag {
		fmt.Println(conf.Version())

	} else if *updateFlag {
		clientVersion, err := c.GetLatestVersion(*tokenFlag)
		fmt.Println(fmt.Sprintf("latest version is %s", clientVersion.Version))
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}
		if clientVersion.Version != conf.Version() {
			fmt.Println("updating...")
			err := c.DoUpdate(*tokenFlag, clientVersion.Checksum)
			if err != nil {
				log.Error().Stack().Err(err).Msg("")
				os.Exit(1)
			}
		} else {
			fmt.Println("already on the latest version")
		}

	} else if strings.TrimSpace(*checksumFlag) != "" {
		payload, err := ioutil.ReadFile(*checksumFlag)
		if err != nil {
			log.Error().Stack().Err(errors.WithStack(err)).Msg("")
			os.Exit(1)
		}

		// Must match the hash function options passed to the go-update lib, see
		// vendor/github.com/inconshreveable/go-update/apply.go
		h := crypto.SHA256
		hash := h.New()
		_, err = hash.Write(payload)
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}

		b := hash.Sum([]byte{})
		// Print hash bytes as hexadecimal (base 16), see https://pkg.go.dev/fmt
		fmt.Printf("%x\n", b)

	} else {
		flag.Usage()
	}

	os.Exit(0)
}
