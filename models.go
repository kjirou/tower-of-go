package main

import (
	"fmt"
	"github.com/kjirou/tower_of_go/utils"
)

type FieldPosition struct {
	x int
	y int
}

func (fieldPosition *FieldPosition) GetY() int {
	return fieldPosition.y
}

func (fieldPosition *FieldPosition) GetX() int {
	return fieldPosition.x
}

func (fieldPosition *FieldPosition) Validate(rowLength int, columnLength int) bool {
	y := fieldPosition.GetY()
	x := fieldPosition.GetX()
	return y >= 0 && y < rowLength && x >= 0 && x < columnLength
}

type FieldObject struct {
	Class string
}

func (fieldObject *FieldObject) IsEmpty() bool {
	return fieldObject.Class == "empty"
}

type FieldElement struct {
	Object   FieldObject
	Position FieldPosition
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

func (field *Field) At(position utils.MatrixPosition) utils.IFieldElement {
	y := position.GetY()
	x := position.GetX()
	// TODO: Error handling.
	if y < 0 || y > field.MeasureRowLength() {
		panic(fmt.Sprintf("That position (Y=%d) does not exist on the field.", y))
	} else if x < 0 || x > field.MeasureColumnLength() {
		panic(fmt.Sprintf("That position (X=%d) does not exist on the field.", x))
	}
	return &(field.matrix[y][x])
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

func (field *Field) MoveObject(from utils.MatrixPosition, to utils.MatrixPosition) error {
	fromElement := field.At(from)
	if fromElement.IsObjectEmpty() {
		return fmt.Errorf("The object to be moved does not exist.")
	}
	toElement := field.At(to)
	if toElement.IsObjectEmpty() == false {
		return fmt.Errorf("An object exists at the destination.")
	}
	toElement.UpdateObjectClass(fromElement.GetObjectClass())
	fromElement.UpdateObjectClass("empty")
	return nil
}

type FourDirection int
const (
	FourDirectionUp FourDirection = iota
	FourDirectionRight
	FourDirectionDown
	FourDirectionLeft
)

func (field *Field) WalkHero(direction FourDirection) error {
	element := field.GetElementOfHero()
	nextY := element.Position.GetY()
	nextX := element.Position.GetX()
	switch direction {
	case FourDirectionUp:
		nextY -= 1
	case FourDirectionRight:
		nextX += 1
	case FourDirectionDown:
		nextY += 1
	case FourDirectionLeft:
		nextX -= 1
	}
	var position utils.MatrixPosition = &element.Position
	var nextPosition utils.MatrixPosition = &FieldPosition{
		y: nextY,
		x: nextX,
	}
	// TODO: Error handling.
	if nextPosition.Validate(field.MeasureRowLength(), field.MeasureColumnLength()) {
		if field.At(nextPosition).IsObjectEmpty() {
			return field.MoveObject(position, nextPosition)
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
				y: rowIndex,
				x: columnIndex,
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

