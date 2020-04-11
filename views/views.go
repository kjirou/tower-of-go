package views

//
// The "views" package creates a layer that avoids to write logics tightly coupled with "termbox".
//

import (
	"fmt"
	"github.com/kjirou/tower_of_go/utils"
	"github.com/nsf/termbox-go"
)

type ScreenCellProps struct {
	Symbol          rune
	ForegroundColor termbox.Attribute
	BackgroundColor termbox.Attribute
}

// TODO: "ForegroundColor"->"Foreground"
// TODO: Prevent public access.
type screenCell struct {
	Symbol          rune
	ForegroundColor termbox.Attribute
	BackgroundColor termbox.Attribute
}

func (screenCell_ *screenCell) render(props *ScreenCellProps) {
	screenCell_.Symbol = props.Symbol
	screenCell_.ForegroundColor = props.ForegroundColor
	screenCell_.BackgroundColor = props.BackgroundColor
}

type screenText struct {
	Position *utils.MatrixPosition
	// ASCII only. Line breaks are not allowed.
	Text string
	Foreground termbox.Attribute
}

func createSequentialScreenTexts(position *utils.MatrixPosition, parts []*screenText) []*screenText {
	texts := make([]*screenText, 0)
	deltaX := 0
	for _, part := range parts {
		pos := &utils.MatrixPosition{
			Y: position.GetY(),
			X: position.GetX() + deltaX,
		}
		deltaX += len(part.Text)
		fg := termbox.ColorWhite
		if part.Foreground != 0 {
			fg = part.Foreground
		}
		text := screenText{
			Position: pos,
			Text: part.Text,
			Foreground: fg,
		}
		texts = append(texts, &text)
	}
	return texts
}

type ScreenProps struct {
	FieldCells [][]*ScreenCellProps
	FloorNumber int
	LankMessage string
	LankMessageForeground termbox.Attribute
	RemainingTime float64
}

type Screen struct {
	matrix [][]*screenCell
	staticTexts []*screenText
}

func (screen *Screen) measureRowLength() int {
	return len(screen.matrix)
}

func (screen *Screen) measureColumnLength() int {
	return len(screen.matrix[0])
}

func (screen *Screen) ForEachCells(
	callback func(
		y int,
		x int,
		symbol rune,
		foreground termbox.Attribute,
		background termbox.Attribute)) {
	for y, row := range screen.matrix {
		for x, cell := range row {
			callback(y, x, cell.Symbol, cell.ForegroundColor, cell.BackgroundColor)
		}
	}
}

func (screen *Screen) Render(props *ScreenProps) {
	rowLength := screen.measureRowLength()
	columnLength := screen.measureColumnLength()

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
			cell := screen.matrix[y][x]
			cell.render(&ScreenCellProps{
				Symbol: symbol,
				ForegroundColor: termbox.ColorWhite,
				BackgroundColor: termbox.ColorBlack,
			})
		}
	}

	// Place the field.
	fieldPosition := &utils.MatrixPosition{Y: 2, X: 2}
	for y, rowProps := range props.FieldCells {
		for x, cellProps := range rowProps {
			cell := screen.matrix[y + fieldPosition.GetY()][x + fieldPosition.GetX()]
			cell.render(cellProps)
		}
	}

	// Prepare texts.
	texts := make([]*screenText, 0)
	texts = append(texts, screen.staticTexts...)
	remainingTimeText := fmt.Sprintf("%4.1f", props.RemainingTime)
	timeText := &screenText{
		Position: &utils.MatrixPosition{Y: 3, X: 25},
		Text: fmt.Sprintf("Time : %s", remainingTimeText),
		Foreground: termbox.ColorWhite,
	}
	texts = append(texts, timeText)
	floorNumberText := &screenText{
		Position: &utils.MatrixPosition{Y: 4, X: 25},
		Text: fmt.Sprintf("Floor: %2d", props.FloorNumber),
		Foreground: termbox.ColorWhite,
	}
	texts = append(texts, floorNumberText)
	if props.LankMessage != "" {
		lankText := &screenText{
			Position: &utils.MatrixPosition{Y: 5, X: 27},
			Text: props.LankMessage,
			Foreground: props.LankMessageForeground,
		}
		texts = append(texts, lankText)
	}

	// Place texts.
	for _, textInstance := range texts {
		for deltaX, character := range textInstance.Text {
			cell := screen.matrix[textInstance.Position.GetY()][textInstance.Position.GetX() + deltaX]
			cell.render(&ScreenCellProps{
				Symbol: character,
				ForegroundColor: textInstance.Foreground,
				BackgroundColor: termbox.ColorBlack,
			})
		}
	}
}

