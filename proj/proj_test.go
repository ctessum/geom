package proj

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
	"testing"
)

func TestAddDef(t *testing.T) {
	err := addDef("EPSG:102018", "+proj=gnom +lat_0=90 +lon_0=0 +x_0=6300000 +y_0=6300000 +ellps=WGS84 +datum=WGS84 +units=m +no_defs")
	if err != nil {
		t.Error(err)
	}
	err = addDef("testmerc", "+proj=merc +lon_0=5.937 +lat_ts=45.027 +ellps=sphere +datum=none")
	if err != nil {
		t.Error(err)
	}
	err = addDef("testmerc2", "+proj=merc +a=6378137 +b=6378137 +lat_ts=0.0 +lon_0=0.0 +x_0=0.0 +y_0=0 +units=m +k=1.0 +nadgrids=@null +no_defs")
	if err != nil {
		t.Error(err)
	}
	err = addDef("esriOnline", `PROJCS["WGS_1984_Web_Mercator_Auxiliary_Sphere",GEOGCS["GCS_WGS_1984",DATUM["D_WGS_1984",SPHEROID["WGS_1984",6378137.0,298.257223563]],PRIMEM["Greenwich",0.0],UNIT["Degree",0.0174532925199433]],PROJECTION["Mercator_Auxiliary_Sphere"],PARAMETER["False_Easting",0.0],PARAMETER["False_Northing",0.0],PARAMETER["Central_Meridian",0.0],PARAMETER["Standard_Parallel_1",0.0],PARAMETER["Auxiliary_Sphere_Type",0.0],UNIT["Meter",1.0]]`)
	if err != nil {
		t.Error(err)
	}
}

func TestUnits(t *testing.T) {
	if defs["testmerc2"].Units != "m" {
		t.Error("should parse units")
	}
}

func closeTo(t *testing.T, a, b, tol float64, prefix string) {
	if 2*math.Abs(a-b)/math.Abs(a+b) > tol {
		t.Errorf("%s: value should be %f but is %f", prefix, b, a)
	}
}

func TestProj2Proj(t *testing.T) {
	// transforming from one projection to another
	sweref99tm, err := Parse("+proj=utm +zone=33 +ellps=GRS80 +towgs84=0,0,0,0,0,0,0 +units=m +no_defs")
	if err != nil {
		t.Error(err)
	}
	rt90, err := Parse("+lon_0=15.808277777799999 +lat_0=0.0 +k=1.0 +x_0=1500000.0 +y_0=0.0 +proj=tmerc +ellps=bessel +units=m +towgs84=414.1,41.3,603.1,-0.855,2.141,-7.023,0 +no_defs")
	if err != nil {
		t.Error(err)
	}
	trans, err := sweref99tm.NewTransformFunc(rt90)
	if err != nil {
		t.Error(err)
	}
	rsltx, rslty, err := trans(319180, 6399862)
	if err != nil {
		t.Error(err)
	}
	closeTo(t, rsltx, 1271137.927154, 0.000001, "x")
	closeTo(t, rslty, 6404230.291456, 0.000001, "y")
}

func TestProj4(t *testing.T) {
	type testPoint struct {
		code   string
		xy, ll []float64
		acc    struct {
			xy, ll float64
		}
	}
	var testPoints []testPoint
	f, err := os.Open("testData.json")
	if err != nil {
		t.Fatal(err)
	}
	d := json.NewDecoder(f)
	err = d.Decode(&testPoints)
	if err != nil {
		t.Fatal(err)
	}
	for _, testPoint := range testPoints {
		wgs84, err := Parse("+longlat")
		if err != nil {
			t.Fatal(err)
		}
		xyAcc := 2.
		llAcc := 6.
		if testPoint.acc.xy != 0 {
			xyAcc = testPoint.acc.xy
		}
		if testPoint.acc.ll != 0 {
			llAcc = testPoint.acc.ll
		}
		xyEPSLN := math.Pow(10, -1*xyAcc)
		llEPSLN := math.Pow(10, -1*llAcc)
		proj, err := Parse(testPoint.code)
		if err != nil {
			t.Errorf("%s: %s", testPoint.code, err)
		}
		trans, err := wgs84.NewTransformFunc(proj)
		if err != nil {
			t.Errorf("%s: %s", testPoint.code, err)
		}
		x, y, err := trans(testPoint.ll[0], testPoint.ll[1])
		if err != nil {
			t.Errorf("%s: %s", testPoint.code, err)
		}
		closeTo(t, x, testPoint.xy[0], xyEPSLN, fmt.Sprintf("%s fwd x", testPoint.code))
		closeTo(t, y, testPoint.xy[1], xyEPSLN, fmt.Sprintf("%s fwd y", testPoint.code))
		trans, err = proj.NewTransformFunc(wgs84)
		if err != nil {
			t.Errorf("%s: %s", testPoint.code, err)
		}
		lon, lat, err := trans(testPoint.xy[0], testPoint.xy[1])
		if err != nil {
			t.Errorf("%s: %s", testPoint.code, err)
		}
		closeTo(t, lon, testPoint.ll[0], llEPSLN, fmt.Sprintf("%s inv x", testPoint.code))
		closeTo(t, lat, testPoint.ll[1], llEPSLN, fmt.Sprintf("%s inv y", testPoint.code))
	}
}

