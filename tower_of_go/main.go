package main

// TODO:
// - gofmt
// - Separate to modules

import (
	"fmt"
	"strings"
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

func (fo FieldObject) IsEmpty() bool {
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

func renderFieldObject(fo FieldObject) string {
	switch fo.Class {
		case "empty":
			return ""
		case "hero":
			return "@"
		case "wall":
			return "#"
		default:
			return "?"
	}
}

func renderFieldElement(fe FieldElement) string {
	output := renderFieldObject(fe.Object)
	if output == "" {
		output = "."
	}
	return output
}

func renderFieldMatrix(fieldMatrix FieldMatrix) string {
	y := fieldMatrix.MeasureY()
	x := fieldMatrix.MeasureX()
	lines := make([]string, y)
	for rowIndex := 0; rowIndex < y; rowIndex++ {
		line := ""
		// TODO: Use mapping method
		for columnIndex := 0; columnIndex < x; columnIndex++ {
			line += renderFieldElement(fieldMatrix[rowIndex][columnIndex])
		}
		lines[rowIndex] = line
	}
	return strings.Join(lines, "\n")
}

func render(state State) string {
	return renderFieldMatrix(state.fieldMatrix)
}

// Main Process
// ------------

func main() {
	state := State{
		fieldMatrix: createFieldMatrix(12, 20),
	}
	state.fieldMatrix[1][2].Object = FieldObject{
		Class: "hero",
	}
	moveErr := state.fieldMatrix.MoveObject(FieldPosition{Y: 1, X: 2}, FieldPosition{Y: 1, X: 5})
	fmt.Println(moveErr)
	output := render(state)
	fmt.Println(output)
}
