package cmd

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetNormalizedCmd(t *testing.T) {
	assert := assert.New(t)

	expected1 := []*string{}
	viper.Set("cmd", []string{})
	actual1 := getNormalizedCmd()
	assert.Equal(expected1, actual1)

	echo := "echo"
	helloWorld := "hello world"
	expected2 := []*string{&echo, &helloWorld}
	viper.Set("cmd", []string{echo, helloWorld})
	actual2 := getNormalizedCmd()
	assert.Equal(expected2, actual2)
}

func TestGetAssignPublicIP(t *testing.T) {
	assert := assert.New(t)

	expected1 := "ENABLED"
	viper.Set("public", true)
	actual1 := getAssignPublicIP()
	assert.Equal(expected1, actual1)

	expected2 := "DISABLED"
	viper.Set("public", false)
	actual2 := getAssignPublicIP()
	assert.Equal(expected2, actual2)
}

func TestGetContainerName(t *testing.T) {
	assert := assert.New(t)

	expected1 := "task-def-name"
	viper.Set("task", "task-def-name")
	actual1 := getContainerName()
	assert.Equal(expected1, actual1)

	expected2 := "container-name"
	viper.Set("task", "task-def-name")
	viper.Set("name", "container-name")
	actual2 := getContainerName()
	assert.Equal(expected2, actual2)
}

func TestGetTaskDefinition(t *testing.T) {
	assert := assert.New(t)

	expected1 := "task-def-name"
	viper.Set("task", "task-def-name")
	actual1 := getTaskDefinition()
	assert.Equal(expected1, actual1)

	expected2 := "task-def-name"
	viper.Set("task", "task-def-name")
	viper.Set("revision", "")
	actual2 := getTaskDefinition()
	assert.Equal(expected2, actual2)

	expected3 := "task-def-name:5"
	viper.Set("task", "task-def-name")
	viper.Set("revision", "5")
	actual3 := getTaskDefinition()
	assert.Equal(expected3, actual3)
}
