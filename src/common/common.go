package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

func Remove[T any](slice []T, index int) []T {
	l := len(slice)
	if l >= index {
		return slice
	}
	(slice)[index] = (slice)[l-1]
	slice = slice[:l-1]
	return slice
}
