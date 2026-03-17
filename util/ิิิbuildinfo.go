package util

import (
	"fmt"
	"sync/atomic"
	"time"
)

var loc, _ = time.LoadLocation("Asia/Bangkok")

var (
	Version       string
	buildTime     string
	startTime     = time.Now()
	totalRequests uint64
)

func Buildtime() string {
	bt, err := time.Parse(time.RFC3339, buildTime)
	if err != nil {
		return "Invalid Time!"
	} else {
		return bt.In(loc).Format("2006/01/02 | 15:04:05")
	}
}

func StartTime() string {
	return startTime.In(loc).Format("2006/01/02 | 15:04:05")
}

func Uptime() string {
	duration := time.Since(startTime)
	h := int(duration.Hours())
	m := int(duration.Minutes()) % 60
	s2 := int(duration.Seconds()) % 60
	return fmt.Sprintf("%dh %dm %ds", h, m, s2)
}

func IncrementRequest() {
	atomic.AddUint64(&totalRequests, 1)
}

func GetTotalRequests() uint64 {
	return atomic.LoadUint64(&totalRequests)
}
