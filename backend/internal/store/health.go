package store

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
)

func (s *SQLStorage) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := s.connPool.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	poolStats := s.connPool.Stat()
	stats["acquire_count"] = strconv.FormatInt(poolStats.AcquireCount(), 10)
	stats["acquired_conns"] = strconv.FormatInt(int64(poolStats.AcquiredConns()), 10)
	stats["acquire_duration"] = poolStats.AcquireDuration().String()
	stats["canceled_acquire_count"] = strconv.FormatInt(poolStats.CanceledAcquireCount(), 10)
	stats["constructing_conns"] = strconv.FormatInt(int64(poolStats.ConstructingConns()), 10)
	stats["empty_acquire_count"] = strconv.FormatInt(poolStats.EmptyAcquireCount(), 10)
	stats["idle_conns"] = strconv.FormatInt(int64(poolStats.IdleConns()), 10)
	stats["max_conns"] = strconv.FormatInt(int64(poolStats.MaxConns()), 10)
	stats["total_conns"] = strconv.FormatInt(int64(poolStats.TotalConns()), 10)

	if poolStats.EmptyAcquireCount() > 1000 {
		stats["message"] = "Connection pool has been exhausted over 1000 times; consider increasing max connections or reviewing query efficiency."
	}

	if poolStats.AcquireDuration() > 500*time.Millisecond {
		stats["message"] = "Average connection acquire time is high; potential connection contention or slow database responses."
	}

	if poolStats.AcquiredConns() > poolStats.MaxConns()*80/100 {
		stats["message"] = "Over 80% of the pool connections are currently in use; approaching saturation."
	}

	if poolStats.ConstructingConns() > 5 {
		stats["message"] = "New connections are being constructed frequently; connection churn may be high."
	}

	if poolStats.IdleConns() == 0 && poolStats.AcquiredConns() > 0 {
		stats["message"] = "No idle connections available while some are in use; pool might be undersized."
	}

	return stats
}
