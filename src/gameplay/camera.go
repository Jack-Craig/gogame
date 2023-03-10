package gameplay

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Camera struct {
	w                                     *World
	offX, offY, screenWidth, screenHeight float32
	geo                                   *ebiten.GeoM
}

func NewCamera(w *World) *Camera {
	return &Camera{
		w:   w,
		geo: &ebiten.GeoM{},
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
	c.offX = -totalX/float32(len(c.w.playerObjects)) + c.screenWidth/2
	c.offY = -totalY/float32(len(c.w.playerObjects)) + c.screenHeight/2
}

// Returns a copy of the render transformation matrix
func (c *Camera) GetRenderOffset() ebiten.GeoM {
	c.geo.Reset()
	c.geo.Translate(float64(c.offX), float64(c.offY))
	return *c.geo
}

// Returns true if the coordinates are within the camera bounds
func (c *Camera) IsInsideCamera(x, y float32) bool {
	topLeftX := -c.offX
	topLeftY := -c.offY
	bottomRightX := -c.offX + c.screenWidth
	bottomRightY := -c.offY + c.screenHeight
	return !(x < topLeftX || x > bottomRightX || y < topLeftY || y > bottomRightY)
}

func (c *Camera) Draw(screen *ebiten.Image) {
	ebitenutil.DrawCircle(screen, float64(c.screenWidth/2), float64(c.screenHeight/2), 5, color.White)
}
