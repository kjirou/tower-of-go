package utils

import(
	"github.com/pkg/errors"
	"math/rand"
)

type MazeCellContent int
const (
	MazeCellContentBreakableWall MazeCellContent = iota
	MazeCellContentEmpty
	MazeCellContentUnbreakableWall
)

type mazeCell struct {
	ClusterIndex int
	Content MazeCellContent
	X int
	Y int
}

func generateRawMazeMatrix(rowLength int, columnLength int) ([][]*mazeCell, error) {
	cells := make([][]*mazeCell, rowLength)

	if rowLength < 3 || columnLength < 3 {
		return cells, errors.Errorf("The number of rows and columns must be at least 3.")
	} else if (rowLength % 2 != 1 || columnLength % 2 != 1) {
		return cells, errors.Errorf("The number of rows and columns should be 2n+1.")
	}

	clusterIndex := 0
	for y := 0; y < rowLength; y++ {
		row := make([]*mazeCell, columnLength)
		for x := 0; x < columnLength; x++ {
			content := MazeCellContentUnbreakableWall
			if (y%2 == 1 && x%2 == 1) {
				content = MazeCellContentEmpty
			} else if (
				y != 0 && y != rowLength-1 &&
				x != 0 && x != columnLength-1 &&
				(y%2 == 0 && x%2 == 1 || y%2 == 1 && x%2 == 0)) {
				content = MazeCellContentBreakableWall
			}
			cell := mazeCell{
				Content: content,
				ClusterIndex: clusterIndex,
				Y: y,
				X: x,
			}
			row[x] = &cell
			clusterIndex++
		}
		cells[y] = row
	}
	return cells, nil
}

// Generate a maze with the clustering method.
//
// The maze generation algorithm referred to the following article.
// https://qiita.com/kaityo256/items/b2e504c100f4274deb42
//
// For example, if set rowLength=5 and columnLength=7 then a maze of the following size is generated.
// #######
// #     #
// #     #
// #     #
// #######
func GenerateMaze(rowLength int, columnLength int) ([][]*mazeCell, error) {
	cells, err := generateRawMazeMatrix(rowLength, columnLength)
	if err != nil {
		return cells, err
	}

	breakableWalls := make([]*mazeCell, 0)
	for _, row := range cells {
		for _, cell := range row {
			if cell.Content == MazeCellContentBreakableWall {
				breakableWalls = append(breakableWalls, cell)
			}
		}
	}

	rand.Shuffle(len(breakableWalls), func (i, j int) {
		breakableWalls[i], breakableWalls[j] = breakableWalls[j], breakableWalls[i]
	})

	for _, breakableWall := range breakableWalls {
		var a *mazeCell
		var b *mazeCell
		upperCell := cells[breakableWall.Y - 1][breakableWall.X]
		//
		// # = MazeCellContentUnbreakableWall
		// * = MazeCellContentBreakableWall
		// @ = MazeCellContentBreakableWall in this loop
		//
		// #*#
		// *a*
		// #@#
		// *b*
		// #*#
		//
		if upperCell.Content == MazeCellContentEmpty {
			a = upperCell
			b = cells[breakableWall.Y + 1][breakableWall.X]
		//
		// #*#*#
		// #b@a#
		// #*#*#
		//
		} else {
			a = cells[breakableWall.Y][breakableWall.X + 1]
			b = cells[breakableWall.Y][breakableWall.X - 1]
		}

		if a.ClusterIndex != b.ClusterIndex {
			aci := a.ClusterIndex
			bci := b.ClusterIndex
			breakableWall.Content = MazeCellContentEmpty
			breakableWall.ClusterIndex = aci
			for _, row := range cells {
				for _, cell := range row {
					if cell.ClusterIndex == bci {
						cell.ClusterIndex = aci
					}
				}
			}
		} else {
			breakableWall.Content = MazeCellContentUnbreakableWall
		}
	}

	return cells, nil
}
