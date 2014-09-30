package carto

import (
	"fmt"
	"github.com/ajstarks/svgo"
	"image/color"
	"io"
	"math"
	"os/exec"
	"sort"
	"strings"
)

const (
	pointsPerIn = 72. // postscript points per inch
)

var (
	Font = "" // Default font to use
)

type Colorlist struct {
	Val, R, G, B        []float64
	HighLimit, LowLimit color.NRGBA
}

var (
	// optimized olors from http://www.cs.unm.edu/~kmorel/documents/ColorMaps/index.html
	// Originally the 255 values were 221's
	Optimized Colorlist = Colorlist{
		[]float64{-1., -0.9375, -0.875, -0.8125,
			-0.75, -0.6875, -0.625, -0.5625, -0.5, -0.4375, -0.375,
			-0.3125, -0.25, -0.1875, -0.125, -0.0625, 0., 0.0625, 0.125,
			0.1875, 0.25, 0.3125, 0.375, 0.4375, 0.5, 0.5625, 0.625,
			0.6875, 0.75, 0.8125, 0.875, 0.9375, 1.},
		[]float64{59., 68, 77, 87, 98, 108, 119, 130, 141, 152,
			163, 174, 184, 194, 204, 213, 255, 229, 236, 241, 245, 247,
			247, 247, 244, 241, 236, 229, 222, 213, 203, 192, 180},
		[]float64{76., 90, 104, 117, 130, 142, 154, 165, 176,
			185, 194, 201, 208, 213, 217, 219, 255, 216, 211, 204, 196, 187,
			177, 166, 154, 141, 127, 112, 96, 80, 62, 40, 4},
		[]float64{192., 204, 215, 225, 234, 241, 247, 251, 254,
			255, 255, 253, 249, 244, 238, 230, 255, 209, 197, 185, 173, 160,
			148, 135, 123, 111, 99, 88, 77, 66, 56, 47, 38},
		color.NRGBA{70., 6, 16, 255},
		color.NRGBA{27., 34, 85, 255}}
	OptimizedGrey Colorlist = Colorlist{
		[]float64{-1., -0.9375, -0.875, -0.8125,
			-0.75, -0.6875, -0.625, -0.5625, -0.5, -0.4375, -0.375,
			-0.3125, -0.25, -0.1875, -0.125, -0.0625, 0., 0.0625, 0.125,
			0.1875, 0.25, 0.3125, 0.375, 0.4375, 0.5, 0.5625, 0.625,
			0.6875, 0.75, 0.8125, 0.875, 0.9375, 1.},
		[]float64{59., 68, 77, 87, 98, 108, 119, 130, 141, 152, 163, 174,
			184, 194, 204, 213, 221, 229, 236, 241, 245, 247, 247, 247, 244, 241,
			236, 229, 222, 213, 203, 192, 180},
		[]float64{76., 90, 104, 117, 130, 142, 154, 165, 176, 185, 194,
			201, 208, 213, 217, 219, 221, 216, 211, 204, 196, 187, 177, 166,
			154, 141, 127, 112, 96, 80, 62, 40, 4},
		[]float64{192., 204, 215, 225, 234, 241, 247, 251, 254, 255,
			255, 253, 249, 244, 238, 230, 221, 209, 197, 185, 173, 160, 148,
			135, 123, 111, 99, 88, 77, 66, 56, 47, 38},
		color.NRGBA{70., 6, 16, 255},
		color.NRGBA{27., 34, 85, 255}}
	Jet Colorlist = Colorlist{
		[]float64{-1, -0.866666666666667, -0.733333333333333, -0.6,
			-0.466666666666667, -0.333333333333333, -0.2, -0.0666666666666668,
			0.0666666666666665, 0.2, 0.333333333333333, 0.466666666666666, 0.6,
			0.733333333333333, 0.866666666666666, 1},
		[]float64{0, 0, 0, 0, 0, 0, 66, 132, 189, 255, 255, 255,
			255, 255, 189, 132},
		[]float64{0, 0, 66, 132, 189, 255, 255, 255, 255, 255, 189,
			132, 66, 0, 0, 0},
		[]float64{189, 255, 255, 255, 255, 255, 189, 132, 66, 0, 0,
			0, 0, 0, 0, 0},
		color.NRGBA{249., 15, 225, 255},
		color.NRGBA{154., 0, 171, 255}}
	JetPosOnly Colorlist = Colorlist{
		[]float64{-1, 0, 0.0666666666666667, 0.133333333333333, 0.2,
			0.266666666666667, 0.333333333333333, 0.4, 0.466666666666667,
			0.533333333333333, 0.6, 0.666666666666667, 0.733333333333333, 0.8,
			0.866666666666667, 0.933333333333333, 1},
		[]float64{0, 0, 0, 0, 0, 0, 0, 66, 132, 189, 255, 255, 255,
			255, 255, 189, 132},
		[]float64{0, 0, 0, 66, 132, 189, 255, 255, 255, 255, 255, 189,
			132, 66, 0, 0, 0},
		[]float64{189, 189, 255, 255, 255, 255, 255, 189, 132, 66, 0,
			0, 0, 0, 0, 0, 0},
		color.NRGBA{249., 15, 225, 255},
		color.NRGBA{154., 0, 171, 255}}
)

