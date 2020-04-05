package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"strings"
)

type ScreenPosition struct {
	X int
	Y int
}

type ScreenElement struct {
	Symbol          rune
	ForegroundColor termbox.Attribute
	BackgroundColor termbox.Attribute
}

func (screenElement *ScreenElement) renderFieldElement(fieldElement *FieldElement) {
	symbol := '.'
	fg := termbox.ColorWhite
	bg := termbox.ColorBlack
	if !fieldElement.Object.IsEmpty() {
		switch fieldElement.Object.Class {
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

func (screen *Screen) MeasureRowLength() int {
	return len(screen.matrix)
}

func (screen *Screen) MeasureColumnLength() int {
	return len(screen.matrix[0])
}

func (screen *Screen) At(position ScreenPosition) *ScreenElement {
	if position.Y < 0 || position.Y > screen.MeasureRowLength() {
		panic(fmt.Sprintf("That position (Y=%d) does not exist on the screen.", position.Y))
	} else if position.X < 0 || position.X > screen.MeasureColumnLength() {
		panic(fmt.Sprintf("That position (X=%d) does not exist on the screen.", position.X))
	}
	return &(screen.matrix[position.Y][position.X])
}

func (screen *Screen) renderField(startPosition ScreenPosition, field *Field) {
	rowLength := field.MeasureRowLength()
	columnLength := field.MeasureColumnLength()
	for y := 0; y < rowLength; y++ {
		for x := 0; x < columnLength; x++ {
			position := ScreenPosition{
				Y: startPosition.Y + y,
				X: startPosition.X + x,
			}
			element := screen.At(position)
			element.renderFieldElement(field.At(FieldPosition{Y: y, X: x}))
		}
	}
}

func (screen *Screen) render(state *State) {
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
	screen.renderField(ScreenPosition{Y: 1, X: 1}, &state.Field)
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

func createScreen(rowLength int, columnLength int) Screen {
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
