// Go language map drawing library
package carto

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"reflect"

	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"github.com/gonum/plot/vg/vgimg"
	"github.com/twpayne/gogeom/geom"
)

// Canvas is a canvas for drawing maps.
type Canvas struct {
	draw.Canvas
	bounds *geom.Bounds // geographic boundaries of map
	scale  float64
}

//type MarkerFunction func(*Canvas, float64, float64, float64) // Function for specifying the shape of the marker for points

//var (
//	Circle MarkerFunction = func(m *Canvas, x, y, markersize vg.Length) {
//		m.GC.ArcTo(x, y, markersize, markersize, 0, 2*math.Pi)
//	}
//	Square MarkerFunction = func(m *Canvas, x, y, markersize vg.Length) {
//		adjMS := markersize / 1.2 // ratio to adjust the markersize
//		// to make the area be the same as the circle
//		m.GC.MoveTo(x-adjMS, y-adjMS)
//		m.GC.LineTo(x+adjMS, y-adjMS)
//		m.GC.LineTo(x+adjMS, y+adjMS)
//		m.GC.LineTo(x-adjMS, y+adjMS)
//		m.GC.LineTo(x-adjMS, y-adjMS)
//	}
//	Triangle MarkerFunction = func(m *Canvas, x, y, markersize vg.Length) {
//		adjMS := markersize / 0.75 // ratio to adjust the markersize
//		// to make the area be the same as the circle
//		cosval := math.Cos(0.125 * math.Pi)
//		sinval := math.Sin(0.125 * math.Pi)
//		m.GC.MoveTo(x-adjMS*cosval, y+adjMS*sinval)
//		m.GC.LineTo(x+adjMS*cosval, y+adjMS*sinval)
//		m.GC.LineTo(x, y-adjMS)
//		m.GC.LineTo(x-adjMS*cosval, y+adjMS*sinval)
//	}
//	Star MarkerFunction = func(m *Canvas, x, y, markersize vg.Length) {
//		adjMS := markersize / 0.75 // ratio to adjust the markersize
//		// to make the area be the same as the circle
//		var alpha = (2 * math.Pi) / 10
//		// works out the angle between each vertex (5 external + 5 internal = 10)
//		var r_concave = adjMS / 2.25 // r_point is the radius to the external point
//		for i := 11; i != 0; i-- {
//			var ra float64
//			if i%2 == 1 {
//				ra = adjMS
//			} else {
//				ra = r_concave
//			}
//			omega := alpha * float64(i) //omega is the angle of the current point
//			//cx and cy are the center point of the star.
//			if i == 11 {
//				m.GC.MoveTo(x+(ra*math.Sin(omega)), y+(ra*math.Cos(omega)))
//			} else {
//				m.GC.LineTo(x+(ra*math.Sin(omega)), y+(ra*math.Cos(omega)))
//			}
//		}
//	}
//)

func NewCanvas(N, S, E, W float64, c *draw.Canvas) *Canvas {
	m := &Canvas{
		*c,
		geom.NewBoundsPoint(geom.Point{W, S}).
			ExtendPoint(geom.Point{E, N}),
		min(float64(c.Max.X-c.Min.X)/(E-W),
			float64(c.Max.Y-c.Min.Y)/(N-S)),
	}
	return m
}

type RasterMap struct {
	Canvas
	I *image.RGBA
}

func NewRasterMap(N, S, E, W float64, width int) *RasterMap {
	const mapWidth = 3.5 // inches, for dpi conversion
	height := int(float64(width) * (N - S) / (E - W))
	I := image.NewRGBA(image.Rect(0, 0, width, height))
	m := &RasterMap{
		Canvas: Canvas{
			Canvas: draw.New(vgimg.NewImage(I, int(float64(width)/mapWidth))),
			bounds: geom.NewBoundsPoint(geom.Point{W, S}).
				ExtendPoint(geom.Point{E, N}),
		},
		I: I,
	}
	m.scale = min(float64(m.Max.X-m.Min.X)/(E-W),
		float64(m.Max.Y-m.Min.Y)/(N-S))
	return m
}

