// Package audit provides config parsing and path checking logic for ppath.
package audit

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/mattn/go-zglob"
)

// Command is a single command section in the config file.
type Command struct {
	Exclude interface{}
	Include interface{}
}

// Precious is an entire config file.
type Precious struct {
	Commands map[string]Command
}

// Paths audits the paths contained in a Precious struct and returns false if
// any of them cannot be found.
func Paths(config *Precious) (bool, error) {
	success := true
	ignoreConfig, err := PpathConfig(".ppath.toml")
	if err != nil {
		return false, err
	}

	for commandName, command := range config.Commands {
		lists := map[string][]string{
			"exclude": patternList(command.Exclude),
			"include": patternList(command.Include),
		}
		for section, list := range lists {
			if len(list) > 0 {
				ok, err := patternsOk(ignoreConfig, commandName, section, list)
				if err != nil {
					return false, err
				}
				if !ok {
					success = false
				}
			}
		}
	}

	return success, nil
}

func patternsOk(ppath *Ppath, commandName, section string, patterns []string) (bool, error) {
	success := true
	for _, pattern := range patterns {
		if patternIgnored(ppath, commandName, pattern) {
			continue
		}
		matches, err := zglob.Glob(pattern)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				success = false
			} else {
				return false, fmt.Errorf("error matching %s pattern '%s' %w", section, pattern, err)
			}
		}

		if len(matches) == 0 {
			success = false
		}
		if !success {
			log.Printf("%s %s pattern %s was not found", commandName, section, pattern)
		}
	}

	return success, nil
}

func patternList(patternGroup interface{}) []string {
	var patterns []string
	if maybeList, ok := patternGroup.([]interface{}); ok {
		for _, pattern := range maybeList {
			str, _ := pattern.(string)
			patterns = append(patterns, str)
		}
	} else if str, ok := patternGroup.(string); ok {
		patterns = append(patterns, str)
	}

	return patterns
}

func patternIgnored(config *Ppath, commandName, pattern string) bool {
	for _, v := range config.Ignore {
		if pattern == v {
			return true
		}
	}

	for name, command := range config.Commands {
		if name != commandName {
			continue
		}
		for _, v := range command.Ignore {
			if pattern == v {
				return true
			}
		}
	}

	return false
}
