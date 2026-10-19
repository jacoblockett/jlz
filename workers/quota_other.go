//go:build !linux

package workers

func cpuQuotaWorkers() int { return 0 }
