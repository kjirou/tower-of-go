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

type FieldObject struct {
	// TODO: Enumerize
	Class string
}

type FieldElement struct {
	Object FieldObject
	X int
	Y int
}

type FieldMatrix [][]FieldElement

type State struct {
	fieldMatrix FieldMatrix
}

func createFieldMatrix(y int, x int) FieldMatrix {
	matrix := make(FieldMatrix, y)
	for rowIndex := 0; rowIndex < y; rowIndex++ {
		row := make([]FieldElement, x)
		for columnIndex := 0; columnIndex < x; columnIndex++ {
			// TODO: Embed into the following FieldElement initialization
			fieldObject := FieldObject{
				Class: "empty",
			}
			row[columnIndex] = FieldElement{
				Y: rowIndex,
				X: columnIndex,
				Object: fieldObject,
			}
		}
		matrix[rowIndex] = row
	}
	return matrix
}

func measureFieldMatrixY(fieldMatrix FieldMatrix) int {
	return len(fieldMatrix)
}

func measureFieldMatrixX(fieldMatrix FieldMatrix) int {
	return len(fieldMatrix[0])
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
	y := measureFieldMatrixY(fieldMatrix)
	x := measureFieldMatrixX(fieldMatrix)
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
	output := render(state)
	fmt.Println(output)
}
