package proj

import (
	"fmt"
	"math"
)

func msfnz(eccent, sinphi, cosphi float64) float64 {
	var con = eccent * sinphi
	return cosphi / (math.Sqrt(1 - con*con))
}

func sign(x float64) float64 {
	if x < 0 {
		return -1
	}
	return 1
}

const (
	twoPi = math.Pi * 2
	// SPI is slightly greater than Math.PI, so values that exceed the -180..180
	// degree range by a tiny amount don't get wrapped. This prevents points that
	// have drifted from their original location along the 180th meridian (due to
	// floating point error) from changing their sign.
	sPi    = 3.14159265359
	halfPi = math.Pi / 2
)

func adjust_lon(x float64) float64 {
	if math.Abs(x) <= sPi {
		return x
	}
	return (x - (sign(x) * twoPi))
}

func tsfnz(eccent, phi, sinphi float64) float64 {
	var con = eccent * sinphi
	var com = 0.5 * eccent
	con = math.Pow(((1 - con) / (1 + con)), com)
	return (math.Tan(0.5*(halfPi-phi)) / con)
}

func phi2z(eccent, ts float64) (float64, error) {
	var eccnth = 0.5 * eccent
	phi := halfPi - 2*math.Atan(ts)
	for i := 0; i <= 15; i++ {
		con := eccent * math.Sin(phi)
		dphi := halfPi - 2*math.Atan(ts*(math.Pow(((1-con)/(1+con)), eccnth))) - phi
		phi += dphi
		if math.Abs(dphi) <= 0.0000000001 {
			return phi, nil
		}
	}
	return math.NaN(), fmt.Errorf("phi2z has no convergence")
}
