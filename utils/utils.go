package utils

type MatrixPosition interface {
	GetX() int
	GetY() int
	Validate(rowLength int, columnLength int) bool
}
