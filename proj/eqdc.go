package proj

import (
	"fmt"
	"math"
)

// EqdC is an Equidistant Conic projection.
func EqdC(this *SR) (forward, inverse Transformer, err error) {
	// Standard Parallels cannot be equal and on opposite sides of the equator
	if math.Abs(this.Lat1+this.Lat2) < epsln {
		return nil, nil, fmt.Errorf("proj: Equidistant Conic parallels cannot be equal and on opposite sides of the equator but are %g and %g", this.Lat1, this.Lat2)
	}
	if math.IsNaN(this.Lat2) {
		this.Lat2 = this.Lat1
	}

	temp := this.B / this.A
	this.Es = 1 - math.Pow(temp, 2)
	this.E = math.Sqrt(this.Es)
	e0 := e0fn(this.Es)
	e1 := e1fn(this.Es)
	e2 := e2fn(this.Es)
	e3 := e3fn(this.Es)

	sinphi := math.Sin(this.Lat1)
	cosphi := math.Cos(this.Lat1)

	ms1 := msfnz(this.E, sinphi, cosphi)
	ml1 := mlfn(e0, e1, e2, e3, this.Lat1)

	var ns float64
	if math.Abs(this.Lat1-this.Lat2) < epsln {
		ns = sinphi
	} else {
		sinphi = math.Sin(this.Lat2)
		cosphi = math.Cos(this.Lat2)
		ms2 := msfnz(this.E, sinphi, cosphi)
		ml2 := mlfn(e0, e1, e2, e3, this.Lat2)
		ns = (ms1 - ms2) / (ml2 - ml1)
	}
	g := ml1 + ms1/ns
	ml0 := mlfn(e0, e1, e2, e3, this.Lat0)
	rh := this.A * (g - ml0)

	/* Equidistant Conic forward equations--mapping lat,long to x,y
	   -----------------------------------------------------------*/
	forward = func(lon, lat float64) (x, y float64, err error) {
		var rh1 float64

		/* Forward equations
		   -----------------*/
		if this.sphere {
			rh1 = this.A * (g - lat)
		} else {
			var ml = mlfn(e0, e1, e2, e3, lat)
			rh1 = this.A * (g - ml)
		}
		var theta = ns * adjust_lon(lon-this.Long0)
		x = this.X0 + rh1*math.Sin(theta)
		y = this.Y0 + rh - rh1*math.Cos(theta)
		return x, y, nil
	}

	/* Inverse equations
	   -----------------*/
	inverse = func(x, y float64) (lon, lat float64, err error) {
		x -= this.X0
		y = rh - y + this.Y0
		var con, rh1 float64
		if ns >= 0 {
			rh1 = math.Sqrt(x*x + y*y)
			con = 1
		} else {
			rh1 = -math.Sqrt(x*x + y*y)
			con = -1
		}
		var theta float64
		if rh1 != 0 {
			theta = math.Atan2(con*x, con*y)
		}

		if this.sphere {
			lon = adjust_lon(this.Long0 + theta/ns)
			lat = adjust_lat(g - rh1/this.A)
			return
		}
		var ml = g - rh1/this.A
		lat, err = imlfn(ml, e0, e1, e2, e3)
		if err != nil {
			return math.NaN(), math.NaN(), err
		}
		lon = adjust_lon(this.Long0 + theta/ns)
		return
	}
	return
}

func init() {
	registerTrans(EqdC, "Equidistant_Conic", "eqdc")
}
