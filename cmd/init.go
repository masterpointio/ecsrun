package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// InitCmd is responsible for creating a `ecsrun.yml`
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a blank `ecsrun.yml` config file in the current directory.",
	Run: func(cmd *cobra.Command, args []string) {
		// Reset viper as it carries over config from root. We want to build our own config.
		viper.Reset()

		initCmd()
	},
}

func initCmd() {

	config := make(map[string]interface{})
	config["default"] = map[string]interface{}{
		"cluster":        "TODO",
		"task":           "TODO",
		"security-group": "TODO",
		"subnet":         "TODO",
		"cmd":            []string{"bash", "-c", "echo", "hello", "world"},
	}

	if err := viper.MergeConfigMap(config); err != nil {
		log.Fatal(err)
	}

	// Write the file.
	if err := viper.SafeWriteConfigAs("./ecsrun.yaml"); err != nil {
		log.Fatal(err)
		red := color.New(color.FgRed)
		red.Printf("ecsrun.yaml already exists.\n")
	}
}
