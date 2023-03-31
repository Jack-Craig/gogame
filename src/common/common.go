package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"os"
)

type Pair struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type BiomeJson struct {
	NextTo          []string `json:"nextTo"`
	SurfaceTiles    int      `json:"surfaceTiles"`
	SubSurfaceTiles int      `json:"subsurfaceTiles"`
	GenAmplitude    uint32   `json:"genAmplitude"`
	GenFrequency    float64  `json:"genFrequency"`
}

type BiomeDataJson struct {
	Biomes map[string]BiomeJson `json:"biomes"`
}

type PlayerDataJson struct {
	Players map[string]struct {
		ImageId int `json:"imageId"`
	} `json:"players"`
}

func LoadJSON[T any](filePath string, container T) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	jsonFileBytes, _ := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(jsonFileBytes, container)
}

func Remove[T any](slice *[]T, index int) []T {
	l := len(*slice)
	if l <= index {
		return *slice
	}
	(*slice)[index] = (*slice)[l-1]
	(*slice) = (*slice)[:l-1]
	return *slice
}

type Vec2 struct {
	X, Y float64
}

func NewVec2(x, y float64) Vec2 {
	return Vec2{x, y}
}

func Add(v1, v2 Vec2) Vec2 {
	return NewVec2(v1.X+v2.X, v1.Y+v2.Y)
}

func Sub(v1, v2 Vec2) Vec2 {
	return NewVec2(v1.X-v2.X, v1.Y-v2.Y)
}

func Neg(v1 Vec2) Vec2 {
	return NewVec2(-v1.X, -v1.Y)
}

func Dot(v1, v2 Vec2) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}

func Normal(v1 Vec2) Vec2 {
	return NewVec2(v1.Y, -v1.X)
}

func Normalize(v1 Vec2) Vec2 {
	magn := math.Sqrt(v1.X*v1.X + v1.Y*v1.Y)
	return NewVec2(v1.X/magn, v1.Y/magn)
}

func MinMaxProjection(x, y, width, height float64, axis Vec2) (float64, float64) {
	// TL
	projection := Dot(NewVec2(x, y), axis)
	min := projection
	max := projection

	// TR
	projection = Dot(NewVec2(x+width, y), axis)
	min = math.Min(min, projection)
	max = math.Max(max, projection)

	// BR
	projection = Dot(NewVec2(x+width, y+height), axis)
	min = math.Min(min, projection)
	max = math.Max(max, projection)

	// BL
	projection = Dot(NewVec2(x, y+height), axis)
	min = math.Min(min, projection)
	max = math.Max(max, projection)

	return min, max
}

func TwoMaxPoints(p1, p2, p3, p4 float64) (float64, float64) {
	// m1 is max1, m2 is max2
	if p1 > p2 && p1 > p3 && p1 > p4 {
		return p1, MaxPoint(p2, p3, p4)
	}
	if p2 > p1 && p2 > p3 && p2 > p4 {
		return p2, MaxPoint(p1, p3, p4)
	}
	if p3 > p2 && p3 > p1 && p3 > p4 {
		return p1, MaxPoint(p1, p2, p4)
	}
	return p4, MaxPoint(p1, p2, p3)
}

func MaxPoint(p1, p2, p3 float64) float64 {
	if p1 > p2 && p1 > p3 {
		return p1
	}
	if p2 > p1 && p2 > p3 {
		return p2
	}
	return p3
}

func TwoMinPoints(p1, p2, p3, p4 float64) (float64, float64) {
	// m1 is max1, m2 is max2
	if p1 < p2 && p1 < p3 && p1 < p4 {
		return p1, MaxPoint(p2, p3, p4)
	}
	if p2 < p1 && p2 < p3 && p2 < p4 {
		return p2, MaxPoint(p1, p3, p4)
	}
	if p3 < p2 && p3 < p1 && p3 < p4 {
		return p1, MaxPoint(p1, p2, p4)
	}
	return p4, MaxPoint(p1, p2, p3)
}

func MinPoint(p1, p2, p3 float64) float64 {
	if p1 < p2 && p1 < p3 {
		return p1
	}
	if p2 < p1 && p2 < p3 {
		return p2
	}
	return p3
}
