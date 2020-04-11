package utils

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
	GetFloorObjectClass() string
	GetObjectClass() string
	GetPosition() IMatrixPosition
	IsObjectEmpty() bool
	UpdateFloorObjectClass(class string)
	UpdateObjectClass(class string)
}

var HeroPosition IMatrixPosition = &MatrixPosition{Y: 1, X: 1}
var UpstairsPosition IMatrixPosition = &MatrixPosition{Y: 11, X: 19}
