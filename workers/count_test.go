package workers

import (
	"runtime"
	"testing"
)

func TestLimitedWorkers(t *testing.T) {
	tests := []struct{ jobs, parallelism, quota, want int }{
		{20, 16, 0, 16}, // No artificial five- or eight-worker cap.
		{20, 16, 9, 9},
		{20, 16, 4, 4},
		{3, 16, 0, 3},
		{20, 1, 0, 1},
		{20, 2, 0, 2},
		{0, 16, 0, 1},
	}
	for _, tt := range tests {
		got := Limited(tt.jobs, tt.parallelism, tt.quota)
		if got != tt.want {
			t.Errorf("Limited(%d,%d,%d)=%d, want %d", tt.jobs, tt.parallelism, tt.quota, got, tt.want)
		}
	}
}

func TestWorkersHonorGOMAXPROCS(t *testing.T) {
	old := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(old)
	if n := Available(20); n != 1 {
		t.Fatalf("GOMAXPROCS=1 produced %d workers", n)
	}
}
