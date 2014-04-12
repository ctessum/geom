package gis

import (
	"encoding/json"
	"github.com/twpayne/gogeom/geom/encoding/geojson"
	"io"
)

type GeoJSONfeature struct {
	Type       string
	Geometry   *geojson.Geometry
	Properties map[string]float64
}
type GeoJSON struct {
	Type     string
	Features []*GeoJSONfeature
}

func LoadGeoJSON(r io.Reader) (*GeoJSON, error) {
	out := new(GeoJSON)
	d := json.NewDecoder(r)
	err := d.Decode(&out)
	return out, err
}

func (g *GeoJSON) Sum(propertyName string) float64 {
	sum := 0.
	for _, f := range g.Features {
		sum += f.Properties[propertyName]
	}
	return sum
}
