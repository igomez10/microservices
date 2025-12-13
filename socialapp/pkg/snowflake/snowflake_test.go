package snowflake

import (
	"sync"
	"testing"
	"time"
)

func TestNewGenerator_ValidNodeID(t *testing.T) {
	tests := []struct {
		name   string
		nodeID int64
	}{
		{"min node ID", 0},
		{"mid node ID", 512},
		{"max node ID", 1023},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGenerator(tt.nodeID)
			if err != nil {
				t.Errorf("NewGenerator(%d) returned error: %v", tt.nodeID, err)
			}
			if g == nil {
				t.Errorf("NewGenerator(%d) returned nil generator", tt.nodeID)
			}
		})
	}
}

func TestNewGenerator_InvalidNodeID(t *testing.T) {
	tests := []struct {
		name   string
		nodeID int64
	}{
		{"negative node ID", -1},
		{"too large node ID", 1024},
		{"way too large", 10000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGenerator(tt.nodeID)
			if err != ErrInvalidNodeID {
				t.Errorf("NewGenerator(%d) expected ErrInvalidNodeID, got: %v", tt.nodeID, err)
			}
			if g != nil {
				t.Errorf("NewGenerator(%d) expected nil generator, got: %v", tt.nodeID, g)
			}
		})
	}
}

func TestNextID_Uniqueness(t *testing.T) {
	g, err := NewGenerator(1)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	seen := make(map[int64]bool)
	count := 10000

	for i := 0; i < count; i++ {
		id, err := g.NextID()
		if err != nil {
			t.Fatalf("NextID() returned error: %v", err)
		}
		if seen[id] {
			t.Errorf("Duplicate ID generated: %d", id)
		}
		seen[id] = true
	}

	if len(seen) != count {
		t.Errorf("Expected %d unique IDs, got %d", count, len(seen))
	}
}

func TestNextID_Monotonic(t *testing.T) {
	g, err := NewGenerator(1)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	var lastID int64
	count := 10000

	for i := 0; i < count; i++ {
		id, err := g.NextID()
		if err != nil {
			t.Fatalf("NextID() returned error: %v", err)
		}
		if id <= lastID {
			t.Errorf("IDs not monotonically increasing: %d <= %d", id, lastID)
		}
		lastID = id
	}
}

func TestNextID_Concurrent(t *testing.T) {
	g, err := NewGenerator(1)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	var wg sync.WaitGroup
	idChan := make(chan int64, 10000)
	goroutines := 10
	idsPerGoroutine := 1000

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				id, err := g.NextID()
				if err != nil {
					t.Errorf("NextID() returned error: %v", err)
					return
				}
				idChan <- id
			}
		}()
	}

	wg.Wait()
	close(idChan)

	seen := make(map[int64]bool)
	for id := range idChan {
		if seen[id] {
			t.Errorf("Duplicate ID generated in concurrent test: %d", id)
		}
		seen[id] = true
	}

	expectedCount := goroutines * idsPerGoroutine
	if len(seen) != expectedCount {
		t.Errorf("Expected %d unique IDs, got %d", expectedCount, len(seen))
	}
}

func TestNextID_Deterministic(t *testing.T) {
	// Use a fixed time function to make ID generation deterministic
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	timeFunc := func() time.Time { return fixedTime }

	g1, err := NewGeneratorWithTime(42, timeFunc)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	g2, err := NewGeneratorWithTime(42, timeFunc)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	// First IDs from both generators should be identical
	id1, _ := g1.NextID()
	id2, _ := g2.NextID()

	if id1 != id2 {
		t.Errorf("Deterministic IDs should match: %d != %d", id1, id2)
	}

	// Verify the ID contains expected components
	timestamp, nodeID, sequence := Decompose(id1)

	if nodeID != 42 {
		t.Errorf("Expected nodeID 42, got %d", nodeID)
	}
	if sequence != 0 {
		t.Errorf("Expected sequence 0, got %d", sequence)
	}
	if !timestamp.Equal(fixedTime) {
		t.Errorf("Expected timestamp %v, got %v", fixedTime, timestamp)
	}
}

func TestNextID_DifferentNodes(t *testing.T) {
	// Same time, different nodes should produce different IDs
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	timeFunc := func() time.Time { return fixedTime }

	g1, _ := NewGeneratorWithTime(1, timeFunc)
	g2, _ := NewGeneratorWithTime(2, timeFunc)

	id1, _ := g1.NextID()
	id2, _ := g2.NextID()

	if id1 == id2 {
		t.Errorf("Different nodes should produce different IDs: %d == %d", id1, id2)
	}

	_, nodeID1, _ := Decompose(id1)
	_, nodeID2, _ := Decompose(id2)

	if nodeID1 != 1 {
		t.Errorf("Expected nodeID 1, got %d", nodeID1)
	}
	if nodeID2 != 2 {
		t.Errorf("Expected nodeID 2, got %d", nodeID2)
	}
}

func TestDecompose(t *testing.T) {
	g, _ := NewGenerator(100)
	id, _ := g.NextID()

	timestamp, nodeID, sequence := Decompose(id)

	if nodeID != 100 {
		t.Errorf("Expected nodeID 100, got %d", nodeID)
	}
	if sequence != 0 {
		t.Errorf("Expected sequence 0 for first ID, got %d", sequence)
	}
	if timestamp.IsZero() {
		t.Error("Expected non-zero timestamp")
	}
	// Timestamp should be recent (within last minute)
	if time.Since(timestamp) > time.Minute {
		t.Errorf("Timestamp seems too old: %v", timestamp)
	}
}

func TestMustNextID_Success(t *testing.T) {
	g, _ := NewGenerator(1)
	
	// Should not panic
	id := g.MustNextID()
	if id == 0 {
		t.Error("MustNextID() returned 0")
	}
}

func TestSequenceIncrement(t *testing.T) {
	// Use fixed time to force sequence increment
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	timeFunc := func() time.Time { return fixedTime }

	g, _ := NewGeneratorWithTime(1, timeFunc)

	// Generate multiple IDs in the same "millisecond"
	id1, _ := g.NextID()
	id2, _ := g.NextID()
	id3, _ := g.NextID()

	_, _, seq1 := Decompose(id1)
	_, _, seq2 := Decompose(id2)
	_, _, seq3 := Decompose(id3)

	if seq1 != 0 {
		t.Errorf("Expected first sequence 0, got %d", seq1)
	}
	if seq2 != 1 {
		t.Errorf("Expected second sequence 1, got %d", seq2)
	}
	if seq3 != 2 {
		t.Errorf("Expected third sequence 2, got %d", seq3)
	}
}

func BenchmarkNextID(b *testing.B) {
	g, _ := NewGenerator(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = g.NextID()
	}
}

func BenchmarkNextID_Concurrent(b *testing.B) {
	g, _ := NewGenerator(1)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = g.NextID()
		}
	})
}

