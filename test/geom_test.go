package test

import (
	"github.com/twpayne/gogeom/geom/encoding/hex"
	"github.com/twpayne/gogeom/geom/encoding/wkb"
	"reflect"
	"testing"
)

func TestHexEncode(t *testing.T) {
	for _, c := range cases {
		if got, err := hex.Encode(c.g, wkb.NDR); err != nil || got != c.hex {
			t.Errorf("hex.Encode(%#v, %#v) == %#v, %#v, want %#v, nil", c.g, wkb.NDR, got, err, c.hex)
		}
	}
}

func TestHexDecode(t *testing.T) {
	for _, c := range cases {
		if got, err := hex.Decode(c.hex); err != nil || !reflect.DeepEqual(got, c.g) {
			t.Errorf("hex.Decode(%#v) == %#v, %#v, want %#v, nil", c.wkb, got, err, c.g)
		}
	}
}

func TestWKBDecode(t *testing.T) {
	for _, c := range cases {
		if got, err := wkb.Encode(c.g, wkb.NDR); err != nil || !reflect.DeepEqual(got, c.wkb) {
			t.Errorf("wkb.Encode(%#v, %#v) == %#v, %#v, want %#v, nil", c.g, wkb.NDR, got, err, c.wkb)
		}
	}
}

func TestWKBEncode(t *testing.T) {
	for _, c := range cases {
		if got, err := wkb.Decode(c.wkb); err != nil || !reflect.DeepEqual(got, c.g) {
			t.Errorf("wkb.Decode(%#v) == %#v, %#v, want %#v, nil", c.wkb, got, err, c.g)
		}
	}
}
