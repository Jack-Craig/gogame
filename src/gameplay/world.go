package gameplay

import (
	"log"

	"github.com/Jack-Craig/gogame/src/graphics"
	"github.com/Jack-Craig/gogame/src/input"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	TILEWIDTH         float32 = 32
	MAXVEL            float32 = 12
	WORLDBUFFERLEN    uint32  = 80
	WORLDBUFFERHEIGHT uint32  = 30
)

type World struct {
	gdl           *graphics.GraphicsDataLoader
	im            *input.InputManager
	camera        *Camera
	gameObjects   []*GameObject
	entityObjects []*Entity
	playerObjects []*Player
	gravity       float32
	worldTiles    [][]*Tile

	// In array coordinates, the start and end of the visible world. Can wrap
	worldFrameStart uint32
	worldFrameEnd   uint32
	// In world coordinates, the left most coordinate
	worldXStart float32
}

func NewWorld(gdl *graphics.GraphicsDataLoader, im *input.InputManager) *World {
	w := &World{}
	w.gdl = gdl
	w.im = im
	w.camera = NewCamera(w)
	for id, playerInput := range *w.im.GetPlayerInputs() {
		log.Printf("Created player: ID: %d\n", id)
		p := NewPlayer(1, 50, 1, TILEWIDTH, TILEWIDTH, w, w.gdl.GetSpriteImage(3), playerInput)
		w.gameObjects = append(w.gameObjects, &p.GameObject)
		w.entityObjects = append(w.entityObjects, &p.Entity)
		w.playerObjects = append(w.playerObjects, p)
	}
	w.gravity = .25
	w.generateWorld()
	return w
}

func (w *World) Update() {
	w.camera.Update()
	// Update ring offets
	if w.camera.screenWidth > 0 && w.worldFrameEnd == 0 {
		w.worldFrameEnd = w.worldFrameStart + uint32(w.camera.screenWidth/TILEWIDTH)
	}
	if w.camera.screenWidth > 0 && !w.camera.IsInsideCamera(w.worldXStart+TILEWIDTH, -1) {
		// Need to advance world buffer
		w.worldXStart += TILEWIDTH
		lastTile := w.worldFrameStart
		w.worldFrameStart += 1
		w.worldFrameStart %= WORLDBUFFERLEN
		w.worldFrameEnd = w.worldFrameStart + uint32(w.camera.screenWidth/TILEWIDTH)
		w.worldFrameEnd %= WORLDBUFFERLEN
		// Need to change coordinates of blocks that just left the screen
		for y := uint32(0); y < WORLDBUFFERHEIGHT; y++ {
			var prevTile uint32
			if lastTile == 0 {
				prevTile = WORLDBUFFERLEN - 1
			} else {
				prevTile = lastTile - 1
			}
			w.worldTiles[y][lastTile].x = w.worldTiles[y][prevTile].x + TILEWIDTH
		}
	}
	for _, player := range w.playerObjects {
		player.Update()
	}
	for _, entity := range w.entityObjects {
		entity.AddVel(0, w.gravity)
		entity.Update()
	}
}

func (w *World) Draw(screen *ebiten.Image) {
	screenBounds := screen.Bounds().Max
	w.camera.screenWidth = float32(screenBounds.X)
	w.camera.screenHeight = float32(screenBounds.Y)

	for _, row := range w.worldTiles {
		for _, tile := range row {
			tile.GameObject.Draw(screen)
		}
	}
	for _, gobj := range w.gameObjects {
		gobj.Draw(screen)
	}
	w.camera.Draw(screen)
}

func (w *World) generateWorld() {
	grass := w.gdl.GetSpriteImage(2)
	sky := w.gdl.GetSpriteImage(4)
	dirt := w.gdl.GetSpriteImage(1)
	for y := float32(0); y < float32(WORLDBUFFERHEIGHT); y++ {
		var row []*Tile
		for x := float32(0); x < float32(WORLDBUFFERLEN); x++ {
			if y == float32(WORLDBUFFERHEIGHT)-6 {
				t := NewTile(2, x*TILEWIDTH, y*TILEWIDTH, w, grass)
				row = append(row, t)
			} else if y > float32(WORLDBUFFERHEIGHT)-6 {
				t := NewTile(2, x*TILEWIDTH, y*TILEWIDTH, w, dirt)
				row = append(row, t)
			} else {
				t := NewTile(3, x*TILEWIDTH, y*TILEWIDTH, w, sky)
				t.isPassable = true
				row = append(row, t)
			}
		}
		w.worldTiles = append(w.worldTiles, row)
	}
}

// Given an x and y in world coordinates, returns true if there is a tile there and false otherwise
func (w *World) IsWorldCollision(x, y float32) bool {
	gridX, gridY := w.worldToBuffer(x, y)
	if int(gridY) >= len(w.worldTiles) {
		return false
	}
	if int(gridX) >= len(w.worldTiles[0]) {
		return false
	}
	tile := w.worldTiles[gridY][gridX]
	return !tile.isPassable
}

func (w *World) worldToBuffer(x, y float32) (uint32, uint32) {
	bufferY := uint32(y / TILEWIDTH)
	gridX := uint32(x / TILEWIDTH)
	bufferX := gridX % WORLDBUFFERLEN
	return bufferX, bufferY
}
