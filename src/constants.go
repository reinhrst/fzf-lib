package fzf

import (
	"time"

	"github.com/reinhrst/fzf-lib/src/util"
)

const (
	// Matcher
	numPartitionsMultiplier = 8
	maxPartitions           = 32
	progressMinDuration     = 200 * time.Millisecond

	// Capacity of each chunk
	chunkSize int = 100

	// Pre-allocated memory slices to minimize GC
	slab16Size int = 100 * 1024 // 200KB * 32 = 12.8MB
	slab32Size int = 2048       // 8KB * 32 = 256KB

	// Do not cache results of low selectivity queries
	queryCacheMax int = chunkSize / 5

	// Not to cache mergers with large lists
	mergerCacheMax int = 100000
)

// fzf events
const (
	EvtSearchProgress util.EventType = iota
	EvtSearchFin
	EvtQuit
)
