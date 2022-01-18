package args

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestArgs(t *testing.T) {
	os.Args = []string{"-a", "abcd", "-b", "cde", "-c"}
	assert.Equal(t, Get("-a"), "abcd")
	assert.Equal(t, Get("-b"), "cde")
	assert.Equal(t, Exists("-c"), true)
	assert.Equal(t, Exists("-d"), false)
}
