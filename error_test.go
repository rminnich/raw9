package main

import "testing"

func TestError(t *testing.T) {
	var b [128]byte

	l := xrerror(b[:], 1, "suck")
	v := rrerror(b[:l])

	if v != "suck" {
		t.Errorf("%q != suck", v)
	}
}

func BenchmarkXRError(b *testing.B) {
	var d [128]byte
	for i := 0; i < b.N; i++ {
		xrerror(d[:], 1, "suck")
	}
}

func BenchmarkRRError(b *testing.B) {
	var d [128]byte
	l := xrerror(d[:], 1, "suck")
	for i := 0; i < b.N; i++ {
		rrerror(d[:l])
	}
}
