package genesis

import (
	"encoding/json"
	"github.com/therealfakemoot/genesis/noise"
)

func flatten(source [][]float64) []float64 {
	r := make([]float64, len(source)*len(source[0]))
	for _, a := range source {
		r = append(r, a...)
	}

	return r
}

// Grid encodes the dimensions of a Map.
type Grid struct {
	X int
	Y int
	Z int
}

// Map is a nested array of float64 values.
type Map struct {
	Grid   Grid
	Points [][]float64
}

// MarshalJSON is used for encoding maps to a JSON payload suitable for use with d3.js .
func (m *Map) MarshalJSON() ([]byte, error) {

	mj := MapJSON{}
	mj.Width = m.Grid.X
	mj.Height = m.Grid.Y
	mj.Values = flatten(m.Points)

	return json.Marshal(mj)
}

// RenderHTML creates an HTML file that displays a contour map of the terrain data.
func (m *Map) RenderHTML(name string) {}

// MapJSON is used for encoding maps to a JSON payload suitable for use with d3.js .
type MapJSON struct {
	Width  int       `json:"width"`
	Height int       `json:"height"`
	Values []float64 `json:"values"`
}

// MapGen will allow for reuse and iterative tweaking of noise generation
// parameters.
type MapGen struct {
	Stretch float64
	Squish  float64
	Noise   genesis.Noise
}

// Generate takes x,y coordinates indicating the maximum dimensions of the
// terrain map to be generated.
func Generate(mg *MapGen, x float64, y float64) Map {
	m := Map{}
	for i := 0.0; i < y; i++ {
		row := make([]float64, int(x))
		m.Points[int(i)] = row
	}

	for xGen := 0.0; xGen < x; xGen++ {
		for yGen := 0.0; yGen < y; yGen++ {
			m.Points[int(xGen)][int(yGen)] = mg.Noise.Eval3(xGen, yGen, 0)
		}
	}

	return m
}
