package genesis

import (
	"context"
)

type Point map[string]interface{}

type Movable interface {
	Move(p *Point) (*Feature, error)
}

type Walkable interface {
	Walk(w *Walkable) (walkFunc WalkFunc)
}

type WalkCtx struct {
	context.Context
}

type WalkFunc func(root *Feature, ctx WalkCtx) error

func NewFeature(locMap map[string]interface{}, args ...interface{}) (*Feature, error) {
	f := new(Feature)

	//for k, v := range args["LocMap"] {
	//locMap[k] = v
	//}
	f.LocMap = locMap

	return f, nil
}

type Feature struct {
	Name     string
	LocMap   Point
	Features []Feature
}

	return nil
}

func (f *Feature) Move(p *Point) (*Feature, error) {
	return moveFeature(f, p)
}

func moveFeature(f Movable, p *Point) (*Feature, error) {
	return nil, nil
}
