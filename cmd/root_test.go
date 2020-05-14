package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spf13/viper"
)

var previousProfile string

func setup() {
	previousProfile = os.Getenv("AWS_PROFILE")
	os.Unsetenv("AWS_PROFILE")
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
}

func teardown() {
	os.Setenv("AWS_PROFILE", previousProfile)
	viper.Reset()
}

func TestExecute(t *testing.T) {
	setup()
	assert := assert.New(t)

	os.Setenv("AWS_ACCESS_KEY_ID", "123")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "SECRET123")

	Execute()

	var accessKey = viper.Get("accessKey")
	assert.Equal("123", accessKey)

	var secretKey = viper.Get("secretKey")
	assert.Equal("SECRET123", secretKey)

	teardown()
}

func TestInitAws(t *testing.T) {
	assert := assert.New(t)
	setup()

	// TODO
	t.Skip("TODO")
	assert.Equal(true, true)

	teardown()
}

func TestGetProfile(t *testing.T) {
	assert := assert.New(t)
	setup()

	var profile1 = getProfile()
	assert.Equal("default", profile1)

	os.Setenv("AWS_PROFILE", "not-default")
	var profile2 = getProfile()
	assert.Equal("not-default", profile2)

	teardown()
}
