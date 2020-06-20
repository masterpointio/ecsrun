package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var testFs = afero.NewMemMapFs()

// Tests
/////////

func TestInitCmd(t *testing.T) {
	assert := assert.New(t)

	viper.SetFs(testFs)

	exists, err := afero.Exists(testFs, "./ecsrun.yaml")
	if err != nil {
		panic(err)
	}
	assert.False(exists)

	initCmd()

	exists, err = afero.Exists(testFs, "./ecsrun.yaml")
	if err != nil {
		panic(err)
	}
	assert.True(exists)

	testFs.Remove("./ecsrun.yaml")
}
