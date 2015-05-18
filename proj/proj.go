/*
Package projgeom is performs geodesic reprojections on
Open GIS Consortium style geometry objects.
It is an interface between
	"github.com/pebbe/go-proj-4/proj"
and
	"github.com/twpayne/gogeom/geom"
*/
package proj

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	"github.com/ctessum/geom"
	"github.com/lukeroth/gdal"
	"github.com/pebbe/go-proj-4/proj"
)

type UnsupportedGeometryError struct {
	Type reflect.Type
}

func (e UnsupportedGeometryError) Error() string {
	return "projgeom: unsupported geometry type: " + e.Type.String()
}

// Project geometry from src to dst projection. inputDegrees and outputDegrees are `true` if
// the input or output geometries is in units of degrees. We need to know this
// because the Proj4 library works in units of radians.
// Because I don't know whether to transform Z values from degrees to radians or
// not, Z values are not supported.
// I also don't know what to do with M values so they are not supported either.
func project(g geom.T, src, dst *proj.Proj, inputDegrees, outputDegrees bool) (geom.T, error) {
	if g == nil {
		return nil, nil
	}
	switch g.(type) {
	case geom.Point:
		point := g.(geom.Point)
		return projectPoint(&point, src, dst, inputDegrees, outputDegrees)
	//case geom.PointZ:
	//	pointZ := g.(geom.PointZ)
	//	return projectPointZ(&pointZ, src, dst)
	//case geom.PointM:
	//	pointM := g.(geom.PointM)
	//	return projectPointM(&pointM, src, dst)
	//case geom.PointZM:
	//	pointZM := g.(geom.PointZM)
	//	return projectPointZM(&pointZM, src, dst)
	case geom.LineString:
		lineString := g.(geom.LineString)
		return projectLineString(lineString, src, dst, inputDegrees,
			outputDegrees)
	//case geom.LineStringZ:
	//	lineStringZ := g.(geom.LineStringZ)
	//	return projectLineStringZ(&lineStringZ, src, dst)
	//case geom.LineStringM:
	//	lineStringM := g.(geom.LineStringM)
	//	return projectLineStringM(&lineStringM, src, dst)
	//case geom.LineStringZM:
	//	lineStringZM := g.(geom.LineStringZM)
	//	return projectLineStringZM(&lineStringZM, src, dst)
	case geom.MultiLineString:
		multiLineString := g.(geom.MultiLineString)
		return projectMultiLineString(multiLineString, src, dst,
			inputDegrees, outputDegrees)
	case geom.Polygon:
		polygon := g.(geom.Polygon)
		return projectPolygon(polygon, src, dst,
			inputDegrees, outputDegrees)
	//case geom.PolygonZ:
	//	polygonZ := g.(geom.PolygonZ)
	//	return projectPolygonZ(&polygonZ, src, dst)
	//case geom.PolygonM:
	//	polygonM := g.(geom.PolygonM)
	//	return projectPolygonM(&polygonM, src, dst)
	//case geom.PolygonZM:
	//	polygonZM := g.(geom.PolygonZM)
	//	return projectPolygonZM(&polygonZM, src, dst), nil
	case geom.MultiPolygon:
		multiPolygon := g.(geom.MultiPolygon)
		return projectMultiPolygon(multiPolygon, src, dst,
			inputDegrees, outputDegrees)
	default:
		return nil, &UnsupportedGeometryError{reflect.TypeOf(g)}
	}
}

type CoordinateTransform struct {
	src, dst                    *proj.Proj
	sameProj                    bool
	inputDegrees, outputDegrees bool
}

func NewCoordinateTransform(src, dst SR) (
	ct *CoordinateTransform, err error) {
	ct = new(CoordinateTransform)
	ct.sameProj = gdal.SpatialReference(src).IsSame(gdal.SpatialReference(dst))
	var srcproj, dstproj string
	if !ct.sameProj {
		srcproj, err = gdal.SpatialReference(src).ToProj4()
		if err != nil && err.Error() != "No Error" {
			return
		}
		ct.inputDegrees = strings.Contains(srcproj, "longlat") ||
			strings.Contains(srcproj, "latlong")
		ct.src, err = proj.NewProj(srcproj)
		if err != nil {
			return
		}

		dstproj, err = gdal.SpatialReference(dst).ToProj4()
		if err != nil && err.Error() != "No Error" {
			return
		}
		ct.outputDegrees = strings.Contains(dstproj, "longlat") ||
			strings.Contains(dstproj, "latlong")
		ct.dst, err = proj.NewProj(dstproj)
		if err != nil {
			return
		}
	}
	return
}

