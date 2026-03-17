package util

import (
	"sync/atomic"
	"time"
)

var (
	Version       string
	BuildTime     string
	StartTime     = time.Now()
	TotalRequests uint64
)

func IncrementRequest() {
	atomic.AddUint64(&TotalRequests, 1)
}

func GetTotalRequests() uint64 {
	return atomic.LoadUint64(&TotalRequests)
}
