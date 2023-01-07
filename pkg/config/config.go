// Code generated with https://github.com/mozey/config DO NOT EDIT

package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// APP_ADDR
var addr string

// APP_EXE
var exe string

// APP_MAX_BYTES_KB
var maxBytesKb string

// APP_MAX_PAYLOAD_MB
var maxPayloadMb string

// APP_NAME
var name string

// APP_DIR
var dir string

// Config fields correspond to config file keys less the prefix
type Config struct {
	addr         string // APP_ADDR
	exe          string // APP_EXE
	maxBytesKb   string // APP_MAX_BYTES_KB
	maxPayloadMb string // APP_MAX_PAYLOAD_MB
	name         string // APP_NAME
	dir          string // APP_DIR
}

// Addr is APP_ADDR
func (c *Config) Addr() string {
	return c.addr
}

// Exe is APP_EXE
func (c *Config) Exe() string {
	return c.exe
}

// MaxBytesKb is APP_MAX_BYTES_KB
func (c *Config) MaxBytesKb() string {
	return c.maxBytesKb
}

// MaxPayloadMb is APP_MAX_PAYLOAD_MB
func (c *Config) MaxPayloadMb() string {
	return c.maxPayloadMb
}

// Name is APP_NAME
func (c *Config) Name() string {
	return c.name
}

// Dir is APP_DIR
func (c *Config) Dir() string {
	return c.dir
}

// SetAddr overrides the value of addr
func (c *Config) SetAddr(v string) {
	c.addr = v
}

// SetExe overrides the value of exe
func (c *Config) SetExe(v string) {
	c.exe = v
}

// SetMaxBytesKb overrides the value of maxBytesKb
func (c *Config) SetMaxBytesKb(v string) {
	c.maxBytesKb = v
}

// SetMaxPayloadMb overrides the value of maxPayloadMb
func (c *Config) SetMaxPayloadMb(v string) {
	c.maxPayloadMb = v
}

// SetName overrides the value of name
func (c *Config) SetName(v string) {
	c.name = v
}

// SetDir overrides the value of dir
func (c *Config) SetDir(v string) {
	c.dir = v
}

// New creates an instance of Config.
// Build with ldflags to set the package vars.
// Env overrides package vars.
// Fields correspond to the config file keys less the prefix.
// The config file must have a flat structure
func New() *Config {
	conf := &Config{}
	SetVars(conf)
	SetEnv(conf)
	return conf
}

// SetVars sets non-empty package vars on Config
func SetVars(conf *Config) {

	if addr != "" {
		conf.addr = addr
	}

	if exe != "" {
		conf.exe = exe
	}

	if maxBytesKb != "" {
		conf.maxBytesKb = maxBytesKb
	}

	if maxPayloadMb != "" {
		conf.maxPayloadMb = maxPayloadMb
	}

	if name != "" {
		conf.name = name
	}

	if dir != "" {
		conf.dir = dir
	}

}

// SetEnv sets non-empty env vars on Config
func SetEnv(conf *Config) {
	var v string

	v = os.Getenv("APP_ADDR")
	if v != "" {
		conf.addr = v
	}

	v = os.Getenv("APP_EXE")
	if v != "" {
		conf.exe = v
	}

	v = os.Getenv("APP_MAX_BYTES_KB")
	if v != "" {
		conf.maxBytesKb = v
	}

	v = os.Getenv("APP_MAX_PAYLOAD_MB")
	if v != "" {
		conf.maxPayloadMb = v
	}

	v = os.Getenv("APP_NAME")
	if v != "" {
		conf.name = v
	}

	v = os.Getenv("APP_DIR")
	if v != "" {
		conf.dir = v
	}

}

// GetMap of all env vars
func (c *Config) GetMap() map[string]string {
	m := make(map[string]string)

	m["APP_ADDR"] = c.addr

	m["APP_EXE"] = c.exe

	m["APP_MAX_BYTES_KB"] = c.maxBytesKb

	m["APP_MAX_PAYLOAD_MB"] = c.maxPayloadMb

	m["APP_NAME"] = c.name

	m["APP_DIR"] = c.dir

	return m
}

// SetEnvBase64 decodes and sets env from the given base64 string
func SetEnvBase64(configBase64 string) (err error) {
	// Decode base64
	decoded, err := base64.StdEncoding.DecodeString(configBase64)
	if err != nil {
		return errors.WithStack(err)
	}
	// UnMarshall json
	configMap := make(map[string]string)
	err = json.Unmarshal(decoded, &configMap)
	if err != nil {
		return errors.WithStack(err)
	}
	// Set config
	for key, value := range configMap {
		err = os.Setenv(key, value)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// LoadFile sets the env from file and returns a new instance of Config
func LoadFile(mode string) (conf *Config, err error) {
	appDir := os.Getenv("APP_DIR")
	p := filepath.Join(appDir, fmt.Sprintf("config.%v.json", mode))
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	configMap := make(map[string]string)
	err = json.Unmarshal(b, &configMap)
	if err != nil {
		return nil, err
	}
	for key, val := range configMap {
		_ = os.Setenv(key, val)
	}
	return New(), nil
}
