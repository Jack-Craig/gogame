package gameplay

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Background for parallax tings
type Background struct {
	world                                        *World
	third, second, first                         *ebiten.Image
	width, height                                float32
	thirdModifier, secondModifier, firstModifier float32
}

func NewBackground(world *World) *Background {
	b := &Background{world: world}
	b.first, b.second, b.third = world.gdl.GetBackgroundImages()
	b.width = float32(b.third.Bounds().Max.X - b.third.Bounds().Min.X)
	b.height = float32(b.third.Bounds().Max.Y - b.third.Bounds().Min.Y)
	b.firstModifier = .05
	b.secondModifier = .2
	b.thirdModifier = .4
	return b
}

func (bg *Background) Draw(screen *ebiten.Image) {
	sizeScale := bg.world.camera.screenHeight / bg.height
	newWidth := sizeScale * bg.width
	newHeight := sizeScale * bg.height
	requiredTiles := 3 * bg.world.camera.screenWidth / newWidth

	cOffX := bg.world.camera.offX
	cOffY := bg.world.camera.offY

	screenTLX := float64(int(cOffX*bg.firstModifier-bg.width/2) % int(newWidth))
	screenTLY := float64(int(cOffY*bg.firstModifier*.25-bg.height/2) % int(newHeight))

	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(sizeScale), float64(sizeScale))
	op.GeoM.Translate(screenTLX, screenTLY+float64(bg.world.camera.screenHeight)*.5)
	for x := 0; x < int(requiredTiles); x++ {
		screen.DrawImage(bg.first, &op)
		op.GeoM.Translate(float64(newWidth), 0)
	}

	screenTLX = float64(int(cOffX*bg.secondModifier-bg.width/2) % int(newWidth))
	screenTLY = float64(int(cOffY*bg.secondModifier*.25-bg.height/2) % int(newHeight))

	op.GeoM.Reset()
	op.GeoM.Scale(float64(sizeScale), float64(sizeScale))
	op.GeoM.Translate(screenTLX, screenTLY+float64(bg.world.camera.screenHeight)*.525)
	for x := 0; x < int(requiredTiles); x++ {
		screen.DrawImage(bg.second, &op)
		op.GeoM.Translate(float64(newWidth), 0)
	}

	screenTLX = float64(int(cOffX*bg.thirdModifier-bg.width/2) % int(newWidth))
	screenTLY = float64(int(cOffY*bg.thirdModifier*.25-bg.height/2) % int(newHeight))

	op.GeoM.Reset()
	op.GeoM.Scale(float64(sizeScale), float64(sizeScale))
	op.GeoM.Translate(screenTLX, screenTLY+float64(bg.world.camera.screenHeight)*.55)
	for x := 0; x < int(requiredTiles); x++ {
		screen.DrawImage(bg.third, &op)
		op.GeoM.Translate(float64(newWidth), 0)
	}
}
