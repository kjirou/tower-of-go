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
		// NOTE: Probably, termbox.SetCell's outputs are asynchronous.
		//       Therefore, multiple executions on the same cell at the same time will nest the output buffers.
		//       As a result, the display will be corrupted.
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
	for {
		// TODO: Expecting 60fps. However, it is behind the real time.
		//       For example, my computer needs 33-36 seconds of real time for 30 seconds of a game.
		//       It becomes more accurate if fps is lowered.
		interval := controller.CalculateIntervalToNextMainLoop(time.Now())
		time.Sleep(interval)

		newState, err := controller.HandleMainLoop(interval)

		if err != nil {
			errMessage, _ := fmt.Printf("%+v", err)
			panic(errMessage)
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
