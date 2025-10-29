package ingestion

import (
	"sync/atomic"
	"time"
)

// Metrics captures ingestion performance indicators for monitoring.
type Metrics struct {
	acceptedCount     atomic.Int64
	rejectedCount     atomic.Int64
	totalLatencyNanos atomic.Int64
}

// ObserveAccepted records a successful ingestion acknowledgement latency.
func (m *Metrics) ObserveAccepted(latency time.Duration) {
	m.acceptedCount.Add(1)
	m.totalLatencyNanos.Add(latency.Nanoseconds())
}

// ObserveRejected increments the rejection counter.
func (m *Metrics) ObserveRejected() {
	m.rejectedCount.Add(1)
}

// Snapshot returns current counter values and average latency.
func (m *Metrics) Snapshot() (accepted, rejected int64, avgLatency time.Duration) {
	accepted = m.acceptedCount.Load()
	rejected = m.rejectedCount.Load()
	total := m.totalLatencyNanos.Load()
	if accepted == 0 {
		return accepted, rejected, 0
	}
	return accepted, rejected, time.Duration(total / accepted)
}
