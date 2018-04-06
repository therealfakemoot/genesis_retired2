package genesis

import (
	"math"
)

/**
 * OpenSimplex Noise in Go.
 * algorithm by Kurt Spencer
 * ported by Owen Raccuglia
 *
 * Based on Java v1.1 (October 5, 2014)
 */

// I'm not sure what these magic numbers are, but they will need
// to be parameterized so users can control how smooth/rough their maps are.
var (
	stretchConstant3D = -1.0 / 6 // (1/Math.sqrt(3+1)-1)/3
	squishConstant3D  = 1.0 / 3  // (Math.sqrt(3+1)-1)/3
)

var normConstant3D float64 = 103

var defaultSeed int64

// Noise is a seeded Noise instance. Reusing a Noise instance (rather than recreating it
// from a known seed) will save some calculation time.
type Noise struct {
	perm            []int16
	permGradIndex3D []int16
}

// New returns a Noise instance with a seed of 0.
func New() *Noise {
	return NewWithSeed(defaultSeed)
}

// NewWithSeed returns a Noise instance with a 64-bit seed. Two Noise instances with the
// same seed will have the same output.
func NewWithSeed(seed int64) *Noise {
	s := Noise{
		perm:            make([]int16, 256),
		permGradIndex3D: make([]int16, 256),
	}

	source := make([]int16, 256)
	for i := range source {
		source[i] = int16(i)
	}

	seed = seed*6364136223846793005 + 1442695040888963407
	seed = seed*6364136223846793005 + 1442695040888963407
	seed = seed*6364136223846793005 + 1442695040888963407
	for i := int32(255); i >= 0; i-- {
		seed = seed*6364136223846793005 + 1442695040888963407
		r := int32((seed + 31) % int64(i+1))
		if r < 0 {
			r += i + 1
		}

		s.perm[i] = source[r]
		s.permGradIndex3D[i] = (s.perm[i] % (int16(len(gradients3D)) / 3)) * 3
		source[r] = source[i]
	}

	return &s
}

// NewWithPerm returns a Noise instance with a specific internal permutation state.
// If you're not sure about this, you probably want NewWithSeed().
func NewWithPerm(perm []int16) *Noise {
	s := Noise{
		perm:            perm,
		permGradIndex3D: make([]int16, 256),
	}

	for i, p := range perm {
		// Since 3D has 24 gradients, simple bitmask won't work, so precompute modulo array.
		s.permGradIndex3D[i] = (p % (int16(len(gradients3D)) / 3)) % 3
	}

	return &s
}

