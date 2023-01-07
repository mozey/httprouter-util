package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mozey/httprouter-util/pkg/client"
	"github.com/mozey/httprouter-util/pkg/config"
	"github.com/mozey/logutil"
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
	tokenFlag := flag.String(
		"token", "", "Auth token")
	versionFlag := flag.Bool(
		"version", false, "Print version")
	updateFlag := flag.Bool(
		"update", false, "Request an update from the server")
	flag.Parse()

	c := client.NewHandler(conf)

	if *versionFlag {
		fmt.Println(conf.Version())

	} else if *updateFlag {
		latestVersion, err := c.GetLatestVersion(*tokenFlag)
		fmt.Println(fmt.Sprintf("latest version is %s", latestVersion))
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			os.Exit(1)
		}
		latestVersion = strings.TrimSpace(
			strings.ReplaceAll(latestVersion, "\n", ""))
		if latestVersion != conf.Version() {
			fmt.Println("updating...")
			err := c.DoUpdate(*tokenFlag)
			if err != nil {
				log.Error().Stack().Err(err).Msg("")
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
