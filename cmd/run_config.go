package cmd

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/spf13/viper"
)

// RunConfig is the main config object used to configure the RunTask.
type RunConfig struct {
	Command                []*string
	Cluster                string
	TaskDefinition         string
	TaskDefinitionName     string
	TaskDefinitionRevision string
	ContainerName          string
	LaunchType             string
	Count                  int64

	SubnetID           string
	SecurityGroupID    string
	AssignPublicIPFlag bool
	AssignPublicIP     string

	Session *session.Session
}

// BuildRunConfig constructs the our primary RunConfig object using the given
// AWS session and the CLI args from viper.
func BuildRunConfig(session *session.Session) *RunConfig {

	// Convert our cmd slice to a slice of pointers

	cmd := getNormalizedCmd()
	taskDef := getTaskDefinition()
	name := getContainerName()
	assignPublicIP := getAssignPublicIP()

	return &RunConfig{
		Command:                cmd,
		Cluster:                viper.GetString("cluster"),
		TaskDefinition:         taskDef,
		TaskDefinitionName:     viper.GetString("task"),
		TaskDefinitionRevision: viper.GetString("revision"),
		ContainerName:          name,
		LaunchType:             viper.GetString("launch-type"),
		Count:                  viper.GetInt64("count"),
		SubnetID:               viper.GetString("subnet"),
		SecurityGroupID:        viper.GetString("security-group"),
		AssignPublicIPFlag:     viper.GetBool("public"),
		AssignPublicIP:         assignPublicIP,
		Session:                session,
	}
}

func getNormalizedCmd() []*string {
	result := []*string{}
	for _, v := range viper.GetStringSlice("cmd") {
		result = append(result, &v)
	}

	return result
}

func getTaskDefinition() string {
	if viper.GetString("revision") != "" {
		return viper.GetString("task") + ":" + viper.GetString("revision")
	}

	return viper.GetString("task")
}

func getContainerName() string {
	if viper.GetString("name") != "" {
		return viper.GetString("name")
	}

	return viper.GetString("task")
}

func getAssignPublicIP() string {
	if viper.GetBool("public") {
		return ecs.AssignPublicIpEnabled
	}

	return ecs.AssignPublicIpDisabled
}
