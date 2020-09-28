package osm

import (
	"os"
	"testing"
)

func TestCheck(t *testing.T) {
	f, err := os.Open("testdata/honolulu_hawaii.osm.pbf")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}
	data, err := ExtractTag(f, "source", true, "Bing")
	if err != nil {
		t.Fatal(err)
	}
	if err = data.Check(); err != nil {
		t.Fatal(err)
	}
}
