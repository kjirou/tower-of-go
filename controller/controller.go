package controller

//
// NOTE: 全体的な設計について。
//
// Inputs   = 経過時間とキー入力が本アプリケーションが認識する外部入力である。
//   |
// Reducers = いち Reducer 関数は、その関数が対応する Input と現在の Models を組み合わせて、
//   |          次の Models を生成して返す。
// Models   = 本アプリケーションの正規化された状態である。
//   |        これ以下の状態は全て Models の写像として生成される。
// Props    = Views 側が要求する Views への更新クエリである。
//   |        Models の写像として生成される。
// Views    = 端末出力を抽象化した層である。
//   |        termbox へ渡すためのセルの矩形の集合を生成することが目的である。
// Outputs  = 基本的には termbox 経由で端末へ出力する。
//            デバッグ用に標準出力をすることもある。
//

import (
	"github.com/kjirou/tower-of-go/models"
	"github.com/kjirou/tower-of-go/utils"
	"github.com/kjirou/tower-of-go/reducers"
	"github.com/kjirou/tower-of-go/views"
	"github.com/nsf/termbox-go"
	"time"
)

func mapFieldElementToScreenCellProps(fieldElement *models.FieldElement) *views.ScreenCellProps {
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
		Foreground: fg,
		Background: bg,
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
			fieldElement, _ := field.At(&utils.MatrixPosition{Y: y, X: x})
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
	inputtedCharacter rune
	inputtedKey termbox.Key
	state  *models.State
	screen *views.Screen
}

func (controller *Controller) GetScreen() *views.Screen {
	return controller.screen
}

func (controller *Controller) SetState(state *models.State) {
	controller.state = state
}

func (controller *Controller) setKeyInputs(ch rune, key termbox.Key) {
	controller.inputtedCharacter = ch
	controller.inputtedKey = key
}

func (controller *Controller) resetKeyInputs() {
	controller.setKeyInputs(0, 0)
}

func (controller *Controller) Dispatch(newState *models.State) {
	controller.SetState(newState)
	screenProps := mapStateModelToScreenProps(controller.state)
	controller.screen.Render(screenProps)
}

func (controller *Controller) HandleMainLoop(interval time.Duration) (*models.State, error) {
	ch := controller.inputtedCharacter
	key := controller.inputtedKey
	controller.resetKeyInputs()

	var newState *models.State
	var err error

	switch {
	// Start or restart a game.
	case ch == 's':
		newState, err = reducers.StartOrRestartGame(*controller.state)
	// Move the hero.
	case key == termbox.KeyArrowUp || ch == 'k':
		newState, err = reducers.WalkHero(*controller.state, reducers.FourDirectionUp)
	case key == termbox.KeyArrowRight || ch == 'l':
		newState, err = reducers.WalkHero(*controller.state, reducers.FourDirectionRight)
	case key == termbox.KeyArrowDown || ch == 'j':
		newState, err = reducers.WalkHero(*controller.state, reducers.FourDirectionDown)
	case key == termbox.KeyArrowLeft || ch == 'h':
		newState, err = reducers.WalkHero(*controller.state, reducers.FourDirectionLeft)
	}

	if err != nil {
		return newState, err
	}

	if newState == nil {
		return reducers.AdvanceTime(*controller.state, interval)
	} else {
		return reducers.AdvanceTime(*newState, interval)
	}
}

func (controller *Controller) HandleKeyPress(ch rune, key termbox.Key) {
	controller.setKeyInputs(ch, key)
}

func CreateController() (*Controller, error) {
	controller := &Controller{}

	state := models.CreateState()
	err := state.SetWelcomeData()
	if err != nil {
		return controller, err
	}

	screen := views.CreateScreen(24, 80)

	controller.resetKeyInputs()
	controller.state = state
	controller.screen = screen
	controller.Dispatch(state)

	return controller, nil
}
