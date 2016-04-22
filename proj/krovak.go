package proj

import (
	"fmt"
	"math"
)

// Krovak is a Krovak projection.
func Krovak(this *SR) (forward, inverse Transformer, err error) {
	this.A = 6377397.155
	this.Es = 0.006674372230614
	this.E = math.Sqrt(this.Es)
	if math.IsNaN(this.Lat0) {
		this.Lat0 = 0.863937979737193
	}
	if math.IsNaN(this.Long0) {
		this.Long0 = 0.7417649320975901 - 0.308341501185665
	}
	/* if scale not set default to 0.9999 */
	if math.IsNaN(this.K0) {
		this.K0 = 0.9999
	}
	this.S45 = 0.785398163397448 /* 45 */
	this.S90 = 2 * this.S45
	this.Fi0 = this.Lat0
	this.E2 = this.Es
	this.E = math.Sqrt(this.E2)
	this.Alfa = math.Sqrt(1 + (this.E2*math.Pow(math.Cos(this.Fi0), 4))/(1-this.E2))
	this.Uq = 1.04216856380474
	this.U0 = math.Asin(math.Sin(this.Fi0) / this.Alfa)
	this.G = math.Pow((1+this.E*math.Sin(this.Fi0))/(1-this.E*math.Sin(this.Fi0)), this.Alfa*this.E/2)
	this.K = math.Tan(this.U0/2+this.S45) / math.Pow(math.Tan(this.Fi0/2+this.S45), this.Alfa) * this.G
	this.K1 = this.K0
	this.N0 = this.A * math.Sqrt(1-this.E2) / (1 - this.E2*math.Pow(math.Sin(this.Fi0), 2))
	this.S0 = 1.37008346281555
	this.N = math.Sin(this.S0)
	this.Ro0 = this.K1 * this.N0 / math.Tan(this.S0)
	this.Ad = this.S90 - this.Uq

	/* ellipsoid */
	/* calculate xy from lat/lon */
	/* Constants, identical to inverse transform function */
	forward = func(lon, lat float64) (x, y float64, err error) {
		var gfi, u, deltav, s, d, eps, ro float64
		delta_lon := adjust_lon(lon - this.Long0)
		/* Transformation */
		gfi = math.Pow(((1 + this.E*math.Sin(lat)) / (1 - this.E*math.Sin(lat))), (this.Alfa * this.E / 2))
		u = 2 * (math.Atan(this.K*math.Pow(math.Tan(lat/2+this.S45), this.Alfa)/gfi) - this.S45)
		deltav = -delta_lon * this.Alfa
		s = math.Asin(math.Cos(this.Ad)*math.Sin(u) + math.Sin(this.Ad)*math.Cos(u)*math.Cos(deltav))
		d = math.Asin(math.Cos(u) * math.Sin(deltav) / math.Cos(s))
		eps = this.N * d
		ro = this.Ro0 * math.Pow(math.Tan(this.S0/2+this.S45), this.N) / math.Pow(math.Tan(s/2+this.S45), this.N)
		y = ro * math.Cos(eps) / 1
		x = ro * math.Sin(eps) / 1

		if !this.Czech {
			y *= -1
			x *= -1
		}
		return
	}

	/* calculate lat/lon from xy */
	inverse = func(x, y float64) (lon, lat float64, err error) {
		var u, deltav, s, d, eps, ro, fi1 float64
		var ok int

		/* Transformation */
		/* revert y, x*/
		x, y = y, x
		if !this.Czech {
			y *= -1
			x *= -1
		}
		ro = math.Sqrt(x*x + y*y)
		eps = math.Atan2(y, x)
		d = eps / math.Sin(this.S0)
		s = 2 * (math.Atan(math.Pow(this.Ro0/ro, 1/this.N)*math.Tan(this.S0/2+this.S45)) - this.S45)
		u = math.Asin(math.Cos(this.Ad)*math.Sin(s) - math.Sin(this.Ad)*math.Cos(s)*math.Cos(d))
		deltav = math.Asin(math.Cos(s) * math.Sin(d) / math.Cos(u))
		x = this.Long0 - deltav/this.Alfa
		fi1 = u
		ok = 0
		var iter = 0
		for {
			if !(ok == 0 && iter < 15) {
				break
			}
			y = 2 * (math.Atan(math.Pow(this.K, -1/this.Alfa)*math.Pow(math.Tan(u/2+this.S45), 1/this.Alfa)*math.Pow((1+this.E*math.Sin(fi1))/(1-this.E*math.Sin(fi1)), this.E/2)) - this.S45)
			if math.Abs(fi1-y) < 0.0000000001 {
				ok = 1
			}
			fi1 = y
			iter++
		}
		if iter >= 15 {
			err = fmt.Errorf("proj.Krovak: iter >= 15")
			return
		}

		return
	}
	return
}

func init() {
	registerTrans(Krovak, "Krovak", "krovak")
}
