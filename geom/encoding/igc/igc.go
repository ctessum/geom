package igc

import (
	"bufio"
	"fmt"
	"github.com/twpayne/gogeom/geom"
	"io"
	"strings"
	"time"
)

type Errors map[int]error

func (es Errors) Error() string {
	ss := make([]string, len(es))
	for i, e := range es {
		ss[i] = e.Error()
	}
	return strings.Join(ss, "\n")
}

func parseDec(s string, start, stop int) (int, error) {
	result := 0
	for i := start; i < stop; i++ {
		if c := s[i]; '0' <= c && c <= '9' {
			result = 10*result + int(c) - '0'
		} else {
			return 0, fmt.Errorf("invalid")
		}
	}
	return result, nil
}

type parser struct {
	pointMs          []geom.PointM
	year, month, day int
}

func newParser() *parser {
	p := new(parser)
	p.year = 2000
	p.month = 1
	p.day = 1
	return p
}

func (p *parser) parseB(line string) error {
	var err error
	var hour, minute, second int
	if hour, err = parseDec(line, 1, 3); err != nil {
		return err
	}
	if minute, err = parseDec(line, 3, 5); err != nil {
		return err
	}
	if second, err = parseDec(line, 5, 7); err != nil {
		return err
	}
	var latDeg, latMin int
	if latDeg, err = parseDec(line, 7, 9); err != nil {
		return err
	}
	if latMin, err = parseDec(line, 9, 14); err != nil {
		return err
	}
	lat := float64(latDeg) + float64(latMin)/60000.
	switch c := line[14]; c {
	case 'N':
	case 'S':
		lat = -lat
	default:
		return fmt.Errorf("unexpected character %v", c)
	}
	var lngDeg, lngMin int
	lngDeg, err = parseDec(line, 15, 18)
	if err != nil {
		return err
	}
	lngMin, err = parseDec(line, 18, 23)
	if err != nil {
		return err
	}
	lng := float64(lngDeg) + float64(lngMin)/60000.
	switch c := line[23]; c {
	case 'E':
	case 'W':
		lng = -lng
	default:
		return fmt.Errorf("unexpected character %v", c)
	}
	date := time.Date(p.year, time.Month(p.month), p.day, hour, minute, second, 0, time.UTC)
	pointM := geom.PointM{lng, lat, float64(date.UnixNano()) / 1e9}
	p.pointMs = append(p.pointMs, pointM)
	return nil
}

func (p *parser) parseH(line string) error {
	switch {
	case strings.HasPrefix(line, "HFDTE"):
		return p.parseHFDTE(line)
	default:
		return nil
	}
}

func (p *parser) parseHFDTE(line string) error {
	var err error
	var day, month, year int
	if day, err = parseDec(line, 5, 7); err != nil {
		return err
	}
	if month, err = parseDec(line, 7, 9); err != nil {
		return err
	}
	if year, err = parseDec(line, 9, 11); err != nil {
		return err
	}
	// FIXME check for invalid dates
	p.day = day
	p.month = month
	if year < 70 {
		p.year = 2000 + year
	} else {
		p.year = 1970 + year
	}
	return nil
}

func (p *parser) parseLine(line string) error {
	switch line[0] {
	case 'B':
		return p.parseB(line)
	case 'H':
		return p.parseH(line)
	default:
		return nil
	}
}

func doParse(r io.Reader) (*parser, Errors) {
	errors := make(Errors)
	p := newParser()
	s := bufio.NewScanner(r)
	line := 0
	for s.Scan() {
		line++
		if err := p.parseLine(s.Text()); err != nil {
			errors[line] = err
		}
	}
	return p, errors
}

func Read(r io.Reader) ([]geom.PointM, error) {
	p, errors := doParse(r)
	if len(errors) == 0 {
		return p.pointMs, nil
	} else {
		return p.pointMs, errors
	}
}
