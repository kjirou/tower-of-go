package models

import (
	"fmt"
	"github.com/kjirou/tower_of_go/utils"
	"testing"
	"time"
	"strings"
)

func TestField(t *testing.T) {
	t.Run("At", func(t *testing.T) {
		field := createField(2, 3)

		t.Run("指定した位置の要素を取得できる", func(t *testing.T) {
			element, _ := field.At(&utils.MatrixPosition{Y: 1, X: 2})
			if element.GetPosition().GetY() != 1 {
				t.Fatal("Y が違う")
			} else if element.GetPosition().GetX() != 2 {
				t.Fatal("X が違う")
			}
		})

		t.Run("存在しない位置を指定したとき", func(t *testing.T) {
			type testCase struct{
				Y int
				X int
			}
			var testCases []testCase
			testCases = append(testCases, testCase{Y: -1, X: 0})
			testCases = append(testCases, testCase{Y: 2, X: 0})
			testCases = append(testCases, testCase{Y: 0, X: -1})
			testCases = append(testCases, testCase{Y: 0, X: 3})
			for _, tc := range testCases {
				tc := tc
				t.Run(fmt.Sprintf("Y=%d,X=%dはエラーを返す", tc.Y, tc.X), func(t *testing.T) {
					_, err := field.At(&utils.MatrixPosition{Y: tc.Y, X: tc.X})
					if err == nil {
						t.Fatal("エラーを返さない")
					}
				})
			}
		})
	})

	t.Run("MoveObject", func(t *testing.T) {
		field := createField(2, 3)
		fromPosition := &utils.MatrixPosition{Y: 0, X: 0}
		toPosition := &utils.MatrixPosition{Y: 1, X: 2}
		fromElement, _ := field.At(fromPosition)
		toElement, _ := field.At(toPosition)

		t.Run("始点の物体が空ではなく、終点の物体が空のとき、物体種別が移動する", func(t *testing.T) {
			fromElement.UpdateObjectClass("wall")
			toElement.UpdateObjectClass("empty")
			field.MoveObject(fromPosition, toPosition)
			if toElement.GetObjectClass() != "wall" {
				t.Fatal("物体種別が移動していない")
			}
		})

		t.Run("始点の物体が空ではなく、終点の物体が空ではないとき、エラーを返す", func(t *testing.T) {
			fromElement.UpdateObjectClass("wall")
			toElement.UpdateObjectClass("wall")
			err := field.MoveObject(fromPosition, toPosition)
			if err == nil {
				t.Fatal("エラーを返さない")
			} else if !strings.Contains(err.Error(), "object exists") {
				t.Fatal("意図したエラーメッセージではない")
			}
		})

		t.Run("始点の物体が空のとき、エラーを返す", func(t *testing.T) {
			fromElement.UpdateObjectClass("empty")
			err := field.MoveObject(fromPosition, toPosition)
			if err == nil {
				t.Fatal("エラーを返さない")
			} else if !strings.Contains(err.Error(), "does not exist") {
				t.Fatal("意図したエラーメッセージではない")
			}
		})
	})

	t.Run("ResetMaze", func(t *testing.T) {
		t.Run("外周1マスは壁になる", func(t *testing.T) {
			field := createField(7, 7)
			field.ResetMaze()
			for y, row := range field.matrix {
				for x, element := range row {
					isTopOrBottomEdge := y == 0 || y == field.MeasureRowLength()-1
					isLeftOrRightEdge := x == 0 || x == field.MeasureColumnLength()-1
					if (isTopOrBottomEdge || isLeftOrRightEdge) && element.GetObjectClass() != "wall" {
						t.Fatalf("Y=%d, X=%d が壁ではない", y, x)
					}
				}
			}
		})

		t.Run("ヒーローが存在していたとき、ヒーローは削除される", func(t *testing.T) {
			field := createField(7, 7)
			element, err := field.At(HeroPosition)
			if err != nil {
				t.Fatal("ヒーローの配置に失敗する")
			}
			element.UpdateObjectClass("hero")
			field.ResetMaze()
			for _, row := range field.matrix {
				for _, element := range row {
					if element.GetObjectClass() == "hero" {
						t.Fatal("ヒーローが存在している")
					}
				}
			}
		})
	})
}

func TestGame(t *testing.T) {
	t.Run("CalculateRemainingTime", func(t *testing.T) {
		game := &Game{}

		t.Run("リセット直後は30を返す", func(t *testing.T) {
			game.Reset()
			executionTime, _ := time.ParseDuration("2s")
			remainingTime := game.CalculateRemainingTime(executionTime)
			if remainingTime.Seconds() != 30 {
				t.Fatal("30ではない")
			}
		})

		t.Run("最小で0を返す", func(t *testing.T) {
			game.Reset()
			startTime, _ := time.ParseDuration("1s")
			game.Start(startTime)
			currentTime, _ := time.ParseDuration("999s")
			remainingTime := game.CalculateRemainingTime(currentTime)
			if remainingTime.Seconds() != 0 {
				t.Fatal("0ではない")
			}
		})
	})

	t.Run("Start", func(t *testing.T) {
		game := &Game{}

		t.Run("It works", func(t *testing.T) {
			executionTime, _ := time.ParseDuration("0s")
			game.Start(executionTime)
			if game.IsStarted() {
				t.Fatal("開始している")
			}
			if game.IsFinished() {
				t.Fatal("終了している")
			}
		})
	})
}
