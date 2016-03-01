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

func (p *Proj) sExpr(v []interface{}) error {
	for i, vv := range v {
		switch vv.(type) {
		case string:
			switch vv.(string) {
			case "PROJCS":
				p.srsCode = v[i+1].(string)
				// we are only interested in PROJCS
				j := findWKTSectionEnd(i, v)
				return p.parseWKTProjCS(v[i+2 : j])
			case "GEOCS":
				// This should only happen if there is no PROJCS.
				p.name = "longlat"
				j := findWKTSectionEnd(i, v)
				if err := p.parseWKTGeoCS(v[i+1 : j]); err != nil {
					return err
				}
			case "LOCAL_CS":
				p.name = "identity"
				p.local = true
			}
		}
	}
	return nil
}

func (p *Proj) parseWKTProjCS(v []interface{}) error {
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

func (p *Proj) parseWKTGeoCS(v []interface{}) error {
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
	p.datumCode = strings.ToLower(v[0].(string))
	p.datumRename()
	return nil
}

func (p *Proj) parseWKTDatum(v []interface{}) error {
	p.datumCode = strings.ToLower(v[0].(string))
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

func (p *Proj) datumRename() {
	if p.datumCode[0:2] == "d_" {
		p.datumCode = p.datumCode[2:len(p.datumCode)]
	}
	if p.datumCode == "new_zealand_geodetic_datum_1949" ||
		p.datumCode == "new_zealand_1949" {
		p.datumCode = "nzgd49"
	}
	if p.datumCode == "wgs_1984" {
		if p.name == "Mercator_Auxiliary_Sphere" {
			p.sphere = true
		}
		p.datumCode = "wgs84"
	}
	if strings.HasSuffix(p.datumCode, "_ferro") {
		p.datumCode = strings.TrimSuffix(p.datumCode, "_ferro")
	}
	if strings.HasSuffix(p.datumCode, "_jakarta") {
		p.datumCode = strings.TrimSuffix(p.datumCode, "_jakarta")
	}
	if strings.Contains(p.datumCode, "belge") {
		p.datumCode = "rnb72"
	}
}

func (p *Proj) parseWKTSpheroid(v []interface{}) error {
	p.ellps = strings.Replace(v[0].(string), "_19", "", -1)
	p.ellps = strings.Replace(p.ellps, "clarke_18", "clrk", -1)
	p.ellps = strings.Replace(p.ellps, "Clarke_18", "clrk", -1)
	if strings.ToLower(p.ellps[0:13]) == "international" {
		p.ellps = "intl"
	}
	a, err := strconv.ParseFloat(v[1].(string), 64)
	if err != nil {
		return fmt.Errorf("in proj.parseWKTSpheroid a: %v", err)
	}
	p.a = a
	p.rf, err = strconv.ParseFloat(v[2].(string), 64)
	if err != nil {
		return fmt.Errorf("in proj.parseWKTSpheroid rf: %v", err)
	}
	if strings.Contains(p.datumCode, "osgb_1936") {
		p.datumCode = "osgb36"
	}
	if math.IsNaN(p.b) {
		p.b = p.a
	}
	return nil
}

func (p *Proj) parseWKTProjection(v []interface{}) {
	p.name = v[0].(string)
}

func (p *Proj) parseWKTParameter(v []interface{}) error {
	name := v[0].(string)
	val, err := strconv.ParseFloat(v[1].(string), 64)
	if err != nil {
		return fmt.Errorf("in proj.parseWKTParameter: %v", err)
	}
	switch name {
	case "Standard_Parallel_1", "standard_parallel_1":
		p.lat0 = d2r(val)
		p.lat1 = d2r(val)
	case "Standard_Parallel_2", "standard_parallel_2":
		p.lat2 = d2r(val)
	case "False_Easting":
		p.x0 = p.toMeter(val)
	case "False_Northing":
		p.y0 = p.toMeter(val)
	case "Central_Meridian":
		p.long0 = d2r(val)
	case "Latitude_Of_Origin":
		p.lat0 = d2r(val)
	case "Central_Parallel":
		p.lat0 = d2r(val)
	case "Scale_Factor", "scale_factor":
		p.k0 = val
	case "Latitude_of_center", "latitude_of_center":
		p.lat0 = d2r(val)
	case "longitude_of_center", "Longitude_Of_Center":
		p.longc = d2r(val)
	case "false_easting":
		p.x0 = p.toMeter(val)
	case "false_northing":
		p.y0 = p.toMeter(val)
	case "central_meridian":
		p.long0 = d2r(val)
	case "latitude_of_origin":
		p.lat0 = d2r(val)
	case "azimuth":
		p.alpha = d2r(val)
	}
	return nil
}

func (p *Proj) parseWKTPrimeM(v []interface{}) error {
	name := strings.ToLower(v[0].(string))
	if name != "greenwich" {
		return fmt.Errorf("in proj.parseWTKPrimeM: prime meridian is %s but"+
			"only greenwich is supported", name)
	}
	return nil
}

func (p *Proj) parseWKTUnit(v []interface{}) error {
	p.units = strings.ToLower(v[0].(string))
	if p.units == "metre" {
		p.units = "meter"
	}
	if len(v) > 1 {
		convert, err := strconv.ParseFloat(v[1].(string), 64)
		if err != nil {
			return fmt.Errorf("in proj.parseWKTUnit: %v", err)
		}
		if p.name == "longlat" {
			p.to_meter = convert * p.a
		} else {
			p.to_meter = convert
		}
	}
	return nil
}

func d2r(input float64) float64 {
	return input * D2R
}

func (p *Proj) toMeter(input float64) float64 {
	return p.to_meter * input
}

var wktregexp *regexp.Regexp

func init() {
	wktregexp = regexp.MustCompile("([A-Z]+)(\\[)")
}

func wkt(wkt string) (*Proj, error) {
	wkt = wktregexp.ReplaceAllString(wkt, "$2\"$1\",")
	fmt.Println(wkt)

	var lisp interface{}
	dec := json.NewDecoder(strings.NewReader(wkt))
	err := dec.Decode(&lisp)
	if err != nil {
		panic(err)
	}
	fmt.Println(lisp)
	o := newProj()
	o.sExpr(lisp.([]interface{}))
	return o, nil
}
