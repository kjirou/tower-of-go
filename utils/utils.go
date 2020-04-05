package utils

type MatrixPosition interface {
	GetX() int
	GetY() int
	Validate(rowLength int, columnLength int) bool
}

type IFieldElement interface {
	GetObjectClass() string
	IsObjectEmpty() bool
	UpdateObjectClass(class string)
}

type IField interface {
	At(position MatrixPosition) IFieldElement
	MeasureColumnLength() int
	MeasureRowLength() int
}
