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

// It may need to make the following processes:
//   https://github.com/nsf/termbox-go/blob/4d2b513ad8bee47a9a5a65b0dee0182049a31916/_demos/keyboard.go#L669
//   (However, details cannot be read...)
// TODO: Replace `ch` type with termbox's `Cell.Ch` type.
func handleKeyPress(controller *Controller, ch rune, key termbox.Key) {
	var newState *models.State
	var err error
	state := controller.GetState()

	switch {
	// Start a game.
	case ch == 's':
		newState, err = reducers.StartGame(*state)
	// Move the hero.
	case key == termbox.KeyArrowUp || ch == 'k':
		newState, err = reducers.WalkHero(*state, utils.FourDirectionUp)
	case key == termbox.KeyArrowRight || ch == 'l':
		newState, err = reducers.WalkHero(*state, utils.FourDirectionRight)
	case key == termbox.KeyArrowDown || ch == 'j':
		newState, err = reducers.WalkHero(*state, utils.FourDirectionDown)
	case key == termbox.KeyArrowLeft || ch == 'h':
		newState, err = reducers.WalkHero(*state, utils.FourDirectionLeft)
	}

	if err != nil {
		panic(err)
	} else if newState != nil {
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

// TODO: 1fpsなのでループの周期によっては2秒待つことになり不自然。runGameLoopにして30fps程度にするのが楽そう。
func runTimer(controller *Controller) {
	interval := time.Second
	for {
		var newState *models.State
		var err error
		state := controller.GetState()
		game := state.GetGame()

		time.Sleep(interval)

		if game.IsStarted() && !game.IsFinished() {
			newState, err = reducers.AlterPlaytime(*state, interval)
		}

		if err != nil {
			panic(err)
		} else if newState != nil {
			controller.Dispatch(newState)
			drawTerminal(controller.GetScreen())
		}
	}
}

func main() {
	// TODO: Look for a tiny CLI argument parser like the "minimist" of Node.js.
	commandLineArgs := os.Args[1:]
	doesRunTermbox := false
	for _, arg := range commandLineArgs {
		if arg == "-t" {
			doesRunTermbox = true
		}
	}

	state := models.CreateState()
	err := state.InitializeDummyData()
	if err != nil {
		panic(err)
	}

	screen := views.CreateScreen(24, 80)

	controller := Controller{
		state: &state,
		screen: &screen,
	}

	controller.Dispatch(&state)

	if doesRunTermbox {
		termboxErr := initializeTermbox()
		if termboxErr != nil {
			panic(termboxErr)
		}
		defer termbox.Close()
		drawTerminal(&screen)
		go runTimer(&controller)
		handleTermboxEvents(&controller)
	} else {
		fmt.Println(screen.AsText())
	}
}
