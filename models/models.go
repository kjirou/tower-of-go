package models

import (
	"fmt"
	"github.com/kjirou/tower_of_go/utils"
	"time"
)

var HeroPosition = &utils.MatrixPosition{Y: 1, X: 1}
var UpstairsPosition = &utils.MatrixPosition{Y: 11, X: 19}

type FieldObject struct {
	Class string
}

type FieldFloorObject struct {
	Class string
}

func (fieldObject *FieldObject) IsEmpty() bool {
	return fieldObject.Class == "empty"
}

type FieldElement struct {
	Object   FieldObject
	FloorObject   FieldFloorObject
	Position *utils.MatrixPosition
}

func (fieldElement *FieldElement) GetPosition() *utils.MatrixPosition {
	return fieldElement.Position
}

func (fieldElement *FieldElement) GetObjectClass() string {
	return fieldElement.Object.Class
}

func (fieldElement *FieldElement) GetFloorObjectClass() string {
	return fieldElement.FloorObject.Class
}

func (fieldElement *FieldElement) IsObjectEmpty() bool {
	return fieldElement.Object.IsEmpty()
}

func (fieldElement *FieldElement) UpdateObjectClass(class string) {
	fieldElement.Object.Class = class
}

func (fieldElement *FieldElement) UpdateFloorObjectClass(class string) {
	fieldElement.FloorObject.Class = class
}

type Field struct {
	matrix [][]*FieldElement
}

func (field *Field) MeasureRowLength() int {
	return len(field.matrix)
}

func (field *Field) MeasureColumnLength() int {
	return len(field.matrix[0])
}

func (field *Field) At(position *utils.MatrixPosition) (*FieldElement, error) {
	y := position.GetY()
	x := position.GetX()
	if y < 0 || y > field.MeasureRowLength()-1 {
		return &FieldElement{}, fmt.Errorf("That position (Y=%d) does not exist on the field.", y)
	} else if x < 0 || x > field.MeasureColumnLength()-1 {
		return &FieldElement{}, fmt.Errorf("That position (X=%d) does not exist on the field.", x)
	}
	return field.matrix[y][x], nil
}

// TODO: Refer `FieldObject.Class` type.
func (field *Field) findElementsByObjectClass(objectClass string) []*FieldElement {
	elements := make([]*FieldElement, 0)
	for _, row := range field.matrix {
		for _, element := range row {
			if element.Object.Class == objectClass {
				element_ := element
				elements = append(elements, element_)
			}
		}
	}
	return elements
}

func (field *Field) GetElementOfHero() *FieldElement {
	elements := field.findElementsByObjectClass("hero")
	// TODO: Error handling.
	if len(elements) == 0 {
		panic("The hero does not exist.")
	} else if len(elements) > 1 {
		panic("There are multiple heroes.")
	}
	return elements[0]
}

func (field *Field) MoveObject(from *utils.MatrixPosition, to *utils.MatrixPosition) error {
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

func (field *Field) ResetMaze() error {
	rowLength := field.MeasureRowLength()
	columnLength := field.MeasureColumnLength()
	mazeCells, err := utils.GenerateMaze(rowLength, columnLength)
	if err != nil {
		return err
	}
	for y, mazeRow := range mazeCells {
		for x, mazeCell := range mazeRow {
			element, err := field.At(&utils.MatrixPosition{Y: y, X: x})
			if err != nil {
				return err
			}
			switch mazeCell.Content {
			case utils.MazeCellContentEmpty:
				element.UpdateObjectClass("empty")
			case utils.MazeCellContentUnbreakableWall:
				element.UpdateObjectClass("wall")
			}
		}
	}
	return nil
}

func createField(y int, x int) *Field {
	matrix := make([][]*FieldElement, y)
	for rowIndex := 0; rowIndex < y; rowIndex++ {
		row := make([]*FieldElement, x)
		for columnIndex := 0; columnIndex < x; columnIndex++ {
			row[columnIndex] = &FieldElement{
				Position: &utils.MatrixPosition{
					Y: rowIndex,
					X: columnIndex,
				},
				Object: FieldObject{
					Class: "empty",
				},
				FloorObject: FieldFloorObject{
					Class: "empty",
				},
			}
		}
		matrix[rowIndex] = row
	}
	return &Field{
		matrix: matrix,
	}
}

type Game struct {
	floorNumber int
	isFinished bool
	// A snapshot of `state.executionTime` when a game has started.
	startedAt time.Duration
}

func (game *Game) Reset() {
	zeroDuration, _ := time.ParseDuration("0s")
	game.startedAt = zeroDuration
	game.floorNumber = 1
	game.isFinished = false
}

func (game *Game) IsStarted() bool {
	zeroDuration, _ := time.ParseDuration("0s")
	return game.startedAt != zeroDuration
}

func (game *Game) IsFinished() bool {
	return game.isFinished
}

func (game *Game) CalculateRemainingTime(executionTime time.Duration) time.Duration {
	oneGameTime, _ := time.ParseDuration("30s")
	if game.IsStarted() {
		playtime := executionTime - game.startedAt
		remainingTime := oneGameTime - playtime
		if remainingTime < 0 {
			zeroTime, _ := time.ParseDuration("0s")
			return zeroTime
		}
		return remainingTime
	}
	return oneGameTime
}

func (game *Game) GetFloorNumber() int{
	return game.floorNumber
}

func (game *Game) IncrementFloorNumber() {
	game.floorNumber += 1
}

func (game *Game) Start(executionTime time.Duration) {
	game.startedAt = executionTime
}

func (game *Game) Finish() {
	game.isFinished = true
}

type State struct {
	// This is the total of main loop intervals.
	// It is different from the real time.
	executionTime time.Duration
	field *Field
	game *Game
}

func (state *State) GetExecutionTime() time.Duration {
	return state.executionTime
}

func (state *State) GetField() *Field {
	return state.field
}

func (state *State) GetGame() *Game {
	return state.game
}

func (state *State) AlterExecutionTime(delta time.Duration) {
	state.executionTime = state.executionTime + delta
}

func (state *State) SetWelcomeData() error {
	field := state.GetField()

	// Place a hero to be the player's alter ego.
	heroFieldElement, err := field.At(HeroPosition)
	if err != nil {
		return err
	}
	heroFieldElement.UpdateObjectClass("hero")

	// Place an upstairs.
	upstairsFieldElement, err := field.At(UpstairsPosition)
	if err != nil {
		return err
	}
	upstairsFieldElement.UpdateFloorObjectClass("upstairs")

	// Place defalt walls.
	fieldRowLength := field.MeasureRowLength()
	fieldColumnLength := field.MeasureColumnLength()
	for y := 0; y < fieldRowLength; y++ {
		for x := 0; x < fieldColumnLength; x++ {
			isTopOrBottomEdge := y == 0 || y == fieldRowLength-1
			isLeftOrRightEdge := x == 0 || x == fieldColumnLength-1
			if isTopOrBottomEdge || isLeftOrRightEdge {
				elem, _ := field.At(&utils.MatrixPosition{Y: y, X: x})
				elem.UpdateObjectClass("wall")
			}
		}
	}

	return nil
}

func CreateState() *State {
	executionTime, _ := time.ParseDuration("0")
	state := &State{
		executionTime: executionTime,
		field: createField(13, 21),
		game: &Game{},
	}
	state.game.Reset()
	return state
}
