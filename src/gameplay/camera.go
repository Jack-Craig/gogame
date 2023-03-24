package gameplay

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	w                                     *World
	offX, offY, screenWidth, screenHeight float32
}

func NewCamera(w *World) *Camera {
	return &Camera{
		w: w,
	}
}

func (c *Camera) Update() {
	if c.screenHeight == 0 || c.screenWidth == 0 {
		return
	}
	var totalX, totalY float32
	for _, player := range c.w.playerObjects {
		totalX += player.x + player.width/2
		totalY += player.y + player.height/2
	}
	newXOffset := -totalX/float32(len(c.w.playerObjects)) + c.screenWidth/2
	// Only able to move right, cant see outside of world on right
	if newXOffset < c.offX && int(c.w.level.worldXEnd) < int(c.w.level.worldWidth) {
		c.offX = newXOffset
	}
	newYOffset := -totalY/float32(len(c.w.playerObjects)) + c.screenHeight*2/3
	if float32(WORLDBUFFERHEIGHT*uint32(TILEWIDTH)) > -newYOffset+c.screenHeight {
		c.offY = newYOffset
	}

}

// Returns a copy of the render transformation matrix
func (c *Camera) GetRenderOffset() (float32, float32) {
	return c.offX, c.offY
}

// Returns true if the coordinates are within the camera bounds
func (c *Camera) IsInsideCamera(x, y float32) bool {
	topLeftX := -c.offX
	topLeftY := -c.offY
	bottomRightX := -c.offX + c.screenWidth
	bottomRightY := -c.offY + c.screenHeight
	if y == -1 {
		return !(x < topLeftX || x > bottomRightX)
	}
	return !(x < topLeftX || x > bottomRightX || y < topLeftY || y > bottomRightY)
}

func (c *Camera) Draw(screen *ebiten.Image) {
}
