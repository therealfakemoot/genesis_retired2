package genesis

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

type TreeAxis struct {
	Max int
	Min int
}

type LocMapFunc func(p *Point) *Point

func featureTree(depth TreeAxis, width TreeAxis, locMapFunc LocMapFunc) *Feature {
	return randomLocFeature()
}

func randomLocFeature() *Feature {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	x := r.Int()
	y := r.Int()
	z := r.Int()

	return dummyFeature(x, y, z)

}

func dummyFeature(x int, y int, z int) *Feature {
	m := map[string]interface{}{
		"x": x,
		"y": y,
		"z": z,
	}

	f, _ := NewFeature(m)
	return f
}

func TestFeatureWalk(t *testing.T) {

	var walkTests = []struct {
		in  *Feature
		out int
	}{
		{randomLocFeature(), 1},
	}

	ctx := &WalkCtx{}

	for _, tt := range walkTests {
		walkCount := 0

		walkFunc := func(r *Feature, ctx WalkCtx) error {
			walkCount++
			return nil
		}

		tt.in.Walk(walkFunc, ctx)

		if walkCount != 1 {
			t.Errorf("Expected walkcount = 1, got %v", walkCount)
		}
	}

}

func TestFeatureMove(t *testing.T) {
	locMap := make(map[string]interface{})
	feature, _ := NewFeature(locMap)

	if reflect.DeepEqual(locMap, feature.LocMap) {
		t.Errorf("Expected %v, got %v", locMap, feature.LocMap)
	}

}
