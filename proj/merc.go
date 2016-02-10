package proj

import (
	"fmt"
	"math"
)

const (
	R2D    = 57.29577951308232088
	FORTPI = math.Pi / 4
)

// Merc is a mercator projection.
func Merc(this *Proj) (forward, inverse TransformFunc) {
	var con = this.b / this.a
	this.es = 1 - con*con
	if math.IsNaN(this.x0) {
		this.x0 = 0
	}
	if math.IsNaN(this.y0) {
		this.y0 = 0
	}
	this.e = math.Sqrt(this.es)
	if !math.IsNaN(this.lat_ts) {
		if this.sphere {
			this.k0 = math.Cos(this.lat_ts)
		} else {
			this.k0 = msfnz(this.e, math.Sin(this.lat_ts), math.Cos(this.lat_ts))
		}
	} else {
		if math.IsNaN(this.k0) {
			if !math.IsNaN(this.k) {
				this.k0 = this.k
			} else {
				this.k0 = 1
			}
		}
	}

	// Mercator forward equations--mapping lat,long to x,y
	forward = func(lon, lat float64) (x, y float64, err error) {
		// convert to radians
		if lat*R2D > 90 || lat*R2D < -90 || lon*R2D > 180 || lon*R2D < -180 {
			err = fmt.Errorf("in proj.Merc forward: invalid longitude (%g) or latitude (%g)", lon, lat)
			return
		}

		if math.Abs(math.Abs(lat)-HALF_PI) <= EPSLN {
			err = fmt.Errorf("in proj.Merc forward, abs(lat)==pi/2")
			return
		}
		if this.sphere {
			x = this.x0 + this.a*this.k0*adjust_lon(lon-this.long0)
			y = this.y0 + this.a*this.k0*math.Log(math.Tan(FORTPI+0.5*lat))
		} else {
			var sinphi = math.Sin(lat)
			var ts = tsfnz(this.e, lat, sinphi)
			x = this.x0 + this.a*this.k0*adjust_lon(lon-this.long0)
			y = this.y0 - this.a*this.k0*math.Log(ts)
		}
		return
	}

	// Mercator inverse equations--mapping x,y to lat/long
	inverse = func(x, y float64) (lon, lat float64, err error) {
		x -= this.x0
		y -= this.y0

		if this.sphere {
			lat = HALF_PI - 2*math.Atan(math.Exp(-y/(this.a*this.k0)))
		} else {
			var ts = math.Exp(-y / (this.a * this.k0))
			lat, err = phi2z(this.e, ts)
			if err != nil {
				return
			}
		}
		lon = adjust_lon(this.long0 + x/(this.a*this.k0))
		return
	}
	return
}

func init() {
	// Register this projection with the corresponding names.
	registerTrans(Merc, "Mercator", "Popular Visualisation Pseudo Mercator",
		"Mercator_1SP", "Mercator_Auxiliary_Sphere", "merc")
}
