package carto

import (
	"bufio"
	"code.google.com/p/draw2d/draw2d"
	"fmt"
	"github.com/pmylund/go-cache"
	"github.com/twpayne/gogeom/geom"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math"
	"reflect"
	"time"
)

type Mapper interface {
	DrawVector(geom.T, color.NRGBA, color.NRGBA, float64, float64)
	Save()
}

type RasterMap struct {
	bounds        *geom.Bounds // geographic boundaries of map
	width, height int          // pixel dimensions of map
	dx, dy        float64
	f             io.Writer
	I             draw.Image
	GC            draw2d.GraphicContext
	MarkerFunc    func(*RasterMap, float64, float64, float64) // function that specifies the shape of the marker for points
}

type MarkerFunction func(*RasterMap, float64, float64, float64) // Function for specifying the shape of the marker for points

var (
	Circle MarkerFunction = func(m *RasterMap, x, y, markersize float64) {
		m.GC.ArcTo(x, y, markersize, markersize, 0, 2*math.Pi)
	}
	Square MarkerFunction = func(m *RasterMap, x, y, markersize float64) {
		adjMS := markersize / 1.2 // ratio to adjust the markersize
		// to make the area be the same as the circle
		m.GC.MoveTo(x-adjMS, y-adjMS)
		m.GC.LineTo(x+adjMS, y-adjMS)
		m.GC.LineTo(x+adjMS, y+adjMS)
		m.GC.LineTo(x-adjMS, y+adjMS)
		m.GC.LineTo(x-adjMS, y-adjMS)
	}
	Triangle MarkerFunction = func(m *RasterMap, x, y, markersize float64) {
		adjMS := markersize / 0.75 // ratio to adjust the markersize
		// to make the area be the same as the circle
		cosval := math.Cos(0.125 * math.Pi)
		sinval := math.Sin(0.125 * math.Pi)
		m.GC.MoveTo(x-adjMS*cosval, y+adjMS*sinval)
		m.GC.LineTo(x+adjMS*cosval, y+adjMS*sinval)
		m.GC.LineTo(x, y-adjMS)
		m.GC.LineTo(x-adjMS*cosval, y+adjMS*sinval)
	}
	Star MarkerFunction = func(m *RasterMap, x, y, markersize float64) {
		adjMS := markersize / 0.75 // ratio to adjust the markersize
		// to make the area be the same as the circle
		var alpha = (2 * math.Pi) / 10
		// works out the angle between each vertex (5 external + 5 internal = 10)
		var r_concave = adjMS / 2.25 // r_point is the radius to the external point
		for i := 11; i != 0; i-- {
			var ra float64
			if i%2 == 1 {
				ra = adjMS
			} else {
				ra = r_concave
			}
			omega := alpha * float64(i) //omega is the angle of the current point
			//cx and cy are the center point of the star.
			if i == 11 {
				m.GC.MoveTo(x+(ra*math.Sin(omega)), y+(ra*math.Cos(omega)))
			} else {
				m.GC.LineTo(x+(ra*math.Sin(omega)), y+(ra*math.Cos(omega)))
			}
		}
	}
)

func NewRasterMap(N, S, E, W float64, width int, f io.Writer) *RasterMap {
	r := new(RasterMap)
	r.f = f
	r.bounds = geom.NewBoundsPoint(geom.Point{W, S})
	r.bounds.ExtendPoint(geom.Point{E, N})
	r.width, r.height = width, int(float64(width)*(N-S)/(E-W))
	r.dx = (E - W) / float64(r.width)
	r.dy = (N - S) / float64(r.height)
	r.I = image.NewRGBA(image.Rect(0, 0, r.width, r.height))
	r.GC = draw2d.NewGraphicContext(r.I)
	r.GC.SetFillRule(draw2d.FillRuleWinding)
	r.MarkerFunc = Circle
	return r
}

// Draw a vector on a raster map when given the geometry,
// stroke and fill colors, the width of the bounding line,
// and the size of the marker (only used for point shapes).
func (r *RasterMap) DrawVector(g geom.T, strokeColor,
	fillColor color.NRGBA, linewidth, markersize float64) {
	// check bounding box
	if g == nil {
		return
	}
	gbounds := g.Bounds(nil)
	if !gbounds.Overlaps(r.bounds) {
		return
	}
	r.GC.SetStrokeColor(strokeColor)
	r.GC.SetFillColor(fillColor)
	r.GC.SetLineWidth(linewidth)
	switch g.(type) {
	case geom.Point:
		p := g.(geom.Point)
		x, y := r.coordinates(p.X, p.Y)
		r.MarkerFunc(r, x, y, markersize)
	//case geom.PointZ:
	//case geom.PointM:
	//case geom.PointZM:
	case geom.MultiPoint:
		for _, p := range g.(geom.MultiPoint).Points {
			x, y := r.coordinates(p.X, p.Y)
			r.GC.MoveTo(x, y)
			r.MarkerFunc(r, x, y, markersize)
		}
	//case geom.MultiPointZ:
	//case geom.MultiPointM:
	//case geom.MultiPointZM:
	case geom.LineString:
		l := g.(geom.LineString)
		for i, p := range l.Points {
			x, y := r.coordinates(p.X, p.Y)
			if i == 0 {
				r.GC.MoveTo(x, y)
			} else {
				r.GC.LineTo(x, y)
			}
		}
	//case geom.LineStringZ:
	//case geom.LineStringM:
	//case geom.LineStringZM:
	case geom.MultiLineString:
		l := g.(geom.MultiLineString)
		for _, ls := range l.LineStrings {
			r.DrawVector(ls, strokeColor,
				fillColor, linewidth, markersize)
		}
	//case geom.MultiLineStringZ:
	//case geom.MultiLineStringM:
	//case geom.MultiLineStringZM:
	case geom.Polygon:
		pg := g.(geom.Polygon)
		for _, ring := range pg.Rings {
			for i, p := range ring {
				x, y := r.coordinates(p.X, p.Y)
				if i == 0 {
					r.GC.MoveTo(x, y)
				} else {
					r.GC.LineTo(x, y)
				}
			}
		}
	//case geom.PolygonZ:
	//case geom.PolygonM:
	//case geom.PolygonZM:
	case geom.MultiPolygon:
		mpg := g.(geom.MultiPolygon)
		for _, pg := range mpg.Polygons {
			r.DrawVector(pg, strokeColor,
				fillColor, linewidth, markersize)
		}
	//case geom.MultiPolygonZ:
	//case geom.MultiPolygonM:
	//case geom.MultiPolygonZM:
	default:
		panic(&UnsupportedGeometryError{reflect.TypeOf(g)})
	}
	r.GC.FillStroke()
}

