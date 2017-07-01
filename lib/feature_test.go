package genesis

import (
	"reflect"
	"testing"
)

func TestFeatureMove(t *testing.T) {
	locMap := make(map[string]interface{})
	feature, _ := NewFeature(locMap)

	if reflect.DeepEqual(locMap, feature.LocMap) {
		t.Errorf("Expected %v, got %v", locMap, feature.LocMap)
	}

}

func main() {
}
