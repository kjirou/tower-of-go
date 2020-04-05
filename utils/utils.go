package utils

type MatrixPosition interface {
	GetX() int
	GetY() int
	Validate(rowLength int, columnLength int) bool
}

type RpgFieldElement interface {
	GetObjectClass() string
	IsObjectEmpty() bool
}
