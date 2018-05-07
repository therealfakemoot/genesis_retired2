package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	genesis "github.com/therealfakemoot/genesis/app"
	l "github.com/therealfakemoot/genesis/log"
	"os"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "genesis",
	Short: "A procedural world generation toolkit",
	Long: `Genesis is an interactive tool for creating, modifying,
rendering, and exporting maps containing extensible Features and
generation parameters.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlags(cmd.Flags())

		viper.BindPFlag("verbose", cmd.Flags().Lookup("verbose"))

		if viper.GetBool("verbose") {
			l.Term.SetLevel(logrus.DebugLevel)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
	PostRun: func(cmd *cobra.Command, args []string) {
		for _, key := range viper.AllKeys() {
			l.Term.WithFields(logrus.Fields{
				key: viper.Get(key),
			}).Debug("Parameter found")
		}

		genesis.DumpSettings(l.Term)
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		l.Term.WithError(err).Info("genesis failed to start")
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.Flags().String("config", "", "config file (default is $HOME/.genesis.yaml)")
	RootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose logging. [POSSIBLE PERFORMANCE IMPLICATIONS]")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	genesis.LoadConfig(l.Term)
}
