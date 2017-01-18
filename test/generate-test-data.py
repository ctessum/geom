#!/usr/bin/python

import random
import sys

import shapely.geometry


R = random.Random(0)


def r():
    return float(R.randint(-1000000, 1000000)) / 1000000


def randomCoord():
    return (r(), r())


def randomCoords(n):
    return [(r(), r()) for i in xrange(n)]


class RandomPoint(shapely.geometry.Point):

    def __init__(self, coord=None):
        if coord is None:
            coord = randomCoord()
        shapely.geometry.Point.__init__(self, coord)

    def goify(self):
        return 'geom.Point{%f, %f}' % (self.x, self.y)


class RandomLineString(shapely.geometry.LineString):

    def __init__(self, coords=None):
        if coords is None:
            coords = randomCoords(R.randint(2, 8))
        shapely.geometry.LineString.__init__(self, coords)

    def goify(self):
        return 'geom.LineString{[]geom.Point{' + ', '.join('{%f, %f}' % c for c in self.coords) + '}}'


class RandomPolygon(shapely.geometry.Polygon):

    def __init__(self, rings=None):
        if rings is None:
            rings = [randomCoords(R.randint(3, 8))] + [randomCoords(R.randint(3, 8)) for i in xrange(R.randint(0, 4))]
        shapely.geometry.Polygon.__init__(self, rings[0], rings[1:])

    def goify(self):
        rings = [self.exterior.coords] + [i.coords for i in self.interiors]
        return 'geom.Polygon{[]geom.Path{' + ', '.join('{' + ', '.join('{%f, %f}' % c for c in ring) + '}' for ring in rings) + '}}'


class RandomMultiPoint(shapely.geometry.MultiPoint):

    def __init__(self):
        shapely.geometry.MultiPoint.__init__(self, [RandomPoint() for i in xrange(R.randint(1, 8))])

    def goify(self):
        return 'geom.MultiPoint{[]geom.Point{' + ', '.join(RandomPoint(g.coords[0]).goify() for g in self.geoms) + '}}'


class RandomMultiLineString(shapely.geometry.MultiLineString):

    def __init__(self):
        shapely.geometry.MultiLineString.__init__(self, [RandomLineString() for i in xrange(R.randint(1, 8))])

    def goify(self):
        return 'geom.MultiLineString{[]geom.LineString{' + ', '.join(RandomLineString(g.coords).goify() for g in self.geoms) + '}}'


class RandomMultiPolygon(shapely.geometry.MultiPolygon):

    def __init__(self):
        shapely.geometry.MultiPolygon.__init__(self, [RandomPolygon() for i in xrange(R.randint(1, 8))])

    def goify(self):
        return 'geom.MultiPolygon{[]geom.Polygon{' + ', '.join(RandomPolygon([g.exterior] + list(g.interiors)).goify() for g in self.geoms) + '}}'


def main(argv):
    # FIXME add GeoJSON support
    print 'package test'
    print
    print 'import ('
    print '\t"github.com/twpayne/gogeom/geom"'
    print ')'
    print
    print 'var cases = []struct {'
    print '\tg   geom.Geom'
    print '\thex string'
    print '\twkb []byte'
    print '\twkt string'
    print '}{'
    for klass in (
            RandomPoint,
            RandomLineString,
            RandomPolygon,
            RandomMultiPoint,
            RandomMultiLineString,
            RandomMultiPolygon):
        for i in xrange(8):
            g = klass()
            print '\t{'
            print '\t\t%s,' % (g.goify(),)
            print '\t\t"%s",' % (g.wkb.encode('hex'),)
            print '\t\t[]byte("%s"),' % (''.join('\\x%02X' % ord(c) for c in g.wkb),)
            print '\t\t"%s",' % (g.wkt,)
            print '\t},'
    print '}'


if __name__ == '__main__':
    sys.exit(main(sys.argv))
