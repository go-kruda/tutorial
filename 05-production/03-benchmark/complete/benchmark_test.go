package main

import (
	"encoding/json"
	"testing"

	"github.com/go-kruda/kruda"
)

// ============================================================
// Benchmark Tests for Kruda Application
// ============================================================
//
// Run these benchmarks with:
//
//   go test -bench=. -benchmem
//
// Understanding the output:
//
//   BenchmarkCreateItem-8   500000   2340 ns/op   512 B/op   8 allocs/op
//   |                   |   |        |            |           |
//   |                   |   |        |            |           +-- heap allocations per op
//   |                   |   |        |            +-- bytes allocated per op
//   |                   |   |        +-- nanoseconds per operation
//   |                   |   +-- number of iterations run
//   |                   +-- GOMAXPROCS
//   +-- benchmark function name
//
// Key metrics to watch:
//   - ns/op    -- lower is better; measures handler throughput
//   - B/op     -- lower is better; measures memory pressure
//   - allocs/op -- lower is better; fewer allocs = less GC pressure
//
// Tips for reliable benchmarks:
//   - Close other applications to reduce noise
//   - Run multiple times: go test -bench=. -benchmem -count=5
//   - Use benchstat to compare results:
//     go install golang.org/x/perf/cmd/benchstat@latest
//   - Pin CPU frequency if possible (disable turbo boost)

// ----------------------------------------------------------------
// Benchmark: HTTP Handler Execution via TestClient
// ----------------------------------------------------------------

// BenchmarkListItems measures the throughput of the list-items
// endpoint. The store is pre-populated with 100 items to simulate
// a realistic workload.
//
// This benchmark answers: "How fast can Kruda serialise a list
// of items to JSON and deliver the HTTP response?"
func BenchmarkListItems(b *testing.B) {
	store := NewItemStore()
	for i := 0; i < 100; i++ {
		store.Create("Item", i*10)
	}
	app := NewApp(store)
	client := kruda.NewTestClient(app)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Get("/items")
	}
}

// BenchmarkCreateItem measures the throughput of the create-item
// endpoint including JSON deserialisation by Kruda and the
// handler's business logic.
//
// This benchmark answers: "How many create requests can Kruda
// handle per second on a single core?"
func BenchmarkCreateItem(b *testing.B) {
	store := NewItemStore()
	app := NewApp(store)
	client := kruda.NewTestClient(app)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Post("/items", map[string]interface{}{
			"name":  "Benchmark Item",
			"price": 999,
		})
	}
}

// BenchmarkHealthCheck measures the throughput of the simplest
// possible handler -- useful as a baseline to understand the
// framework overhead of Kruda's handler system.
//
// This benchmark answers: "What is the minimum overhead of a
// Kruda HTTP handler?"
func BenchmarkHealthCheck(b *testing.B) {
	store := NewItemStore()
	app := NewApp(store)
	client := kruda.NewTestClient(app)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Get("/health")
	}
}

// ----------------------------------------------------------------
// Benchmark: JSON Serialisation (baseline comparison)
// ----------------------------------------------------------------

// BenchmarkJSONMarshalItemResponse measures raw JSON serialisation
// speed for a single ItemResponse. Compare this with
// BenchmarkCreateItem to see how much overhead the handler layer
// adds on top of pure JSON encoding.
func BenchmarkJSONMarshalItemResponse(b *testing.B) {
	item := ItemResponse{
		ID:    1,
		Name:  "Benchmark Item",
		Price: 999,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(item)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkJSONMarshalItemSlice measures JSON serialisation speed
// for a slice of 100 items. This isolates the JSON encoding cost
// from the handler logic measured in BenchmarkListItems.
func BenchmarkJSONMarshalItemSlice(b *testing.B) {
	items := make([]ItemResponse, 100)
	for i := range items {
		items[i] = ItemResponse{
			ID:    i + 1,
			Name:  "Item",
			Price: i * 10,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(items)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkJSONUnmarshalCreateInput measures raw JSON
// deserialisation speed for a CreateItemInput. Compare with
// BenchmarkCreateItem to understand how much of the handler
// time is spent parsing JSON vs. business logic.
func BenchmarkJSONUnmarshalCreateInput(b *testing.B) {
	data := []byte(`{"name":"Benchmark Item","price":999}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var req CreateItemInput
		if err := json.Unmarshal(data, &req); err != nil {
			b.Fatal(err)
		}
	}
}

// ----------------------------------------------------------------
// Benchmark: Parallel Handler Execution
// ----------------------------------------------------------------

// BenchmarkListItems_Parallel measures read throughput under
// concurrent load with a pre-populated store.
//
// Compare ns/op with BenchmarkListItems to see how well the
// application scales across CPU cores.
func BenchmarkListItems_Parallel(b *testing.B) {
	store := NewItemStore()
	for i := 0; i < 100; i++ {
		store.Create("Item", i*10)
	}
	app := NewApp(store)
	client := kruda.NewTestClient(app)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			client.Get("/items")
		}
	})
}

// BenchmarkCreateItem_Parallel measures handler throughput under
// concurrent load. This uses b.RunParallel to simulate multiple
// goroutines hitting the handler simultaneously -- closer to
// real-world server behaviour.
//
// Compare ns/op with BenchmarkCreateItem to see the impact of
// mutex contention in the ItemStore.
func BenchmarkCreateItem_Parallel(b *testing.B) {
	store := NewItemStore()
	app := NewApp(store)
	client := kruda.NewTestClient(app)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			client.Post("/items", map[string]interface{}{
				"name":  "Parallel Item",
				"price": 500,
			})
		}
	})
}
