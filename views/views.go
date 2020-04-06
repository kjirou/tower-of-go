package views

import (
	"fmt"
	"github.com/kjirou/tower_of_go/utils"
	"github.com/nsf/termbox-go"
	"strings"
)

type ScreenElement struct {
	Symbol          rune
	ForegroundColor termbox.Attribute
	BackgroundColor termbox.Attribute
}

func (screenElement *ScreenElement) renderWithFieldElement(fieldElement utils.IFieldElement) {
	symbol := '.'
	fg := termbox.ColorWhite
	bg := termbox.ColorBlack
	if !fieldElement.IsObjectEmpty() {
		switch fieldElement.GetObjectClass() {
		case "hero":
			symbol = '@'
			fg = termbox.ColorMagenta
		case "wall":
			symbol = '#'
			fg = termbox.ColorYellow
		default:
			symbol = '?'
		}
	}
	screenElement.Symbol = symbol
	screenElement.ForegroundColor = fg
	screenElement.BackgroundColor = bg
}

//
// A layer that avoids to write logics tightly coupled with "termbox".
//
type Screen struct {
	matrix [][]ScreenElement
}

func (screen *Screen) GetMatrix() [][]ScreenElement {
	return screen.matrix
}

func (screen *Screen) MeasureRowLength() int {
	return len(screen.matrix)
}

func (screen *Screen) MeasureColumnLength() int {
	return len(screen.matrix[0])
}

func (screen *Screen) At(position utils.IMatrixPosition) *ScreenElement {
	y := position.GetY()
	x := position.GetX()
	// TODO: Error handling.
	if y < 0 || y > screen.MeasureRowLength() {
		panic(fmt.Sprintf("That position (Y=%d) does not exist on the screen.", y))
	} else if x < 0 || x > screen.MeasureColumnLength() {
		panic(fmt.Sprintf("That position (X=%d) does not exist on the screen.", x))
	}
	return &(screen.matrix[y][x])
}

func (screen *Screen) renderField(startPosition utils.IMatrixPosition, field utils.IField) {
	rowLength := field.MeasureRowLength()
	columnLength := field.MeasureColumnLength()
	for y := 0; y < rowLength; y++ {
		for x := 0; x < columnLength; x++ {
			var screenElementPosition utils.IMatrixPosition = &utils.MatrixPosition{
				Y: startPosition.GetY() + y,
				X: startPosition.GetX() + x,
			}
			element := screen.At(screenElementPosition)
			var fieldElementPosition utils.IMatrixPosition = &utils.MatrixPosition{Y: y, X: x}
			// TODO: Error handling.
			var fieldElement, _ = field.At(fieldElementPosition)
			element.renderWithFieldElement(fieldElement)
		}
	}
}

func (screen *Screen) Render(state utils.IState) {
	rowLength := screen.MeasureRowLength()
	columnLength := screen.MeasureColumnLength()

	// Set borders.
	for y := 0; y < rowLength; y++ {
		for x := 0; x < columnLength; x++ {
			isTopOrBottomEdge := y == 0 || y == rowLength-1
			isLeftOrRightEdge := x == 0 || x == columnLength-1
			symbol := ' '
			switch {
			case isTopOrBottomEdge && isLeftOrRightEdge:
				symbol = '+'
			case isTopOrBottomEdge && !isLeftOrRightEdge:
				symbol = '-'
			case !isTopOrBottomEdge && isLeftOrRightEdge:
				symbol = '|'
			}
			screen.matrix[y][x].Symbol = symbol
		}
	}

	// Place the field.
	var fieldPosition utils.IMatrixPosition = &utils.MatrixPosition{Y: 1, X: 1}
	screen.renderField(fieldPosition, state.GetField())
}

func (screen *Screen) AsText() string {
	rowLength := screen.MeasureRowLength()
	columnLength := screen.MeasureColumnLength()
	lines := make([]string, rowLength)
	for y := 0; y < rowLength; y++ {
		line := make([]rune, columnLength)
		for x := 0; x < columnLength; x++ {
			line[x] = screen.matrix[y][x].Symbol
		}
		lines[y] = string(line)
	}
	return strings.Join(lines, "\n")
}

func CreateScreen(rowLength int, columnLength int) Screen {
	matrix := make([][]ScreenElement, rowLength)
	for rowIndex := 0; rowIndex < rowLength; rowIndex++ {
		row := make([]ScreenElement, columnLength)
		for columnIndex := 0; columnIndex < columnLength; columnIndex++ {
			row[columnIndex] = ScreenElement{
				Symbol:          '_',
				ForegroundColor: termbox.ColorWhite,
				BackgroundColor: termbox.ColorBlack,
			}
		}
		matrix[rowIndex] = row
	}
	return Screen{
		matrix: matrix,
	}
}