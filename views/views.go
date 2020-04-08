package views

import (
	"fmt"
	"github.com/kjirou/tower_of_go/utils"
	"github.com/nsf/termbox-go"
	"strings"
)

type ScreenElement struct {
	Symbol          rune
	ForegroundColor termbox.Attribute
	BackgroundColor termbox.Attribute
}

func (screenElement *ScreenElement) renderWithFieldElement(fieldElement utils.IFieldElement) {
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
	screenElement.Symbol = symbol
	screenElement.ForegroundColor = fg
	screenElement.BackgroundColor = bg
}

type ScreenText struct {
	Position utils.IMatrixPosition
	// ASCII only. Line breaks are not allowed.
	Text string
	Foreground termbox.Attribute
}

func createSequentialScreenTexts(position utils.IMatrixPosition, parts []*ScreenText) []*ScreenText {
	texts := make([]*ScreenText, 0)
	deltaX := 0
	for _, part := range parts {
		var pos utils.IMatrixPosition = &utils.MatrixPosition{
			Y: position.GetY(),
			X: position.GetX() + deltaX,
		}
		deltaX += len(part.Text)
		fg := termbox.ColorWhite
		if part.Foreground != 0 {
			fg = part.Foreground
		}
		text := ScreenText {
			Position: pos,
			Text: part.Text,
			Foreground: fg,
		}
		texts = append(texts, &text)
	}
	return texts
}

//
// A layer that avoids to write logics tightly coupled with "termbox".
//
type Screen struct {
	matrix [][]ScreenElement
	staticTexts []*ScreenText
}

func (screen *Screen) GetMatrix() [][]ScreenElement {
	return screen.matrix
}

func (screen *Screen) MeasureRowLength() int {
	return len(screen.matrix)
}

func (screen *Screen) MeasureColumnLength() int {
	return len(screen.matrix[0])
}

func (screen *Screen) At(position utils.IMatrixPosition) *ScreenElement {
	y := position.GetY()
	x := position.GetX()
	// TODO: Error handling.
	if y < 0 || y > screen.MeasureRowLength() {
		panic(fmt.Sprintf("That position (Y=%d) does not exist on the screen.", y))
	} else if x < 0 || x > screen.MeasureColumnLength() {
		panic(fmt.Sprintf("That position (X=%d) does not exist on the screen.", x))
	}
	return &(screen.matrix[y][x])
}

func (screen *Screen) renderField(startPosition utils.IMatrixPosition, field utils.IField) {
	rowLength := field.MeasureRowLength()
	columnLength := field.MeasureColumnLength()
	for y := 0; y < rowLength; y++ {
		for x := 0; x < columnLength; x++ {
			var screenElementPosition utils.IMatrixPosition = &utils.MatrixPosition{
				Y: startPosition.GetY() + y,
				X: startPosition.GetX() + x,
			}
			element := screen.At(screenElementPosition)
			var fieldElementPosition utils.IMatrixPosition = &utils.MatrixPosition{Y: y, X: x}
			// TODO: Error handling.
			var fieldElement, _ = field.At(fieldElementPosition)
			element.renderWithFieldElement(fieldElement)
		}
	}
}

func (screen *Screen) Render(state utils.IState) {
	game := state.GetGame()

	rowLength := screen.MeasureRowLength()
	columnLength := screen.MeasureColumnLength()

	// Pad elements with blanks.
	// Set borders.
	for y := 0; y < rowLength; y++ {
		for x := 0; x < columnLength; x++ {
			isTopOrBottomEdge := y == 0 || y == rowLength-1
			isLeftOrRightEdge := x == 0 || x == columnLength-1
			symbol := ' '
			switch {
			case isTopOrBottomEdge && isLeftOrRightEdge:
				symbol = '+'
			case isTopOrBottomEdge && !isLeftOrRightEdge:
				symbol = '-'
			case !isTopOrBottomEdge && isLeftOrRightEdge:
				symbol = '|'
			}
			screen.matrix[y][x].Symbol = symbol
		}
	}

	// Place the field.
	var fieldPosition utils.IMatrixPosition = &utils.MatrixPosition{Y: 2, X: 2}
	screen.renderField(fieldPosition, state.GetField())

	// Prepare texts.
	texts := make([]*ScreenText, 0)
	texts = append(texts, screen.staticTexts...)
	var timeTextPosition utils.IMatrixPosition = &utils.MatrixPosition{Y: 3, X: 25}
	remainingTime := game.CalculateRemainingTime(state.GetExecutionTime()).Seconds()
	remainingTimeText := fmt.Sprintf("%4.1f", remainingTime)
	timeText := ScreenText{
		Position: timeTextPosition,
		Text: fmt.Sprintf("Time : %s", remainingTimeText),
		Foreground: termbox.ColorWhite,
	}
	texts = append(texts, &timeText)
	var floorNumberTextPosition utils.IMatrixPosition = &utils.MatrixPosition{Y: 4, X: 25}
	floorNumberText := ScreenText{
		Position: floorNumberTextPosition,
		Text: fmt.Sprintf("Floor: %2d", game.GetFloorNumber()),
		Foreground: termbox.ColorWhite,
	}
	texts = append(texts, &floorNumberText)

	// Place texts.
	for _, textInstance := range texts {
		for deltaX, character := range textInstance.Text {
			element := &screen.matrix[textInstance.Position.GetY()][textInstance.Position.GetX() + deltaX]
			element.Symbol = character
			element.ForegroundColor = textInstance.Foreground
		}
	}
}

