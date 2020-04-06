package utils

import(
	"time"
)

type FourDirection int
const (
	FourDirectionUp FourDirection = iota
	FourDirectionRight
	FourDirectionDown
	FourDirectionLeft
)

type IMatrixPosition interface {
	GetX() int
	GetY() int
	Validate(rowLength int, columnLength int) bool
}

type MatrixPosition struct {
	X int
	Y int
}
func (matrixPosition *MatrixPosition) GetX() int {
	return matrixPosition.X
}
func (matrixPosition *MatrixPosition) GetY() int {
	return matrixPosition.Y
}
func (matrixPosition *MatrixPosition) Validate(rowLength int, columnLength int) bool {
	y := matrixPosition.GetY()
	x := matrixPosition.GetX()
	return y >= 0 && y < rowLength && x >= 0 && x < columnLength
}

type IFieldElement interface {
	GetObjectClass() string
	GetPosition() IMatrixPosition
	IsObjectEmpty() bool
	UpdateObjectClass(class string)
}

type IField interface {
	At(position IMatrixPosition) (IFieldElement, error)
	GetElementOfHero() IFieldElement
	MeasureColumnLength() int
	MeasureRowLength() int
	MoveObject(from IMatrixPosition, to IMatrixPosition) error
}

type IGame interface {
	AlterPlaytime(delta time.Duration)
	GetPlaytimeAsSeconds() int
	GetPlaytimeAsString() string
	Initialize()
	IsStarted() bool
	IsFinished() bool
//	Finish()
//	GetTime() int
//	Start()
}

type IState interface {
	GetField() IField
	GetGame() IGame
}
