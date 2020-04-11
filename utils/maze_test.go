package utils

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestGenerateRawMazeMatrix(t *testing.T) {
	t.Run("行数が3未満のときはエラーを返す", func(t *testing.T) {
		_, err := generateRawMazeMatrix(2, 3)
		if err == nil {
			t.Fatal("エラーを返さない")
		}
	})

	t.Run("列数が3未満のときはエラーを返す", func(t *testing.T) {
		_, err := generateRawMazeMatrix(3, 2)
		if err == nil {
			t.Fatal("エラーを返さない")
		}
	})

	t.Run("行数が4のときはエラーを返す", func(t *testing.T) {
		_, err := generateRawMazeMatrix(4, 3)
		if err == nil {
			t.Fatal("エラーを返さない")
		}
	})

	t.Run("列数が4のときはエラーを返す", func(t *testing.T) {
		_, err := generateRawMazeMatrix(3, 4)
		if err == nil {
			t.Fatal("エラーを返さない")
		}
	})

	t.Run("指定した行数と同じ長さの行列を返す", func(t *testing.T) {
		cells, _ := generateRawMazeMatrix(3, 5)
		if len(cells) != 3 {
			t.Fatal("行数が指定した値と異なる")
		}
	})

	t.Run("1行目が指定した列数と同じ長さの行列を返す", func(t *testing.T) {
		cells, _ := generateRawMazeMatrix(3, 5)
		if len(cells[0]) != 5 {
			t.Fatal("1行目の列数が指定した値と異なる")
		}
	})

	t.Run("各セルのContent", func(t *testing.T) {
		t.Run("外周1セルは全て壊せない壁である", func(t *testing.T) {
			cells, _ := generateRawMazeMatrix(5, 5)
			testCases := [][]int{
				[]int{0, 0},
				[]int{0, 1},
				[]int{0, 2},
				[]int{0, 3},
				[]int{0, 4},
				[]int{1, 0},
				[]int{1, 4},
				[]int{2, 0},
				[]int{2, 4},
				[]int{3, 0},
				[]int{3, 4},
				[]int{4, 0},
				[]int{4, 1},
				[]int{4, 2},
				[]int{4, 3},
				[]int{4, 4},
			}
			for _, testCase := range testCases {
				y := testCase[0]
				x := testCase[1]
				t.Run(fmt.Sprintf("Y=%d, X=%d は壊せない壁である", y, x), func(t *testing.T) {
					if cells[y][x].Content != MazeCellContentUnbreakableWall {
						t.Fatal("壊せない壁ではない")
					}
				})
			}
		})

		t.Run("Y=2n+1, X=2n+1 の位置のセルは空である", func(t *testing.T) {
			cells, _ := generateRawMazeMatrix(5, 5)
			testCases := [][]int{
				[]int{1, 1},
				[]int{1, 3},
				[]int{3, 1},
				[]int{3, 3},
			}
			for _, testCase := range testCases {
				y := testCase[0]
				x := testCase[1]
				t.Run(fmt.Sprintf("Y=%d, X=%d は空である", y, x), func(t *testing.T) {
					if cells[y][x].Content != MazeCellContentEmpty {
						t.Fatal("空ではない")
					}
				})
			}
		})

		t.Run("外周1セルを除く Y=2n+1, X=2n または Y=2n, X=2n+1 の位置のセルは壊せる壁である", func(t *testing.T) {
			cells, _ := generateRawMazeMatrix(5, 5)
			testCases := [][]int{
				[]int{1, 2},
				[]int{2, 1},
				[]int{2, 3},
				[]int{3, 2},
			}
			for _, testCase := range testCases {
				y := testCase[0]
				x := testCase[1]
				t.Run(fmt.Sprintf("Y=%d, X=%d は壊せる壁である", y, x), func(t *testing.T) {
					if cells[y][x].Content != MazeCellContentBreakableWall {
						t.Fatal("壊せる壁ではない")
					}
				})
			}
		})
	})

	t.Run("各セルのClusterIndex", func(t *testing.T) {
		cells, _ := generateRawMazeMatrix(5, 5)

		t.Run("重複していない", func(t *testing.T) {
			clusterIndexes := make(map[int]bool)
			for _, row := range cells {
				for _, cell := range row {
					clusterIndex := cell.ClusterIndex
					_, ok := clusterIndexes[clusterIndex]
					if ok {
						t.Fatal("重複した値がある")
					}
					clusterIndexes[clusterIndex] = true
				}
			}
		})
	})
}

