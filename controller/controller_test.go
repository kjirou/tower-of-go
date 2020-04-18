package controller

import (
	"testing"
	"time"
)

func TestController_CalculateIntervalToNextMainLoop_NotTD(t *testing.T) {
	controller := &Controller{}

	t.Run("以前に一度もこの関数を呼び出していないとき、lastMainLoopRanAtはゼロ値である", func(t *testing.T) {
		if !controller.lastMainLoopRanAt.IsZero() {
				t.Fatal("ゼロ値ではない")
		}
	})

	t.Run("以前に一度もこの関数を呼び出していないとき、16666マイクロ秒を返す", func(t *testing.T) {
		lastMainLoopRanAt := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
		interval := controller.CalculateIntervalToNextMainLoop(lastMainLoopRanAt)
		if interval != time.Microsecond*16666 {
				t.Fatal("時間が違う")
		}
	})

	t.Run("前回呼び出したときに指定した引数をlastMainLoopRanAtへ格納する", func(t *testing.T) {
		lastMainLoopRanAt := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
		controller.CalculateIntervalToNextMainLoop(lastMainLoopRanAt)
		if !controller.lastMainLoopRanAt.Equal(lastMainLoopRanAt) {
				t.Fatal("前回呼び出した引数と時間が違う")
		}
	})

	t.Run("前回呼び出した時刻と今回の時刻の差が16666マイクロ秒未満のとき、16666マイクロ秒を返す", func(t *testing.T) {
		lastMainLoopRanAt := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
		controller.CalculateIntervalToNextMainLoop(lastMainLoopRanAt)
		now := time.Date(2000, time.January, 1, 0, 0, 0, 16665999, time.UTC)
		interval := controller.CalculateIntervalToNextMainLoop(now)
		if interval != time.Microsecond*16666 {
				t.Fatal("時間が違う")
		}
	})

	t.Run("前回呼び出した時刻と今回の時刻の差がとても長いとき、8333マイクロ秒を返す", func(t *testing.T) {
		lastMainLoopRanAt := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
		controller.CalculateIntervalToNextMainLoop(lastMainLoopRanAt)
		now := time.Date(2000, time.January, 1, 0, 0, 1, 0, time.UTC)
		interval := controller.CalculateIntervalToNextMainLoop(now)
		if interval != time.Microsecond*8333 {
				t.Fatal("時間が違う")
		}
	})

	t.Run("前回呼び出した時刻と今回の時刻の差が16666マイクロ秒を超えるとき、超えた分を差し引いた時間を返す", func(t *testing.T) {
		lastMainLoopRanAt := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
		controller.CalculateIntervalToNextMainLoop(lastMainLoopRanAt)
		now := time.Date(2000, time.January, 1, 0, 0, 0, 17666000, time.UTC)
		interval := controller.CalculateIntervalToNextMainLoop(now)
		if interval != time.Nanosecond*15666000 {
				t.Fatal("時間が違う")
		}
	})
}
