package genesis

import (
	"encoding/json"
	"fmt"
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

// Map describes the topographical layout of a map.
type Map struct {
	Grid   Grid
	Points [][]float64
}

// MarshalJSON is used for encoding maps to a JSON payload suitable for use with d3.js .
func (m Map) MarshalJSON() ([]byte, error) {

	mj := MapJSON{}
	mj.Width = m.Grid.X
	mj.Height = m.Grid.Y
	mj.Values = flatten(m.Points)

	return json.Marshal(mj)
}

func (m Map) String() string {

	s := ""

	grid, _ := json.Marshal(m.Grid)
	s += string(grid)
	s += "\n"
	for _, r := range m.Points {
		s += fmt.Sprintf("%v", r) + "\n"
	}

	return s
}

// RenderHTML creates an HTML file that displays a contour map of the terrain data.
func (m *Map) RenderHTML(name string) {}

// MapJSON is used for encoding maps to a JSON payload suitable for use with d3.js .
type MapJSON struct {
	Width  int       `json:"width"`
	Height int       `json:"height"`
	Values []float64 `json:"values"`
}
