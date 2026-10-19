package jlz

import "runtime"

// availableWorkers chooses the work parallelism automatically. GOMAXPROCS
// accounts for CPU affinity and explicit process limits; cpuQuotaWorkers
// additionally honors Linux container CPU quotas on Go versions before 1.25.
// No preset cap is imposed on machines with more usable CPUs.
func availableWorkers(jobs int) int {
	return limitedWorkers(jobs, runtime.GOMAXPROCS(0), cpuQuotaWorkers())
}

func limitedWorkers(jobs, gomaxprocs, quota int) int {
	if jobs < 1 {
		return 1
	}
	if gomaxprocs < 1 {
		gomaxprocs = 1
	}
	workers := min(jobs, gomaxprocs)
	if quota > 0 && workers > quota {
		workers = quota
	}
	return workers
}
