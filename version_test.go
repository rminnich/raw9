package main

import "testing"

func TestVersion(t *testing.T) {
	var b [128]byte

	l := xversion(b[:], "9P2000.L", 0xcafebeef)
	v, m := rversion(b[:l])

	if v != "9P2000.L" {
		t.Errorf("%q != 9P2000.L", v)
	}
	if m != 0xcafebeef {
		t.Errorf("%#x != 0xcafebeef", m)
	}
}

func BenchmarkTVersion(b *testing.B) {
	var d [128]byte
	for i := 0; i < b.N; i++ {
		xversion(d[:], "9P2000.L", 0xcafebeef)
	}
}

func BenchmarkRVersion(b *testing.B) {
	var d [128]byte
	l := xversion(d[:], "9P2000.L", 0xcafebeef)
	for i := 0; i < b.N; i++ {
		rversion(d[:l])
	}
}
