//go:build !linux

package main

func cpuQuotaWorkers() int { return 0 }
