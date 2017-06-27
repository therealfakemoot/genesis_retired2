package main

import "fmt"

type Point map[string]int

type Mover interface {
	Move(p *Point) (*Location, error)
}

func main() {
	fmt.Println("vim-go")
}
