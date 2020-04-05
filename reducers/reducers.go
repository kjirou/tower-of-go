package reducers

import(
	"github.com/kjirou/tower_of_go/models"
	"github.com/kjirou/tower_of_go/utils"
)

// TODO: Generalize the interface between functions
func WalkHero(state models.State, direction utils.FourDirection) (*models.State, error) {
	field := state.GetField()
	element := field.GetElementOfHero()
	nextY := element.GetPosition().GetY()
	nextX := element.GetPosition().GetX()
	switch direction {
	case utils.FourDirectionUp:
		nextY -= 1
	case utils.FourDirectionRight:
		nextX += 1
	case utils.FourDirectionDown:
		nextY += 1
	case utils.FourDirectionLeft:
		nextX -= 1
	}
	var position utils.IMatrixPosition = element.GetPosition()
	var nextPosition utils.IMatrixPosition = &utils.MatrixPosition{
		Y: nextY,
		X: nextX,
	}
	if nextPosition.Validate(field.MeasureRowLength(), field.MeasureColumnLength()) {
		element, err := field.At(nextPosition)
		if err != nil {
			return &state, err
		} else if element.IsObjectEmpty() {
			err := field.MoveObject(position, nextPosition)
			return &state, err
		}
	}
	return &state, nil
}
