package main

import (
	"fmt"
	"strings"
)

// Point represents a point on a 2D plane.
type Point struct {
	X int
	Y int
}

func (p Point) Add(o Point) Point {
	p.X += o.X
	p.Y += o.Y
	return p
}

func (p Point) Scale(m int) Point {
	p.X *= m
	p.Y *= m
	return p
}

// Quantity represents a number with two fractional digits.
type Quantity int

func (q Quantity) Int() int {
	return q / 100
}
