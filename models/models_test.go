package models

import (
	"fmt"
	"github.com/kjirou/tower_of_go/utils"
	"testing"
	"strings"
)

func TestField(t *testing.T) {
	t.Run("At", func(t *testing.T) {
		field := createField(2, 3)

		t.Run("指定した位置の要素を取得できる", func(t *testing.T) {
			var position utils.IMatrixPosition = &utils.MatrixPosition{Y: 1, X: 2}
			element, _ := field.At(position)
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
					var position utils.IMatrixPosition = &utils.MatrixPosition{Y: tc.Y, X: tc.X}
					_, err := field.At(position)
					if err == nil {
						t.Fatal("エラーを返さない")
					}
				})
			}
		})
	})

	t.Run("MoveObject", func(t *testing.T) {
		field := createField(2, 3)
		var fromPosition utils.IMatrixPosition = &utils.MatrixPosition{Y: 0, X: 0}
		var toPosition utils.IMatrixPosition = &utils.MatrixPosition{Y: 1, X: 2}
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
}

func TestGame(t *testing.T) {
	t.Run("GetPlaytimeAsSeconds", func(t *testing.T) {
		game := &Game{}

		t.Run("初期化直後は0を返す", func(t *testing.T) {
			game.Initialize()
			if game.GetPlaytimeAsSeconds() != 0 {
				t.Fatal("0ではない")
			}
		})
	})

	t.Run("GetPlaytimeAsString", func(t *testing.T) {
		game := &Game{}

		t.Run("初期化直後は\"0\"を返す", func(t *testing.T) {
			game.Initialize()
			if game.GetPlaytimeAsString() != "0" {
				t.Fatal("\"0\"ではない")
			}
		})
	})

	t.Run("Start", func(t *testing.T) {
		game := &Game{}

		t.Run("It works", func(t *testing.T) {
			game.Start()
			if game.IsStarted() != true {
				t.Fatal("開始していない")
			}
			if game.IsFinished() != false {
				t.Fatal("終了していない")
			}
		})
	})
}