func TestWKT(t *testing.T) {
	err := addDef("EPSG:4269", `GEOGCS["NAD83",DATUM["North_American_Datum_1983",SPHEROID["GRS 1980",6378137,298.257222101,AUTHORITY["EPSG","7019"]],AUTHORITY["EPSG","6269"]],PRIMEM["Greenwich",0,AUTHORITY["EPSG","8901"]],UNIT["degree",0.01745329251994328,AUTHORITY["EPSG","9122"]],AUTHORITY["EPSG","4269"]]`)
	if err != nil {
		t.Fatal(err)
	}
	if defs["EPSG:4269"].ToMeter != 6378137*0.01745329251994328 {
		t.Errorf("should provide the correct conversion factor for WKT GEOGCS projections")
	}

	err = addDef("EPSG:4279", `GEOGCS["OS(SN)80",DATUM["OS_SN_1980",SPHEROID["Airy 1830",6377563.396,299.3249646,AUTHORITY["EPSG","7001"]],AUTHORITY["EPSG","6279"]],PRIMEM["Greenwich",0,AUTHORITY["EPSG","8901"]],UNIT["degree",0.01745329251994328,AUTHORITY["EPSG","9122"]],AUTHORITY["EPSG","4279"]]`)
	if err != nil {
		t.Fatal(err)
	}
	if defs["EPSG:4279"].ToMeter != 6377563.396*0.01745329251994328 {
		t.Errorf("should provide the correct conversion factor for WKT GEOGCS projections")
	}
}

func TestErrors(t *testing.T) {
	_, err := Parse("fake one")
	if err != nil || !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("should throw an error for an unknown ref")
	}
}

func TestDatum(t *testing.T) {
	err := addDef("EPSG:5514", "+proj=krovak +lat_0=49.5 +lon_0=24.83333333333333 +alpha=30.28813972222222 +k=0.9999 +x_0=0 +y_0=0 +ellps=bessel +pm=greenwich +units=m +no_defs +towgs84=570.8,85.7,462.8,4.998,1.587,5.261,3.56")
	if err != nil {
		t.Fatal(err)
	}
	wgs84, err := Parse("WGS84")
	if err != nil {
		t.Fatal(err)
	}
	to, err := Parse("EPSG:5514")
	if err != nil {
		t.Fatal(err)
	}
	trans, err := wgs84.NewTransformFunc(to)
	if err != nil {
		t.Fatal(err)
	}
	x, y, err := trans(12.806988, 49.452262)
	if err != nil {
		t.Fatal(err)
	}
	closeTo(t, x, -868208.61, 1.e-8, "Longitude of point from WGS84")
	closeTo(t, y, -1095793.64, 1.e-9, "Latitude of point from WGS84")
	trans2, err := wgs84.NewTransformFunc(to)
	if err != nil {
		t.Fatal(err)
	}
	x2, y2, err := trans2(12.806988, 49.452262)
	if err != nil {
		t.Fatal(err)
	}
	closeTo(t, x2, -868208.61, 1.e-8, "Longitude 2nd of point from WGS84")
	closeTo(t, y2, -1095793.64, 1.e-9, "Latitude of 2nd point from WGS84")
}
