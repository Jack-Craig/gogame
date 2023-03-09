package main

import (
	"log"

	"github.com/Jack-Craig/gogame/src/input"
)

func main() {
	im := input.NewInputManager()
	im.InitiateJoyConConnections()
	inputMap := im.GetPlayerInputs()
	for {
		for id, pi := range *inputMap {
			xAxis, yAxis := pi.GetAxes()
			log.Printf("Controller [%d]: [%f, %f], Button: %t\n", id, xAxis, yAxis, pi.IsButtonPressed(input.JoyConX))
		}
	}
}