type ColorMapType int

const (
	Linear    ColorMapType = iota // Linear color gradient
	LinCutoff                     // linear with a discontinuity at a percentile
	// specified by "CutPercentile"
)

type ColorMap struct {
	cutoff            float64
	maxval            float64
	minval            float64
	cutptlist         []float64
	Type              ColorMapType
	CutPercentile     float64 // Percentile at which discontinuity occurs for "LinCutoff" type.
	NumDivisions      int     // "Number of tick marks on legend.
	rulestring        string
	colorstops        []float64
	stopcolors        []color.NRGBA
	LegendWidth       float64 // width of legend in inches
	LegendHeight      float64 // height of legend in inches
	LineWidth         float64 // width of lines in legend in points
	FontSize          float64 // font size in points.
	Font              string  // Name of the font to use in legend
	FontColor         string
	EdgeColor         string // Color for legend outline
	BackgroundColor   string
	BackgroundOpacity float64
	negativeOutlier   bool
	positiveOutlier   bool
	ColorScheme       Colorlist
}

// Initialize new color map.
func NewColorMap(Type ColorMapType) (c *ColorMap) {
	c = new(ColorMap)
	c.cutptlist = make([]float64, 0)
	c.Type = Type
	c.CutPercentile = 99.
	c.ColorScheme = Optimized
	c.NumDivisions = 9
	c.colorstops = make([]float64, 0)
	c.stopcolors = make([]color.NRGBA, 0)
	c.LegendWidth = 3.70
	c.LegendHeight = c.LegendWidth * 0.1067
	c.LineWidth = 0.5
	c.FontSize = 7.
	c.negativeOutlier = false
	c.positiveOutlier = false
	c.Font = Font
	c.FontColor = "black"
	c.EdgeColor = "black"
	c.BackgroundColor = "white"
	c.BackgroundOpacity = 0.
	return
}

func (c *ColorMap) AddArray(data []float64) {
	var max, min float64
	for i := 0; i < len(data); i++ {
		if data[i] > max {
			max = data[i]
		}
		if data[i] < min {
			min = data[i]
		}
	}
	if max*1.00001 > c.maxval {
		c.maxval = max * 1.00001
	}
	if min*1.00001 < c.minval {
		c.minval = min * 1.00001
	}
	if c.Type == LinCutoff {
		tmpAbs := make([]float64, len(data))
		for i := 0; i < len(data); i++ {
			tmpAbs[i] = math.Abs(data[i])
		}
		sort.Float64s(tmpAbs)
		cutpt := tmpAbs[roundInt(c.CutPercentile/100.*
			float64(len(data)))-1]
		c.cutptlist = append(c.cutptlist, cutpt)
	}
}

func (c *ColorMap) AddGeoJSON(g *GeoJSON, propertyName string) {
	vals := make([]float64, len(g.Features))
	for i, f := range g.Features {
		vals[i] = f.Properties[propertyName]
	}
	c.AddArray(vals)
}

func (c *ColorMap) AddMap(data map[string]float64) {
	vals := make([]float64, len(data))
	i := 0
	for _, val := range data {
		vals[i] = val
		i++
	}
	c.AddArray(vals)
}

func (c *ColorMap) AddArrayServer(datachan chan []float64,
	finished chan int) {
	var data []float64
	for {
		data = <-datachan
		if data == nil {
			break
		}
		c.AddArray(data)
	}
	finished <- 0
}

