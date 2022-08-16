package config

import (
	"flag"
	"fmt"
)

//Flags that can be used by user for app configuration
var (
	Studio   = flag.String("studio", "", "studio name to find his games")
	maxConns = flag.Int("max-conns", 1, "max connections for database")
	minConns = flag.Int("min-conns", 1, "max connections for database")
	initDB   = flag.Bool("init-db", false, "init database data")
)

//AppConfig contains configuration parameters defined by user or set by default
type AppConfig struct {
	MaxConns int
	MinConns int
	InitDB   bool
}

//Method Validate() validates set configuration and returns an error if it fails
func (c *AppConfig) Validate() error {
	if c.MaxConns < 1 || c.MaxConns > 50 {
		return fmt.Errorf("Max connections are limited from 1 to 50")
	}
	if c.MinConns < 1 || c.MinConns > 50 {
		return fmt.Errorf("Min connections are limited from 1 to 50")
	}

	return nil
}

// Use method NewAppConfig() to create a new AppConfig
func NewAppConfig() (*AppConfig, error) {
	flag.Parse()
	config := &AppConfig{*maxConns, *minConns, *initDB}
	return config, config.Validate()
}