// 指定セルから上下左右に一歩ずつ空セルを探して進み、踏んだセルリストを返す。
// 走査が枝分かれしても、後戻りしなければ全ての空セルを踏んで返せる。
// また誤った迷路の生成をして、道が循環しているところがある場合、無限ループで落ちる。
func exploreMaze(
	cells [][]*mazeCell, currentCell *mazeCell, beforeCell *mazeCell, steppedCells []*mazeCell) []*mazeCell {
	steppedCells = append(steppedCells, currentCell)
	fourDirections := []struct{
		deltaY int
		deltaX int
	}{
		{deltaY: -1, deltaX: 0},
		{deltaY: 0, deltaX: 1},
		{deltaY: 1, deltaX: 0},
		{deltaY: 0, deltaX: -1},
	}
	for _, direction := range fourDirections {
		nextCell := cells[currentCell.Y + direction.deltaY][currentCell.X + direction.deltaX]
		if (nextCell != beforeCell && nextCell.Content == MazeCellContentEmpty) {
			steppedCells = exploreMaze(cells, nextCell, currentCell, steppedCells)
		}
	}
	return steppedCells
}

func TestGenerateMaze(t *testing.T) {
	var seed int64 = time.Now().UnixNano()

	t.Run("クラスタリングによる迷路を生成していること", func(t *testing.T) {
		testCases := []struct{
			columnLength int
			rowLength int
		}{
			{rowLength: 3, columnLength: 3},
			{rowLength: 5, columnLength: 3},
			{rowLength: 3, columnLength: 5},
			{rowLength: 7, columnLength: 7},
			{rowLength: 13, columnLength: 21},
			{rowLength: 21, columnLength: 13},
			{rowLength: 21, columnLength: 21},
			{rowLength: 51, columnLength: 51},
		}
		for _, testCase := range testCases {
			title := fmt.Sprintf("行%d*列%dの迷路を生成するとき", testCase.rowLength, testCase.columnLength)
			t.Run(title, func(t *testing.T) {
				seed++
				rand.Seed(seed)
				cells, _ := GenerateMaze(testCase.rowLength, testCase.columnLength)

				t.Run("行と1行目の列の数が指定した値と等しい", func(t *testing.T) {
					if len(cells) != testCase.rowLength {
						t.Fatal("行の数が違う")
					}
					if len(cells[0]) != testCase.columnLength {
						t.Fatal("列の数が違う")
					}
				})

				t.Run("壊せる壁は存在しない", func(t *testing.T) {
					for _, row := range cells {
						for _, cell := range row {
							if cell.Content == MazeCellContentBreakableWall {
								t.Fatalf("Y=%d,X=%d は壊せる壁である", cell.Y, cell.X)
							}
						}
					}
				})

				t.Run("正しい迷路である", func(t *testing.T) {
					noBeforeMazeCell := mazeCell{}
					steppedCells := exploreMaze(cells, cells[1][1], &noBeforeMazeCell, make([]*mazeCell, 0))

					emptyCellCount := 0
					for _, row := range cells {
						for _, cell := range row {
							if cell.Content == MazeCellContentEmpty {
								emptyCellCount++
							}
						}
					}
					steppedEmptyCellCount := 0
					for _, cell := range steppedCells {
						if cell.Content == MazeCellContentEmpty {
							steppedEmptyCellCount++
						}
					}
					if emptyCellCount != steppedEmptyCellCount {
						t.Fatal("全ての空セルが結合されていない")
					}
				})
			})
		}
	})
}
