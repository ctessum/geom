package martinez

// signdedArea returns the signed area of the triangle (p0, p1, p2).
func signedArea(p0, p1, p2 Point) float64 {
	return (p0.X()-p2.X())*(p1.Y()-p2.Y()) - (p1.X()-p2.X())*(p0.Y()-p2.Y())
}
