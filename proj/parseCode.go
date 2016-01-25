package proj

import "fmt"

func testWKT(code string) bool {
	var codeWords = []string{"GEOGCS", "GEOCCS", "PROJCS", "LOCAL_CS"}
	for _, c := range codeWords {
		if code == c {
			return true
		}
	}
	return false
}
func testProj(code string) bool {
	return code[0] == '+'
}
func parse(code string) (*Proj, error) {
	//check to see if this is a WKT string
	if p, ok := defs[code]; ok {
		return p, nil
	} else if testWKT(code) {
		return wkt(code)
	} else if testProj(code) {
		return projString(code)
	}
	return nil, fmt.Errorf("unsupported projection definition %s; only proj4 and "+
		"WKT are supported", code)
}
