package geom

type GeoJSONGeom interface {
	Geom
	GeoJSON() map[string]interface{}
}

func (point Point) geoJSONCoordinates() []float64 {
	return []float64{point.X, point.Y}
}

func (pointZ PointZ) geoJSONCoordinates() []float64 {
	return []float64{pointZ.X, pointZ.Y, pointZ.Z}
}

func (linearRing LinearRing) geoJSONCoordinates() [][]float64 {
	result := make([][]float64, len(linearRing))
	for i, point := range linearRing {
		result[i] = point.geoJSONCoordinates()
	}
	return result
}

func (linearRingZ LinearRingZ) geoJSONCoordinates() [][]float64 {
	result := make([][]float64, len(linearRingZ))
	for i, pointZ := range linearRingZ {
		result[i] = pointZ.geoJSONCoordinates()
	}
	return result
}

func (linearRings LinearRings) geoJSONCoordinates() [][][]float64 {
	result := make([][][]float64, len(linearRings))
	for i, linearRing := range linearRings {
		result[i] = linearRing.geoJSONCoordinates()
	}
	return result
}

func (linearRingZs LinearRingZs) geoJSONCoordinates() [][][]float64 {
	result := make([][][]float64, len(linearRingZs))
	for i, linearRingZ := range linearRingZs {
		result[i] = linearRingZ.geoJSONCoordinates()
	}
	return result
}

func (point Point) GeoJSON() map[string]interface{} {
	return map[string]interface{}{
		"type":        "Point",
		"coordinates": point.geoJSONCoordinates(),
	}
}

func (pointZ PointZ) GeoJSON() map[string]interface{} {
	return map[string]interface{}{
		"type":        "Point",
		"coordinates": pointZ.geoJSONCoordinates(),
	}
}

func (lineString LineString) GeoJSON() map[string]interface{} {
	return map[string]interface{}{
		"type":        "LineString",
		"coordinates": lineString.Points.geoJSONCoordinates(),
	}
}

func (lineStringZ LineStringZ) GeoJSON() map[string]interface{} {
	return map[string]interface{}{
		"type":        "LineString",
		"coordinates": lineStringZ.Points.geoJSONCoordinates(),
	}
}

func (polygon Polygon) GeoJSON() map[string]interface{} {
	return map[string]interface{}{
		"type":        "Polygon",
		"coordinates": polygon.Rings.geoJSONCoordinates(),
	}
}

func (polygonZ PolygonZ) GeoJSON() map[string]interface{} {
	return map[string]interface{}{
		"type":        "Polygon",
		"coordinates": polygonZ.Rings.geoJSONCoordinates(),
	}
}
