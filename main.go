package main

import (
	"fmt"
	"github.com/kjirou/tower_of_go/controller"
	"github.com/kjirou/tower_of_go/views"
	"github.com/nsf/termbox-go"
	"math/rand"
	"os"
	"time"
	"strings"
)

func drawTerminal(screen *views.Screen) {
	for y, row := range screen.GetMatrix() {
		for x, element := range row {
			// TODO: おそらく termbox.SetCell() の終了は非同期で、
			//         同 cell へ高速で重複して出力した場合に互いの出力バッファが入れ子になってしまう。
			//       たまに、端末の画面に ANSI の破片らしき文字列がが出力されることがあることが根拠。
			//       現在は、メインループによる再描画とキー操作による再描画が重なることがあり、
			//         それで高速に重複して同じ cell へ出力することがある。
			//       解決案は、キー操作による状態更新もメインループで解決するようにすることで、
			//         これでおそらくはほとんど発生しなくなると思う。
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

func handleTermboxEvents(controller *controller.Controller) {
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
				drawTerminal(controller.GetScreen())
			}
		}
	}
}

func runMainLoop(controller *controller.Controller) {
	// About 60fps.
	// TODO: Some delay from real time.
	interval := time.Millisecond * 17
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
		handleTermboxEvents(controller)
	}
}
