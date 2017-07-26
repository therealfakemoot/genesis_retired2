package genesis

import (
	"context"
)

// Point describes a single location in 'the world'.
//
// This type is deliberately open-ended to allow consumers of
// the Genesis API and structs to define their own location systems.
// For example, one system may express a Feature's location in terms of
// spatial ( X, Y, Z ) coordinates AND a time-component, allowing their map
// data to express information about how a specific region changes over time,
// like a castle being constructed, populated, and eventually falling into decay.
type Point map[string]interface{}

// Movable describes a Feature that may be in one place at a particular moment
// but another at some other point in time. This will allow describing Features
// or Points of Interest that migrate over time ( a roaming caravan, a seasonal
// festival, or even armies ).
type Movable interface {
	Move(p *Point) (*Feature, error)
}

// Walkable describes a Feature which contains other Features and implements a
// means for traversing its contents. The specifics of how this is achieved is
// up to the individual type and could be depth-first, breadth-first and is not
// guaranteed by this API.
type Walkable interface {
	Walk(w *Walkable) (walkFunc WalkFunc)
}

// WalkCtx is used when satisfying the Walkable interface. It will be used to pass
// search parameters ( such as "ignore man-made structures" or "only find magical Features" )
// down the call chain if further invocations of Walk() are made.
//
// WalkCtx also carries the Matches channel, onto which matching Features will be pushed for
// consumption by who or whatever may need them.
type WalkCtx struct {
	context.Context
	Matches chan interface{}
}

// WalkFunc is the callback passed to Walk() methods. It receives a root Feature
// which exposes its contained Features directly (or via its own Walk() method)
// and the WalkCtx value, allowing it to push desired Features into the Matches channel for
// later processing.
type WalkFunc func(root *Feature, ctx WalkCtx) error

// NewFeature is a helper function mainly intended for use during development and for testing
// the core implementation of Genesis. Will most likely end up not being exported in the 'final'
// product.
func NewFeature(locMap map[string]interface{}, args ...interface{}) (*Feature, error) {
	f := new(Feature)

	//for k, v := range args["LocMap"] {
	//locMap[k] = v
	//}
	f.LocMap = locMap

	return f, nil
}

// Feature is the stock implementation of the Walkable/Movable interfaces.
//
// Feature should be sufficient for most of the maps and points-of-interest, as it is
// highly generalized and places few restrictions on the data which you can store on/in it.
type Feature struct {
	Name     string
	LocMap   Point
	Features []Feature
}

// Walk is the demo implemntation of Walkable. It will iterate over all
// child Features on a given root Feature.
func (f *Feature) Walk(w WalkFunc, ctx *WalkCtx) error {
	for _, f := range f.Features {
		w(&f, *ctx)
	}
	return nil
}

// Move allows a Feature to be moved. The conceptual significance of such a change
// is entirely defined by the caller's design intent.
func (f *Feature) Move(p *Point) (*Feature, error) {
	return moveFeature(f, p)
}

func moveFeature(f Movable, p *Point) (*Feature, error) {
	return nil, nil
}