// Eval3 returns a random noise value in three dimensions.
func (s *Noise) Eval3(x, y, z float64) float64 {
	// Place input coordinates on simplectic honeycomb.
	stretchOffset := (x + y + z) * stretchConstant3D
	xs := float64(x + stretchOffset)
	ys := float64(y + stretchOffset)
	zs := float64(z + stretchOffset)

	// Floor to get simplectic honeycomb coordinates of rhombohedron (stretched cube) super-cell origin.
	xsb := int32(math.Floor(xs))
	ysb := int32(math.Floor(ys))
	zsb := int32(math.Floor(zs))

	// Skew out to get actual coordinates of rhombohedron origin. We'll need these later.
	squishOffset := float64(xsb+ysb+zsb) * squishConstant3D
	xb := float64(xsb) + squishOffset
	yb := float64(ysb) + squishOffset
	zb := float64(zsb) + squishOffset

	// Compute simplectic honeycomb coordinates relative to rhombohedral origin.
	xins := xs - float64(xsb)
	yins := ys - float64(ysb)
	zins := zs - float64(zsb)

	// Sum those together to get a value that determines which region we're in.
	inSum := xins + yins + zins

	// Positions relative to origin point.
	dx0 := x - xb
	dy0 := y - yb
	dz0 := z - zb

	// We'll be defining these inside the next block and using them afterwards.
	var dxExt0, dyExt0, dzExt0 float64
	var dxExt1, dyExt1, dzExt1 float64
	var xsvExt0, ysvExt0, zsvExt0 int32
	var xsvExt1, ysvExt1, zsvExt1 int32

	value := float64(0)
	if inSum <= 1 { // We're inside the tetrahedron (3-Simplex) at (0,0,0)

		// Determine which two of (0,0,1), (0,1,0), (1,0,0) are closest.
		aPoint := byte(0x01)
		bPoint := byte(0x02)
		aScore := xins
		bScore := yins
		if aScore >= bScore && zins > bScore {
			bScore = zins
			bPoint = 0x04
		} else if aScore < bScore && zins > aScore {
			aScore = zins
			aPoint = 0x04
		}

		// Now we determine the two lattice points not part of the tetrahedron that may contribute.
		// This depends on the closest two tetrahedral vertices, including (0,0,0)
		wins := 1 - inSum
		if wins > aScore || wins > bScore { // (0,0,0) is one of the closest two tetrahedral vertices.
			var c byte // Our other closest vertex is the closest out of a and b.
			if bScore > aScore {
				c = bPoint
			} else {
				c = aPoint
			}

			if (c & 0x01) == 0 {
				xsvExt0 = xsb - 1
				xsvExt1 = xsb
				dxExt0 = dx0 + 1
				dxExt1 = dx0
			} else {
				xsvExt1 = xsb + 1
				xsvExt0 = xsvExt1
				dxExt1 = dx0 - 1
				dxExt0 = dxExt1
			}

			if (c & 0x02) == 0 {
				ysvExt1 = ysb
				ysvExt0 = ysvExt1
				dyExt1 = dy0
				dyExt0 = dyExt1
				if (c & 0x01) == 0 {
					ysvExt1--
					dyExt1++
				} else {
					ysvExt0--
					dyExt0++
				}
			} else {
				ysvExt1 = ysb + 1
				ysvExt0 = ysvExt1
				dyExt1 = dy0 - 1
				dyExt0 = dyExt1
			}

			if (c & 0x04) == 0 {
				zsvExt0 = zsb
				zsvExt1 = zsb - 1
				dzExt0 = dz0
				dzExt1 = dz0 + 1
			} else {
				zsvExt1 = zsb + 1
				zsvExt0 = zsvExt1
				dzExt1 = dz0 - 1
				dzExt0 = dzExt1
			}
		} else { // (0,0,0) is not one of the closest two tetrahedral vertices.
			c := aPoint | bPoint // Our two extra vertices are determined by the closest two.

			if (c & 0x01) == 0 {
				xsvExt0 = xsb
				xsvExt1 = xsb - 1
				dxExt0 = dx0 - 2*squishConstant3D
				dxExt1 = dx0 + 1 - squishConstant3D
			} else {
				xsvExt1 = xsb + 1
				xsvExt0 = xsvExt1
				dxExt0 = dx0 - 1 - 2*squishConstant3D
				dxExt1 = dx0 - 1 - squishConstant3D
			}

			if (c & 0x02) == 0 {
				ysvExt0 = ysb
				ysvExt1 = ysb - 1
				dyExt0 = dy0 - 2*squishConstant3D
				dyExt1 = dy0 + 1 - squishConstant3D
			} else {
				ysvExt1 = ysb + 1
				ysvExt0 = ysvExt1
				dyExt0 = dy0 - 1 - 2*squishConstant3D
				dyExt1 = dy0 - 1 - squishConstant3D
			}

			if (c & 0x04) == 0 {
				zsvExt0 = zsb
				zsvExt1 = zsb - 1
				dzExt0 = dz0 - 2*squishConstant3D
				dzExt1 = dz0 + 1 - squishConstant3D
			} else {
				zsvExt1 = zsb + 1
				zsvExt0 = zsvExt1
				dzExt0 = dz0 - 1 - 2*squishConstant3D
				dzExt1 = dz0 - 1 - squishConstant3D
			}
		}

		// Contribution (0,0,0)
		attn0 := 2 - dx0*dx0 - dy0*dy0 - dz0*dz0
		if attn0 > 0 {
			attn0 *= attn0
			value += attn0 * attn0 * s.extrapolate3(xsb+0, ysb+0, zsb+0, dx0, dy0, dz0)
		}

		// Contribution (1,0,0)
		dx1 := dx0 - 1 - squishConstant3D
		dy1 := dy0 - 0 - squishConstant3D
		dz1 := dz0 - 0 - squishConstant3D
		attn1 := 2 - dx1*dx1 - dy1*dy1 - dz1*dz1
		if attn1 > 0 {
			attn1 *= attn1
			value += attn1 * attn1 * s.extrapolate3(xsb+1, ysb+0, zsb+0, dx1, dy1, dz1)
		}

		// Contribution (0,1,0)
		dx2 := dx0 - 0 - squishConstant3D
		dy2 := dy0 - 1 - squishConstant3D
		dz2 := dz1
		attn2 := 2 - dx2*dx2 - dy2*dy2 - dz2*dz2
		if attn2 > 0 {
			attn2 *= attn2
			value += attn2 * attn2 * s.extrapolate3(xsb+0, ysb+1, zsb+0, dx2, dy2, dz2)
		}

		// Contribution (0,0,1)
		dx3 := dx2
		dy3 := dy1
		dz3 := dz0 - 1 - squishConstant3D
		attn3 := 2 - dx3*dx3 - dy3*dy3 - dz3*dz3
		if attn3 > 0 {
			attn3 *= attn3
			value += attn3 * attn3 * s.extrapolate3(xsb+0, ysb+0, zsb+1, dx3, dy3, dz3)
		}
	} else if inSum >= 2 { // We're inside the tetrahedron (3-Simplex) at (1,1,1)

		// Determine which two tetrahedral vertices are the closest, out of (1,1,0), (1,0,1), (0,1,1) but not (1,1,1).
		aPoint := byte(0x06)
		aScore := xins
		bPoint := byte(0x05)
		bScore := yins
		if aScore <= bScore && zins < bScore {
			bScore = zins
			bPoint = 0x03
		} else if aScore > bScore && zins < aScore {
			aScore = zins
			aPoint = 0x03
		}

		// Now we determine the two lattice points not part of the tetrahedron that may contribute.
		// This depends on the closest two tetrahedral vertices, including (1,1,1)
		wins := 3 - inSum
		if wins < aScore || wins < bScore { // (1,1,1) is one of the closest two tetrahedral vertices.
			var c byte // Our other closest vertex is the closest out of a and b.
			if bScore < aScore {
				c = bPoint
			} else {
				c = aPoint
			}

			if (c & 0x01) != 0 {
				xsvExt0 = xsb + 2
				xsvExt1 = xsb + 1
				dxExt0 = dx0 - 2 - 3*squishConstant3D
				dxExt1 = dx0 - 1 - 3*squishConstant3D
			} else {
				xsvExt1 = xsb
				xsvExt0 = xsvExt1
				dxExt1 = dx0 - 3*squishConstant3D
				dxExt0 = dxExt1
			}

			if (c & 0x02) != 0 {
				ysvExt1 = ysb + 1
				ysvExt0 = ysvExt1
				dyExt1 = dy0 - 1 - 3*squishConstant3D
				dyExt0 = dyExt1
				if (c & 0x01) != 0 {
					ysvExt1++
					dyExt1--
				} else {
					ysvExt0++
					dyExt0--
				}
			} else {
				ysvExt1 = ysb
				ysvExt0 = ysvExt1
				dyExt1 = dy0 - 3*squishConstant3D
				dyExt0 = dyExt1
			}

			if (c & 0x04) != 0 {
				zsvExt0 = zsb + 1
				zsvExt1 = zsb + 2
				dzExt0 = dz0 - 1 - 3*squishConstant3D
				dzExt1 = dz0 - 2 - 3*squishConstant3D
			} else {
				zsvExt1 = zsb
				zsvExt0 = zsvExt1
				dzExt1 = dz0 - 3*squishConstant3D
				dzExt0 = dzExt1
			}
		} else { // (1,1,1) is not one of the closest two tetrahedral vertices.
			c := aPoint & bPoint // Our two extra vertices are determined by the closest two.

			if (c & 0x01) != 0 {
				xsvExt0 = xsb + 1
				xsvExt1 = xsb + 2
				dxExt0 = dx0 - 1 - squishConstant3D
				dxExt1 = dx0 - 2 - 2*squishConstant3D
			} else {
				xsvExt1 = xsb
				xsvExt0 = xsvExt1
				dxExt0 = dx0 - squishConstant3D
				dxExt1 = dx0 - 2*squishConstant3D
			}

			if (c & 0x02) != 0 {
				ysvExt0 = ysb + 1
				ysvExt1 = ysb + 2
				dyExt0 = dy0 - 1 - squishConstant3D
				dyExt1 = dy0 - 2 - 2*squishConstant3D
			} else {
				ysvExt1 = ysb
				ysvExt0 = ysvExt1
				dyExt0 = dy0 - squishConstant3D
				dyExt1 = dy0 - 2*squishConstant3D
			}

			if (c & 0x04) != 0 {
				zsvExt0 = zsb + 1
				zsvExt1 = zsb + 2
				dzExt0 = dz0 - 1 - squishConstant3D
				dzExt1 = dz0 - 2 - 2*squishConstant3D
			} else {
				zsvExt1 = zsb
				zsvExt0 = zsvExt1
				dzExt0 = dz0 - squishConstant3D
				dzExt1 = dz0 - 2*squishConstant3D
			}
		}

		// Contribution (1,1,0)
		dx3 := dx0 - 1 - 2*squishConstant3D
		dy3 := dy0 - 1 - 2*squishConstant3D
		dz3 := dz0 - 0 - 2*squishConstant3D
		attn3 := 2 - dx3*dx3 - dy3*dy3 - dz3*dz3
		if attn3 > 0 {
			attn3 *= attn3
			value += attn3 * attn3 * s.extrapolate3(xsb+1, ysb+1, zsb+0, dx3, dy3, dz3)
		}

		// Contribution (1,0,1)
		dx2 := dx3
		dy2 := dy0 - 0 - 2*squishConstant3D
		dz2 := dz0 - 1 - 2*squishConstant3D
		attn2 := 2 - dx2*dx2 - dy2*dy2 - dz2*dz2
		if attn2 > 0 {
			attn2 *= attn2
			value += attn2 * attn2 * s.extrapolate3(xsb+1, ysb+0, zsb+1, dx2, dy2, dz2)
		}

		// Contribution (0,1,1)
		dx1 := dx0 - 0 - 2*squishConstant3D
		dy1 := dy3
		dz1 := dz2
		attn1 := 2 - dx1*dx1 - dy1*dy1 - dz1*dz1
		if attn1 > 0 {
			attn1 *= attn1
			value += attn1 * attn1 * s.extrapolate3(xsb+0, ysb+1, zsb+1, dx1, dy1, dz1)
		}

		// Contribution (1,1,1)
		dx0 = dx0 - 1 - 3*squishConstant3D
		dy0 = dy0 - 1 - 3*squishConstant3D
		dz0 = dz0 - 1 - 3*squishConstant3D
		attn0 := 2 - dx0*dx0 - dy0*dy0 - dz0*dz0
		if attn0 > 0 {
			attn0 *= attn0
			value += attn0 * attn0 * s.extrapolate3(xsb+1, ysb+1, zsb+1, dx0, dy0, dz0)
		}
	} else { // We're inside the octahedron (Rectified 3-Simplex) in between.
		var aScore, bScore float64
		var aPoint, bPoint byte
		var aIsFurtherSide, bIsFurtherSide bool

		// Decide between point (0,0,1) and (1,1,0) as closest
		p1 := xins + yins
		if p1 > 1 {
			aScore = p1 - 1
			aPoint = 0x03
			aIsFurtherSide = true
		} else {
			aScore = 1 - p1
			aPoint = 0x04
			aIsFurtherSide = false
		}

		// Decide between point (0,1,0) and (1,0,1) as closest
		p2 := xins + zins
		if p2 > 1 {
			bScore = p2 - 1
			bPoint = 0x05
			bIsFurtherSide = true
		} else {
			bScore = 1 - p2
			bPoint = 0x02
			bIsFurtherSide = false
		}

		// The closest out of the two (1,0,0) and (0,1,1) will replace the furthest out of the two decided above, if closer.
		p3 := yins + zins
		if p3 > 1 {
			score := p3 - 1
			if aScore <= bScore && aScore < score {
				aScore = score
				aPoint = 0x06
				aIsFurtherSide = true
			} else if aScore > bScore && bScore < score {
				bScore = score
				bPoint = 0x06
				bIsFurtherSide = true
			}
		} else {
			score := 1 - p3
			if aScore <= bScore && aScore < score {
				aScore = score
				aPoint = 0x01
				aIsFurtherSide = false
			} else if aScore > bScore && bScore < score {
				bScore = score
				bPoint = 0x01
				bIsFurtherSide = false
			}
		}

		// Where each of the two closest points are determines how the extra two vertices are calculated.
		if aIsFurtherSide == bIsFurtherSide {
			if aIsFurtherSide { // Both closest points on (1,1,1) side

				// One of the two extra points is (1,1,1)
				dxExt0 = dx0 - 1 - 3*squishConstant3D
				dyExt0 = dy0 - 1 - 3*squishConstant3D
				dzExt0 = dz0 - 1 - 3*squishConstant3D
				xsvExt0 = xsb + 1
				ysvExt0 = ysb + 1
				zsvExt0 = zsb + 1

				// Other extra point is based on the shared axis.
				c := aPoint & bPoint
				if (c & 0x01) != 0 {
					dxExt1 = dx0 - 2 - 2*squishConstant3D
					dyExt1 = dy0 - 2*squishConstant3D
					dzExt1 = dz0 - 2*squishConstant3D
					xsvExt1 = xsb + 2
					ysvExt1 = ysb
					zsvExt1 = zsb
				} else if (c & 0x02) != 0 {
					dxExt1 = dx0 - 2*squishConstant3D
					dyExt1 = dy0 - 2 - 2*squishConstant3D
					dzExt1 = dz0 - 2*squishConstant3D
					xsvExt1 = xsb
					ysvExt1 = ysb + 2
					zsvExt1 = zsb
				} else {
					dxExt1 = dx0 - 2*squishConstant3D
					dyExt1 = dy0 - 2*squishConstant3D
					dzExt1 = dz0 - 2 - 2*squishConstant3D
					xsvExt1 = xsb
					ysvExt1 = ysb
					zsvExt1 = zsb + 2
				}
			} else { // Both closest points on (0,0,0) side

				// One of the two extra points is (0,0,0)
				dxExt0 = dx0
				dyExt0 = dy0
				dzExt0 = dz0
				xsvExt0 = xsb
				ysvExt0 = ysb
				zsvExt0 = zsb

				// Other extra point is based on the omitted axis.
				c := aPoint | bPoint
				if (c & 0x01) == 0 {
					dxExt1 = dx0 + 1 - squishConstant3D
					dyExt1 = dy0 - 1 - squishConstant3D
					dzExt1 = dz0 - 1 - squishConstant3D
					xsvExt1 = xsb - 1
					ysvExt1 = ysb + 1
					zsvExt1 = zsb + 1
				} else if (c & 0x02) == 0 {
					dxExt1 = dx0 - 1 - squishConstant3D
					dyExt1 = dy0 + 1 - squishConstant3D
					dzExt1 = dz0 - 1 - squishConstant3D
					xsvExt1 = xsb + 1
					ysvExt1 = ysb - 1
					zsvExt1 = zsb + 1
				} else {
					dxExt1 = dx0 - 1 - squishConstant3D
					dyExt1 = dy0 - 1 - squishConstant3D
					dzExt1 = dz0 + 1 - squishConstant3D
					xsvExt1 = xsb + 1
					ysvExt1 = ysb + 1
					zsvExt1 = zsb - 1
				}
			}
		} else { // One point on (0,0,0) side, one point on (1,1,1) side
			var c1, c2 byte
			if aIsFurtherSide {
				c1 = aPoint
				c2 = bPoint
			} else {
				c1 = bPoint
				c2 = aPoint
			}

			// One contribution is a permutation of (1,1,-1)
			if (c1 & 0x01) == 0 {
				dxExt0 = dx0 + 1 - squishConstant3D
				dyExt0 = dy0 - 1 - squishConstant3D
				dzExt0 = dz0 - 1 - squishConstant3D
				xsvExt0 = xsb - 1
				ysvExt0 = ysb + 1
				zsvExt0 = zsb + 1
			} else if (c1 & 0x02) == 0 {
				dxExt0 = dx0 - 1 - squishConstant3D
				dyExt0 = dy0 + 1 - squishConstant3D
				dzExt0 = dz0 - 1 - squishConstant3D
				xsvExt0 = xsb + 1
				ysvExt0 = ysb - 1
				zsvExt0 = zsb + 1
			} else {
				dxExt0 = dx0 - 1 - squishConstant3D
				dyExt0 = dy0 - 1 - squishConstant3D
				dzExt0 = dz0 + 1 - squishConstant3D
				xsvExt0 = xsb + 1
				ysvExt0 = ysb + 1
				zsvExt0 = zsb - 1
			}

			// One contribution is a permutation of (0,0,2)
			dxExt1 = dx0 - 2*squishConstant3D
			dyExt1 = dy0 - 2*squishConstant3D
			dzExt1 = dz0 - 2*squishConstant3D
			xsvExt1 = xsb
			ysvExt1 = ysb
			zsvExt1 = zsb
			if (c2 & 0x01) != 0 {
				dxExt1 -= 2
				xsvExt1 += 2
			} else if (c2 & 0x02) != 0 {
				dyExt1 -= 2
				ysvExt1 += 2
			} else {
				dzExt1 -= 2
				zsvExt1 += 2
			}
		}

		// Contribution (1,0,0)
		dx1 := dx0 - 1 - squishConstant3D
		dy1 := dy0 - 0 - squishConstant3D
		dz1 := dz0 - 0 - squishConstant3D
		attn1 := 2 - dx1*dx1 - dy1*dy1 - dz1*dz1
		if attn1 > 0 {
			attn1 *= attn1
			value += attn1 * attn1 * s.extrapolate3(xsb+1, ysb+0, zsb+0, dx1, dy1, dz1)
		}

		// Contribution (0,1,0)
		dx2 := dx0 - 0 - squishConstant3D
		dy2 := dy0 - 1 - squishConstant3D
		dz2 := dz1
		attn2 := 2 - dx2*dx2 - dy2*dy2 - dz2*dz2
		if attn2 > 0 {
			attn2 *= attn2
			value += attn2 * attn2 * s.extrapolate3(xsb+0, ysb+1, zsb+0, dx2, dy2, dz2)
		}

		// Contribution (0,0,1)
		dx3 := dx2
		dy3 := dy1
		dz3 := dz0 - 1 - squishConstant3D
		attn3 := 2 - dx3*dx3 - dy3*dy3 - dz3*dz3
		if attn3 > 0 {
			attn3 *= attn3
			value += attn3 * attn3 * s.extrapolate3(xsb+0, ysb+0, zsb+1, dx3, dy3, dz3)
		}

		// Contribution (1,1,0)
		dx4 := dx0 - 1 - 2*squishConstant3D
		dy4 := dy0 - 1 - 2*squishConstant3D
		dz4 := dz0 - 0 - 2*squishConstant3D
		attn4 := 2 - dx4*dx4 - dy4*dy4 - dz4*dz4
		if attn4 > 0 {
			attn4 *= attn4
			value += attn4 * attn4 * s.extrapolate3(xsb+1, ysb+1, zsb+0, dx4, dy4, dz4)
		}

		// Contribution (1,0,1)
		dx5 := dx4
		dy5 := dy0 - 0 - 2*squishConstant3D
		dz5 := dz0 - 1 - 2*squishConstant3D
		attn5 := 2 - dx5*dx5 - dy5*dy5 - dz5*dz5
		if attn5 > 0 {
			attn5 *= attn5
			value += attn5 * attn5 * s.extrapolate3(xsb+1, ysb+0, zsb+1, dx5, dy5, dz5)
		}

		// Contribution (0,1,1)
		dx6 := dx0 - 0 - 2*squishConstant3D
		dy6 := dy4
		dz6 := dz5
		attn6 := 2 - dx6*dx6 - dy6*dy6 - dz6*dz6
		if attn6 > 0 {
			attn6 *= attn6
			value += attn6 * attn6 * s.extrapolate3(xsb+0, ysb+1, zsb+1, dx6, dy6, dz6)
		}
	}

	// First extra vertex
	attnExt0 := 2 - dxExt0*dxExt0 - dyExt0*dyExt0 - dzExt0*dzExt0
	if attnExt0 > 0 {
		attnExt0 *= attnExt0
		value += attnExt0 * attnExt0 * s.extrapolate3(xsvExt0, ysvExt0, zsvExt0, dxExt0, dyExt0, dzExt0)
	}

	// Second extra vertex
	attnExt1 := 2 - dxExt1*dxExt1 - dyExt1*dyExt1 - dzExt1*dzExt1
	if attnExt1 > 0 {
		attnExt1 *= attnExt1
		value += attnExt1 * attnExt1 * s.extrapolate3(xsvExt1, ysvExt1, zsvExt1, dxExt1, dyExt1, dzExt1)
	}

	return value / normConstant3D
}

