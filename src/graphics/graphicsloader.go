package graphics

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	TILESIZE = 16
)

type SpriteID uint32

const (
	DirtTile SpriteID = iota // 0
	GrassTile
	RockTile
	UserTile
	Background1
	Background2 // 5
	Background3
	Bullet
	Final
)

type spriteData struct {
	Frame            struct{ X, Y, W, H int }
	Rotated, Trimmed bool
	SpriteSourceSize struct{ X, Y, W, H int }
	SourceSize       struct{ W, H int }
}
type mapData struct {
	Frames map[string]spriteData
}

// Loads buffered sprite sheet into memory
// Loads sprit sheet map into memory
// Serves requests from sprite_id to spritesheet coordinates
type GraphicsDataLoader struct {
	spriteSheet *ebiten.Image
	spriteMap   map[SpriteID]*ebiten.Image
	font        font.Face
}

func NewGraphicsDataLoader(path string) *GraphicsDataLoader {
	gdl := &GraphicsDataLoader{}

	// Load spriteImage
	spriteImageFile, err := os.Open(path + "/spritesheet.png")
	if err != nil {
		log.Fatal(err)
	}
	defer spriteImageFile.Close()
	spriteImage, _, err := image.Decode(spriteImageFile)
	if err != nil {
		log.Fatal(err)
	}
	gdl.spriteSheet = ebiten.NewImageFromImage(spriteImage)

	// Load spriteMap
	spriteMapFile, err := os.Open(path + "/spritesheet.json")
	if err != nil {
		log.Fatal(err)
	}
	defer spriteMapFile.Close()
	spriteMapBytes, _ := ioutil.ReadAll(spriteMapFile)
	var md mapData
	json.Unmarshal(spriteMapBytes, &md)
	gdl.spriteMap = make(map[SpriteID]*ebiten.Image)
	for cur := DirtTile; cur < Final; cur++ {
		mapKey := fmt.Sprintf("%d.png", cur)
		sd := md.Frames[mapKey]
		im := gdl.spriteSheet.SubImage(image.Rect(sd.Frame.X, sd.Frame.Y, sd.Frame.X+sd.Frame.W, sd.Frame.Y+sd.Frame.H)).(*ebiten.Image)
		//log.Printf("SubImage: %p, ImageID: %d\n", im, cur)
		gdl.spriteMap[cur] = im
	}
	// Load font
	/**
	fontFile, err := os.Open("res/font.ttf")
	if err != nil {
		log.Fatal(err)
	}
	defer fontFile.Close()
	fontBytes, _ := ioutil.ReadAll(spriteMapFile)
	parsedFont, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}
	*/
	parsedFont, _ := opentype.Parse(fonts.MPlus1pRegular_ttf)
	font, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	gdl.font = text.FaceWithLineHeight(font, 54)
	return gdl
}

func (gdl *GraphicsDataLoader) GetSpriteImage(spriteId SpriteID) *ebiten.Image {
	return gdl.spriteMap[spriteId]
}

func (gdl *GraphicsDataLoader) GetFont() *font.Face {
	return &gdl.font
}
