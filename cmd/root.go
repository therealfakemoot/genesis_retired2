package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	genesis "github.com/therealfakemoot/genesis/app"
	"os"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "genesis",
	Short: "A procedural world generation toolkit",
	Long: `Genesis is an interactive tool for creating, modifying,
rendering, and exporting maps containing extensible Features and
generation parameters.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.genesis.yaml)")
	RootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose logging. [POSSIBLE PERFORMANCE IMPLICATIONS]")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	l := logrus.New()
	l.Out = os.Stdout

	genesis.LoadConfig(l)
	//logCtx := log.WithFields(log.Fields{
	//"name": "butts",
	//})

	l.WithFields(logrus.Fields{
		"Settings": viper.AllSettings(),
	}).Info("All settings")

	viper.Set("logger", l)
}