func (r *RasterMap) WriteTo(f io.Writer) error {
	return png.Encode(f, r.I)
}

// Draw a vector on a raster map when given the geometry,
// stroke and fill colors, the width of the bounding line,
// and the markerGlyph, which specifies the shape of the marker
// (only used for point shapes).
func (m *Canvas) DrawVector(g geom.T, fillColor color.NRGBA,
	lineStyle draw.LineStyle, markerGlyph draw.GlyphStyle) error {
	// check bounding box
	if g == nil {
		return nil
	}
	gbounds := g.Bounds(nil)
	if !gbounds.Overlaps(m.bounds) {
		return nil
	}
	m.SetLineStyle(lineStyle)
	switch g.(type) {
	case geom.Point:
		pTemp := g.(geom.Point)
		p := m.coordinates(pTemp)
		m.DrawGlyph(markerGlyph, p)
	//case geom.PointZ:
	//case geom.PointM:
	//case geom.PointZM:
	case geom.MultiPoint:
		for _, pTemp := range g.(geom.MultiPoint).Points {
			p := m.coordinates(pTemp)
			m.DrawGlyph(markerGlyph, p)
		}
	//case geom.MultiPointZ:
	//case geom.MultiPointM:
	//case geom.MultiPointZM:
	case geom.LineString:
		l := g.(geom.LineString)
		var path vg.Path
		for i, pTemp := range l.Points {
			p := m.coordinates(pTemp)
			if i == 0 {
				path.Move(p.X, p.Y)
			} else {
				path.Line(p.X, p.Y)
			}
		}
		m.Stroke(path)
	//case geom.LineStringZ:
	//case geom.LineStringM:
	//case geom.LineStringZM:
	case geom.MultiLineString:
		l := g.(geom.MultiLineString)
		for _, ls := range l.LineStrings {
			m.DrawVector(ls, fillColor, lineStyle, markerGlyph)
		}
	//case geom.MultiLineStringZ:
	//case geom.MultiLineStringM:
	//case geom.MultiLineStringZM:
	case geom.Polygon:
		pg := g.(geom.Polygon)
		for _, ring := range pg.Rings {
			var path vg.Path
			for i, pTemp := range ring {
				p := m.coordinates(pTemp)
				if i == 0 {
					path.Move(p.X, p.Y)
				} else {
					path.Line(p.X, p.Y)
				}
			}
			path.Close()
			m.Push() // save stroke color
			m.SetColor(fillColor)
			m.Fill(path)
			m.Pop() // retrieve stroke color
			m.Stroke(path)
		}
	//case geom.PolygonZ:
	//case geom.PolygonM:
	//case geom.PolygonZM:
	case geom.MultiPolygon:
		mpg := g.(geom.MultiPolygon)
		for _, pg := range mpg.Polygons {
			m.DrawVector(pg, fillColor, lineStyle, markerGlyph)
		}
	//case geom.MultiPolygonZ:
	//case geom.MultiPolygonM:
	//case geom.MultiPolygonZM:
	default:
		return &UnsupportedGeometryError{reflect.TypeOf(g)}
	}
	return nil
}

// transform geographic coordinates to raster map coordinates
func (m *Canvas) coordinates(pIn geom.Point) draw.Point {
	var pOut draw.Point
	pOut.X = m.Min.X + vg.Length(m.scale*(pIn.X-m.bounds.Min.X))
	pOut.Y = m.Min.Y + vg.Length(m.scale*(pIn.Y-m.bounds.Min.Y))
	return pOut
}

