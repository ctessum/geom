package shp

import (
	"fmt"

	"github.com/ctessum/geom"
)

// This example shows how to read in information from a shapefile, write
// it out again, then read it in a second time to ensure that it was
// properly written.
func Example() {
	type record struct {
		// The geometry data will be stored here.
		geom.Polygon

		// The "value" attribute will be stored here. It would also work to
		// just name this field "Value".
		Val float64 `shp:"value"`
	}

	d, err := NewDecoder("testdata/triangles.shp")
	if err != nil {
		panic(err)
	}

	e, err := NewEncoder("testdata/testout.shp", record{})
	if err != nil {
		panic(err)
	}

	for {
		var rec record
		// Decode a record from the input file.
		if !d.DecodeRow(&rec) {
			break
		}
		// Encode the record to the output file
		if err = e.Encode(rec); err != nil {
			panic(err)
		}
		fmt.Printf("polygon area %.3g, value %g\n", rec.Polygon.Area(), rec.Val)
	}
	// Check to see if any errors occured during decoding.
	if err = d.Error(); err != nil {
		panic(err)
	}
	d.Close()
	e.Close()

	// Read the data back in from the output file.
	d, err = NewDecoder("testdata/testout.shp")
	if err != nil {
		panic(err)
	}
	for {
		var rec record
		// Decode a record from the input file.
		if !d.DecodeRow(&rec) {
			break
		}
		fmt.Printf("polygon area %.3g, value %g\n", rec.Polygon.Area(), rec.Val)
	}
	// Output:
	// polygon area 2.3, value 6
	// polygon area 2.17, value 1
	// polygon area 2.3, value 6
	// polygon area 2.17, value 1
}
