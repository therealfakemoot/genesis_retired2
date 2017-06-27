package genesis

type Point map[string]int

type Mover interface {
	Move(p *Point) (*Location, error)
}

func main() {
}
