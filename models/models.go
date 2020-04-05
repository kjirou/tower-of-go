package models

import (
	"fmt"
	"github.com/kjirou/tower_of_go/utils"
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

func (field *Field) GetElementOfHero() *FieldElement {
	elements := field.FindElementsByObjectClass("hero")
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
	} else if toElement.IsObjectEmpty() == false {
		return fmt.Errorf("An object exists at the destination.")
	}
	toElement.UpdateObjectClass(fromElement.GetObjectClass())
	fromElement.UpdateObjectClass("empty")
	return nil
}

func (field *Field) WalkHero(direction utils.FourDirection) error {
	element := field.GetElementOfHero()
	nextY := element.Position.GetY()
	nextX := element.Position.GetX()
	switch direction {
	case utils.FourDirectionUp:
		nextY -= 1
	case utils.FourDirectionRight:
		nextX += 1
	case utils.FourDirectionDown:
		nextY += 1
	case utils.FourDirectionLeft:
		nextX -= 1
	}
	var position utils.IMatrixPosition = &element.Position
	var nextPosition utils.IMatrixPosition = &utils.MatrixPosition{
		Y: nextY,
		X: nextX,
	}
	if nextPosition.Validate(field.MeasureRowLength(), field.MeasureColumnLength()) {
		element, err := field.At(nextPosition)
		if err != nil {
			return err
		} else if element.IsObjectEmpty() {
			return field.MoveObject(position, nextPosition)
		}
	}
	return nil
}

type State struct {
	field Field
}

func (state *State) GetField() utils.IField {
	var field utils.IField = &state.field
	return field
}

func createField(y int, x int) Field {
	matrix := make([][]FieldElement, y)
	for rowIndex := 0; rowIndex < y; rowIndex++ {
		row := make([]FieldElement, x)
		for columnIndex := 0; columnIndex < x; columnIndex++ {
			// TODO: Embed into the following FieldElement initialization
			fieldPosition := utils.MatrixPosition{
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

func CreateState() State {
	state := State{
		field: createField(12, 20),
	}
	return state
}
