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
type Transformer func(*SR) (forward, inverse TransformFunc, err error)

var projections map[string]Transformer

// SR holds information about a spatial reference (projection).
type SR struct {
	Name, Title                                                string
	SRSCode                                                    string
	DatumCode                                                  string
	Rf                                                         float64
	Lat0, Lat1, Lat2, LatTS                                    float64
	Long0, Long1, Long2, LongC                                 float64
	Alpha                                                      float64
	X0, Y0, K0                                                 float64
	A, A2, B, B2                                               float64
	Ra                                                         bool
	Zone                                                       float64
	UTMSouth                                                   bool
	DatumParams                                                []float64
	ToMeter                                                    float64
	Units                                                      string
	FromGreenwich                                              float64
	NADGrids                                                   string
	Axis                                                       string
	local                                                      bool
	sphere                                                     bool
	Ellps                                                      string
	EllipseName                                                string
	Es                                                         float64
	E                                                          float64
	K                                                          float64
	Ep2                                                        float64
	DatumName                                                  string
	NoDefs                                                     bool
	datum                                                      *datum
	NS                                                         float64
	F0                                                         float64
	RH                                                         float64
	S45, S90, Fi0, E2, Alfa, Uq, U0, G, K1, N0, S0, N, Ro0, Ad float64
	Czech                                                      bool
}

// newProj initializes a SR object and sets fields to default values.
func newSR() *SR {
	p := new(SR)
	// Initialize floats to NaN.
	v := reflect.ValueOf(p).Elem()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := f.Type().Kind()
		if ft == reflect.Float64 {
			f.SetFloat(math.NaN())
		}
	}
	p.ToMeter = 1.
	return p
}

func registerTrans(proj Transformer, names ...string) {
	if projections == nil {
		projections = make(map[string]Transformer)
	}
	for _, n := range names {
		projections[strings.ToLower(n)] = proj
	}
}

// TransformFuncs returns forward and inverse transformation functions for
// this projection.
func (sr *SR) TransformFuncs() (forward, inverse TransformFunc, err error) {
	t, ok := projections[strings.ToLower(sr.Name)]
	if !ok {
		err = fmt.Errorf("in proj.Proj.TransformFuncs, could not find "+
			"transformer for %s", sr.Name)
		return
	}
	forward, inverse, err = t(sr)
	return
}
