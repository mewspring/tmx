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
	cleanData := strings.TrimSpace(data.RawData)
	buf, err := base64.StdEncoding.DecodeString(cleanData)
	if err != nil {
		return err
	}
	switch data.Compression {
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
		return fmt.Errorf("decodeDataBase64: compression '%s' not yet implemented.", data.Compression)
	}
	// We should have one GID for each tile.
	if len(buf)/4 != cols*rows {
		return fmt.Errorf("decodeDataBase64: wrong number of GIDs. Got %d, wanted %d.", len(buf)/4, cols*rows)
	}
	i := 0
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			gid := binary.LittleEndian.Uint32(buf[i*4:])
			data.gids[col][row] = GID(gid)
			if err != nil {
				return err
			}
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
		return fmt.Errorf("decodeDataCsv: wrong number of GIDs. Got %d, wanted %d.", len(rawGIDs), cols*rows)
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
		return fmt.Errorf("decodeDataXml: wrong number of GIDs. Got %d, wanted %d.", len(data.Tiles), cols*rows)
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

// GetGID returns the global tile ID at a given coordinate.
func (l *Layer) GetGID(col, row int) GID {
	return l.Data.gids[col][row]
}
