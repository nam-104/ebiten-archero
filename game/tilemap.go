package game

import (
	"encoding/json"
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// TilemapLayerJSON đại diện cho 1 layer trong map
type TilemapLayerJSON struct {
	Data   []int `json:"data"`
	Width  int   `json:"width"`
	Height int   `json:"height"`
}

// TilemapJSON chứa toàn bộ layer của map
type TilemapJSON struct {
	Layers []TilemapLayerJSON `json:"layers"`
	Width  int                `json:"width"`
	Height int                `json:"height"`
	TileW  int                `json:"tilewidth"`
	TileH  int                `json:"tileheight"`
}

// NewTilemapJSON đọc file map JSON (Tiled) và parse
func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tilemap TilemapJSON
	if err := json.Unmarshal(contents, &tilemap); err != nil {
		return nil, err
	}

	return &tilemap, nil
}

// DrawTilemap vẽ map với camera offset
func DrawTilemap(screen *ebiten.Image, tilemap *TilemapJSON, tileset *ebiten.Image, cameraX, cameraY float64) {
	if tilemap == nil || tileset == nil {
		return
	}

	opts := ebiten.DrawImageOptions{}
	tileW := tilemap.TileW
	tileH := tilemap.TileH
	tilesPerRow := tileset.Bounds().Dx() / tileW

	for _, layer := range tilemap.Layers {
		for idx, id := range layer.Data {
			if id == 0 {
				continue
			}

			x := idx % layer.Width
			y := idx / layer.Width

			dstX := float64(x*tileW) - cameraX
			dstY := float64(y*tileH) - cameraY

			srcX := (id - 1) % tilesPerRow
			srcY := (id - 1) / tilesPerRow

			srcRect := image.Rect(
				srcX*tileW,
				srcY*tileH,
				srcX*tileW+tileW,
				srcY*tileH+tileH,
			)

			opts.GeoM.Translate(dstX, dstY)
			screen.DrawImage(tileset.SubImage(srcRect).(*ebiten.Image), &opts)
			opts.GeoM.Reset()
		}
	}
}
