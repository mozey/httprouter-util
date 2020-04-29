package main

import "github.com/mozey/config/pkg/config"

// Compiled with ldflags
var AppDir string

func main() {
	config.Run(AppDir)
}