// get color for input value. Must run c.Set() first.
func (cm *ColorMap) GetColor(v float64) color.NRGBA {
	var R, G, B uint8
	c := cm.stopcolors
	cv := cm.colorstops
	for i := 1; i < len(cv); i++ {
		if math.Abs(v-cv[i])/math.Abs(cv[i]) < 0.0001 {
			return c[i]
		} else if cv[i] > v {
			valFrac := (v - cv[i-1]) / (cv[i] - cv[i-1])
			R = round(float64(c[i-1].R) + (float64(c[i].R)-float64(c[i-1].R))*
				valFrac)
			G = round(float64(c[i-1].G) + (float64(c[i].G)-float64(c[i-1].G))*
				valFrac)
			B = round(float64(c[i-1].B) + (float64(c[i].B)-float64(c[i-1].B))*
				valFrac)
			return color.NRGBA{R, G, B, 255}
		}
	}
	if math.IsNaN(v) || math.IsInf(v, 0) {
		fmt.Printf("Problem interpolating: %v value\n", v)
		return color.NRGBA{255, 0, 174, 255}
	}
	if len(cm.colorstops) == 0 {
		return color.NRGBA{255, 255, 255, 255}
	}
	fmt.Println("x=", v, "xArray=", cm.colorstops)
	panic("Problem interpolating: x value is larger than xArray")
}

// Given an array of x values and an array of y values, find the y value at a // given x using linear interpolation. xArray must be monotonically increasing.
func (cl *Colorlist) interpolate(v float64) color.NRGBA {
	var R, G, B uint8
	for i, val := range cl.Val {
		if math.Abs(v-val)/math.Abs(val) < 0.0001 {
			R = round(cl.R[i])
			G = round(cl.G[i])
			B = round(cl.B[i])
			return color.NRGBA{R, G, B, 255}
		} else if val > v {
			R = round(cl.R[i-1] + (cl.R[i]-cl.R[i-1])*
				(v-cl.Val[i-1])/(cl.Val[i]-cl.Val[i-1]))
			G = round(cl.G[i-1] + (cl.G[i]-cl.G[i-1])*
				(v-cl.Val[i-1])/(cl.Val[i]-cl.Val[i-1]))
			B = round(cl.B[i-1] + (cl.B[i]-cl.B[i-1])*
				(v-cl.Val[i-1])/(cl.Val[i]-cl.Val[i-1]))
			return color.NRGBA{R, G, B, 255}
		}
	}
	fmt.Println("x=", v, "xArray=", cl.Val)
	panic("Problem interpolating: x value is larger than xArray")
}

// round float to an integer
func round(x float64) uint8 {
	return uint8(x + 0.5)
}

// round float to an integer
func roundInt(x float64) int {
	return int(x + 0.5)
}

func newsvgcolor(cIn color.NRGBA) (cOut svg.Offcolor) {
	cOut.Offset = uint8(0)
	cOut.Opacity = float64(cIn.A) / 255
	cOut.Color = fmt.Sprintf("rgb(%v,%v,%v)", cIn.R, cIn.G, cIn.B)
	return
}

// Figure out rules for color map
func (c *ColorMap) Set() {

	var linmin, linmax, absmax float64
	cutpt := average(c.cutptlist)
	if c.minval*-1 > c.maxval {
		absmax = c.minval * -1
	} else {
		absmax = c.maxval
	}

	if c.Type == LinCutoff && cutpt < absmax && cutpt != 0 {
		linmin = cutpt * -1
		linmax = cutpt
	} else {
		linmin = absmax * -1
		linmax = absmax
	}

	c.colorstops = make([]float64, 0)
	c.stopcolors = make([]color.NRGBA, 0)

	if absmax == 0. {
		return
	}

	if c.Type == LinCutoff && cutpt*-1 > c.minval && cutpt != 0 {
		c.colorstops = append(c.colorstops, absmax*-1)
		c.stopcolors = append(c.stopcolors, c.ColorScheme.LowLimit)
		c.negativeOutlier = true
	}
	if (c.minval-linmin)/(c.minval+linmin) > 0.001 {
		c.colorstops = append(c.colorstops, c.minval)
		c.stopcolors = append(c.stopcolors, c.ColorScheme.interpolate(-1.))
	}

	interval := (linmax - linmin) / float64(c.NumDivisions+1)
	for val := linmin; (val-linmax)/linmax < 0.001; val += interval {
		if (val-c.minval)/linmax > -0.0001 &&
			(val-c.maxval)/linmax < 0.0001 {
			c.colorstops = append(c.colorstops, val)
			c.stopcolors = append(c.stopcolors, c.ColorScheme.interpolate(val/linmax))
		}
	}
	if (c.maxval-linmax)/(c.maxval+linmax) < 0.001 {
		c.colorstops = append(c.colorstops, c.maxval)
		c.stopcolors = append(c.stopcolors, c.ColorScheme.interpolate(1.))
	}

	if c.Type == LinCutoff && cutpt < c.maxval && cutpt != 0 {
		c.colorstops = append(c.colorstops, absmax)
		c.stopcolors = append(c.stopcolors, c.ColorScheme.HighLimit)
		c.positiveOutlier = true
	}
}

