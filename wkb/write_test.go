package wkb

import (
	"bytes"
	"encoding/binary"
	"testing"
)

type writeTestCase struct {
	g         Geom
	byteOrder binary.ByteOrder
	want      string
}

func testWrite(t *testing.T, tc writeTestCase) {
	w := bytes.NewBuffer(nil)
	if err := Write(w, tc.byteOrder, tc.g); err != nil {
		t.Errorf("Write(%q, %q, %q) == %q, want nil", w, tc.byteOrder, tc.g, err)
	}
	if got := w.String(); got != tc.want {
		t.Errorf("expected Write(%q, %q, %q) to write %q, got %q", w, tc.byteOrder, tc.g, tc.want, got)
	}
}

func TestWritePointZM(t *testing.T) {
	testWrite(t, writeTestCase{PointZM{1, 2, 3, 4}, NDR, "\x01\xb9\x0b\x00\x00\x00\x00\x00\x00\x00\x00\xf0?\x00\x00\x00\x00\x00\x00\x00@\x00\x00\x00\x00\x00\x00\x08@\x00\x00\x00\x00\x00\x00\x10@"})
}
