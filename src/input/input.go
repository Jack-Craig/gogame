package input

import (
	"log"
	"math"
	"sync"

	"github.com/Jack-Craig/gogame/src/common"
	"github.com/flynn/hid"
	"github.com/nobonobo/joycon"
)

type JoyConButton uint32

const (
	JoyConHome JoyConButton = iota
	JoyConSign
	JoyConStick
	JoyConX
	JoyConY
	JoyConA
	JoyConB
	JoyConSideTrigger
	JoyConSideBumper
	JoyConTriggerLeft
	JoyConTriggerRight
)

var maskMap map[JoyConButton]common.Pair = map[JoyConButton]common.Pair{
	JoyConHome:         {X: 0x2000, Y: 0x1000},
	JoyConSign:         {X: 0x100, Y: 0x200},
	JoyConStick:        {X: 0x800, Y: 0x400},
	JoyConX:            {X: 0x4, Y: 0x20000},
	JoyConY:            {X: 0x1, Y: 0x80000},
	JoyConA:            {X: 0x8, Y: 0x40000},
	JoyConB:            {X: 0x2, Y: 0x10000},
	JoyConSideTrigger:  {X: 0x800000, Y: 0x80},
	JoyConSideBumper:   {X: 0x40, Y: 0x400000},
	JoyConTriggerLeft:  {X: 0x200000, Y: 0x20},
	JoyConTriggerRight: {X: 0x100000, Y: 0x10},
}

type InputManager struct {
	playerInputs map[uint32]*PlayerInput
}

func NewInputManager() *InputManager {
	return &InputManager{
		playerInputs: make(map[uint32]*PlayerInput),
	}
}

func (im *InputManager) InitiateJoyConConnections() {
	leftJoyCons, _ := joycon.Search(joycon.JoyConL)
	rightJoyCons, _ := joycon.Search(joycon.JoyConR)
	if rightJoyCons == nil && leftJoyCons == nil {
		log.Fatalf("No left or right joycons \n")
	}
	var id uint32
	for _, d := range leftJoyCons {
		im.pairJoyCon(d, &id, true)
	}
	for _, d := range rightJoyCons {
		im.pairJoyCon(d, &id, false)
	}
}

func (im *InputManager) pairJoyCon(device *hid.DeviceInfo, id *uint32, isLeft bool) {
	jc, err := joycon.NewJoycon(device.Path, false)
	if err != nil {
		log.Fatalln(err)
	}
	playerInput := &PlayerInput{
		id:     *id,
		isLeft: isLeft,
	}
	im.playerInputs[*id] = playerInput
	go func(joycon *joycon.Joycon, pi *PlayerInput) {
		for {
			pi.SetControlState(<-joycon.State())
		}
	}(jc, playerInput)
	*id++
}

func (im *InputManager) GetPlayerInputs() *map[uint32]*PlayerInput {
	return &im.playerInputs
}

// PlayerInput is given to player objects for them to take controls
// Input is a joycon state, output is axes and buttons (L/R agnostic)
type PlayerInput struct {
	mut          sync.Mutex
	id           uint32
	isLeft       bool
	xAxis, yAxis float32
	buttons      uint32
}

func (pi *PlayerInput) SetControlState(state joycon.State) {
	pi.mut.Lock()
	defer pi.mut.Unlock()
	if pi.isLeft {
		pi.xAxis = -state.LeftAdj.X
		pi.yAxis = -state.LeftAdj.Y
	} else {
		pi.xAxis = state.RightAdj.X
		pi.yAxis = state.RightAdj.Y
	}
	if math.Abs(float64(pi.xAxis)) > 10 || math.Abs(float64(pi.xAxis)) < .1 {
		pi.xAxis = 0
	}
	if math.Abs(float64(pi.yAxis)) > 10 || math.Abs(float64(pi.yAxis)) < .1 {
		pi.yAxis = 0
	}
	pi.buttons = state.Buttons
}

func (pi *PlayerInput) GetAxes() (float32, float32) {
	pi.mut.Lock()
	defer pi.mut.Unlock()
	return pi.xAxis, pi.yAxis
}

func (pi *PlayerInput) IsButtonPressed(button JoyConButton) bool {
	pi.mut.Lock()
	defer pi.mut.Unlock()
	mapPair := maskMap[button]
	return (uint32(mapPair.X)&pi.buttons) != 0 || (uint32(mapPair.Y)&pi.buttons) != 0
}
