package hex

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/twpayne/gogeom/geom"
	"github.com/twpayne/gogeom/geom/encoding/wkb"
)

func Marshal(g geom.T, byteOrder binary.ByteOrder) (string, error) {
	wkb, err := wkb.Marshal(g, byteOrder)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(wkb), nil
}

func Unmarshal(s string) (geom.T, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return wkb.Unmarshal(data)
}
