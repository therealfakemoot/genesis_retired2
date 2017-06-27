package genesis

type Feature interface {
	Walk(func(*Feature) error) error
	Set(string, interface{}) error
}

type Location struct {
	Name     string
	Location Point
	Features []Feature
}

func (l *Location) New() (Location, error) {
	return Location{}, nil
}

func (l *Location) Move(p *Point) (*Location, error) {
	return moveLocation(l, p)
}

func (l *Location) Set(key string, value interface{}) error {
	return nil
}

func moveLocation(l Mover, p *Point) (*Location, error) {
	return nil, nil
}
