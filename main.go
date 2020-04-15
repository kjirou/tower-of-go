package main

import (
	"flag"
	"fmt"
	"github.com/kjirou/tower-of-go/controller"
	"github.com/kjirou/tower-of-go/views"
	"github.com/nsf/termbox-go"
	"math/rand"
	"time"
)

func drawTerminal(screen *views.Screen) {
	screen.ForEachCells(func (y int, x int, symbol rune, fg termbox.Attribute, bg termbox.Attribute) {
		// おそらく termbox.SetCell() の終了は非同期で、
		//   同 cell へ高速で重複して出力した場合に互いの出力バッファが入れ子になる。
		// メインループとキー操作それぞれで本関数を実行していたら、
		//   頻繁に壊れた ANSI の破片のような文字列が出力されていたことからの推測である。
		// メインループのみで描画していれば、問題なさそう。
		termbox.SetCell(x, y, symbol, fg, bg)
	})
	termbox.Flush()
}

func convertScreenToText(screen *views.Screen) string {
	output := ""
	lastY := 0
	screen.ForEachCells(func (y int, x int, symbol rune, fg termbox.Attribute, bg termbox.Attribute) {
		if y != lastY {
			output += "\n"
			lastY = y
		}
		output += string(symbol)
	})
	return output
}

func runMainLoop(controller *controller.Controller) {
	// TODO: 現実時間よりゲーム時間の進みが遅い。
	//       手元環境だと、ゲーム時間 30 秒に対して現実時間 36 秒というずれになっている。
	//
	//       HandleMainLoop が同期的に止めている時間かと思ったが、HandleMainLoop 呼び出しの直前から
	//         for の最後までを time.Now() の差分で計測したら、800 マイクロ秒程度だった。
	//
	//       一方で time.Sleep の呼びだしの時間を計測すると、17000-20000マイクロ秒程度あった。
	//       こちらのずれの方が大きい。このずれの理由はまだ不明。
	//
	//       なお、fps を落とすとその分現実時間に近づく。手元環境だと 25fps ならほぼ現実時間に等しくなる。

	for {
		// Expecting 60fps. However, it is behind the real time.
		// For example, my computer needs 33-36 seconds of real time for 30 seconds of a game.
		interval := controller.CalculateIntervalToNextMainLoop(time.Now())
		time.Sleep(interval)

		newState, err := controller.HandleMainLoop(interval)

		if err != nil {
			panic(err)
		} else if newState != nil {
			controller.Dispatch(newState)
			drawTerminal(controller.GetScreen())
		}
	}
}

func main() {
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "Runs with debug mode.")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	controller, createControllerErr := controller.CreateController()
	if createControllerErr != nil {
		panic(createControllerErr)
	}

	if debugMode {
		fmt.Println(convertScreenToText(controller.GetScreen()))
	} else {
		termboxErr := termbox.Init()
		if termboxErr != nil {
			panic(termboxErr)
		}
		termbox.SetInputMode(termbox.InputEsc)
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		defer termbox.Close()
		drawTerminal(controller.GetScreen())
		go runMainLoop(controller)
		// Observe termbox events.
		didQuitApplication := false
		for !didQuitApplication {
			event := termbox.PollEvent()
			switch event.Type {
			case termbox.EventKey:
				// Quit the application. Only this operation is resolved with priority.
				if event.Key == termbox.KeyEsc || event.Key == termbox.KeyCtrlC || event.Key == termbox.KeyCtrlQ {
					didQuitApplication = true
					break
				}
				controller.HandleKeyPress(event.Ch, event.Key)
			}
		}
	}
}
