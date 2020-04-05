package main

// TODO:
// - Write unit test for each packages.

import (
	"fmt"
	"github.com/kjirou/tower_of_go/models"
	"github.com/kjirou/tower_of_go/utils"
	"github.com/kjirou/tower_of_go/views"
	"github.com/nsf/termbox-go"
	"os"
)

func drawTerminal(screen *views.Screen) {
	for y, row := range screen.GetMatrix() {
		for x, element := range row {
			termbox.SetCell(x, y, element.Symbol, element.ForegroundColor, element.BackgroundColor)
		}
	}
	termbox.Flush()
}

func initializeTermbox(screen *views.Screen) error {
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
func handleKeyPress(state *models.State, screen *views.Screen, ch rune, key termbox.Key) {
	var err error
	field := state.GetField()
	stateChanged := false

	// Move the hero.
	// TODO: Consider arrow keys.
	if ch == 'k' {
		err = field.WalkHero(utils.FourDirectionUp)
		stateChanged = true
	} else if ch == 'l' {
		err = field.WalkHero(utils.FourDirectionRight)
		stateChanged = true
	} else if ch == 'j' {
		err = field.WalkHero(utils.FourDirectionDown)
		stateChanged = true
	} else if ch == 'h' {
		err = field.WalkHero(utils.FourDirectionLeft)
		stateChanged = true
	}

	if err != nil {
		panic(err)
	}

	if stateChanged {
		screen.Render(state)
		drawTerminal(screen)
	}
}

func handleTermboxEvents(state *models.State, screen *views.Screen) {
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

			handleKeyPress(state, screen, event.Ch, event.Key)
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

	screen := views.CreateScreen(24+2, 80+2)
	screen.Render(&state)

	if doesRunTermbox {
		termboxErr := initializeTermbox(&screen)
		if termboxErr != nil {
			panic(termboxErr)
		}
		defer termbox.Close()
		drawTerminal(&screen)
		handleTermboxEvents(&state, &screen)
	} else {
		fmt.Println(screen.AsText())
	}
}
