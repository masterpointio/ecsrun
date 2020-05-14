package cmd

import (
	// "fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/ecs"
)

// ECSClient is the wrapper around the aws-sdk ECS client and its various structs / methods.
type ECSClient interface {
	BuildRunTaskInput() (*ecs.RunTaskInput, error)
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

func (c *ecsClient) BuildRunTaskInput() (*ecs.RunTaskInput, error) {

	taskDefition := c.getTaskDefinition()
	assignPublicIP := c.getAssignPublicIp()

	runInput := &ecs.RunTaskInput{
		Cluster:        &c.config.Cluster,
		TaskDefinition: &taskDefition,
		Count:          &c.config.Count,
		LaunchType:     &c.config.LaunchType,
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				AssignPublicIp: &assignPublicIP,
				SecurityGroups: []*string{&c.config.SecurityGroupID},
				Subnets:        []*string{&c.config.SubnetID},
			},
		},
		Overrides: &ecs.TaskOverride{
			ContainerOverrides: []*ecs.ContainerOverride{
				{
					Command: c.config.Command,
					Name:    &def,
				},
			},
		},
	}

	return runInput, nil
}

func (c *ecsClient) RunTask(runTaskInput *ecs.RunTaskInput) (*ecs.RunTaskOutput, error) {

	output, err := client.RunTask(runInput)
	if err != nil {
		log.Fatal("Received error when invoking RunTask.", err)
		log.Fatal("Error: ", err)
		os.Exit(1)
	}

	log.Info("Output: ", output)
	return output, err
}

func (c *ecsClient) getTaskDefinition() string {
	if c.config.TaskDefinitionRevision != nil {
		return c.config.TaskDefinitionName + ":" + c.config.TaskDefinitionRevision
	}

	return c.config.TaskDefinitionName
}

func (c *ecsClient) getAssignPublicIP() string {
	if c.config.AssignPublicIP {
		return ecs.AssignPublicIpEnabled
	}

	return ecs.AssignPublicIpDisabled
}