package main

import (
	"fmt"
)

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

type FourDirection int
const (
	FourDirectionUp FourDirection = iota
	FourDirectionRight
	FourDirectionDown
	FourDirectionLeft
)

func (field *Field) WalkHero(direction FourDirection) error {
	element := field.GetElementOfHero()
	nextPosition := element.Position
	switch direction {
	case FourDirectionUp:
		nextPosition.Y -= 1
	case FourDirectionRight:
		nextPosition.X += 1
	case FourDirectionDown:
		nextPosition.Y += 1
	case FourDirectionLeft:
		nextPosition.X -= 1
	}
	// TODO: Error handling.
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

