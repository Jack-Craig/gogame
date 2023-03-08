package input

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type InputManager struct {
	gamepadIDs map[ebiten.GamepadID]struct{}
}

func (im *InputManager) Update() {

}
