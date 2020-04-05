package main

// TODO:
// - Separate the main package to sub local packages.
//   e.g.) Screeen -> views.Screen
// - Summarize "go run a.go b.go (...more local files!)" into "go run (--option) main.go".

import (
	"fmt"
	"github.com/kjirou/tower_of_go/utils"
	"github.com/nsf/termbox-go"
	"os"
)

func drawTerminal(screen *Screen) {
	for y, row := range screen.matrix {
		for x, element := range row {
			termbox.SetCell(x, y, element.Symbol, element.ForegroundColor, element.BackgroundColor)
		}
	}
	termbox.Flush()
}

func initializeTermbox(screen *Screen) error {
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
func handleKeyPress(state *State, screen *Screen, ch rune, key termbox.Key) {
	var err error
	field := &state.Field
	stateChanged := false

	// Move the hero.
	// TODO: Consider arrow keys.
	if ch == 'k' {
		err = field.WalkHero(FourDirectionUp)
		stateChanged = true
	} else if ch == 'l' {
		err = field.WalkHero(FourDirectionRight)
		stateChanged = true
	} else if ch == 'j' {
		err = field.WalkHero(FourDirectionDown)
		stateChanged = true
	} else if ch == 'h' {
		err = field.WalkHero(FourDirectionLeft)
		stateChanged = true
	}

	if err != nil {
		panic(err)
	}

	if stateChanged {
		screen.render(state)
		drawTerminal(screen)
	}
}

func handleTermboxEvents(state *State, screen *Screen) {
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

	state := State{
		Field: createField(12, 20),
	}

	// Dummy data
	var heroPosition utils.MatrixPosition = &FieldPosition{y: 2, x: 5}
	state.Field.At(heroPosition).UpdateObjectClass("hero")
	fieldRowLength := state.Field.MeasureRowLength()
	fieldColumnLength := state.Field.MeasureColumnLength()
	for y := 0; y < fieldRowLength; y++ {
		for x := 0; x < fieldColumnLength; x++ {
			isTopOrBottomEdge := y == 0 || y == fieldRowLength-1
			isLeftOrRightEdge := x == 0 || x == fieldColumnLength-1
			if isTopOrBottomEdge || isLeftOrRightEdge {
				var wallPosition utils.MatrixPosition = &FieldPosition{y: y, x: x}
				state.Field.At(wallPosition).UpdateObjectClass("wall")
			}
		}
	}

	screen := createScreen(24+2, 80+2)
	screen.render(&state)

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
