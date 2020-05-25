package cmd

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/sirupsen/logrus"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var awsSession *session.Session

var log = logrus.New()

var rootCmd = &cobra.Command{
	Use:   "escrun",
	Short: "Easily run one-off tasks against a ECS Cluster",
	Long: `
ecsrun is a CLI tool that allows users to run one-off administrative tasks
using their existing Task Definitions.`,

	Run: func(cmd *cobra.Command, args []string) {
		config := BuildRunConfig(awsSession)

		ecsClient := NewEcsClient(config)

		input := ecsClient.BuildRunTaskInput()
		output, err := ecsClient.RunTask(input)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("RunTask output: ", output)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initVerbose, initAws)

	log.SetOutput(os.Stderr)

	// Basic Flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.escrun)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	// AWS Cred / Environment Flags
	rootCmd.PersistentFlags().String("cred", "", "AWS credentials file (default is $HOME/.aws/.credentials)")
	rootCmd.PersistentFlags().StringP("profile", "p", "", "AWS profile to target (default is AWS_PROFILE or 'default')")
	rootCmd.PersistentFlags().String("region", "", `AWS region to target (default is AWS_REGION or pulled from $HOME/.aws/.credentials)`)

	// Task Flags
	rootCmd.PersistentFlags().StringP("cluster", "c", "", "The ECS Cluster to run the task in.")
	rootCmd.PersistentFlags().StringP("task", "t", "", "The name of the ECS Task Definition to use.")
	rootCmd.PersistentFlags().StringP("revision", "r", "", "The Task Definition revision to use.")
	rootCmd.PersistentFlags().StringP("name", "n", "", "The name of the container in the Task Definition.")
	rootCmd.PersistentFlags().StringP("launch-type", "l", "FARGATE", "The launch type to run as. Currently only Fargate is supported.")
	rootCmd.PersistentFlags().StringSlice("cmd", []string{}, "The comma separated command override to apply.")
	rootCmd.PersistentFlags().Int64("count", 1, "The number of tasks to launch for the given cmd.")

	// Network Flags
	rootCmd.PersistentFlags().StringP("subnet", "s", "", "The Subnet ID that the task should be launched in.")
	rootCmd.PersistentFlags().StringP("security-group", "g", "", "The Security Group ID that the task should be associated with.")
	rootCmd.PersistentFlags().Bool("public", false, "Assigns a public IP to the task if included. (default is false)")

	// Require specific flags
	rootCmd.MarkPersistentFlagRequired("cluster")
	rootCmd.MarkPersistentFlagRequired("task")
	rootCmd.MarkPersistentFlagRequired("cmd")
	rootCmd.MarkPersistentFlagRequired("subnet")
	rootCmd.MarkPersistentFlagRequired("security-group")

	// Bind em All
	viper.BindPFlags(rootCmd.PersistentFlags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		// Search config in home directory with name ".escrun" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".escrun")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	}
}

func initVerbose() {
	verbose, err := rootCmd.PersistentFlags().GetBool("verbose")
	if err != nil {
		log.Fatal("Unable to pull verbose flag.", err)
	}

	if verbose {
		log.Info("Enabling verbose output.")
		log.SetLevel(logrus.DebugLevel)
	}
}

func initAws() {
	profile := getProfile()
	viper.Set("profile", profile)

	// Create our AWS session object for AWS API Usage
	sesh, err := initAwsSession(profile)
	if err != nil {
		log.Fatal("Unable to init AWS Session. Check your credentials and profile.", err)
	}

	region := viper.GetString("region")
	if region == "" {
		region = *sesh.Config.Region
	}

	// Override our Session's region in case it was set.
	sesh.Config.WithRegion(region)

	// Set our awsSession for later use.
	// TODO: What's the proper way to do this... This seems weird.
	awsSession = sesh
}

func getProfile() string {
	var profile = viper.GetString("profile")
	if profile == "" {
		profile = "default"
		if os.Getenv("AWS_PROFILE") != "" {
			profile = os.Getenv("AWS_PROFILE")
		}
	}

	log.Debug("Using AWS Profile: " + profile)
	return profile
}

func initAwsSession(profile string) (*session.Session, error) {
	credFile, err := rootCmd.PersistentFlags().GetString("cred")
	if err != nil {
		log.Fatal("Not able to get credFile from cmd.", err)
	}

	log.Debug("Cred File: " + credFile)

	var sesh *session.Session

	if credFile != "" {
		sesh, err = session.NewSession(&aws.Config{
			Credentials: credentials.NewSharedCredentials(credFile, profile),
		})
	} else {
		sesh, err = session.NewSessionWithOptions(session.Options{
			Profile:           profile,
			SharedConfigState: session.SharedConfigEnable,
			Config: aws.Config{
				CredentialsChainVerboseErrors: aws.Bool(true),
			},
		})
	}

	return sesh, err
}
