package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	assert := assert.New(t)
	Execute()
	_ = assert
}

func TestInitAws(t *testing.T) {
	assert := assert.New(t)
}

func TestGetProfile(t *testing.T) {
	assert := assert.New(t)

	profile = getProfile()
	assert.Equal(profile, "default")

	os.Setenv("AWS_PROFILE", "not-default")
	profile = getProfile()
	assert.Equal(profile, "not-default")
}
