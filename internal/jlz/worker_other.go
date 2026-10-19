//go:build !linux

package jlz

func cpuQuotaWorkers() int { return 0 }
