package main

// TODO:
// - go fmt
// - Separate to modules
// - Why did diffs in the go.mod/go.sub have increased? Probably only `go run` was executed.

import (
	"fmt"
	"os"
	"strings"
	"github.com/doronbehar/termbox-go"
)

// Model
// -----

type FieldPosition struct {
	X int
	Y int
}

type FieldObject struct {
	// TODO: Enumerize
	Class string
}

func (fo *FieldObject) IsEmpty() bool {
	return fo.Class == "empty"
}

type FieldElement struct {
	Object FieldObject
	Position FieldPosition
}

type FieldMatrix [][]FieldElement

func (fm *FieldMatrix) MeasureY() int {
	return len(*fm)
}

func (fm *FieldMatrix) MeasureX() int {
	return len((*fm)[0])
}

func (fm *FieldMatrix) At(fp FieldPosition) (*FieldElement, error) {
	// TODO: Is it correct? Should it return nil?
	notFound := FieldElement{}
	if fp.Y < 0 || fp.Y > fm.MeasureY() {
		return &notFound, fmt.Errorf("That position (Y=%d) does not exist on the field-matrix.", fp.Y)
	} else if fp.X < 0 || fp.X > fm.MeasureX() {
		return &notFound, fmt.Errorf("That position (X=%d) does not exist on the field-matrix.", fp.X)
	}
	return &((*fm)[fp.Y][fp.X]), nil
}

func (fm *FieldMatrix) MoveObject(from FieldPosition, to FieldPosition) error {
	fromElement, fromErr := fm.At(from)
	if fromErr != nil {
		return fromErr
	}
	if fromElement.Object.IsEmpty() {
		return fmt.Errorf("The object to be moved does not exist.")
	}
	toElement, toErr := fm.At(to)
	if toErr != nil {
		return toErr
	}
	if toElement.Object.IsEmpty() == false {
		return fmt.Errorf("An object exists at the destination.")
	}
	toElement.Object = fromElement.Object
	fromElement.Object = FieldObject{
		Class: "empty",
	}
	return nil
}

type State struct {
	fieldMatrix FieldMatrix
}

func createFieldMatrix(y int, x int) FieldMatrix {
	matrix := make(FieldMatrix, y)
	for rowIndex := 0; rowIndex < y; rowIndex++ {
		row := make([]FieldElement, x)
		for columnIndex := 0; columnIndex < x; columnIndex++ {
			// TODO: Embed into the following FieldElement initialization
			fieldPosition := FieldPosition{
				Y: rowIndex,
				X: columnIndex,
			}
			fieldObject := FieldObject{
				Class: "empty",
			}
			row[columnIndex] = FieldElement{
				Position: fieldPosition,
				Object: fieldObject,
			}
		}
		matrix[rowIndex] = row
	}
	return matrix
}

// View
// ----

// TODO: Combine them into one `map[string]rune`.
const blankRune rune = 0x0020  // " "
const sharpRune rune = 0x0023  // "#"
const plusRune rune = 0x002b  // "+"
const hyphenRune rune = 0x002d  // "-"
const dotRune rune = 0x002e  // "."
const atRune rune = 0x0040  // "@"
const virticalBarRune rune = 0x007C  // "|"

type ScreenElement struct {
	character rune
	//foregroundColor
	//backgroundColor
}

// A layer that avoid to write logics tightly coupled with "termbox".
type Screen struct {
	matrix [][]ScreenElement
}

func (s *Screen) MeasureRowLength() int {
	return len(s.matrix)
}

func (s *Screen) MeasureColumnLength() int {
	return len(s.matrix[0])
}

func (s *Screen) AsText() string {
	rowLength := s.MeasureRowLength()
	columnLength := s.MeasureColumnLength()
	lines := make([]string, rowLength)
	for rowIndex := 0; rowIndex < rowLength; rowIndex++ {
		line := make([]rune, columnLength)
		// TODO: Use mapping method
		for columnIndex := 0; columnIndex < columnLength; columnIndex++ {
			line[columnIndex] = s.matrix[rowIndex][columnIndex].character
		}
		lines[rowIndex] = string(line)
	}
	return strings.Join(lines, "\n")
}

