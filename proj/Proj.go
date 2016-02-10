package proj

import (
	"fmt"
	"math"
	"reflect"
	"strings"
)

// TransformFunc takes input coordinates and returns output coordinates and an error.
type TransformFunc func(float64, float64) (float64, float64, error)

// A Transformer creates forward and inverse TransformFuncs from a projection.
type Transformer func(*Proj) (forward, inverse TransformFunc)

var projections map[string]Transformer

// Proj holds information about a spatial projection.
type Proj struct {
	projName                   string
	srsCode                    string
	datumCode                  string
	rf                         float64
	lat0, lat1, lat2, lat_ts   float64
	long0, long1, long2, longc float64
	alpha                      float64
	x0, y0, k0                 float64
	a, a2, b, b2               float64
	R_A                        bool
	zone                       int64
	utmSouth                   bool
	datum_params               []float64
	to_meter                   float64
	units                      string
	from_greenwich             float64
	nadGrids                   string
	axis                       string
	local                      bool
	sphere                     bool
	ellps                      string
	ellipseName                string
	es                         float64
	e                          float64
	k                          float64
	ep2                        float64
	datumName                  string
	datum                      *datum
}

// newProj initializes a Proj object and sets fields to default values.
func newProj() *Proj {
	p := new(Proj)
	// Initialize floats to NaN.
	v := reflect.ValueOf(p).Elem()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := f.Type().Kind()
		if ft == reflect.Float64 {
			f.SetFloat(math.NaN())
		}
	}
	p.to_meter = 1.
	return p
}

func registerTrans(proj Transformer, names ...string) {
	if projections == nil {
		projections = make(map[string]Transformer)
	}
	for _, n := range names {
		projections[strings.ToLower(n)] = Merc
	}
}

// TransformFuncs returns forward and inverse transformation functions for
// this projection.
func (p *Proj) TransformFuncs() (forward, inverse TransformFunc, err error) {
	t, ok := projections[strings.ToLower(p.projName)]
	if !ok {
		err = fmt.Errorf("in proj.Proj.TransformFuncs, could not find "+
			"transformer for %s", p.projName)
	}
	forward, inverse = t(p)
	return
}
