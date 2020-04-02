package main

// TODO:
// - gofmt

import (
	"fmt"
	"strings"
)

type FieldElement struct {
	Symbol string
	X int
	Y int
}

type FieldMatrix [][]FieldElement

func createFieldMatrix(y int, x int) FieldMatrix {
	matrix := make(FieldMatrix, y)
	for rowIndex := 0; rowIndex < y; rowIndex++ {
		row := make([]FieldElement, x)
		for columnIndex := 0; columnIndex < x; columnIndex++ {
			row[columnIndex] = FieldElement{
				Y: rowIndex,
				X: columnIndex,
				Symbol: ".",
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

func renderFieldElement(fieldElement FieldElement) string {
	return fieldElement.Symbol
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

func main() {
	fieldMatrix := createFieldMatrix(12, 20)
	output := renderFieldMatrix(fieldMatrix)
	fmt.Println(output)
}