func CreateScreen(rowLength int, columnLength int) *Screen {
	matrix := make([][]*screenCell, rowLength)
	for rowIndex := 0; rowIndex < rowLength; rowIndex++ {
		row := make([]*screenCell, columnLength)
		for columnIndex := 0; columnIndex < columnLength; columnIndex++ {
			row[columnIndex] = &screenCell{
				Symbol:          '_',
				ForegroundColor: termbox.ColorWhite,
				BackgroundColor: termbox.ColorBlack,
			}
		}
		matrix[rowIndex] = row
	}

	staticTexts := make([]*screenText, 0)

	titleText := &screenText{
		Position: &utils.MatrixPosition{Y: 0, X: 2},
		Text: "[ A Tower of Go ]",
		Foreground: termbox.ColorWhite,
	}
	staticTexts = append(staticTexts, titleText)

	urlText := &screenText{
		Position: &utils.MatrixPosition{Y: 22, X: 41},
		Text: "https://github.com/kjirou/tower_of_go",
		Foreground: termbox.ColorWhite | termbox.AttrUnderline,
	}
	staticTexts = append(staticTexts, urlText)

	operationTitleText := &screenText{
		Position: &utils.MatrixPosition{Y: 11, X: 25},
		Text: "[ Operations ]",
		Foreground: termbox.ColorWhite,
	}
	staticTexts = append(staticTexts, operationTitleText)

	var sKeyHelpTextParts = make([]*screenText, 0)
	sKeyHelpTextParts = append(sKeyHelpTextParts, &screenText{Text: "\""})
	sKeyHelpTextParts = append(sKeyHelpTextParts, &screenText{Text: "s", Foreground: termbox.ColorYellow})
	sKeyHelpTextParts = append(sKeyHelpTextParts, &screenText{Text: "\" ... Start or restart a new game."})
	sKeyHelpTexts := createSequentialScreenTexts(&utils.MatrixPosition{Y: 12, X: 25}, sKeyHelpTextParts)
	staticTexts = append(staticTexts, sKeyHelpTexts...)

	var moveKeysHelpTextParts = make([]*screenText, 0)
	moveKeysHelpTextParts =
		append(moveKeysHelpTextParts, &screenText{Text: "Arrow keys", Foreground: termbox.ColorYellow})
	moveKeysHelpTextParts = append(moveKeysHelpTextParts, &screenText{Text: " or \""})
	moveKeysHelpTextParts =
		append(moveKeysHelpTextParts, &screenText{Text: "k,l,j,h", Foreground: termbox.ColorYellow})
	moveKeysHelpTextParts = append(moveKeysHelpTextParts, &screenText{Text: "\" ... Move the player."})
	staticTexts = append(
		staticTexts,
		createSequentialScreenTexts(&utils.MatrixPosition{Y: 13, X: 25}, moveKeysHelpTextParts)...
	)

	description1Text := &screenText{
		Position: &utils.MatrixPosition{Y: 17, X: 3},
		Text: "Move the player in the upper left to reach the stairs in the lower right.",
		Foreground: termbox.ColorWhite,
	}
	staticTexts = append(staticTexts, description1Text)

	description2Text := &screenText{
		Position: &utils.MatrixPosition{Y: 18, X: 3},
		Text: "The score is the number of floors that can be reached within 30 seconds.",
		Foreground: termbox.ColorWhite,
	}
	staticTexts = append(staticTexts, description2Text)

	return &Screen{
		matrix: matrix,
		staticTexts: staticTexts,
	}
}
