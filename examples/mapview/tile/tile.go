// Package tile implements handling of tiles and tilesets.
package tile

import (
	"image"
	"image/draw"
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
// startId as the first tile id.
//
// Note: If possible the added tiles will share pixels with the provided sprite
// sheet.
func (tileset Tileset) AddTiles(spriteSheet image.Image, startId, tileWidth, tileHeight int, tileOffset image.Point) {
	sub, ok := spriteSheet.(subImager)
	if !ok {
		sub = &subFallback{
			Image: spriteSheet,
		}
	}
	r := sub.Bounds()
	id := startId
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

// subImager is an interface that adds the SubImage method to the image.Image
// interface.
type subImager interface {
	image.Image
	SubImage(r image.Rectangle) image.Image
}

// subFallback provides a SubImage method for images that lacks it.
type subFallback struct {
	image.Image
}

// SubImage returns an image representing the portion of the image src visible
// through r. The returned value doesn't shares pixels with the original image.
func (src *subFallback) SubImage(r image.Rectangle) image.Image {
	dstRect := image.Rect(0, 0, r.Dx(), r.Dy())
	dst := image.NewRGBA(dstRect)
	draw.Draw(dst, dstRect, src, r.Min, draw.Over)
	return dst
}
