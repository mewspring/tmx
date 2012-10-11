package tmx

import "bytes"
import "compress/gzip"
import "compress/zlib"
import "encoding/base64"
import "encoding/binary"
import "encoding/xml"
import "fmt"
import "io"
import "io/ioutil"
import "os"
import "strconv"
import "strings"

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
		err = l.decodeData(m.Width, m.Height)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

// decodeData decodes the GIDs that are stored in the <data> XML-tag of a layer.
// It will handle the various encodings and compression methods.
func (l *Layer) decodeData(cols, rows int) (err error) {
	if l.gids != nil {
		// data has already been decoded.
		return nil
	}
	// alloc
	l.gids = make([][]GID, cols)
	for i := range l.gids {
		l.gids[i] = make([]GID, rows)
	}
	// decode
	switch l.Data.Encoding {
	case "base64":
		err = l.decodeDataBase64(cols, rows)
		if err != nil {
			return err
		}
	case "csv":
		err = l.decodeDataCsv(cols, rows)
		if err != nil {
			return err
		}
	case "": // XML encoding
		err = l.decodeDataXml(cols, rows)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("decodeData: encoding '%s' not yet implemented.", l.Data.Encoding)
	}
	return nil
}

// decodeDataBase64 decodes the GIDs that are stored as a base64-encoded array
// of unsigned 32-bit integers, using little-endian byte ordering. This array
// may be compressed using gzip or zlib.
func (l *Layer) decodeDataBase64(cols, rows int) (err error) {
	cleanData := strings.TrimSpace(l.Data.RawData)
	buf, err := base64.StdEncoding.DecodeString(cleanData)
	if err != nil {
		return err
	}
	switch l.Data.Compression {
	case "gzip":
		r := bytes.NewBuffer(buf)
		z, err := gzip.NewReader(r)
		if err != nil {
			return err
		}
		defer z.Close()
		buf, err = ioutil.ReadAll(z)
		if err != nil {
			return err
		}
	case "zlib":
		r := bytes.NewBuffer(buf)
		z, err := zlib.NewReader(r)
		if err != nil {
			return err
		}
		defer z.Close()
		buf, err = ioutil.ReadAll(z)
		if err != nil {
			return err
		}
	case "": // no compression.
		break
	default:
		return fmt.Errorf("decodeDataBase64: compression '%s' not yet implemented.", l.Data.Compression)
	}
	// We should have one GID for each tile.
	if len(buf)/4 != cols*rows {
		return fmt.Errorf("decodeDataBase64: wrong number of GIDs. Got %d, wanted %d.", len(buf)/4, cols*rows)
	}
	i := 0
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			gid := binary.LittleEndian.Uint32(buf[i*4:])
			l.gids[col][row] = GID(gid)
			if err != nil {
				return err
			}
			i++
		}
	}
	return nil
}

// decodeDataCvs decodes the GIDs that are stored as comma-separated values.
func (l *Layer) decodeDataCsv(cols, rows int) (err error) {
	cleanData := strings.Map(clean, l.Data.RawData)
	rawGIDs := strings.Split(cleanData, ",")
	// We should have one GID for each tile.
	if len(rawGIDs) != cols*rows {
		return fmt.Errorf("decodeDataCsv: wrong number of GIDs. Got %d, wanted %d.", len(rawGIDs), cols*rows)
	}
	i := 0
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			gid, err := strconv.Atoi(rawGIDs[i])
			if err != nil {
				return err
			}
			l.gids[col][row] = GID(gid)
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

// decodeDataCvs decodes the GIDs that are stored in the <tile> XML-tags' 'gid'
// attribute.
func (l *Layer) decodeDataXml(cols, rows int) (err error) {
	if len(l.Data.Tiles) != cols*rows {
		return fmt.Errorf("decodeDataXml: wrong number of GIDs. Got %d, wanted %d.", len(l.Data.Tiles), cols*rows)
	}
	i := 0
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			l.gids[col][row] = l.Data.Tiles[i].GID
			i++
		}
	}
	return nil
}

// GetGID returns the global tile ID at a given coordinate.
func (l *Layer) GetGID(col, row int) GID {
	return l.gids[col][row]
}