func (ct *CoordinateTransform) Reproject(g geom.T) (geom.T, error) {
	if ct.sameProj {
		return g, nil
	}
	g2, err := project(g, ct.src, ct.dst,
		ct.inputDegrees, ct.outputDegrees)
	return g2, err
}

// ReadPrj reads an ESRI '.prj' projection file and
// creates a corresponding spatial reference.
func ReadPrj(f io.Reader) (SR, error) {
	sr := gdal.CreateSpatialReference("")
	prj, err := ioutil.ReadAll(f)
	if err != nil {
		return SR{}, err
	}
	err = sr.FromWKT(string(prj))
	if err != nil && err.Error() == "No Error" {
		err = nil
	}
	return SR(sr), err
}

// FromProj4 converts a Proj4 string into the corresponding
// spatial reference
func FromProj4(proj4 string) (SR, error) {
	sr := gdal.CreateSpatialReference("")
	err := sr.FromProj4(proj4)
	if err != nil && err.Error() == "No Error" {
		err = nil
	}
	return SR(sr), err
}

// SR is a spatial reference or coordinate system to use in
// reprojections.
type SR gdal.SpatialReference

type ParsedProj4 struct {
	SRID          int
	Proj          string
	Lat_1         float64
	Lat_2         float64
	Lat_0         float64
	Lon_0         float64
	EarthRadius_a float64
	EarthRadius_b float64
	To_meter      float64
}

func (p *ParsedProj4) Equals(p2 *ParsedProj4) bool {
	switch p.Proj {
	case "lcc":
		return (p.Proj == p2.Proj &&
			p.Lat_1 == p2.Lat_1 &&
			p.Lat_2 == p2.Lat_2 &&
			p.Lat_0 == p2.Lat_0 &&
			p.Lon_0 == p2.Lon_0 &&
			p.EarthRadius_a == p2.EarthRadius_a &&
			p.EarthRadius_b == p2.EarthRadius_b &&
			p.To_meter == p2.To_meter)
	case "longlat":
		return (p.Proj == p2.Proj)
	case "merc":
		return (p.Proj == p2.Proj &&
			p.Lon_0 == p2.Lon_0 &&
			p.EarthRadius_a == p2.EarthRadius_a &&
			p.EarthRadius_b == p2.EarthRadius_b)
	}
	panic("Unsupported projection type " + p.Proj)
	return false
}

func parseHelper(s1, s2 string) string {
	return strings.Split(strings.Split(strings.ToLower(s1),
		"+"+strings.ToLower(s2)+"=")[1], " ")[0]
}
func parseHelperFloat(s1, s2 string) float64 {
	f, _ := strconv.ParseFloat(parseHelper(s1, s2), 64)
	return f
}

func ParseProj4(proj4 string) *ParsedProj4 {
	p := new(ParsedProj4)
	p.Proj = parseHelper(proj4, "proj")
	switch p.Proj {
	case "lcc":
		p.Lat_1 = parseHelperFloat(proj4, "lat_1")
		p.Lat_2 = parseHelperFloat(proj4, "lat_2")
		p.Lat_0 = parseHelperFloat(proj4, "lat_0")
		p.Lon_0 = parseHelperFloat(proj4, "lon_0")
		p.To_meter = parseHelperFloat(proj4, "to_meter")
		p.EarthRadius_a = parseHelperFloat(proj4, "a")
		p.EarthRadius_b = parseHelperFloat(proj4, "b")
	case "merc":
		p.Lon_0 = parseHelperFloat(proj4, "lon_0")
		p.EarthRadius_a = parseHelperFloat(proj4, "a")
		p.EarthRadius_b = parseHelperFloat(proj4, "b")
	}
	return p
}

func (p *ParsedProj4) ToString() string {
	var s string
	switch p.Proj {
	case "lcc":
		s = fmt.Sprintf("+proj=lcc +lat_1=%f +lat_2=%f +lat_0=%f +lon_0=%f "+
			"+x_0=0 +y_0=0 +a=%f +b=%f +to_meter=%v",
			p.Lat_1, p.Lat_2, p.Lat_0, p.Lon_0,
			p.EarthRadius_a, p.EarthRadius_b, p.To_meter)
	case "longlat":
		s = "+proj=longlat +datum=WGS84 +no_defs"
	case "merc":
		s = fmt.Sprintf("+proj=merc +a=%f +b=%f +lat_ts=0.0 +lon_0=%f "+
			"+x_0=0.0 +y_0=0 +k=1.0 +units=m +nadgrids=@null +wktext  +no_defs",
			p.EarthRadius_a, p.EarthRadius_b, p.Lon_0)
	default:
		panic(fmt.Errorf("Unknown proj4 projection `%v'.", p.Proj))
	}
	return s
}
