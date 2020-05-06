/*
Copyright © 2020 Matt Gowie <matt@masterpoint.io>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"testing"

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
using their existing Task Definitions.
TODO: Supply more info here.`,

	Run: func(cmd *cobra.Command, args []string) {
		log.Info("SHRED!")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initVerbose, initAws)

	log.SetOutput(os.Stderr)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.escrun.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringP("cred", "c", "", "aws credentials file (default is $HOME/.aws/.credentials)")
	rootCmd.PersistentFlags().StringP("profile", "p", "", "aws profile to target (default is AWS_PROFILE or 'default')")
	rootCmd.PersistentFlags().StringP("region", "r", "", `aws region to target (default is AWS_REGION or pulled from $HOME/.aws/.credentials)`)

	viper.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile"))
	viper.BindPFlag("region", rootCmd.PersistentFlags().Lookup("region"))
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
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".escrun" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".escrun")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initVerbose() {
	verbose, err := rootCmd.PersistentFlags().GetBool("verbose")
	if err != nil {
		log.Fatal("Unable to pull verbose flag.")
		log.Fatal(err)
		os.Exit(1)
	}

	if verbose {
		log.Info("Enabling verbose output.")
		log.SetLevel(logrus.DebugLevel)
	}
}

func initAws() {
	profile := getProfile(nil)

	// Create our AWS session object for AWS API Usage
	sesh, err := initAwsSession(profile)
	if err != nil {
		log.Fatal("Unable to init AWS Session. Check your credentials and profile.")
		log.Fatal(err)
		os.Exit(1)
	}

	cred, err := sesh.Config.Credentials.Get()
	if err != nil {
		log.Fatal("Unable to get credentials from Session. Check your credentials and profile.")
		log.Fatal(err)
		os.Exit(1)
	}

	region := viper.GetString("region")
	if region == "" {
		region = *sesh.Config.Region
	}
	// Override our Session's region in case it was set.
	sesh.Config.WithRegion(region)

	viper.Set("profile", profile)
	viper.Set("region", region)
	viper.Set("accesskey", cred.AccessKeyID)
	viper.Set("secretkey", cred.SecretAccessKey)
}

func getProfile(t *testing.T) string {
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
		log.Fatal("Not able to get credFile from cmd.")
		log.Fatal(err)
		os.Exit(1)
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
