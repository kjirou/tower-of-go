package utils

type FourDirection int
const (
	FourDirectionUp FourDirection = iota
	FourDirectionRight
	FourDirectionDown
	FourDirectionLeft
)

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
	WalkHero(direction FourDirection) error
}

type IState interface {
	GetField() IField
}
