// Package snowflake provides a distributed unique ID generator based on Twitter's Snowflake algorithm.
// IDs are 64-bit integers composed of:
// - 41 bits: timestamp (milliseconds since custom epoch) - ~69 years of IDs
// - 10 bits: node/machine ID (0-1023) - supports distributed generation
// - 12 bits: sequence number (0-4095) - up to 4096 IDs per millisecond per node
package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	// Custom epoch: 2020-01-01 00:00:00 UTC
	epoch int64 = 1577836800000

	// Bit lengths
	nodeBits     = 10
	sequenceBits = 12

	// Max values
	maxNodeID   = -1 ^ (-1 << nodeBits)     // 1023
	maxSequence = -1 ^ (-1 << sequenceBits) // 4095

	// Bit shifts
	nodeShift      = sequenceBits
	timestampShift = nodeBits + sequenceBits
)

var (
	ErrInvalidNodeID  = errors.New("node ID must be between 0 and 1023")
	ErrClockMovedBack = errors.New("clock moved backwards, refusing to generate ID")
)

// Generator generates unique snowflake IDs.
type Generator struct {
	mu            sync.Mutex
	nodeID        int64
	lastTimestamp int64
	sequence      int64

	// timeFunc allows overriding time for testing
	timeFunc func() time.Time
}

// IDGenerator interface for dependency injection and mocking.
type IDGenerator interface {
	NextID() (int64, error)
}

// Ensure Generator implements IDGenerator
var _ IDGenerator = (*Generator)(nil)

// NewGenerator creates a new snowflake ID generator with the given node ID.
// Node ID must be between 0 and 1023 (inclusive).
func NewGenerator(nodeID int64) (*Generator, error) {
	if nodeID < 0 || nodeID > maxNodeID {
		return nil, ErrInvalidNodeID
	}

	return &Generator{
		nodeID:   nodeID,
		timeFunc: time.Now,
	}, nil
}

// NewGeneratorWithTime creates a generator with a custom time function (for testing).
func NewGeneratorWithTime(nodeID int64, timeFunc func() time.Time) (*Generator, error) {
	g, err := NewGenerator(nodeID)
	if err != nil {
		return nil, err
	}
	g.timeFunc = timeFunc
	return g, nil
}

// NextID generates the next unique snowflake ID.
// Returns an error if the clock has moved backwards.
func (g *Generator) NextID() (int64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	timestamp := g.currentMillis()

	if timestamp < g.lastTimestamp {
		return 0, ErrClockMovedBack
	}

	if timestamp == g.lastTimestamp {
		// Same millisecond, increment sequence
		g.sequence = (g.sequence + 1) & maxSequence
		if g.sequence == 0 {
			// Sequence overflow, wait for next millisecond
			timestamp = g.waitNextMillis(timestamp)
		}
	} else {
		// New millisecond, reset sequence
		g.sequence = 0
	}

	g.lastTimestamp = timestamp

	id := ((timestamp - epoch) << timestampShift) |
		(g.nodeID << nodeShift) |
		g.sequence

	return id, nil
}

// MustNextID generates the next unique snowflake ID.
// Panics if there's an error (e.g., clock moved backwards).
func (g *Generator) MustNextID() int64 {
	id, err := g.NextID()
	if err != nil {
		panic(err)
	}
	return id
}

// currentMillis returns the current time in milliseconds since Unix epoch.
func (g *Generator) currentMillis() int64 {
	return g.timeFunc().UnixMilli()
}

// waitNextMillis blocks until the next millisecond.
func (g *Generator) waitNextMillis(lastTimestamp int64) int64 {
	timestamp := g.currentMillis()
	for timestamp <= lastTimestamp {
		time.Sleep(100 * time.Microsecond)
		timestamp = g.currentMillis()
	}
	return timestamp
}

// Decompose extracts the timestamp, node ID, and sequence from a snowflake ID.
func Decompose(id int64) (timestamp time.Time, nodeID int64, sequence int64) {
	timestampMs := (id >> timestampShift) + epoch
	nodeID = (id >> nodeShift) & maxNodeID
	sequence = id & maxSequence
	timestamp = time.UnixMilli(timestampMs)
	return
}
