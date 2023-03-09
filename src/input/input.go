package input

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InputManager struct {
	gamepadIDs     map[ebiten.GamepadID]struct{}
	axes           map[ebiten.GamepadID][]float64
	pressedButtons map[ebiten.GamepadID][]int
}

func (im *InputManager) Update() {
	if im.gamepadIDs == nil {
		log.Println("Input Manager Init")
		im.gamepadIDs = make(map[ebiten.GamepadID]struct{})
	}
	// Register newly connected gamepads
	newGamepads := []ebiten.GamepadID{}
	newGamepads = inpututil.AppendJustConnectedGamepadIDs(newGamepads)
	for _, id := range newGamepads {
		log.Printf("gamepad connected: id: %d\n", id)
		im.gamepadIDs[id] = struct{}{}
		ebiten.VibrateGamepad(id, &ebiten.VibrateGamepadOptions{
			Duration:        1000 * time.Millisecond,
			WeakMagnitude:   1,
			StrongMagnitude: 1,
		})
	}
	// Unregister newly disconnected gamepads
	for id := range im.gamepadIDs {
		if inpututil.IsGamepadJustDisconnected(id) {
			log.Printf("gamepad disconnected: id: %d\n", id)
			delete(im.gamepadIDs, id)
		}
	}

	// Update input data
	im.axes = make(map[ebiten.GamepadID][]float64)
	im.pressedButtons = make(map[ebiten.GamepadID][]int)
	for id := range im.gamepadIDs {
		// Axis data
		//horiz := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal)
		//vert := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical)

		//log.Printf("Horizontal: %f, Vertical: %f\n", horiz, vert)
		//	im.axes[id] = append(im.axes[id], ebiten.GamepadAxisValue(id, a))

		// Button data
		maxButton := ebiten.GamepadButton(ebiten.GamepadButtonCount(id))
		for b := ebiten.GamepadButton(id); b < maxButton; b++ {
			/**
			if ebiten.IsGamepadButtonPressed(id, b) {
				im.pressedButtons[id] = append(im.pressedButtons[id], int(b))
			}
			*/

			if inpututil.IsGamepadButtonJustPressed(id, b) {
				log.Printf("button pressed: id: %d, button: %d\n", id, b)
			}
		}
		// IDK
		for b := ebiten.StandardGamepadButton(0); b <= ebiten.StandardGamepadButtonMax; b++ {
			if inpututil.IsStandardGamepadButtonJustPressed(id, b) {
				log.Printf("standard button pressed: id: %d, button: %d", id, b)
			}
		}
	}
}
