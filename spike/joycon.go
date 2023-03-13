//package a

import (
	"log"
	"math"

	"github.com/nobonobo/joycon"
)

type JoyConButton uint32

// Enum, where int value is the bit set when the button is pressed (See ButtonPressed())
const (
	JoyConHome         JoyConButton = 13
	JoyConSign                      = 8
	JoyConStick                     = 11
	JoyConX                         = 16
	JoyConY                         = 18
	JoyConA                         = 19
	JoyConB                         = 17
	JoyConSideTrigger               = 23
	JoyConSideBumper                = 22
	JoyConTriggerLeft               = 21
	JoyConTriggerRight              = 20
)

type JoyConInput struct {
	id      int
	xAxis   float32
	yAxis   float32
	buttons uint32
}

func (jci *JoyConInput) ButtonPressed(button JoyConButton) bool {
	expectedHighlightedBit := uint32(button)
	mask := uint32(math.Pow(2, float64(expectedHighlightedBit)))
	return mask&jci.buttons != 0
}

func main() {
	// All devices must be paired (bluetooth) before this
	leftJoyCons, _ := joycon.Search(joycon.JoyConL)

	rightJoyCons, _ := joycon.Search(joycon.JoyConR)

	if rightJoyCons == nil && leftJoyCons == nil {
		log.Fatalf("No left or right joycons \n")
	}

	id := 0
	inputChannel := make(chan JoyConInput)
	for _, d := range leftJoyCons {
		id++
		jc, err := joycon.NewJoycon(d.Path, false)
		if err != nil {
			log.Fatalln(err)
		}
		RegisterJoyCon(id, true, jc, inputChannel)
	}
	for _, d := range rightJoyCons {
		id++
		jc, err := joycon.NewJoycon(d.Path, false)
		if err != nil {
			log.Fatalln(err)
		}
		RegisterJoyCon(id, false, jc, inputChannel)
	}
	var last uint32
	for {
		inp := <-inputChannel
		if inp.buttons != last {
			last = inp.buttons
			log.Printf("ID: %d, [%f, %f], %d\n", inp.id, inp.xAxis, inp.yAxis, uint32(math.Log2(float64(inp.buttons))))
			if inp.ButtonPressed(JoyConX) {
				log.Println("X Pressed")
			}
		}
	}

}

func RegisterJoyCon(id int, isLeft bool, jc *joycon.Joycon, ic chan<- JoyConInput) {
	go func(id int, joycon *joycon.Joycon) {
		for {
			state := <-joycon.State()
			if isLeft {
				ic <- JoyConInput{
					id:      id,
					xAxis:   state.LeftAdj.X,
					yAxis:   state.LeftAdj.Y,
					buttons: state.Buttons,
				}
			} else {
				ic <- JoyConInput{
					id:      id,
					xAxis:   state.RightAdj.X,
					yAxis:   state.RightAdj.Y,
					buttons: state.Buttons,
				}
			}
		}

	}(id, jc)
}
