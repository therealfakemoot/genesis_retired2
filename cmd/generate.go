package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	noise "github.com/therealfakemoot/genesis/noise"
	terrain "github.com/therealfakemoot/genesis/terrain"
	"io/ioutil"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		n := noise.NewWithSeed(4074849863)
		mg := terrain.MapGen{
			Stretch: -1.0 / 6,
			Squish:  1 / 3,
			Noise:   n,
		}

		w := float64(viper.GetInt("width"))
		h := float64(viper.GetInt("height"))

		terrainMap := mg.Generate(w, h)

		jsonBytes, _ := json.Marshal(terrainMap)

		o := viper.GetString("output")
		ioutil.WriteFile(o+".json", jsonBytes, 0644)
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(terrainCmd)

	generateCmd.Flags().StringP("output", "o", "", "Name of output file")
	generateCmd.Flags().Int("width", 1000, "Horizontal width of generated map")
	generateCmd.Flags().Int("height", 1000, "Vertical height of generated map")

	viper.BindPFlags(generateCmd.Flags())
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
