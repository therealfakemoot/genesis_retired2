package genesis

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	l "github.com/therealfakemoot/genesis/log"
	noise "github.com/therealfakemoot/genesis/noise"
	Q "github.com/therealfakemoot/go-quantize"
)

// MapGen will allow for reuse and iterative tweaking of noise generation
// parameters.
type MapGen struct {
	Stretch float64
	Squish  float64
	Noise   *noise.Noise
}

// Generate takes x,y coordinates indicating the maximum dimensions of the
// terrain map to be generated.
func (mg *MapGen) Generate(x, y, sampleScale, thresholdScale float64) Map {
	m := Map{}
	m.Grid = Grid{X: int(x), Y: int(y), Z: 0}
	points := make([][]float64, int(y))

	d := Q.Domain{
		Max:  viper.GetFloat64("Terrain.Domain.Max"),
		Min:  viper.GetFloat64("Terrain.Domain.Min"),
		Step: viper.GetFloat64("Terrain.Domain.Step"),
	}

	l.Term.WithFields(logrus.Fields{
		"domain": fmt.Sprintf("%+v", d),
	}).Debug("Domain")

	for yGen := 0.0; yGen < y; yGen++ {
		row := make([]float64, int(x))
		for xGen := 0.0; xGen < x; xGen++ {
			row[int(xGen)] = mg.Noise.Eval3(xGen*sampleScale, yGen*sampleScale, 0)
		}
		quantized := d.Quantize(row)

		points[int(yGen)] = quantized

		l.Term.WithFields(logrus.Fields{
			"Raw Row":       row,
			"Quantized Row": quantized,
		}).Debug(fmt.Sprintf("Row %0.f", yGen))
	}

	m.Points = points

	return m
}
