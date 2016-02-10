package proj

import (
	"fmt"
	"strconv"
	"strings"
)

const D2R = 0.01745329251994329577

func projString(defData string) (*Proj, error) {
	self := newProj()
	paramObj := make(map[string]string)
	for _, a := range strings.Split(defData, "+") {
		a = strings.TrimSpace(a)
		split := strings.Split(a, "=")
		split = append(split, "true") // TODO: Not sure why this is done.
		paramObj[strings.ToLower(split[0])] = split[1]
	}
	var err error
	for paramName, paramVal := range paramObj {
		switch paramName {
		case "proj":
			self.projName = paramVal
		case "datum":
			self.datumCode = paramVal
		case "rf":
			self.rf, err = strconv.ParseFloat(paramVal, 64)
		case "lat_0":
			self.lat0, err = strconv.ParseFloat(paramVal, 64)
			self.lat0 *= D2R
		case "lat_1":
			self.lat1, err = strconv.ParseFloat(paramVal, 64)
			self.lat1 *= D2R
		case "lat_2":
			self.lat2, err = strconv.ParseFloat(paramVal, 64)
			self.lat2 *= D2R
		case "lat_ts":
			self.lat_ts, err = strconv.ParseFloat(paramVal, 64)
			self.lat_ts *= D2R
		case "lon_0":
			self.long0, err = strconv.ParseFloat(paramVal, 64)
			self.long0 *= D2R
		case "lon_1":
			self.long1, err = strconv.ParseFloat(paramVal, 64)
			self.long1 *= D2R
		case "lon_2":
			self.long2, err = strconv.ParseFloat(paramVal, 64)
			self.long2 *= D2R
		case "alpha":
			self.alpha, err = strconv.ParseFloat(paramVal, 64)
			self.alpha *= D2R
		case "lonc":
			self.longc, err = strconv.ParseFloat(paramVal, 64)
			self.longc *= D2R
		case "x_0":
			self.x0, err = strconv.ParseFloat(paramVal, 64)
		case "y_0":
			self.y0, err = strconv.ParseFloat(paramVal, 64)
		case "k_0", "k":
			self.k0, err = strconv.ParseFloat(paramVal, 64)
		case "a":
			self.a, err = strconv.ParseFloat(paramVal, 64)
		case "b":
			self.b, err = strconv.ParseFloat(paramVal, 64)
		case "r_a":
			self.R_A = true
		case "zone":
			self.zone, err = strconv.ParseInt(paramVal, 10, 64)
		case "south":
			self.utmSouth = true
		case "towgs84":
			split := strings.Split(paramVal, ",")
			self.datum_params = make([]float64, len(split))
			for i, s := range split {
				self.datum_params[i], err = strconv.ParseFloat(s, 64)
				if err != nil {
					return nil, err
				}
			}
		case "to_meter":
			self.to_meter, err = strconv.ParseFloat(paramVal, 64)
		case "units":
			self.units = paramVal
			if u, ok := units[paramVal]; ok {
				self.to_meter = u.to_meter
			}
		case "from_greenwich":
			self.from_greenwich, err = strconv.ParseFloat(paramVal, 64)
			self.from_greenwich *= D2R
		case "pm":
			if pm, ok := primeMeridian[paramVal]; ok {
				self.from_greenwich = pm
			} else {
				self.from_greenwich, err = strconv.ParseFloat(paramVal, 64)
				self.from_greenwich *= D2R
			}
		case "nadgrids":
			if paramVal == "@null" {
				self.datumCode = "none"
			} else {
				self.nadGrids = paramVal
			}
		case "axis":
			legalAxis := "ewnsud"
			if len(paramVal) == 3 && strings.Index(legalAxis, paramVal[0:1]) != -1 &&
				strings.Index(legalAxis, paramVal[1:2]) != -1 &&
				strings.Index(legalAxis, paramVal[2:3]) != -1 {
				self.axis = paramVal
			}
		default:
			err = fmt.Errorf("proj: invalid field %s", paramName)
		}
		if err != nil {
			return nil, err
		}
	}
	if self.datumCode != "WGS84" {
		self.datumCode = strings.ToLower(self.datumCode)
	}
	return self, nil
}
