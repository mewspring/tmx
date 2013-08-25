// Package tmx implements access to Tiled's tmx (Tile Map XML) files.
//
// Documentation and specification of the TMX file format has been based on the
// information available at:
//    https://github.com/bjorn/tiled/wiki/TMX-Map-Format
package tmx

import (
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Open reads the provided tmx file and returns a parsed Map, based on the TMX
// file format.
func Open(tmxPath string) (m *Map, err error) {
	fr, err := os.Open(tmxPath)
	if err != nil {
		return nil, err
	}
	defer fr.Close()
	return NewFile(fr)
}

// NewFile reads from the provided io.Reader and returns a parsed Map, based on
// the TMX file format.
func NewFile(r io.Reader) (m *Map, err error) {
	d := xml.NewDecoder(r)
	m = new(Map)
	err = d.Decode(m)
	if err != nil {
		return nil, err
	}
	for _, l := range m.Layers {
		err = l.Data.decode(m.Width, m.Height)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

// decode decodes the GIDs that are stored in the <data> XML-tag of a layer. It
// will handle the various encodings and compression methods.
func (data *Data) decode(cols, rows int) (err error) {
	if data.gids != nil {
		// data has already been decoded.
		return nil
	}
	// alloc
	data.gids = make([][]GID, cols)
	for i := range data.gids {
		data.gids[i] = make([]GID, rows)
	}
	// decode
	switch data.Encoding {
	case "base64":
		err = data.decodeBase64(cols, rows)
		if err != nil {
			return err
		}
	case "csv":
		err = data.decodeCsv(cols, rows)
		if err != nil {
			return err
		}
	case "": // XML encoding
		err = data.decodeXml(cols, rows)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("decodeData: encoding '%s' not yet implemented.", data.Encoding)
	}
	return nil
}

// decodeBase64 decodes the GIDs that are stored as a base64-encoded array of
// unsigned 32-bit integers, using little-endian byte ordering. This array may
// be compressed using gzip or zlib.
func (data *Data) decodeBase64(cols, rows int) (err error) {
	s := strings.TrimSpace(data.RawData)
	r := base64.NewDecoder(base64.StdEncoding, strings.NewReader(s))
	switch data.Compression {
	case "gzip":
		z, err := gzip.NewReader(r)
		if err != nil {
			return err
		}
		defer z.Close()
		r = z
	case "zlib":
		z, err := zlib.NewReader(r)
		if err != nil {
			return err
		}
		defer z.Close()
		r = z
	case "": // no compression.
		break
	default:
		return fmt.Errorf("decodeBase64: compression '%s' not yet implemented.", data.Compression)
	}
	buf, err := ioutil.ReadAll(r)
	// We should have one GID for each tile.
	if len(buf)/4 != cols*rows {
		return fmt.Errorf("decodeBase64: wrong number of GIDs. Got %d, wanted %d.", len(buf)/4, cols*rows)
	}
	i := 0
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			gid := binary.LittleEndian.Uint32(buf[i*4:])
			data.gids[col][row] = GID(gid)
			i++
		}
	}
	return nil
}

// decodeCvs decodes the GIDs that are stored as comma-separated values.
func (data *Data) decodeCsv(cols, rows int) (err error) {
	cleanData := strings.Map(clean, data.RawData)
	rawGIDs := strings.Split(cleanData, ",")
	// We should have one GID for each tile.
	if len(rawGIDs) != cols*rows {
		return fmt.Errorf("decodeCsv: wrong number of GIDs. Got %d, wanted %d.", len(rawGIDs), cols*rows)
	}
	i := 0
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			gid, err := strconv.Atoi(rawGIDs[i])
			if err != nil {
				return err
			}
			data.gids[col][row] = GID(gid)
			i++
		}
	}
	return nil
}

// clean cleans the csv data from superfluous runes.
func clean(r rune) rune {
	if r >= '0' && r <= '9' || r == ',' {
		return r
	}
	// skip rune.
	return -1
}

// decodeXml decodes the GIDs that are stored in the <tile> XML-tags' 'gid'
// attribute.
func (data *Data) decodeXml(cols, rows int) (err error) {
	if len(data.Tiles) != cols*rows {
		return fmt.Errorf("decodeXml: wrong number of GIDs. Got %d, wanted %d.", len(data.Tiles), cols*rows)
	}
	i := 0
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			data.gids[col][row] = data.Tiles[i].GID
			i++
		}
	}
	return nil
}

// GetGID returns the global tile ID at a given coordinate, after clearing the
// flip flags.
func (l *Layer) GetGID(col, row int) int {
	return l.Data.gids[col][row].GlobalTileID()
}

// GetRawGID returns the global tile ID at a given coordinate, without clearing
// the flip flags.
func (l *Layer) GetRawGID(col, row int) GID {
	return l.Data.gids[col][row]
}

// GlobalTileID returns the GID after clearing the flip flags.
func (gid GID) GlobalTileID() int {
	return int(gid &^ FlagFlip)
}

// IsDiagonalFlip returns true if the GID is flipped diagonally.
func (gid GID) IsDiagonalFlip() bool {
	if gid&FlagDiagonalFlip != 0 {
		return true
	}
	return false
}

// IsVerticalFlip returns true if the GID is flipped vertically.
func (gid GID) IsVerticalFlip() bool {
	if gid&FlagVerticalFlip != 0 {
		return true
	}
	return false
}

// IsHorizontalFlip returns true if the GID is flipped horizontally.
func (gid GID) IsHorizontalFlip() bool {
	if gid&FlagHorizontalFlip != 0 {
		return true
	}
	return false
}

// IsFlip returns true if the GID is flipped.
func (gid GID) IsFlip() bool {
	if gid&FlagFlip != 0 {
		return true
	}
	return false
}
