package genesis

type Point map[string]interface{}

type Mover interface {
	Move(p *Point) (*Feature, error)
}
