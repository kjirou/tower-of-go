package main

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
	"strings"
)

func mapFieldElementToScreenCellProps(fieldElement utils.IFieldElement) *views.ScreenCellProps {
	symbol := '.'
	fg := termbox.ColorWhite
	bg := termbox.ColorBlack
	if !fieldElement.IsObjectEmpty() {
		switch fieldElement.GetObjectClass() {
		case "hero":
			symbol = '@'
			fg = termbox.ColorMagenta
		case "wall":
			symbol = '#'
			fg = termbox.ColorYellow
		default:
			symbol = '?'
		}
	} else {
		switch fieldElement.GetFloorObjectClass() {
		case "upstairs":
			symbol = '<'
			fg = termbox.ColorGreen
		}
	}
	return &views.ScreenCellProps{
		Symbol: symbol,
		ForegroundColor: fg,
		BackgroundColor: bg,
	}
}

func mapStateModelToScreenProps(state *models.State) *views.ScreenProps {
	game := state.GetGame()
	field := state.GetField()

	// Cells of the field.
	fieldRowLength := field.MeasureRowLength()
	fieldColumnLength := field.MeasureColumnLength()
	fieldCells := make([][]*views.ScreenCellProps, fieldRowLength)
	for y := 0; y < fieldRowLength; y++ {
		cellsRow := make([]*views.ScreenCellProps, fieldColumnLength)
		for x := 0; x < fieldColumnLength; x++ {
			var fieldElementPosition utils.IMatrixPosition = &utils.MatrixPosition{Y: y, X: x}
			// TODO: Error handling.
			var fieldElement, _ = field.At(fieldElementPosition)
			cellsRow[x] = mapFieldElementToScreenCellProps(fieldElement)
		}
		fieldCells[y] = cellsRow
	}

	// Lank message.
	lankMessage := ""
	lankMessageForeground := termbox.ColorWhite
	if game.IsFinished() {
		score := game.GetFloorNumber()
		switch {
			case score == 3:
				lankMessage = "Good!"
				lankMessageForeground = termbox.ColorGreen
			case score == 4:
				lankMessage = "Excellent!"
				lankMessageForeground = termbox.ColorGreen
			case score == 5:
				lankMessage = "Marvelous!"
				lankMessageForeground = termbox.ColorGreen
			case score >= 6:
				lankMessage = "Gopher!!"
				lankMessageForeground = termbox.ColorCyan
			default:
				lankMessage = "No good..."
		}
	}

	return &views.ScreenProps{
		FieldCells: fieldCells,
		RemainingTime: game.CalculateRemainingTime(state.GetExecutionTime()).Seconds(),
		FloorNumber: game.GetFloorNumber(),
		LankMessage: lankMessage,
		LankMessageForeground: lankMessageForeground,
	}
}

type Controller struct {
	state  *models.State
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
	screenProps := mapStateModelToScreenProps(controller.state)
	controller.screen.Render(screenProps)
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
			if event.Key == termbox.KeyEsc || event.Key == termbox.KeyCtrlC || event.Key == termbox.KeyCtrlQ {
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
		state:  &state,
		screen: &screen,
	}

	controller.Dispatch(&state)

	if debugMode {
		fmt.Println(convertScreenToText(&screen))
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
