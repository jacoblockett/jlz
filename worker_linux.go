//go:build linux

package main

import (
	"os"
	"strconv"
	"strings"
)

// cpuQuotaWorkers returns the cgroup CPU quota rounded up to whole workers,
// or 0 when no quota is configured or it cannot be read.
func cpuQuotaWorkers() int {
	if data, err := os.ReadFile("/sys/fs/cgroup/cpu.max"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) == 2 && fields[0] != "max" {
			q, e1 := strconv.ParseInt(fields[0], 10, 64)
			p, e2 := strconv.ParseInt(fields[1], 10, 64)
			if e1 == nil && e2 == nil && q > 0 && p > 0 {
				return quotaToWorkers(q, p)
			}
		}
	}
	qdata, qe := os.ReadFile("/sys/fs/cgroup/cpu/cpu.cfs_quota_us")
	pdata, pe := os.ReadFile("/sys/fs/cgroup/cpu/cpu.cfs_period_us")
	if qe == nil && pe == nil {
		q, e1 := strconv.ParseInt(strings.TrimSpace(string(qdata)), 10, 64)
		p, e2 := strconv.ParseInt(strings.TrimSpace(string(pdata)), 10, 64)
		if e1 == nil && e2 == nil && q > 0 && p > 0 {
			return quotaToWorkers(q, p)
		}
	}
	return 0
}

func quotaToWorkers(quota, period int64) int {
	// Division avoids overflow even when values approach MaxInt64.
	n := quota / period
	if quota%period != 0 {
		n++
	}
	if n > int64(^uint(0)>>1) {
		return int(^uint(0) >> 1)
	}
	if n < 1 {
		return 1
	}
	return int(n)
}