func (c *ColorMap) Legend(w io.Writer, label string) (err error) {
	const dpi = 300.
	const topPad = 0.    // points
	const bottomPad = 2. // points
	const unitsPad = 2.5 // pad between units and bar
	const labelPad = 2.  // pad between bar and label
	const wPad = 10.     // points
	pts2px := dpi / pointsPerIn
	fontHeight := c.FontSize * pts2px
	strokeWidth := round(c.LineWidth * pts2px)
	var fontConfig string
	if c.Font == "" {
		fontConfig = fmt.Sprintf("text-anchor:middle;font-size:%dpx;fill:%v",
			int(fontHeight), c.FontColor)
	} else {
		fontConfig = fmt.Sprintf(
			"text-anchor:middle;font-size:%dpx;fill:%v;font-family:%v",
			int(fontHeight), c.FontColor, c.Font)
	}
	width := roundInt(c.LegendWidth * dpi)
	height := roundInt(c.LegendHeight * dpi)
	barWstart := roundInt(wPad * pts2px)
	barWidth := roundInt(float64(width) - 2*wPad*pts2px)
	labelX := roundInt(float64(width) * 0.5)
	labelY := roundInt(topPad*pts2px + fontHeight)
	unitsYunder := roundInt(float64(height) - bottomPad*pts2px)
	unitsYover := labelY
	barHstart := roundInt((topPad+labelPad)*pts2px + fontHeight)
	barHeight := roundInt(float64(height) - 2*fontHeight -
		(topPad+labelPad+bottomPad)*pts2px)

	legendcolors := make([]svg.Offcolor, len(c.stopcolors))
	numstops := len(c.stopcolors) - 1
	ticklocs := make([]float64, numstops+1)
	loc := 0.
	for i, _ := range c.colorstops {
		ticklocs[i] = loc
		if c.negativeOutlier && i == 0 ||
			c.positiveOutlier && i == numstops-1 {
			loc += 0.33
		} else {
			loc += 1.
		}
	}
	for i, val := range ticklocs {
		ticklocs[i] = val/ticklocs[numstops]*float64(barWidth) + float64(barWstart)
	}
	for i, val := range ticklocs {
		legendcolors[i] = newsvgcolor(c.stopcolors[i])
		legendcolors[i].Offset = uint8(round(val /
			ticklocs[numstops] * 100))
	}

	g := svg.New(w)
	g.Start(width, height)
	g.Rect(0, 0, width, height, fmt.Sprintf("fill:%v;fill-opacity:%v",
		c.BackgroundColor, c.BackgroundOpacity))
	g.Def()
	g.LinearGradient("cmap", 0, 0, 100, 0, legendcolors)
	g.DefEnd()
	g.Rect(barWstart, barHstart, barWidth, barHeight,
		fmt.Sprintf("fill:url(#cmap);stroke:%v;stroke-width:%d", c.EdgeColor,
			strokeWidth))
	for i, tickloc := range ticklocs {
		val := c.colorstops[i]
		var valStr string
		if math.Abs(val) < absmax(c.maxval, c.minval)*1.e-10 {
			valStr = "0"
		} else {
			valStr = strings.Replace(strings.Replace(fmt.Sprintf("%3.2g", val),
				"e+0", "e", -1), "e-0", "e-", -1)
		}
		if c.negativeOutlier && i == 0 ||
			c.positiveOutlier && i == numstops {
			g.Text(roundInt(tickloc), unitsYover, valStr, fontConfig)
		} else {
			g.Text(roundInt(tickloc), unitsYunder, valStr, fontConfig)
		}
	}
	g.Text(labelX, labelY, label, fontConfig)
	g.End()
	return
}

func max(a, b float64) float64 {
	if a > b {
		return a
	} else {
		return b
	}
}
func absmax(a, b float64) float64 {
	absa := math.Abs(a)
	absb := math.Abs(b)
	if absa > absb {
		return absa
	} else {
		return absb
	}
}

func ConvertSVGToPNG(filename string) {
	cmd := exec.Command("convert", "-density", "600",
		filename, strings.Replace(filename, "svg", "png", -1))
	out, err := cmd.CombinedOutput()
	output := fmt.Sprintf("%s", out)
	if err != nil {
		panic(fmt.Errorf(output))
	}
}

func average(a []float64) (avg float64) {
	for _, val := range a {
		avg += val
	}
	avg /= float64(len(a))
	return
}
