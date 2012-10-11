package tmx

// Flip flags stored in the highest three bits of the GID.
const (
	FlagDiagonalFlip   = 0x20000000
	FlagVerticalFlip   = 0x40000000
	FlagHorizontalFlip = 0x80000000
	FlagFlip           = FlagDiagonalFlip | FlagVerticalFlip | FlagHorizontalFlip
)

// A Map contains all the map information stored in tmx files.
//
// The TileWidth and TileHeight properties determine the general grid size of
// the map. The individual tiles may have different sizes. Larger tiles will
// extend at the top and right (anchored to the bottom left).
type Map struct {
	// The TMX format version, generally 1.0.
	Version string `xml:"version,attr"`
	// Map orientation. Tiled supports "orthogonal" and "isometric" at the
	// moment.
	Orientation string `xml:"orientation,attr"`
	// The map width in tiles.
	Width int `xml:"width,attr"`
	// The map height in tiles.
	Height int `xml:"height,attr"`
	// The width of a tile.
	TileWidth int `xml:"tilewidth,attr"`
	// The height of a tile.
	TileHeight int `xml:"tileheight,attr"`
	// Properties associated with the map.
	Properties []Property `xml:"properties>property"`
	// Tilesets associated with the map.
	Tilesets []Tileset `xml:"tileset"`
	// Layers associated with the map.
	Layers []Layer `xml:"layer"`
	// Object layers associated with the map.
	ObjectLayers []ObjectLayer `xml:"objectgroup"`
}

// A Property is a name, value pair.
type Property struct {
	// The name of the property.
	Name string `xml:"name,attr"`
	// The value of the property.
	Value string `xml:"value,attr"`
}

/// ### [ todo ] ###
///    - Source: load info from TSX files.
/// ### [/ todo ] ###

// A Tileset is a sprite sheet of tiles.
type Tileset struct {
	// FirstGID is the first global tile ID of the tileset and it maps to the
	// first tile in the tilset.
	FirstGID int `xml:"firstgid,attr"`
	// Source refers to an external TSX (Tile Set XML) file. The TSX file has the
	// same structure as the Tileset described here, but without the firstgid and
	// source attributes, since they are map specific.
	Source string `xml:"source,attr"`
	// The name of the tileset.
	Name string `xml:"name,attr"`
	// The (maximum) width of the tiles in the tileset.
	TileWidth int `xml:"tilewidth,attr"`
	// The (maximum) height of the tiles in the tileset.
	TileHeight int `xml:"tileheight,attr"`
	// The spacing in pixels between the tiles in the tileset (applies to the
	// tileset image).
	Spacing int `xml:"spacing,attr"`
	// The margin around the tiles in the tileset (applies to the tileset image).
	Margin int `xml:"margin,attr"`
	// Tile offset associated with the tileset.
	TileOffset TileOffset `xml:"tileoffset"`
	// Properties associated with the tileset.
	Properties []Property `xml:"properties>property"`
	// The image associated with the tileset.
	Image Image `xml:"image"`
	// TilesInfo contains information about the tiles within a tileset.
	TilesInfo []TileInfo `xml:"tile"`
}

// A TileOffset specifies an offset in pixels, to be applied when drawing a tile
// from the related tileset.
type TileOffset struct {
	// Horizontal offset in pixels
	X int `xml:"x,attr"`
	// Vertical offset in pixels (positive is down)
	Y int `xml:"y,attr"`
}

// A single image is associated with each tileset. It is cut into smaller tiles
// based on the attributes defined in the tileset.
type Image struct {
	// Source refers to the tileset image file.
	Source string `xml:"source,attr"`
	// Trans defines a specific color that is treated as transparent (example
	// value: "FF00FF" for magenta).
	Trans string `xml:"trans,attr"`
	// The image width in pixels (optional, used for tile index correction when
	// the image changes).
	Width int `xml:"width,attr"`
	// The image height in pixels (optional).
	Height int `xml:"height,attr"`
}

// TileInfo contains information about a tile within a tileset.
type TileInfo struct {
	// The local tile ID within its tileset.
	Id int `xml:"id,attr"`
	// Properties associated with the tile.
	Properties []Property `xml:"properties>property"`
}

/// ### [ todo ] ###
///   - Visible: default value true
///   - Opacity: default value 1.0
/// ### [/ todo ] ###

