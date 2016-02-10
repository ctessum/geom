package proj

import "math"

const (
	EPSLN = 1.0e-10
	// ellipoid pj_set_ell.c
	SIXTH = 0.1666666666666666667
	/* 1/6 */
	RA4 = 0.04722222222222222222
	/* 17/360 */
	RA6 = 0.02215608465608465608
)

func (json *Proj) deriveConstants() {
	// DGR 2011-03-20 : nagrids -> nadgrids
	if json.datumCode != "" && json.datumCode != "none" {
		datumDef, ok := datumDefs[json.datumCode]
		if ok {
			json.datum_params = datumDef.towgs84
			json.ellps = datumDef.ellipse
			if datumDef.datumName != "" {
				json.datumName = datumDef.datumName
			} else {
				json.datumName = json.datumCode
			}
		}
	}
	if math.IsNaN(json.a) { // do we have an ellipsoid?
		ellipse, ok := ellipsoidDefs[json.ellps]
		if !ok {
			ellipse = ellipsoidDefs["WGS84"]
		}
		json.a, json.b, json.rf = ellipse.a, ellipse.b, ellipse.rf
		json.ellipseName = ellipse.ellipseName
	}
	if !math.IsNaN(json.rf) && math.IsNaN(json.b) {
		json.b = (1.0 - 1.0/json.rf) * json.a
	}
	if json.rf == 0 || math.Abs(json.a-json.b) < EPSLN {
		json.sphere = true
		json.b = json.a
	}
	json.a2 = json.a * json.a               // used in geocentric
	json.b2 = json.b * json.b               // used in geocentric
	json.es = (json.a2 - json.b2) / json.a2 // e ^ 2
	json.e = math.Sqrt(json.es)             // eccentricity
	if json.R_A {
		json.a *= 1 - json.es*(SIXTH+json.es*(RA4+json.es*RA6))
		json.a2 = json.a * json.a
		json.b2 = json.b * json.b
		json.es = 0
	}
	json.ep2 = (json.a2 - json.b2) / json.b2 // used in geocentric
	if math.IsNaN(json.k0) {
		json.k0 = 1.0 //default value
	}
	//DGR 2010-11-12: axis
	if json.axis == "" {
		json.axis = "enu"
	}

	if json.datum == nil {
		json.datum = json.getDatum()
	}
}
