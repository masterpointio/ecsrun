package cmd

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
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

type ecsClientFake struct {
	client ecsiface.ECSAPI
	config *RunConfig
}

func (c *ecsClientFake) BuildRunTaskInput() *ecs.RunTaskInput {
	return &ecs.RunTaskInput{}
}

func (c *ecsClientFake) RunTask(runTaskInput *ecs.RunTaskInput) (*ecs.RunTaskOutput, error) {
	return &ecs.RunTaskOutput{}, nil
}

func newEcsClientFake(c *RunConfig) ECSClient {
	return &ecsClientFake{}
}

func TestExecute(t *testing.T) {
	setup()
	assert := assert.New(t)

	os.Setenv("AWS_PROFILE", "go-tester")
	os.Setenv("AWS_ACCESS_KEY_ID", "123")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "SECRET123")
	os.Setenv("ECSRUN_CMD", "echo, hello, world")
	os.Setenv("ECSRUN_CLUSTER", "shred")
	os.Setenv("ECSRUN_SECURITY_GROUP", "sg-1")
	os.Setenv("ECSRUN_SUBNET", "public-subnet-1")
	os.Setenv("ECSRUN_TASK", "task")
	os.Setenv("ECSRUN_VERBOSE", "true")

	Execute(newEcsClientFake, VersionInfo{})

	var sesh = viper.Get("session")
	cred, _ := sesh.(*session.Session).Config.Credentials.Get()
	assert.Equal("123", cred.AccessKeyID)

	teardown()
}

func TestVersion(t *testing.T) {
	setup()
	assert := assert.New(t)

	vInfo := VersionInfo{
		Version: "0.1.0",
		Commit:  "238958943",
		Date:    "06/19/20",
		BuiltBy: "MDG",
	}

	result := vInfo.String()

	assert.Contains(result, "0.1.0")
	assert.Contains(result, "238958943")
	assert.Contains(result, "06/19/20")
	assert.Contains(result, "MDG")

	teardown()
}

func TestInitAws(t *testing.T) {
	assert := assert.New(t)
	setup()

	viper.Set("region", "random-region")

	initAws()
	sesh1 := viper.Get("session").(*session.Session)
	assert.Equal(sesh1.Config.Region, aws.String("random-region"))

	viper.Reset()
	os.Setenv("AWS_REGION", "us-west-47")

	initAws()
	sesh2 := viper.Get("session").(*session.Session)
	assert.Equal(sesh2.Config.Region, aws.String("us-west-47"))

	teardown()
}

func TestConfigFile(t *testing.T) {
	assert := assert.New(t)
	setup()

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
