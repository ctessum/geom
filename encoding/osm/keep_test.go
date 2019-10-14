package osm

import (
	"testing"

	"github.com/paulmach/osm"
)

func TestKeepTags(t *testing.T) {
	k := KeepTags(map[string][]string{"x": []string{}})
	r := &osm.Node{
		Tags: []osm.Tag{
			{
				Key:   "x",
				Value: "y",
			},
		},
	}
	keep := k(nil, r)
	if !keep {
		t.Errorf("keep should be true but is false")
	}
}
