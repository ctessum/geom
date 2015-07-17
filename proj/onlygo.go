// +build onlygo

package proj

import (
	"fmt"
	"io"

	"github.com/ctessum/geom"
)

type CoordinateTransform struct{}

func NewCoordinateTransform(src, dst SR) (ct *CoordinateTransform, err error) {
	fmt.Println("Not implemented in pure go. To use this function " +
		"compile without the `onlygo` tag")
	return
}
func (ct *CoordinateTransform) Reproject(g geom.T) (geom.T, error) {
	fmt.Println("Not implemented in pure go. To use this function " +
		"compile without the `onlygo` tag")
	var err error
	return geom.Point{}, err
}

type ParsedProj4 struct{}

func ParseProj4(proj4 string) *ParsedProj4 {
	fmt.Println("Not implemented in pure go. To use this function " +
		"compile without the `onlygo` tag")
	return nil
}
func (p *ParsedProj4) Equals(p2 *ParsedProj4) bool {
	fmt.Println("Not implemented in pure go. To use this function " +
		"compile without the `onlygo` tag")
	return false
}
func (p *ParsedProj4) ToString() string {
	fmt.Println("Not implemented in pure go. To use this function " +
		"compile without the `onlygo` tag")
	return ""
}

type SR struct{}

func FromProj4(proj4 string) (SR, error) {
	fmt.Println("Not implemented in pure go. To use this function " +
		"compile without the `onlygo` tag")
	var err error
	return SR{}, err
}
func ReadPrj(f io.Reader) (SR, error) {
	fmt.Println("Not implemented in pure go. To use this function " +
		"compile without the `onlygo` tag")
	var err error
	return SR{}, err
}

type UnsupportedGeometryError struct{}

func (e UnsupportedGeometryError) Error() string {
	fmt.Println("Not implemented in pure go. To use this function " +
		"compile without the `onlygo` tag")
	return ""
}