func (r *RasterMap) Save() {
	b := bufio.NewWriter(r.f)
	err := png.Encode(b, r.I)
	if err != nil {
		panic(err)
	}
	err = b.Flush()
	if err != nil {
		panic(err)
	}
}

// transform geographic coordinates to raster map coordinates
func (r *RasterMap) coordinates(X, Y float64) (
	x, y float64) {
	x = (X - r.bounds.Min.X) / r.dx
	y = float64(r.height) - 1. - (Y-r.bounds.Min.Y)/r.dy
	return
}

type MapData struct {
	Cmap      *ColorMap
	Shapes    []geom.T
	Data      []float64
	tileCache *cache.Cache
	DrawEdges bool
	EdgeWidth float64
}

func NewMapData(numShapes int, Type ColorMapType) *MapData {
	m := new(MapData)
	m.Cmap = NewColorMap(Type)
	m.Shapes = make([]geom.T, numShapes)
	m.Data = make([]float64, numShapes)
	m.tileCache = cache.New(1*time.Hour, 10*time.Minute)
	m.EdgeWidth = 0.5
	return m
}

func (m *MapData) WriteGoogleMapTile(w io.Writer, zoom, x, y int) error {
	// Check if image is already in the cache.
	cacheKey := fmt.Sprintf("%v_%v_%v", zoom, x, y)
	if img, found := m.tileCache.Get(cacheKey); found {
		err := png.Encode(w, img.(image.Image))
		if err != nil {
			return err
		}
		return nil
	}
	//strokeColor := color.NRGBA{0, 0, 0, 255}
	N, S, E, W := GetGoogleTileBounds(zoom, x, y)
	maptile := NewRasterMap(N, S, E, W, 256, w)

	var strokeColor color.NRGBA
	for i, shp := range m.Shapes {
		fillColor := m.Cmap.GetColor(m.Data[i])
		if m.DrawEdges {
			strokeColor = color.NRGBA{0, 0, 0, 255}
		} else {
			strokeColor = fillColor
		}
		// use the fill color for both the fill and the stroke
		// to avoid unsightly gaps between shapes.
		maptile.DrawVector(shp, strokeColor, fillColor, m.EdgeWidth, 0)
	}
	err := png.Encode(w, maptile.I)
	if err != nil {
		return err
	}
	m.tileCache.Set(cacheKey, maptile.I, 0)
	return nil
}

func GetGoogleTileBounds(zoom, x, y int) (N, S, E, W float64) {
	const originShift = math.Pi * 6378137. // for mercator projection
	// get boundaries in lat/lon
	n := math.Pow(2, float64(zoom))
	W_lon := float64(x)/n*360.0 - 180.0
	E_lon := float64(x+1)/n*360.0 - 180.0
	N_rad := math.Atan(math.Sinh(math.Pi * (1 - 2*float64(y)/n)))
	N_lat := N_rad * 180.0 / math.Pi
	S_rad := math.Atan(math.Sinh(math.Pi * (1 - 2*float64(y+1)/n)))
	S_lat := S_rad * 180.0 / math.Pi
	// convert to Mercator meters
	W = W_lon * originShift / 180.0
	E = E_lon * originShift / 180.0
	N = math.Log(math.Tan((90+N_lat)*math.Pi/360.0)) /
		(math.Pi / 180.0) * originShift / 180.0
	S = math.Log(math.Tan((90+S_lat)*math.Pi/360.0)) /
		(math.Pi / 180.0) * originShift / 180.0
	return
}

// convert from long/lat to google mercator (or EPSG:4326 to EPSG:900913)
func Degrees2meters(lon, lat float64) (x, y float64) {
	x = lon * 20037508.34 / 180.
	y = math.Log(math.Tan((90.+lat)*math.Pi/360.)) / (math.Pi / 180.)
	y *= 20037508.34 / 180.
	return x, y
}

type UnsupportedGeometryError struct {
	Type reflect.Type
}

func (e UnsupportedGeometryError) Error() string {
	return "Unsupported geometry type: " + e.Type.String()
}

// Convenience function for making a simple map.
func DrawShapes(f io.Writer, strokeColor, fillColor []color.NRGBA,
	linewidth, markersize float64, shapes ...geom.T) {
	bounds := geom.NewBounds()
	for _, s := range shapes {
		if s != nil {
			bounds = s.Bounds(bounds)
		}
	}
	m := NewRasterMap(bounds.Max.Y, bounds.Min.Y,
		bounds.Max.X, bounds.Min.X, 500, f)
	for i, s := range shapes {
		m.DrawVector(s, strokeColor[i], fillColor[i], linewidth, markersize)
	}
	m.Save()
}
