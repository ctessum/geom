package proj

import (
	"fmt"
	"math"
	"reflect"
	"strings"
)

// A Transformer takes input coordinates and returns output coordinates and an error.
type Transformer func(X, Y float64) (x, y float64, err error)

// A TransformerFunc creates forward and inverse Transformers from a projection.
type TransformerFunc func(*SR) (forward, inverse Transformer, err error)

var projections map[string]TransformerFunc

// SR holds information about a spatial reference (projection).
type SR struct {
	Name, Title                string
	SRSCode                    string
	DatumCode                  string
	Rf                         float64
	Lat0, Lat1, Lat2, LatTS    float64
	Long0, Long1, Long2, LongC float64
	Alpha                      float64
	X0, Y0, K0, K              float64
	A, A2, B, B2               float64
	Ra                         bool
	Zone                       float64
	UTMSouth                   bool
	DatumParams                []float64
	ToMeter                    float64
	Units                      string
	FromGreenwich              float64
	NADGrids                   string
	Axis                       string
	local                      bool
	sphere                     bool
	Ellps                      string
	EllipseName                string
	Es                         float64
	E                          float64
	Ep2                        float64
	DatumName                  string
	NoDefs                     bool
	datum                      *datum
	Czech                      bool
}

// NewSR initializes a SR object and sets fields to default values.
func NewSR() *SR {
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

func registerTrans(proj TransformerFunc, names ...string) {
	if projections == nil {
		projections = make(map[string]TransformerFunc)
	}
	for _, n := range names {
		projections[strings.ToLower(n)] = proj
	}
}

// Transformers returns forward and inverse transformation functions for
// this projection.
func (sr *SR) Transformers() (forward, inverse Transformer, err error) {
	t, ok := projections[strings.ToLower(sr.Name)]
	if !ok {
		err = fmt.Errorf("in proj.Proj.TransformFuncs, could not find "+
			"transformer for %s", sr.Name)
		return
	}
	forward, inverse, err = t(sr)
	return
}