func createScreen(rowLength int, columnLength int) Screen {
	matrix := make([][]ScreenElement, rowLength)
	for rowIndex := 0; rowIndex < rowLength; rowIndex++ {
		row := make([]ScreenElement, columnLength)
		for columnIndex := 0; columnIndex < columnLength; columnIndex++ {
			row[columnIndex] = ScreenElement{
				character: dotRune,
			}
		}
		matrix[rowIndex] = row
	}
	return Screen{
		matrix: matrix,
	}
}

func renderFieldObject(fo *FieldObject) string {
	switch fo.Class {
		case "hero":
			return "@"
		case "wall":
			return "#"
		default:
			return "?"
	}
}

func renderFieldElement(fe *FieldElement) string {
	if fe.Object.IsEmpty() {
		return "."
	}
	return renderFieldObject(&fe.Object)
}

func renderFieldMatrix(fieldMatrix FieldMatrix) string {
	y := fieldMatrix.MeasureY()
	x := fieldMatrix.MeasureX()
	lines := make([]string, y)
	for rowIndex := 0; rowIndex < y; rowIndex++ {
		line := ""
		// TODO: Use mapping method
		for columnIndex := 0; columnIndex < x; columnIndex++ {
			line += renderFieldElement(&(fieldMatrix[rowIndex][columnIndex]))
		}
		lines[rowIndex] = line
	}
	return strings.Join(lines, "\n")
}

func render(screen *Screen, state *State) error {
	rowLength := screen.MeasureRowLength()
	columnLength := screen.MeasureColumnLength()

	// Set borders on the screen.
	for y := 0; y < rowLength; y++ {
		for x := 0; x < columnLength; x++ {
			isTopOrBottomEdge := y == 0 || y == rowLength - 1
			isLeftOrRightEdge := x == 0 || x == columnLength - 1
			character := blankRune
			switch {
			case isTopOrBottomEdge && isLeftOrRightEdge:
				character = plusRune
			case isTopOrBottomEdge && !isLeftOrRightEdge:
				character = hyphenRune
			case !isTopOrBottomEdge && isLeftOrRightEdge:
				character = virticalBarRune
			}
			screen.matrix[y][x].character = character
		}
	}

	return nil
}

// Main Process
// ------------

func runTermbox(initialOutput string) error {
	termboxErr := termbox.Init()
	if termboxErr != nil {
		return termboxErr
	}

	termbox.SetInputMode(termbox.InputEsc)

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	termbox.SetCell(0, 0, plusRune, termbox.ColorWhite, termbox.ColorBlack)

	termbox.Flush()

	return nil
}

func main() {
	// TODO: Look for a tiny CLI argument parser like the "minimist" of Node.js.
	commandLineArgs := os.Args[1:]
	doesRunTermbox := false
	for _, arg := range commandLineArgs {
		if arg == "-t" {
			doesRunTermbox = true
		}
	}

	state := State{
		fieldMatrix: createFieldMatrix(12, 20),
	}
	state.fieldMatrix[1][2].Object = FieldObject{
		Class: "hero",
	}

	screen := createScreen(24 + 2, 80 + 2)

	state.fieldMatrix.MoveObject(FieldPosition{Y: 1, X: 2}, FieldPosition{Y: 1, X: 5})

	render(&screen, &state)

	if doesRunTermbox {
		termboxErr := runTermbox("")
		if termboxErr != nil {
			panic(termboxErr)
		}
		// TODO: Can it move into the runTermbox?
		didQuitApplication := false
		for didQuitApplication == false {
			event := termbox.PollEvent()
			fmt.Println(event.Type)
			switch event.Type {
			case termbox.EventKey:
				if event.Key == termbox.KeyCtrlC || event.Key == termbox.KeyCtrlQ {
					didQuitApplication = true
				}
			}
		}
		defer termbox.Close()
	} else {
		fmt.Println(screen.AsText())
	}
}
