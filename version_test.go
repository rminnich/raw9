package raw9

import "testing"

func TestVersion(t *testing.T) {
	var b [128]byte

	l := xtversion(b[:], "9P2000.L", 0xcafebeef)
	v, m := rtversion(b[:l])

	if v != "9P2000.L" {
		t.Errorf("%q != 9P2000.L", v)
	}
	if m != 0xcafebeef {
		t.Errorf("%#x != 0xcafebeef", m)
	}

	l = xrversion(b[:], "9P2000.L", 0xcafebeef)
	v, m = rrversion(b[:l])

	if v != "9P2000.L" {
		t.Errorf("%q != 9P2000.L", v)
	}
	if m != 0xcafebeef {
		t.Errorf("%#x != 0xcafebeef", m)
	}
}

func BenchmarkXTVersion(b *testing.B) {
	var d [128]byte
	for i := 0; i < b.N; i++ {
		xtversion(d[:], "9P2000.L", 0xcafebeef)
	}
}

func BenchmarkRTVersion(b *testing.B) {
	var d [128]byte
	l := xtversion(d[:], "9P2000.L", 0xcafebeef)
	for i := 0; i < b.N; i++ {
		rtversion(d[:l])
	}
}

func BenchmarkXRVersion(b *testing.B) {
	var d [128]byte
	for i := 0; i < b.N; i++ {
		xtversion(d[:], "9P2000.L", 0xcafebeef)
	}
}

func BenchmarkRRVersion(b *testing.B) {
	var d [128]byte
	l := xtversion(d[:], "9P2000.L", 0xcafebeef)
	for i := 0; i < b.N; i++ {
		rtversion(d[:l])
	}
}