// A Map can contain any number of layers. Each layer contains information about
// which global tile ID any given coordinate has.
type Layer struct {
	// The name of the layer.
	Name string `xml:"name,attr"`
	// Visible specifies whether the layer is shown (true) or hidden (false).
	Visible bool `xml:"visible,attr"`
	// The opacity of the layer as a value from 0.0 to 1.0.
	Opacity float64 `xml:"opacity,attr"`
	// Properties associated with the layer.
	Properties []Property `xml:"properties>property"`
	// Data contains the information about the tile GIDs associated with a layer.
	//
	// Note: Data should not be accessed directly. Use the GetGID method instead
	// to obtain the GID at a given coordinate.
	Data *Data `xml:"data"`
}

// GID corresponds to a global tile ID.
//
// Note: The highest three bits of the GID are used to store flip flags. These
// must be cleared before using the GID as a global tile ID. Either use the
// convenience methods or clear the flip flags manually before using the GID
// value.
type GID uint32

// Data contains the information about the tile GIDs associated with a layer.
//
// Note: Data should not be accessed directly. Use the GetGID method instead to
// obtain the GID at a given coordinate.
type Data struct {
	// Encoding specifies the encoding method used for the RawData. Options
	// include "base64", "csv" and "" for XML encoding.
	Encoding string `xml:"encoding,attr"`
	// Compression specifies the compression method used for the RawData. Options
	// include "gzip", "zlib" and "" for no compression.
	Compression string `xml:"compression,attr"`
	// RawData contains the raw data of tile GIDs, which can be represented in
	// several different ways as specified by Encoding and Compression.
	RawData string `xml:",innerxml"`
	// Tiles associated with the layer.
	Tiles []Tile `xml:"tile"`
	// gids contains the decoded tile GIDs arranged by col and row.
	gids [][]GID
}

// A Tile contains the GID of a single tile on a tile layer.
type Tile struct {
	// The global tile ID.
	GID GID `xml:"gid,attr"`
}

/// ### [ todo ] ###
///   - Visible: default value true
/// ### [/ todo ] ###

// A Map can contain any number of object layers. Each object layer contains
// information about different objects on the map.
//
// While tile layers are very suitable for anything repetitive aligned to the
// tile grid, sometimes you want to annotate your map with other information,
// not necessarily aligned to the grid. Hence the objects have their coordinates
// and size in pixels, but you can still easily align that to the grid when you
// want to.
type ObjectLayer struct {
	// The name of the object layer.
	Name string `xml:"name,attr"`
	// Visible specifies whether the layer is shown (true) or hidden (false).
	Visible bool `xml:"visible,attr"`
	// The opacity of the layer as a value from 0.0 to 1.0.
	Opacity float64 `xml:"opacity,attr"`
	// Objects associated with the object layer.
	Objects []Object `xml:"object"`
}

// An Object can be positioned anywhere on the map, and is not necessarily
// aligned to the grid.
//
// You generally use objects to add custom information to your tile map, such
// as spawn points, warps, exits, etc.
type Object struct {
	// The name of the object.
	Name string `xml:"name,attr"`
	// The type of the object.
	Type string `xml:"type,attr"`
	// The x coordinate of the object in pixels.
	X int `xml:"x,attr"`
	// The y coordinate of the object in pixels.
	Y int `xml:"y,attr"`
	// The width of the object in pixels.
	Width int `xml:"width,attr"`
	// The height of the object in pixels.
	Height int `xml:"height,attr"`
	// GID is a reference to a global tile ID.
	//
	// When the object has a GID set, then it is represented by the image of the
	// tile with that global tile ID. Currently that means Width and Height are
	// ignored for such objects. The image alignment currently depends on the map
	// orientation. In orthogonal orientation it's aligned to the bottom-left
	// while in isometric it's aligned to the bottom-center.
	GID GID `xml:"gid,attr"`
	// Properties associated with the object.
	Properties []Property `xml:"properties>property"`
	// A Polygon associated with the object.
	Polygon Polygon `xml:"polygon"`
	// A Polyline associated with the object.
	Polyline Polyline `xml:"polyline"`
}

// A Polygon object is made up of a space-delimited list of x,y coordinates. The
// origin for these coordinates is the location of the parent object. By
// default, the first point is created as 0,0 denoting that the point will
// originate exactly where the object is placed.
type Polygon struct {
	// Points contains a list of x,y coordinates in pixels.
	Points string `xml:"points,attr"`
}

// A Polyline object is made up of a space-delimited list of x,y coordinates.
// The origin for these coordinates is the location of the parent object. By
// default, the first point is created as 0,0 denoting that the point will
// originate exactly where the object is placed.
type Polyline struct {
	// Points contains a list of x,y coordinates in pixels.
	Points string `xml:"points,attr"`
}
