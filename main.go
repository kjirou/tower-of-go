package main

// TODO:
// - Separate the main package to sub local packages.

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
	"strings"
)

// Model
// -----

type FieldPosition struct {
	X int
	Y int
}

func (fieldPosition *FieldPosition) Validate(rowLength int, columnLength int) bool {
	return fieldPosition.Y >= 0 &&
		fieldPosition.Y < rowLength &&
		fieldPosition.X >= 0 &&
		fieldPosition.X < columnLength
}

type FieldObject struct {
	// TODO: Enumerize
	Class string
}

func (fieldObject *FieldObject) IsEmpty() bool {
	return fieldObject.Class == "empty"
}

type FieldElement struct {
	Object   FieldObject
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
				element_ := element
				elements = append(elements, &element_)
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

// TODO: Enumerize
func (field *Field) WalkHero(direction string) error {
	element := field.GetElementOfHero()
	nextPosition := element.Position
	switch direction {
	case "up":
		nextPosition.Y -= 1
	case "right":
		nextPosition.X += 1
	case "down":
		nextPosition.Y += 1
	case "left":
		nextPosition.X -= 1
	}
	if nextPosition.Validate(field.MeasureRowLength(), field.MeasureColumnLength()) {
		if field.At(nextPosition).Object.IsEmpty() {
			return field.MoveObject(element.Position, nextPosition)
		}
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
				Object:   fieldObject,
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

// A layer that avoids to write logics tightly coupled with "termbox".
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
	for rowIndex := 0; rowIndex < rowLength; rowIndex++ {
		line := make([]rune, columnLength)
		// TODO: Use mapping method
		for columnIndex := 0; columnIndex < columnLength; columnIndex++ {
			line[columnIndex] = screen.matrix[rowIndex][columnIndex].Symbol
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

// Main Process
// ------------

func drawTerminal(screen *Screen) {
	for y, row := range screen.matrix {
		for x, element := range row {
			termbox.SetCell(x, y, element.Symbol, element.ForegroundColor, element.BackgroundColor)
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

// It may need to make the following processes:
//   https://github.com/nsf/termbox-go/blob/4d2b513ad8bee47a9a5a65b0dee0182049a31916/_demos/keyboard.go#L669
//   (However, details cannot be read...)
// TODO: Use termbox's types
func handleKeyPress(state *State, screen *Screen, ch rune, key termbox.Key) {
	var err error
	field := &state.Field
	stateChanged := false

	// Move the hero.
	// TODO: Consider arrow keys.
	if ch == 'k' {
		err = field.WalkHero("up")
		stateChanged = true
	} else if ch == 'l' {
		err = field.WalkHero("right")
		stateChanged = true
	} else if ch == 'j' {
		err = field.WalkHero("down")
		stateChanged = true
	} else if ch == 'h' {
		err = field.WalkHero("left")
		stateChanged = true
	}

	if err != nil {
		panic(err)
	}

	if stateChanged {
		screen.render(state)
		drawTerminal(screen)
	}
}

func handleTermboxEvents(state *State, screen *Screen) {
	didQuitApplication := false

	for !didQuitApplication {
		event := termbox.PollEvent()
		switch event.Type {
		case termbox.EventKey:
			// Quit the application. Only this operation is resolved with priority.
			if event.Key == termbox.KeyCtrlC || event.Key == termbox.KeyCtrlQ {
				didQuitApplication = true
				break
			}

			handleKeyPress(state, screen, event.Ch, event.Key)
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
			isTopOrBottomEdge := y == 0 || y == fieldRowLength-1
			isLeftOrRightEdge := x == 0 || x == fieldColumnLength-1
			if isTopOrBottomEdge || isLeftOrRightEdge {
				state.Field.At(FieldPosition{Y: y, X: x}).Object = FieldObject{
					Class: "wall",
				}
			}
		}
	}

	screen := createScreen(24+2, 80+2)
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
