package main

import (
	"fmt"
	"github.com/kjirou/tower_of_go/controller"
	"github.com/kjirou/tower_of_go/views"
	"github.com/nsf/termbox-go"
	"math/rand"
	"os"
	"strings"
	"time"
)

func drawTerminal(screen *views.Screen) {
	for y, row := range screen.GetMatrix() {
		for x, element := range row {
			// おそらく termbox.SetCell() の終了は非同期で、
			//   同 cell へ高速で重複して出力した場合に互いの出力バッファが入れ子になる。
			// メインループとキー操作それぞれで本関数を実行していたら、
			//   頻繁に壊れた ANSI の破片のような文字列が出力されていたことからの推測である。
			// メインループのみで描画していれば、問題なさそう。
			termbox.SetCell(x, y, element.Symbol, element.ForegroundColor, element.BackgroundColor)
		}
	}
	termbox.Flush()
}

func convertScreenToText(screen *views.Screen) string {
	lines := make([]string, 0)
	for _, row := range screen.GetMatrix() {
		line := ""
		for _, element := range row {
			line += string(element.Symbol)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func runMainLoop(controller *controller.Controller) {
	// About 60fps.
	// TODO: Some delay from real time.
	interval := time.Microsecond * 16666
	for {
		time.Sleep(interval)

		newState, stateChanged, err := controller.HandleMainLoop(interval)

		if err != nil {
			panic(err)
		} else if newState != nil && stateChanged {
			controller.Dispatch(newState)
			drawTerminal(controller.GetScreen())
		}
	}
}

func observeTermboxEvents(controller *controller.Controller) {
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

			newState, stateChanged, err := controller.HandleKeyPress(event.Ch, event.Key)

			if err != nil {
				panic(err)
			} else if newState != nil && stateChanged {
				controller.Dispatch(newState)
			}
		}
	}
}

func main() {
	// TODO: Look for a tiny CLI argument parser like the "minimist" of Node.js.
	commandLineArgs := os.Args[1:]
	debugMode := false
	for _, arg := range commandLineArgs {
		if arg == "--debug-mode" || arg == "-d" {
			debugMode = true
		}
	}

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
		observeTermboxEvents(controller)
	}
}
