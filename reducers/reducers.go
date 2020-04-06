package reducers

import(
	"fmt"
	"github.com/kjirou/tower_of_go/models"
	"github.com/kjirou/tower_of_go/utils"
	"time"
)

// TODO: Generalize the interface between reducer functions.

func StartGame(state models.State) (*models.State, error) {
	state.GetGame().Start()
	return &state, nil
}

func AlterPlaytime(state models.State, delta time.Duration) (*models.State, error) {
	game := state.GetGame()
	if !game.IsStarted() {
		return &state, fmt.Errorf("The game has not started.")
	} else if game.IsFinished() {
		return &state, fmt.Errorf("The game is over.")
	}
	game.AlterPlaytime(delta)
	return &state, nil
}

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
