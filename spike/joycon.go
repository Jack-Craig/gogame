package main

import (
	"log"

	"github.com/nobonobo/joycon"
)

type JoyConInput struct {
	id    int
	xAxis float32
	yAxis float32
}

func main() {
	// All devices must be paired (bluetooth) before this
	leftJoyCons, err := joycon.Search(joycon.JoyConL)
	if err != nil {
		log.Fatalln(err)
	}

	rightJoyCons, err := joycon.Search(joycon.JoyConR)
	if err != nil {
		log.Fatalln(err)
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

	for {
		inp := <-inputChannel
		log.Printf("ID: %d, [%f, %f]\n", inp.id, inp.xAxis, inp.yAxis)
	}

}

func RegisterJoyCon(id int, isLeft bool, jc *joycon.Joycon, ic chan<- JoyConInput) {
	go func(id int, joycon *joycon.Joycon) {
		for {
			state := <-joycon.State()
			if isLeft {
				ic <- JoyConInput{
					id:    id,
					xAxis: state.LeftAdj.X,
					yAxis: state.LeftAdj.Y,
				}
			} else {
				ic <- JoyConInput{
					id:    id,
					xAxis: state.RightAdj.X,
					yAxis: state.RightAdj.Y,
				}
			}
		}

	}(id, jc)
}
