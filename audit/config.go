// Package audit provides some Ppath config parsing.
package audit

import (
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
)

// Ppath is an entire config file.
type Ppath struct {
	Ignore   []string
	Commands map[string]struct {
		Ignore []string
	}
}

// PpathConfig parses a precious.toml and returns a Ppath struct with the bits
// we care about. Currently it will return an empty config if a config file
// cannot be found.
func PpathConfig(filename string) (*Ppath, error) {
	var config Ppath
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return &config, nil
		}
		return nil, errors.Wrapf(err, "cannot stat %s", filename)
	}
	dat, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"reading ppath config file %s",
			filename,
		)
	}
	if err := toml.Unmarshal(dat, &config); err != nil {
		errors.Wrapf(err, "unmarshal toml in %s", filename)
	}

	return &config, nil
}

// PreciousConfig parses a precious.toml and returns a Precious struct
// with the bits we care about.
func PreciousConfig(filename string) (*Precious, error) {
	dat, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(
			err,
			"reading precious config file %s",
			filename,
		)
	}
	var config Precious
	if err := toml.Unmarshal(dat, &config); err != nil {
		log.Fatalf("unmarshal %v", err)
	}

	return &config, nil
}
