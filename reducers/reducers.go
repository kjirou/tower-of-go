package reducers

import(
	"github.com/kjirou/tower_of_go/models"
	"github.com/kjirou/tower_of_go/utils"
	"time"
)

// TODO: Generalize the interface between reducer functions.

func StartGame(state models.State) (*models.State, bool, error) {
	state.GetGame().Start()
	return &state, true, nil
}

func AdvanceTime(state models.State, delta time.Duration) (*models.State, bool, error) {
	game := state.GetGame()
	game.AlterPlaytime(delta)
	return &state, true, nil
}

func WalkHero(state models.State, direction utils.FourDirection) (*models.State, bool, error) {
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
			return &state, false, err
		} else if element.IsObjectEmpty() {
			err := field.MoveObject(position, nextPosition)
			return &state, err == nil, err
		}
	}
	return &state, true, nil
}
