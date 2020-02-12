package goutils

import "testing"

func TestMaxInt(t *testing.T) {
	a, b, c := 1, 2, 3
	if max := MaxInt(a, b, c); max != c {
		t.Logf("except max is %d, but actual is %d", c, max)
		t.FailNow()
	}
}

func TestMinInt(t *testing.T) {
	a, b, c := 1, 2, 3
	if min := MinInt(a, b, c); min != a {
		t.Logf("except min is %d, but actual is %d", a, min)
		t.FailNow()
	}
}
