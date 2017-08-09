package shp

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/ctessum/geom"
)

func TestEncoder_polygon(t *testing.T) {

	const testFile = "testdata/test_output"

	type polygon struct {
		geom.Polygon
	}

	p := polygon{
		Polygon: geom.Polygon{
			geom.Path{
				geom.Point{X: -104, Y: 42},
				geom.Point{X: -104, Y: 44},
				geom.Point{X: -100, Y: 44},
				geom.Point{X: -100, Y: 42},
				geom.Point{X: -104, Y: 42},
			},
			geom.Path{
				geom.Point{X: -103, Y: 45},
				geom.Point{X: -102, Y: 45},
				geom.Point{X: -102, Y: 44},
				geom.Point{X: -103, Y: 44},
				geom.Point{X: -103, Y: 45},
			},
		},
	}

	shape, err := NewEncoder(testFile+".shp", polygon{})
	if err != nil {
		t.Fatalf("error creating output shapefile: %v", err)
	}
	if err = shape.Encode(p); err != nil {
		fmt.Printf("error writing output shapefile: %v", err)
	}
	shape.Close()

	// Load geometries.
	d, err := NewDecoder(testFile + ".shp")
	if err != nil {
		panic(err)
	}

	var p2 polygon
	d.DecodeRow(&p2)

	if len(p.Polygon) != len(p2.Polygon) {
		t.Fatalf("polygons have different numbers of rings: %d != %d", len(p.Polygon), len(p2.Polygon))
	}
	if !reflect.DeepEqual(p.Polygon[0], p2.Polygon[0]) {
		t.Errorf("ring 0 is different: %+v != %+v", p.Polygon[0], p2.Polygon[0])
	}
	if !reflect.DeepEqual(p.Polygon[1], p2.Polygon[1]) {
		t.Errorf("ring 1 is different: %+v != %+v", p.Polygon[1], p2.Polygon[1])
	}

	// Check to see if any errors occured during decoding.
	if err := d.Error(); err != nil {
		t.Fatalf("error decoding shapefile: %v", err)
	}
	d.Close()
	os.Remove(testFile + ".shp")
	os.Remove(testFile + ".shx")
	os.Remove(testFile + ".dbf")
}
