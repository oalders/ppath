package audit

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	if err := os.Chdir(dir); err != nil {
		panic(err)
	}
}

func TestPpathConfigDoesNotExist(t *testing.T) {
	c, err := PpathConfig("testdata/.noppath.toml")
	assert.NoError(t, err)

	var empty []string
	assert.Equal(t, empty, c.Ignore)
	assert.Equal(
		t,
		empty,
		c.Commands["omegasort-gitignore"].Ignore,
	)
}

func TestPpathConfig(t *testing.T) {
	c, err := PpathConfig("testdata/.ppath.toml")
	assert.NoError(t, err)

	assert.Equal(t, []string{`**/node_modules/**/*`}, c.Ignore)
	assert.Equal(
		t,
		[]string{`**/foo/**/*`},
		c.Commands["omegasort-gitignore"].Ignore,
	)
}

func TestPreciousConfigDoesNotExist(t *testing.T) {
	c, err := PreciousConfig("testdata/precious-not-found.toml")
	assert.Error(t, err)
	assert.Empty(t, c)
}

func TestPreciousFailConfig(t *testing.T) {
	c, err := PreciousConfig("testdata/precious-fail.toml")
	assert.NoError(t, err)

	assert.Equal(t,
		Command{
			Exclude: "baz",
			Include: `**/.gitignore`,
		},
		c.Commands["omegasort-gitignore"],
	)

	assert.Equal(t,
		Command{
			Exclude: []interface{}{"foo", "bar"},
			Include: []interface{}{`**/*.go`},
		},
		c.Commands["golangci-lint"],
	)
}

func TestPreciousSuccessConfig(t *testing.T) {
	c, err := PreciousConfig("testdata/precious.toml")
	assert.NoError(t, err)

	assert.Equal(t,
		Command{
			Exclude: nil,
			Include: `**/.gitignore`,
		},
		c.Commands["omegasort-gitignore"],
	)

	assert.Equal(t,
		Command{
			Exclude: nil,
			Include: []interface{}{`**/*.go`},
		},
		c.Commands["golangci-lint"],
	)
}