func (s *Noise) extrapolate3(xsb, ysb, zsb int32, dx, dy, dz float64) float64 {
	index := s.permGradIndex3D[(int32(s.perm[(int32(s.perm[xsb&0xFF])+ysb)&0xFF])+zsb)&0xFF]
	return float64(gradients3D[index])*dx + float64(gradients3D[index+1])*dy + float64(gradients3D[index+2])*dz
}

// Gradients for 3D. They approximate the directions to the
// vertices of a rhombicuboctahedron from the center, skewed so
// that the triangular and square facets can be inscribed inside
// circles of the same radius.
var gradients3D = []int8{
	-11, 4, 4, -4, 11, 4, -4, 4, 11,
	11, 4, 4, 4, 11, 4, 4, 4, 11,
	-11, -4, 4, -4, -11, 4, -4, -4, 11,
	11, -4, 4, 4, -11, 4, 4, -4, 11,
	-11, 4, -4, -4, 11, -4, -4, 4, -11,
	11, 4, -4, 4, 11, -4, 4, 4, -11,
	-11, -4, -4, -4, -11, -4, -4, -4, -11,
	11, -4, -4, 4, -11, -4, 4, -4, -11,
}

//func main() {
//	n := NewWithSeed(54524524)
//	for x := 0.0; x < 5.0; x++ {
//		for y := 0.0; y < 5.0; y++ {
//			fmt.Println(n.Eval3(x, y, 0))
//		}
//	}
//}
