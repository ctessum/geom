package wkb

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
)

func testReadWrite(t *testing.T, g Geom) {
	for _, byteOrder := range []binary.ByteOrder{XDR, NDR} {
		w := bytes.NewBuffer(nil)
		err := Write(w, byteOrder, g)
		if err != nil {
			t.Errorf("Write(%q, %q, %q) == %q, want nil", w, byteOrder, g, err)
		}
		if got, err := Read(w); !reflect.DeepEqual(got, g) || err != nil {
			t.Errorf("Read(%q) == %q, %q, want %q, nil", w, got, err, g)
		}
	}
}

func TestReadWritePoint(t *testing.T) {
	testReadWrite(t, Point{1, 2})
}

func TestReadWritePointZ(t *testing.T) {
	testReadWrite(t, PointZ{1, 2, 3})
}

func TestReadWritePointM(t *testing.T) {
	testReadWrite(t, PointM{1, 2, 3})
}

func TestReadWritePointZM(t *testing.T) {
	testReadWrite(t, PointZM{1, 2, 3, 4})
}

func TestReadWriteLineStringZM(t *testing.T) {
	testReadWrite(t, LineStringZM{[]PointZM{{1, 2, 3, 4}, {5, 6, 7, 8}}})
}