func (screen *Screen) AsText() string {
	rowLength := screen.MeasureRowLength()
	columnLength := screen.MeasureColumnLength()
	lines := make([]string, rowLength)
	for y := 0; y < rowLength; y++ {
		line := make([]rune, columnLength)
		for x := 0; x < columnLength; x++ {
			line[x] = screen.matrix[y][x].Symbol
		}
		lines[y] = string(line)
	}
	return strings.Join(lines, "\n")
}

func CreateScreen(rowLength int, columnLength int) Screen {
	matrix := make([][]ScreenElement, rowLength)
	for rowIndex := 0; rowIndex < rowLength; rowIndex++ {
		row := make([]ScreenElement, columnLength)
		for columnIndex := 0; columnIndex < columnLength; columnIndex++ {
			row[columnIndex] = ScreenElement{
				Symbol:          '_',
				ForegroundColor: termbox.ColorWhite,
				BackgroundColor: termbox.ColorBlack,
			}
		}
		matrix[rowIndex] = row
	}

	staticTexts := make([]*ScreenText, 0)

	var titleTextPosition utils.IMatrixPosition = &utils.MatrixPosition{Y: 0, X: 2}
	titleText := ScreenText{
		Position: titleTextPosition,
		Text: "[ A Tower of Go ]",
		Foreground: termbox.ColorWhite,
	}
	staticTexts = append(staticTexts, &titleText)

	var urlTextPosition utils.IMatrixPosition = &utils.MatrixPosition{Y: 22, X: 41}
	urlText := ScreenText{
		Position: urlTextPosition,
		Text: "https://github.com/kjirou/tower_of_go",
		Foreground: termbox.ColorWhite | termbox.AttrUnderline,
	}
	staticTexts = append(staticTexts, &urlText)

	var sKeyHelpTextPosition utils.IMatrixPosition = &utils.MatrixPosition{Y: 16, X: 2}
	var sKeyHelpTextParts = make([]*ScreenText, 0)
	sKeyHelpTextParts = append(sKeyHelpTextParts, &ScreenText{Text: "\""})
	sKeyHelpTextParts = append(sKeyHelpTextParts, &ScreenText{Text: "s", Foreground: termbox.ColorYellow})
	sKeyHelpTextParts = append(sKeyHelpTextParts, &ScreenText{Text: "\" ... Start a new game."})
	sKeyHelpTexts := createSequentialScreenTexts(sKeyHelpTextPosition, sKeyHelpTextParts)
	staticTexts = append(staticTexts, sKeyHelpTexts...)

	var rKeyHelpTextPosition utils.IMatrixPosition = &utils.MatrixPosition{Y: 17, X: 2}
	var rKeyHelpTextParts = make([]*ScreenText, 0)
	rKeyHelpTextParts = append(rKeyHelpTextParts, &ScreenText{Text: "\""})
	rKeyHelpTextParts = append(rKeyHelpTextParts, &ScreenText{Text: "r", Foreground: termbox.ColorYellow})
	rKeyHelpTextParts = append(rKeyHelpTextParts, &ScreenText{Text: "\" ... Reset the current game."})
	rKeyHelpTexts := createSequentialScreenTexts(rKeyHelpTextPosition, rKeyHelpTextParts)
	staticTexts = append(staticTexts, rKeyHelpTexts...)

	var moveKeysHelpTextPosition utils.IMatrixPosition = &utils.MatrixPosition{Y: 18, X: 2}
	var moveKeysHelpTextParts = make([]*ScreenText, 0)
	moveKeysHelpTextParts =
		append(moveKeysHelpTextParts, &ScreenText{Text: "Arrow keys", Foreground: termbox.ColorYellow})
	moveKeysHelpTextParts = append(moveKeysHelpTextParts, &ScreenText{Text: " or \""})
	moveKeysHelpTextParts =
		append(moveKeysHelpTextParts, &ScreenText{Text: "k,l,j,h", Foreground: termbox.ColorYellow})
	moveKeysHelpTextParts = append(moveKeysHelpTextParts, &ScreenText{Text: "\" ... Move the player."})
	staticTexts = append(
		staticTexts,
		createSequentialScreenTexts(moveKeysHelpTextPosition, moveKeysHelpTextParts)...
	)

	return Screen{
		matrix: matrix,
		staticTexts: staticTexts,
	}
}
