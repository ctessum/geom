package igc

import (
	"github.com/twpayne/gogeom/geom"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test(t *testing.T) {
	var testCases = []struct {
		igc string
		g   []geom.PointM
	}{
		{
			"HFDTE290813\nB1000294358694N00628792EA0147501551000\n",
			[]geom.PointM{{
				float64(6) + float64(28792)/60000.,
				float64(43) + float64(58694)/60000.,
				float64(time.Date(2013, 8, 29, 10, 0, 29, 0, time.UTC).UnixNano()) / 1e9,
			}},
		},
	}

	for _, tc := range testCases {
		if got, err := Read(strings.NewReader(tc.igc)); err != nil || !reflect.DeepEqual(got, tc.g) {
			t.Errorf("Read(%q) == %q, %q, want %q, nil", tc.igc, got, err, tc.g)
		}
	}
}

func TestErrors(t *testing.T) {

	var testCases = []string{
		//"HFDTE290813\nB1000294358694N00628792EA0147501551000\n",
		"HFDTE290813\nB1X00294358694N00628792EA0147501551000\n",
		"HFDTE290813\nB100X294358694N00628792EA0147501551000\n",
		"HFDTE290813\nB10002X4358694N00628792EA0147501551000\n",
		"HFDTE290813\nB1000294X58694N00628792EA0147501551000\n",
		"HFDTE290813\nB100029435869XN00628792EA0147501551000\n",
		"HFDTE290813\nB1000294358694X00628792EA0147501551000\n",
		"HFDTE290813\nB1000294358694N00X28792EA0147501551000\n",
		"HFDTE290813\nB1000294358694N0062879XEA0147501551000\n",
		"HFDTE290813\nB1000294358694N00628792XA0147501551000\n",
	}

	for _, tc := range testCases {
		if _, err := Read(strings.NewReader(tc)); err == nil {
			t.Errorf("Read(%q) == <don't-care>, nil, want <don't-care>, !nil", tc, err)
		}
	}
}
