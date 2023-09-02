package monitor

import "testing"

const PID = 998

func TestGetProcessInfo(t *testing.T) {
	p, err := getProcessInfo(PID)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(p)
}

func BenchmarkGetProcessInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pi, err := getProcessInfo(PID)
		if err != nil {
			b.Error(err)
			return
		}
		pi.Recycle()
	}
}
