package spherical

import (
	"github.com/twpayne/gogeom/geom"
	"math"
)

func CosineDist(p1, p2 Point) float64 {
	d := math.Sin(p1.Y)*math.Sin(p2.Y) + math.Cos(p1.Y)*math.Cos(p2.Y)*math.Cos(p1.Y-p2.Y)
	if d < 1 {
		return math.Acos(d)
	} else {
		return 0
	}
}

func HaversineDist(p1, p2 Point) float64 {
	halfDeltaX := (p1.X - p2.X) / 2
	halfDeltaY := (p1.Y - p2.Y) / 2
	a := math.Sin(halfDeltaY)*math.Sin(halfDeltaY) + math.Sin(halfDeltaX)*math.Sin(halfDeltaX)*math.Cos(p1.Y)*math.Cos(p2.Y)
	return 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}
