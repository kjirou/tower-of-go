package main

// TODO:
// - gofmt

import (
	"fmt"
	"strings"
)

type FieldObject struct {
	// TODO: Enumerize
	Class string
}

func (fo FieldObject) render() string {
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

type FieldElement struct {
	Object FieldObject
	X int
	Y int
}

func (fe FieldElement) render() string {
	output := fe.Object.render()
	if output == "" {
		output = "."
	}
	return output
}

type FieldMatrix [][]FieldElement

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

func renderFieldMatrix(fieldMatrix FieldMatrix) string {
	y := measureFieldMatrixY(fieldMatrix)
	x := measureFieldMatrixX(fieldMatrix)
	lines := make([]string, y)
	for rowIndex := 0; rowIndex < y; rowIndex++ {
		line := ""
		// TODO: Use mapping method
		for columnIndex := 0; columnIndex < x; columnIndex++ {
			line += fieldMatrix[rowIndex][columnIndex].render()
		}
		lines[rowIndex] = line
	}
	return strings.Join(lines, "\n")
}

func main() {
	fieldMatrix := createFieldMatrix(12, 20)
	output := renderFieldMatrix(fieldMatrix)
	fmt.Println(output)
}
