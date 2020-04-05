package models

import (
	"fmt"
	"github.com/kjirou/tower_of_go/utils"
	"testing"
)

func TestField(t *testing.T) {
	t.Run("At", func(t *testing.T) {
		field := createField(2, 3)

		t.Run("指定した位置の要素を取得できる", func(t *testing.T) {
			var position utils.IMatrixPosition = &utils.MatrixPosition{Y: 1, X: 2}
			element, err := field.At(position)
			if err != nil {
				t.Error(err)
			} else if element.GetPosition().GetY() != 1 {
				t.Error("Y が違う")
			} else if element.GetPosition().GetX() != 2 {
				t.Error("X が違う")
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
						t.Error("エラーを返さない")
					}
				})
			}
		})
	})
}
