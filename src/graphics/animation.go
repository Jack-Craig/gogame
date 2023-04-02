package graphics

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	DEFAULTFRAMEDELAY = 100
)

type Animation struct {
	frames          []*ebiten.Image
	curFrame        int
	frameDelayMs    int64
	lastFrameDrawMs int64
}

func NewAnimation(gdl *GraphicsDataLoader, sprites []SpriteID) *Animation {
	animation := Animation{}
	for _, spriteId := range sprites {
		animation.frames = append(animation.frames, gdl.GetSpriteImage(spriteId))
	}
	animation.frameDelayMs = DEFAULTFRAMEDELAY
	return &animation
}

func (a *Animation) Draw(screen *ebiten.Image, ops *ebiten.DrawImageOptions) {
	screen.DrawImage(a.frames[a.curFrame], ops)
	timeNow := time.Now().UnixMilli()
	if timeNow > a.lastFrameDrawMs+a.frameDelayMs {
		a.curFrame = (a.curFrame + 1) % len(a.frames)
		a.lastFrameDrawMs = timeNow
	}
}
