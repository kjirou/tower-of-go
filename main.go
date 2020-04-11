package main

// TODO:
// - Write unit test for each packages.

import (
	"fmt"
	"github.com/kjirou/tower_of_go/models"
	"github.com/kjirou/tower_of_go/reducers"
	"github.com/kjirou/tower_of_go/utils"
	"github.com/kjirou/tower_of_go/views"
	"github.com/nsf/termbox-go"
	"math/rand"
	"os"
	"time"
)

type Controller struct {
	state *models.State
	screen *views.Screen
}

func (controller *Controller) GetState() *models.State {
	return controller.state
}

func (controller *Controller) GetScreen() *views.Screen {
	return controller.screen
}

func (controller *Controller) Dispatch(newState *models.State) {
	controller.state = newState
	controller.screen.Render(controller.state)
}

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

func initializeTermbox() error {
	termboxErr := termbox.Init()
	if termboxErr != nil {
		return termboxErr
	}
	termbox.SetInputMode(termbox.InputEsc)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	return nil
}

// TODO: Replace `ch` type with termbox's `Cell.Ch` type.
func handleKeyPress(controller *Controller, ch rune, key termbox.Key) {
	var newState *models.State
	var stateChanged bool = false
	var err error
	state := controller.GetState()

	switch {
	// Start or restart a game.
	case ch == 's':
		newState, stateChanged, err = reducers.StartOrRestartGame(*state)
	// Move the hero.
	case key == termbox.KeyArrowUp || ch == 'k':
		newState, stateChanged, err = reducers.WalkHero(*state, utils.FourDirectionUp)
	case key == termbox.KeyArrowRight || ch == 'l':
		newState, stateChanged, err = reducers.WalkHero(*state, utils.FourDirectionRight)
	case key == termbox.KeyArrowDown || ch == 'j':
		newState, stateChanged, err = reducers.WalkHero(*state, utils.FourDirectionDown)
	case key == termbox.KeyArrowLeft || ch == 'h':
		newState, stateChanged, err = reducers.WalkHero(*state, utils.FourDirectionLeft)
	}

	if err != nil {
		panic(err)
	} else if newState != nil && stateChanged {
		controller.Dispatch(newState)
		drawTerminal(controller.GetScreen())
	}
}

func handleTermboxEvents(controller *Controller) {
	didQuitApplication := false

	for !didQuitApplication {
		event := termbox.PollEvent()
		switch event.Type {
		case termbox.EventKey:
			// Quit the application. Only this operation is resolved with priority.
			if event.Key == termbox.KeyCtrlC || event.Key == termbox.KeyCtrlQ {
				didQuitApplication = true
				break
			}

			handleKeyPress(controller, event.Ch, event.Key)
		}
	}
}

func runMainLoop(controller *Controller) {
	// About 60fps.
	// TODO: Some delay from real time.
	interval := time.Millisecond * 17
	for {
		var newState *models.State
		var stateChanged bool = false
		var err error
		state := controller.GetState()

		time.Sleep(interval)

		newState, stateChanged, err = reducers.AdvanceTime(*state, interval)

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

	state := models.CreateState()
	err := state.SetWelcomeData()
	if err != nil {
		panic(err)
	}

	screen := views.CreateScreen(24, 80)

	controller := Controller{
		state: &state,
		screen: &screen,
	}

	controller.Dispatch(&state)

	if debugMode {
		fmt.Println(screen.AsText())
	} else {
		termboxErr := initializeTermbox()
		if termboxErr != nil {
			panic(termboxErr)
		}
		defer termbox.Close()
		drawTerminal(&screen)
		go runMainLoop(&controller)
		handleTermboxEvents(&controller)
	}
}
