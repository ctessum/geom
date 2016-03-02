package proj

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func findWKTSectionEnd(i int, v []interface{}) int {
	// If there is another string, that means that
	// this section is over.
	for j := i; j < len(v); j++ {
		switch v[j].(type) {
		case string:
			return j
		}
	}
	return len(v)
}

func (p *SR) sExpr(v []interface{}) error {
	for i, vv := range v {
		switch vv.(type) {
		case string:
			switch vv.(string) {
			case "PROJCS":
				p.SRSCode = v[i+1].(string)
				// we are only interested in PROJCS
				j := findWKTSectionEnd(i, v)
				return p.parseWKTProjCS(v[i+2 : j])
			case "GEOCS":
				// This should only happen if there is no PROJCS.
				p.Name = "longlat"
				j := findWKTSectionEnd(i, v)
				if err := p.parseWKTGeoCS(v[i+1 : j]); err != nil {
					return err
				}
			case "LOCAL_CS":
				p.Name = "identity"
				p.local = true
			}
		}
	}
	return nil
}

func (p *SR) parseWKTProjCS(v []interface{}) error {
	for _, vv := range v {
		vvv := vv.([]interface{})
		switch vvv[0].(type) {
		case string:
			s := vvv[0].(string)
			switch s {
			case "GEOCS":
				p.parseWKTGeoCS(vvv[1:len(vvv)])
			case "PRIMEM":
				if err := p.parseWKTPrimeM(vvv[1:len(vvv)]); err != nil {
					return err
				}
			case "PROJECTION":
				p.parseWKTProjection(vvv[1:len(vvv)])
			case "PARAMETER":
				p.parseWKTParameter(vvv[1:len(vvv)])
			case "UNIT":
				if err := p.parseWKTUnit(vvv[1:len(vvv)]); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (p *SR) parseWKTGeoCS(v []interface{}) error {
	for _, vv := range v[1:len(v)] {
		vvv := vv.([]interface{})
		switch vvv[0].(type) {
		case string:
			s := vvv[0].(string)
			switch s {
			case "DATUM":
				return p.parseWKTDatum(vvv[1:len(v)])
			}
		}
	}
	// didn't find a datum, so the datum name is the GEOCS name.
	p.DatumCode = strings.ToLower(v[0].(string))
	p.datumRename()
	return nil
}

func (p *SR) parseWKTDatum(v []interface{}) error {
	p.DatumCode = strings.ToLower(v[0].(string))
	p.datumRename()
	for _, vv := range v[1:len(v)] {
		vvv := vv.([]interface{})
		switch vvv[0].(type) {
		case string:
			s := vvv[0].(string)
			switch s {
			case "SPHEROID":
				if err := p.parseWKTSpheroid(vvv[1:len(vvv)]); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (p *SR) datumRename() {
	if p.DatumCode[0:2] == "d_" {
		p.DatumCode = p.DatumCode[2:len(p.DatumCode)]
	}
	if p.DatumCode == "new_zealand_geodetic_datum_1949" ||
		p.DatumCode == "new_zealand_1949" {
		p.DatumCode = "nzgd49"
	}
	if p.DatumCode == "wgs_1984" {
		if p.Name == "Mercator_Auxiliary_Sphere" {
			p.sphere = true
		}
		p.DatumCode = "wgs84"
	}
	if strings.HasSuffix(p.DatumCode, "_ferro") {
		p.DatumCode = strings.TrimSuffix(p.DatumCode, "_ferro")
	}
	if strings.HasSuffix(p.DatumCode, "_jakarta") {
		p.DatumCode = strings.TrimSuffix(p.DatumCode, "_jakarta")
	}
	if strings.Contains(p.DatumCode, "belge") {
		p.DatumCode = "rnb72"
	}
}

func (p *SR) parseWKTSpheroid(v []interface{}) error {
	p.Ellps = strings.Replace(v[0].(string), "_19", "", -1)
	p.Ellps = strings.Replace(p.Ellps, "clarke_18", "clrk", -1)
	p.Ellps = strings.Replace(p.Ellps, "Clarke_18", "clrk", -1)
	if strings.ToLower(p.Ellps[0:13]) == "international" {
		p.Ellps = "intl"
	}
	a, err := strconv.ParseFloat(v[1].(string), 64)
	if err != nil {
		return fmt.Errorf("in proj.parseWKTSpheroid a: %v", err)
	}
	p.A = a
	p.Rf, err = strconv.ParseFloat(v[2].(string), 64)
	if err != nil {
		return fmt.Errorf("in proj.parseWKTSpheroid rf: %v", err)
	}
	if strings.Contains(p.DatumCode, "osgb_1936") {
		p.DatumCode = "osgb36"
	}
	if math.IsNaN(p.B) {
		p.B = p.A
	}
	return nil
}

func (p *SR) parseWKTProjection(v []interface{}) {
	p.Name = v[0].(string)
}

func (p *SR) parseWKTParameter(v []interface{}) error {
	name := v[0].(string)
	val, err := strconv.ParseFloat(v[1].(string), 64)
	if err != nil {
		return fmt.Errorf("in proj.parseWKTParameter: %v", err)
	}
	switch name {
	case "Standard_Parallel_1", "standard_parallel_1":
		p.Lat0 = d2r(val)
		p.Lat1 = d2r(val)
	case "Standard_Parallel_2", "standard_parallel_2":
		p.Lat2 = d2r(val)
	case "False_Easting":
		p.X0 = p.toMeter(val)
	case "False_Northing":
		p.Y0 = p.toMeter(val)
	case "Central_Meridian":
		p.Long0 = d2r(val)
	case "Latitude_Of_Origin":
		p.Lat0 = d2r(val)
	case "Central_Parallel":
		p.Lat0 = d2r(val)
	case "Scale_Factor", "scale_factor":
		p.K0 = val
	case "Latitude_of_center", "latitude_of_center":
		p.Lat0 = d2r(val)
	case "longitude_of_center", "Longitude_Of_Center":
		p.LongC = d2r(val)
	case "false_easting":
		p.X0 = p.toMeter(val)
	case "false_northing":
		p.Y0 = p.toMeter(val)
	case "central_meridian":
		p.Long0 = d2r(val)
	case "latitude_of_origin":
		p.Lat0 = d2r(val)
	case "azimuth":
		p.Alpha = d2r(val)
	}
	return nil
}

func (p *SR) parseWKTPrimeM(v []interface{}) error {
	name := strings.ToLower(v[0].(string))
	if name != "greenwich" {
		return fmt.Errorf("in proj.parseWTKPrimeM: prime meridian is %s but"+
			"only greenwich is supported", name)
	}
	return nil
}

func (p *SR) parseWKTUnit(v []interface{}) error {
	p.Units = strings.ToLower(v[0].(string))
	if p.Units == "metre" {
		p.Units = "meter"
	}
	if len(v) > 1 {
		convert, err := strconv.ParseFloat(v[1].(string), 64)
		if err != nil {
			return fmt.Errorf("in proj.parseWKTUnit: %v", err)
		}
		if p.Name == "longlat" {
			p.ToMeter = convert * p.A
		} else {
			p.ToMeter = convert
		}
	}
	return nil
}

func d2r(input float64) float64 {
	return input * deg2rad
}

func (p *SR) toMeter(input float64) float64 {
	return p.ToMeter * input
}

var wktregexp *regexp.Regexp

func init() {
	wktregexp = regexp.MustCompile("([A-Z]+)(\\[)")
}

func wkt(wkt string) (*SR, error) {
	wkt = wktregexp.ReplaceAllString(wkt, "$2\"$1\",")
	fmt.Println(wkt)

	var lisp interface{}
	dec := json.NewDecoder(strings.NewReader(wkt))
	err := dec.Decode(&lisp)
	if err != nil {
		panic(err)
	}
	fmt.Println(lisp)
	o := newSR()
	o.sExpr(lisp.([]interface{}))
	return o, nil
}
