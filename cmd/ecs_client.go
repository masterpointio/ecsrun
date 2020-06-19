package cmd

import (
	// "fmt"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
)

// ECSClient is the wrapper around the aws-sdk ECS client and its various structs / methods.
type ECSClient interface {
	BuildRunTaskInput() *ecs.RunTaskInput
	RunTask(runTaskInput *ecs.RunTaskInput) (*ecs.RunTaskOutput, error)
}

type ecsClient struct {
	client ecsiface.ECSAPI
	config *RunConfig
}

// NewEcsClient creates a new ecsClient for the given RunConfig
func NewEcsClient(config *RunConfig) ECSClient {
	client := ecs.New(config.Session)

	return newClient(client, config)
}

func newClient(client ecsiface.ECSAPI, config *RunConfig) ECSClient {
	return &ecsClient{
		client: client,
		config: config,
	}
}

func (c *ecsClient) BuildRunTaskInput() *ecs.RunTaskInput {
	return &ecs.RunTaskInput{
		Cluster:        &c.config.Cluster,
		TaskDefinition: &c.config.TaskDefinition,
		Count:          &c.config.Count,
		LaunchType:     &c.config.LaunchType,
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				AssignPublicIp: &c.config.AssignPublicIP,
				SecurityGroups: []*string{&c.config.SecurityGroupID},
				Subnets:        []*string{&c.config.SubnetID},
			},
		},
		Overrides: &ecs.TaskOverride{
			ContainerOverrides: []*ecs.ContainerOverride{
				{
					Command: c.config.Command,
					Name:    &c.config.ContainerName,
				},
			},
		},
	}
}

func (c *ecsClient) RunTask(runTaskInput *ecs.RunTaskInput) (*ecs.RunTaskOutput, error) {
	output, err := c.client.RunTask(runTaskInput)
	if err != nil {
		log.Fatal("Received error when invoking RunTask.", err)
	}

	return output, err
}
