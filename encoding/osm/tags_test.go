package osm

import (
	"os"
	"reflect"
	"testing"
)

func TestCountTags(t *testing.T) {
	f, err := os.Open("testdata/honolulu_hawaii.osm.pbf")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}
	tags, err := CountTags(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 19918 {
		t.Errorf("Wrong number of tags %d", len(tags))
	}

	tags2 := tags.Filter(func(t *TagCount) bool {
		return t.Key == "highway" && t.Value == "residential"
	})
	tableWant := [][]string{
		[]string{"Key", "Value", "Total", "Node", "Closed way", "Open way", "Relation"},
		[]string{"highway", "residential", "6839", "0", "55", "6784", "0"}}
	tableHave := tags2.Table()
	if !reflect.DeepEqual(tableWant, tableHave) {
		t.Error("tables don't match")
	}
	dt := (*tags2)[0].DominantType()
	if dt != OpenWay {
		t.Errorf("dominant type should be %d but is %d", OpenWay, dt)
	}
}

func TestData_CountTags(t *testing.T) {
	f, err := os.Open("testdata/honolulu_hawaii.osm.pbf")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}
	data, err := Extract(f, func(_ *Data, _ interface{}) bool { return true })
	if err != nil {
		t.Fatal(err)
	}
	tags := data.CountTags()
	if len(tags) != 19918 {
		t.Errorf("Wrong number of tags %d", len(tags))
	}

	tags2 := tags.Filter(func(t *TagCount) bool {
		return t.Key == "highway" && t.Value == "residential"
	})
	tableWant := [][]string{
		[]string{"Key", "Value", "Total", "Node", "Closed way", "Open way", "Relation"},
		[]string{"highway", "residential", "6839", "0", "55", "6784", "0"}}
	tableHave := tags2.Table()
	if !reflect.DeepEqual(tableWant, tableHave) {
		t.Error("tables don't match")
	}
	dt := (*tags2)[0].DominantType()
	if dt != OpenWay {
		t.Errorf("dominant type should be %d but is %d", OpenWay, dt)
	}
}
