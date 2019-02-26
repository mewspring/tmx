package mapview

import "image"

import (
	"github.com/mewkiz/pkg/imgutil"
	"github.com/mewspring/tmx"
	"github.com/mewspring/tmx/examples/mapview/tile"
)

// GetTileset returns the combined tileset of a given tmx map.
func GetTileset(m *tmx.Map, dir string) (tileset tile.Tileset, err error) {
	tileset = tile.NewTileset()
	for _, ts := range m.Tilesets {
		spritePath := dir + "/" + ts.Image.Source
		spriteSheet, err := imgutil.ReadFile(spritePath)
		if err != nil {
			return nil, err
		}
		tileOffset := image.Pt(ts.TileOffset.X, ts.TileOffset.Y)
		tileset.AddTiles(spriteSheet, ts.FirstGID, ts.TileWidth, ts.TileHeight, tileOffset)
	}
	return tileset, nil
}
