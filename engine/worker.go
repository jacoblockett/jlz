package engine

import "github.com/jacoblockett/jlz/workers"

// Match the previous internal symbols so existing engine call sites and tests
// preserve their behavior while worker policy belongs to its own package.
func availableWorkers(jobs int) int { return workers.Available(jobs) }
func limitedWorkers(jobs, gomaxprocs, quota int) int { return workers.Limited(jobs, gomaxprocs, quota) }
