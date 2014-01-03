// Package mapview implements functions for drawing image representations of tmx
// maps.
package mapview

import (
	"image"
	"image/draw"

	"github.com/mewmew/tmx"
	"github.com/mewmew/tmx/examples/mapview/tile"
)

// A View corresponds to an image representation of a map.
type View struct {
	// Image corresponds to the image on which the map's tiles are drawn.
	draw.Image
	// cols corresponds to the number of columns in the map.
	cols int
	// rows corresponds to the number of rows in the map.
	rows int
	// tileWidth corresponds to the standard tile width in pixels.
	tileWidth int
	// tileHeight corresponds to the standard tile height in pixels.
	tileHeight int
	// delta is the differance between the map's standard tile height and the
	// maximum tile height of all tilesets.
	delta int
	// layers associated with the map.
	layers []tmx.Layer
	// tileset is a map from a tile ID to a tile image.
	tileset tile.Tileset
}

// NewView returns a new view of the map. The tileset sprite sheet is loaded
// relative to the tmx dir.
func NewView(m *tmx.Map, dir string) (view *View, err error) {
	view = &View{
		cols:       m.Width,
		rows:       m.Height,
		tileWidth:  m.TileWidth,
		tileHeight: m.TileHeight,
		delta:      getDelta(m),
		layers:     m.Layers,
	}
	// Each map is (cols+rows)/2 number of tiles in width and height.
	i := (view.cols + view.rows) / 2
	width := i * view.tileWidth
	height := i*view.tileHeight + view.delta
	view.Image = image.NewRGBA(image.Rect(0, 0, width, height))
	view.tileset, err = GetTileset(m, dir)
	if err != nil {
		return nil, err
	}
	return view, nil
}

// getDelta returns the differance between the map's standard tile height and
// the maximum tile height of all tilesets.
func getDelta(m *tmx.Map) int {
	var max int
	for _, ts := range m.Tilesets {
		if max < ts.TileHeight {
			max = ts.TileHeight
		}
	}
	return max - m.TileHeight
}

// GetCellRect returns the image.Rectangle of the cell at the provided
// coordinates.
//
// The calculations are based on the map coordinate system illustrated below:
//
//                    (0, 0)
//
//             +---------------+
//             |      r /\ c   |
//             |     o /\/\ o  |
//             |    w /\/\/\ l |
//             |     /\/\/\/\  |
//             |    /\/\/\/\/\ |
//             |   /\/\/\/\/\/\|
//             |  /\/\/\/\/\/\/|   (5, 0)
//             | /\/\/\/\/\/\/ |
//             |/\/\/\/\/\/\/  |
//    (0, 8)   |\/\/\/\/\/\/   |
//             | \/\/\/\/\/    |
//             |  \/\/\/\/     |
//             |   \/\/\/      |
//             |    \/\/       |
//             |     \/        |
//             +---------------+
//
//                 (5, 8)
func (view *View) GetCellRect(col, row int) image.Rectangle {
	halfTileWidth := view.tileWidth / 2
	halfTileHeight := view.tileHeight / 2

	// X offset to cell (0, 0):
	x := (view.rows - 1) * halfTileWidth
	// Adjust x offset based on col:
	x += col * halfTileWidth
	// Adjust x offset based on row:
	x -= row * halfTileWidth

	// Y offset to cell (0, 0):
	y := 0
	// Adjust y offset based on col:
	y += col * halfTileHeight
	// Adjust y offset based on row:
	y += row * halfTileHeight

	return image.Rect(x, y, x+view.tileWidth, y+view.tileHeight)
}

// GetTileRect returns the image.Rectangle of the tile at the provided
// coordinates.
func (view *View) GetTileRect(col, row int, tileBounds image.Rectangle) image.Rectangle {
	rect := view.GetCellRect(col, row)
	// Adjust max x based on tile width.
	rect.Max.X += tileBounds.Dx() - view.tileWidth
	// Adjust min y based on tile height.
	rect.Min.Y -= tileBounds.Dy() - view.tileHeight
	return rect
}

// Draw draws the image representation of the map to the view image.
func (view *View) Draw() {
	for _, layer := range view.layers {
		if layer.Name == "collision" {
			continue
		}
		for row := 0; row < view.rows; row++ {
			for col := 0; col < view.cols; col++ {
				gid := layer.GetGID(col, row)
				tile, ok := view.tileset[gid]
				if !ok {
					continue
				}
				sr := tile.Bounds()
				tileRect := view.GetTileRect(col, row, sr)
				tileRect = tileRect.Add(tile.Offset)
				tileRect = tileRect.Add(image.Pt(0, view.delta))
				draw.Draw(view, tileRect, tile, sr.Min, draw.Over)
			}
		}
	}
}
