// Package tile implements handling of tiles and tilesets.
package tile

import (
	"image"

	"github.com/mewkiz/pkg/imgutil"
)

// Tileset is a map from a tile ID to a tile image.
type Tileset map[int]Tile

// Tile corresponds to a tile image.
type Tile struct {
	image.Image
	// Offset to be applied when drawing the tile.
	Offset image.Point
}

// NewTileset returns a new tileset.
func NewTileset() (tileset Tileset) {
	tileset = make(Tileset)
	return tileset
}

/// ### [ todo ] ###
///   - handle Margin?
///   - handle Spacing?
/// ### [/ todo ] ###

// AddTiles adds tiles to the tileset based on a provided sprite sheet, using
// startID as the first tile id.
//
// Note: If possible the added tiles will share pixels with the provided sprite
// sheet.
func (tileset Tileset) AddTiles(spriteSheet image.Image, startID, tileWidth, tileHeight int, tileOffset image.Point) {
	sub := imgutil.SubFallback(spriteSheet)
	r := sub.Bounds()
	id := startID
	for y := r.Min.Y; y < r.Max.Y; y += tileHeight {
		for x := r.Min.X; x < r.Max.X; x += tileWidth {
			tileRect := image.Rect(x, y, x+tileWidth, y+tileHeight)
			tile := Tile{
				Image:  sub.SubImage(tileRect),
				Offset: tileOffset,
			}
			tileset[id] = tile
			id++
		}
	}
}
