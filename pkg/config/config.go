// Code generated with https://github.com/mozey/config DO NOT EDIT

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// APP_ADDR
var addr string

// APP_DEV
var dev string

// APP_DIR
var dir string

// APP_PROXY
var proxy string

// Config fields correspond to config file keys less the prefix
type Config struct {
	addr  string // APP_ADDR
	dev   string // APP_DEV
	dir   string // APP_DIR
	proxy string // APP_PROXY
}

// Addr is APP_ADDR
func (c *Config) Addr() string {
	return c.addr
}

// Dev is APP_DEV
func (c *Config) Dev() string {
	return c.dev
}

// Dir is APP_DIR
func (c *Config) Dir() string {
	return c.dir
}

// Proxy is APP_PROXY
func (c *Config) Proxy() string {
	return c.proxy
}

// SetAddr overrides the value of addr
func (c *Config) SetAddr(v string) {
	c.addr = v
}

// SetDev overrides the value of dev
func (c *Config) SetDev(v string) {
	c.dev = v
}

// SetDir overrides the value of dir
func (c *Config) SetDir(v string) {
	c.dir = v
}

// SetProxy overrides the value of proxy
func (c *Config) SetProxy(v string) {
	c.proxy = v
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

	if dev != "" {
		conf.dev = dev
	}

	if dir != "" {
		conf.dir = dir
	}

	if proxy != "" {
		conf.proxy = proxy
	}

}

// SetEnv sets non-empty env vars on Config
func SetEnv(conf *Config) {
	var v string

	v = os.Getenv("APP_ADDR")
	if v != "" {
		conf.addr = v
	}

	v = os.Getenv("APP_DEV")
	if v != "" {
		conf.dev = v
	}

	v = os.Getenv("APP_DIR")
	if v != "" {
		conf.dir = v
	}

	v = os.Getenv("APP_PROXY")
	if v != "" {
		conf.proxy = v
	}

}

// LoadFile sets the env from file and returns a new instance of Config
func LoadFile(mode string) (conf *Config, err error) {
	appDir := os.Getenv("APP_DIR")
	p := fmt.Sprintf("%v/config.%v.json", appDir, mode)
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
