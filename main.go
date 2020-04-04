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

type Field struct {
	matrix [][]FieldElement
}

func (field *Field) MeasureRowLength() int {
	return len(field.matrix)
}

func (field *Field) MeasureColumnLength() int {
	return len(field.matrix[0])
}

func (field *Field) At(position FieldPosition) *FieldElement {
	if position.Y < 0 || position.Y > field.MeasureRowLength() {
		panic(fmt.Sprintf("That position (Y=%d) does not exist on the field.", position.Y))
	} else if position.X < 0 || position.X > field.MeasureColumnLength() {
		panic(fmt.Sprintf("That position (X=%d) does not exist on the field.", position.X))
	}
	return &(field.matrix[position.Y][position.X])
}

// TODO: Refer `FieldObject.Class` type.
func (field *Field) FindElementsByObjectClass(objectClass string) []*FieldElement {
	elements := make([]*FieldElement, 0)
	for _, row := range field.matrix {
		for _, element := range row {
			if element.Object.Class == objectClass {
				elements = append(elements, &element)
			}
		}
	}
	return elements
}

func (field *Field) GetElementOfHero() *FieldElement {
	elements := field.FindElementsByObjectClass("hero")
	if len(elements) == 0 {
		panic("The hero does not exist.")
	} else if len(elements) > 1 {
		panic("There are multiple heroes.")
	}
	return elements[0]
}

func (field *Field) MoveObject(from FieldPosition, to FieldPosition) error {
	fromElement := field.At(from)
	if fromElement.Object.IsEmpty() {
		return fmt.Errorf("The object to be moved does not exist.")
	}
	toElement := field.At(to)
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
	Field Field
}

func createField(y int, x int) Field {
	matrix := make([][]FieldElement, y)
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
	return Field{
		matrix: matrix,
	}
}

// View
// ----

// TODO: Combine them into one `map[string]rune`.
const blankRune rune = 0x0020  // " "
const sharpRune rune = 0x0023  // "#"
const plusRune rune = 0x002b  // "+"
const hyphenRune rune = 0x002d  // "-"
const dotRune rune = 0x002e  // "."
const questionRune rune = 0x003f  // "?"
const atRune rune = 0x0040  // "@"
const virticalBarRune rune = 0x007C  // "|"

type ScreenPosition struct {
	X int
	Y int
}

type ScreenElement struct {
	character rune
	//foregroundColor
	//backgroundColor
}

func (se *ScreenElement) renderFieldElement(fieldElement *FieldElement) {
	symbol := dotRune
	if !fieldElement.Object.IsEmpty() {
		switch fieldElement.Object.Class {
			case "hero":
				symbol = atRune
			case "wall":
				symbol = sharpRune
			default:
				symbol = questionRune
		}
	}
	se.character = symbol
}

// A layer that avoids to write logics tightly coupled with "termbox".
type Screen struct {
	matrix [][]ScreenElement
}

func (s *Screen) MeasureRowLength() int {
	return len(s.matrix)
}

func (s *Screen) MeasureColumnLength() int {
	return len(s.matrix[0])
}

func (s *Screen) At(position ScreenPosition) *ScreenElement {
	if position.Y < 0 || position.Y > s.MeasureRowLength() {
		panic(fmt.Sprintf("That position (Y=%d) does not exist on the screen.", position.Y))
	} else if position.X < 0 || position.X > s.MeasureColumnLength() {
		panic(fmt.Sprintf("That position (X=%d) does not exist on the screen.", position.X))
	}
	return &(s.matrix[position.Y][position.X])
}

func (s *Screen) renderField(startPosition ScreenPosition, field Field) {
	rowLength := field.MeasureRowLength()
	columnLength := field.MeasureColumnLength()
	for y := 0; y < rowLength; y++ {
		for x := 0; x < columnLength; x++ {
			position := ScreenPosition{
				Y: startPosition.Y + y,
				X: startPosition.X + x,
			}
			element := s.At(position)
			element.renderFieldElement(field.At(FieldPosition{Y: y, X: x}))
		}
	}
}

func (s *Screen) render(state *State) {
	rowLength := s.MeasureRowLength()
	columnLength := s.MeasureColumnLength()

	// Set borders.
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
			s.matrix[y][x].character = character
		}
	}

	// Place the field.
	s.renderField(ScreenPosition{Y: 1, X: 1}, state.Field)
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
				character: questionRune,
			}
		}
		matrix[rowIndex] = row
	}
	return Screen{
		matrix: matrix,
	}
}

// Main Process
// ------------

func drawTerminal(screen *Screen) {
	for y, row := range screen.matrix {
		for x, screenElement := range row {
			termbox.SetCell(x, y, screenElement.character, termbox.ColorWhite, termbox.ColorBlack)
		}
	}
	termbox.Flush()
}

func initializeTermbox(screen *Screen) error {
	termboxErr := termbox.Init()
	if termboxErr != nil {
		return termboxErr
	}
	termbox.SetInputMode(termbox.InputEsc)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	return nil
}

func handleTermboxEvents(state *State, screen *Screen) {
	didQuitApplication := false

	for !didQuitApplication {
		event := termbox.PollEvent()
		switch event.Type {
		case termbox.EventKey:
			if event.Key == termbox.KeyCtrlC || event.Key == termbox.KeyCtrlQ {
				didQuitApplication = true
			}
		}
	}
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
		Field: createField(12, 20),
	}

	// Dummy data
	state.Field.At(FieldPosition{Y: 2, X: 5}).Object = FieldObject{
		Class: "hero",
	}
	fieldRowLength := state.Field.MeasureRowLength()
	fieldColumnLength := state.Field.MeasureColumnLength()
	for y := 0; y < fieldRowLength; y++ {
		for x := 0; x < fieldColumnLength; x++ {
			isTopOrBottomEdge := y == 0 || y == fieldRowLength - 1
			isLeftOrRightEdge := x == 0 || x == fieldColumnLength - 1
			if isTopOrBottomEdge || isLeftOrRightEdge {
				state.Field.At(FieldPosition{Y: y, X: x}).Object = FieldObject{
					Class: "wall",
				}
			}
		}
	}

	screen := createScreen(24 + 2, 80 + 2)
	screen.render(&state)

	if doesRunTermbox {
		termboxErr := initializeTermbox(&screen)
		if termboxErr != nil {
			panic(termboxErr)
		}
		defer termbox.Close()
		drawTerminal(&screen)
		handleTermboxEvents(&state, &screen)
	} else {
		fmt.Println(screen.AsText())
	}
}
