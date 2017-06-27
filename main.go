package main

import (
	"fmt"
	"math/rand"

	"github.com/tbogdala/noisey"
)

func main() {
	r := rand.New(rand.NewSource(int64(1)))

	noiseGen := noisey.NewPerlinGenerator(r)

	v := noiseGen.Get3D(0.4, 0.2, 0.0)
	fmt.Println(v)
}
