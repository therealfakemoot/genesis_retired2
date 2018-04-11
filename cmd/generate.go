package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	l "github.com/therealfakemoot/genesis/log"
	terrain "github.com/therealfakemoot/genesis/map/terrain"
	noise "github.com/therealfakemoot/genesis/noise"
	"os"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:       "generate",
	Short:     "A brief description of your command",
	ValidArgs: []string{"all", "terrain", "feature"},
	Args:      cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			l.Term.Error("Please provide a valid layer.")
		}

		layer := args[0]

		switch layer {
		case "all":
			l.Term.Info("Full-map generation not implemented.")
		case "terrain":
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

			err := os.Mkdir(o, 0755)

			if err != nil {
				l.Term.WithError(err).Error("Failed to create map directory.")
			}

			terrainJSONFile, err := os.OpenFile(o+"/terrain.json", os.O_RDWR|os.O_CREATE, 0644)
			defer terrainJSONFile.Close()

			if err != nil {
				l.Term.WithError(err).Error("Failed to open " + o + "/terrain.json")
			}

			terrainJSONFile.Write(jsonBytes)

			terrainHTMLFile, err := os.OpenFile(o+"/terrain.html", os.O_RDWR|os.O_CREATE, 0644)

			defer terrainHTMLFile.Close()

			if err != nil {
				l.Term.WithError(err).Error("Failed to open " + o + "/terrain.json")
			}

			terrain.RenderTopoHTML(terrainHTMLFile)

		case "feature":
			l.Term.Info("Feature generation not implemented.")
		}
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)

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
