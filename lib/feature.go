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
	l := new(Feature)

	//for k, v := range args["LocMap"] {
	//locMap[k] = v
	//}
	l.LocMap = locMap

	return l, nil
}

type Feature struct {
	Name     string
	LocMap   Point
	Features []Feature
}

func (l *Feature) Walk(w WalkFunc) error {
	return nil
}

func (l *Feature) Move(p *Point) (*Feature, error) {
	return moveFeature(l, p)
}

func moveFeature(l Movable, p *Point) (*Feature, error) {
	return nil, nil
}
