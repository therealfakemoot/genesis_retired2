package genesis

import (
	noise "github.com/therealfakemoot/genesis/noise"
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
func (mg *MapGen) Generate(x float64, y float64) Map {
	m := Map{}
	m.Grid = Grid{X: int(x), Y: int(y), Z: 0}
	points := make([][]float64, int(y))
	for i := 0.0; i < y; i++ {
		row := make([]float64, int(x))
		points[int(i)] = row
		// m.Points[int(i)] = row
	}
	m.Points = points

	for xGen := 0.0; xGen < x; xGen++ {
		for yGen := 0.0; yGen < y; yGen++ {
			m.Points[int(xGen)][int(yGen)] = mg.Noise.Eval3(xGen, yGen, 0)
		}
	}

	return m
}
