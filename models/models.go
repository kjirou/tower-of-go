package models

import (
	"fmt"
	"github.com/kjirou/tower_of_go/utils"
	"time"
)

type FieldObject struct {
	Class string
}

func (fieldObject *FieldObject) IsEmpty() bool {
	return fieldObject.Class == "empty"
}

type FieldElement struct {
	Object   FieldObject
	Position utils.MatrixPosition
}

func (fieldElement *FieldElement) GetPosition() utils.IMatrixPosition {
	var position utils.IMatrixPosition = &fieldElement.Position
	return position
}

func (fieldElement *FieldElement) GetObjectClass() string {
	return fieldElement.Object.Class
}

func (fieldElement *FieldElement) IsObjectEmpty() bool {
	return fieldElement.Object.IsEmpty()
}

func (fieldElement *FieldElement) UpdateObjectClass(class string) {
	fieldElement.Object.Class = class
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

func (field *Field) At(position utils.IMatrixPosition) (utils.IFieldElement, error) {
	y := position.GetY()
	x := position.GetX()
	if y < 0 || y > field.MeasureRowLength()-1 {
		return &FieldElement{}, fmt.Errorf("That position (Y=%d) does not exist on the field.", y)
	} else if x < 0 || x > field.MeasureColumnLength()-1 {
		return &FieldElement{}, fmt.Errorf("That position (X=%d) does not exist on the field.", x)
	}
	return &(field.matrix[y][x]), nil
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

func (field *Field) GetElementOfHero() utils.IFieldElement {
	elements := field.FindElementsByObjectClass("hero")
	// TODO: Error handling.
	if len(elements) == 0 {
		panic("The hero does not exist.")
	} else if len(elements) > 1 {
		panic("There are multiple heroes.")
	}
	return elements[0]
}

func (field *Field) MoveObject(from utils.IMatrixPosition, to utils.IMatrixPosition) error {
	fromElement, err := field.At(from)
	if err != nil {
		return err
	} else if fromElement.IsObjectEmpty() {
		return fmt.Errorf("The object to be moved does not exist.")
	}
	toElement, err := field.At(to)
	if err != nil {
		return err
	} else if !toElement.IsObjectEmpty() {
		return fmt.Errorf("An object exists at the destination.")
	}
	toElement.UpdateObjectClass(fromElement.GetObjectClass())
	fromElement.UpdateObjectClass("empty")
	return nil
}

func createField(y int, x int) Field {
	matrix := make([][]FieldElement, y)
	for rowIndex := 0; rowIndex < y; rowIndex++ {
		row := make([]FieldElement, x)
		for columnIndex := 0; columnIndex < x; columnIndex++ {
			row[columnIndex] = FieldElement{
				Position: utils.MatrixPosition{
					Y: rowIndex,
					X: columnIndex,
				},
				Object: FieldObject{
					Class: "empty",
				},
			}
		}
		matrix[rowIndex] = row
	}
	return Field{
		matrix: matrix,
	}
}

type Game struct {
	isFinished bool
	isStarted bool
	// This is the total of main loop intervals.
	// It is different from the real time.
	playtime time.Duration
}

func (game *Game) Initialize() {
	game.isStarted = false
	game.isFinished = false
	duration, _ := time.ParseDuration("0")
	game.playtime = duration
}

func (game *Game) IsStarted() bool {
	return game.isStarted
}

func (game *Game) IsFinished() bool {
	return game.isFinished
}

func (game *Game) GetPlaytimeAsSeconds() int {
	return int(game.playtime.Seconds())
}

func (game *Game) GetPlaytimeAsString() string {
	return fmt.Sprintf("%d", game.GetPlaytimeAsSeconds())
}

func (game *Game) Start() {
	game.isStarted = true
	game.isFinished = false
}

func (game *Game) AlterPlaytime(duration time.Duration) {
	game.playtime = game.playtime + duration
}

type State struct {
	field Field
	game utils.IGame
}

func (state *State) GetField() utils.IField {
	var field utils.IField = &state.field
	return field
}

func (state *State) GetGame() utils.IGame {
	return state.game
}

func (state *State) InitializeDummyData() error {
	field := state.GetField()

	// Hero
	var heroPosition utils.IMatrixPosition = &utils.MatrixPosition{Y: 2, X: 5}
	heroFieldElement, err := field.At(heroPosition)
	if err != nil {
		return err
	}
	heroFieldElement.UpdateObjectClass("hero")

	// Walls
	fieldRowLength := field.MeasureRowLength()
	fieldColumnLength := field.MeasureColumnLength()
	for y := 0; y < fieldRowLength; y++ {
		for x := 0; x < fieldColumnLength; x++ {
			isTopOrBottomEdge := y == 0 || y == fieldRowLength-1
			isLeftOrRightEdge := x == 0 || x == fieldColumnLength-1
			if isTopOrBottomEdge || isLeftOrRightEdge {
				var wallPosition utils.IMatrixPosition = &utils.MatrixPosition{Y: y, X: x}
				var elem, _ = field.At(wallPosition)
				elem.UpdateObjectClass("wall")
			}
		}
	}

	return nil
}

func CreateState() State {
	var game utils.IGame = &Game{}
	state := State{
		field: createField(12, 20),
		game: game,
	}
	game.Initialize()
	return state
}
