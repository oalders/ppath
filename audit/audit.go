// Package audit provides config parsing and path checking logic for ppath.
package audit

import (
	"log"
	"os"

	"github.com/mattn/go-zglob"
	"github.com/pkg/errors"
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

type matchCache map[string]bool

// Paths audits the paths contained in a Precious struct and returns false if
// any of them cannot be found.
func Paths(config *Precious) (bool, error) {
	success := true
	ignoreConfig, err := PpathConfig(".ppath.toml")
	if err != nil {
		return false, err
	}

	seen := make(matchCache)
	for commandName, command := range config.Commands {
		lists := map[string][]string{
			"exclude": patternList(command.Exclude),
			"include": patternList(command.Include),
		}
		for section, list := range lists {
			if len(list) > 0 {
				ok, err := patternsOk(seen, ignoreConfig, commandName, section, list)
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

func patternsOk(seen matchCache, ppath *Ppath, commandName, section string, patterns []string) (bool, error) {
	success := true
	for _, pattern := range patterns {
		matched, exists := seen[pattern]

		if exists && matched {
			continue
		}

		// For our purposes found and ignored are the same thing.
		if patternIgnored(ppath, commandName, pattern) {
			seen[pattern] = true
			continue
		}

		if !exists {
			matches, err := zglob.Glob(pattern)

			if err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return false, errors.Wrapf(
						err,
						"matching %s pattern '%s'",
						section,
						pattern,
					)
				}
			} else if len(matches) > 0 {
				seen[pattern] = true
				continue
			}
		}

		seen[pattern] = false
		success = false
		log.Printf("%s %s pattern %s was not found", commandName, section, pattern)
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