// Make a new raster map from raster data.
// It is assumed that the outer axis in the data is the Y-axis
// (north-south) and the inner axis is the X-axis (west-east)
// (i.e., len(data)==nx*ny &&  val[j,i] = data[j*nx+i]).
func NewCanvasFromRaster(S, W, dy, dx float64, ny, nx int,
	data []float64, cmap *ColorMap,
	flipVertical, flipHorizontal bool) *RasterMap {
	N := S + float64(ny)*dy
	E := W + float64(nx)*dx
	r := NewRasterMap(N, S, E, W, nx)
	if !flipVertical && !flipHorizontal {
		for i := 0; i < nx; i++ {
			for j := 0; j < ny; j++ {
				val := data[j*nx+i]
				r.I.Set(i, j, cmap.GetColor(val))
			}
		}
	} else if flipVertical && !flipHorizontal {
		for i := 0; i < nx; i++ {
			for j := 0; j < ny; j++ {
				val := data[j*nx+i]
				r.I.Set(i, ny-1-j, cmap.GetColor(val))
			}
		}
	} else if !flipVertical && flipHorizontal {
		for i := 0; i < nx; i++ {
			for j := 0; j < ny; j++ {
				val := data[j*nx+i]
				r.I.Set(nx-1-i, j, cmap.GetColor(val))
			}
		}
	} else if flipVertical && flipHorizontal {
		for i := 0; i < nx; i++ {
			for j := 0; j < ny; j++ {
				val := data[j*nx+i]
				r.I.Set(nx-1-i, ny-1-j, cmap.GetColor(val))
			}
		}
	}
	return r
}

type MapData struct {
	Cmap      *ColorMap
	Shapes    []geom.T
	Data      []float64
	DrawEdges bool
	draw.LineStyle
}

func NewMapData(numShapes int, Type ColorMapType) *MapData {
	m := new(MapData)
	m.Cmap = NewColorMap(Type)
	m.Shapes = make([]geom.T, numShapes)
	m.Data = make([]float64, numShapes)
	m.LineStyle = draw.LineStyle{Width: 1. * vg.Millimeter}
	return m
}

func (m *MapData) WriteGoogleMapTile(w io.Writer, zoom, x, y int) error {
	//strokeColor := color.NRGBA{0, 0, 0, 255}
	N, S, E, W := getGoogleTileBounds(zoom, x, y)
	maptile := NewRasterMap(N, S, E, W, 256)

	var markerGlyph draw.GlyphStyle
	xLen := (maptile.Max.X - maptile.Min.X)
	yLen := (maptile.Max.Y - maptile.Min.Y)
	markerGlyph.Radius = vg.Length(0.01 * math.Sqrt(float64(xLen*xLen+yLen*yLen)))
	markerGlyph.Shape = draw.RingGlyph{}

	for i, shp := range m.Shapes {
		fillColor := m.Cmap.GetColor(m.Data[i])
		if m.DrawEdges {
			markerGlyph.Color = color.NRGBA{0, 0, 0, 255}
			m.LineStyle.Color = color.NRGBA{0, 0, 0, 255}
		} else {
			markerGlyph.Color = fillColor
			m.LineStyle.Color = fillColor
		}
		// use the fill color for both the fill and the stroke
		// to avoid unsightly gaps between shapes.
		maptile.DrawVector(shp, fillColor, m.LineStyle, markerGlyph)
	}
	return maptile.WriteTo(w)
}

func getGoogleTileBounds(zoom, x, y int) (N, S, E, W float64) {
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
	linewidth, markersize vg.Length, shapes ...geom.T) error {
	bounds := geom.NewBounds()
	for _, s := range shapes {
		if s != nil {
			bounds = s.Bounds(bounds)
		}
	}
	m := NewRasterMap(bounds.Max.Y, bounds.Min.Y,
		bounds.Max.X, bounds.Min.X, 500)
	var markerGlyph draw.GlyphStyle
	xLen := (m.Max.X - m.Min.X)
	yLen := (m.Max.Y - m.Min.Y)
	markerGlyph.Radius = 0.01 * vg.Length(math.Sqrt(float64(xLen*xLen+yLen*yLen)))
	markerGlyph.Shape = draw.RingGlyph{}
	lineStyle := draw.LineStyle{Width: 1. * vg.Millimeter}
	for i, s := range shapes {
		markerGlyph.Color = strokeColor[i]
		lineStyle.Color = strokeColor[i]
		err := m.DrawVector(s, fillColor[i], lineStyle, markerGlyph)
		if err != nil {
			return err
		}
	}
	return m.WriteTo(f)
}
