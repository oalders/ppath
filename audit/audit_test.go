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

func TestPaths(t *testing.T) {
	config, err := PreciousConfig("testdata/precious.toml")
	assert.NoError(t, err)
	success, err := Paths(config)
	assert.NoError(t, err)
	assert.True(t, success)
}
