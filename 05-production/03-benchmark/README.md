# ⚡ Section 05-03 — Benchmark: Measure Kruda Application Performance

⏱️ Estimated time: **30 minutes**

Welcome to the final lesson of the Tutorial! In this lesson you will learn how to **benchmark** a Kruda application using Go's built-in benchmarking tool -- measuring handler throughput, JSON serialisation speed, and concurrent load performance to ensure your application is production-ready.

---

## Target Learning Outcomes

- Understand the concept of **Go benchmarking** and why it matters for production applications
- Write benchmark functions with `func Benchmark*(b *testing.B)`
- Measure handler execution throughput (ns/op) using `kruda.NewTestClient`
- Measure memory allocation (B/op) and heap allocation count (allocs/op)
- Compare handler throughput with raw JSON serialisation
- Run parallel benchmarks to simulate concurrent load

---

## What You Will Learn

By the end of this lesson you will be able to:

- Write Go benchmark functions in `*_test.go` files
- Run benchmarks with `go test -bench=. -benchmem`
- Read benchmark results: **ns/op**, **B/op**, **allocs/op**
- Use `b.ResetTimer()` to exclude setup time from measurements
- Write parallel benchmarks with `b.RunParallel()`
- Compare results with `benchstat`
- Identify bottlenecks from memory allocation patterns
- Use `kruda.NewTestClient` to benchmark the full HTTP handler stack

---

## Prerequisites

| Tool | Version |
|---|---|
| Go | 1.25+ |
| Git | Latest |
| Text Editor / IDE | VS Code, GoLand, or your preferred editor |

> If you haven't completed Section 05-02, consider going back first -- see [Section 05-02 -- Docker Deploy](../02-docker-deploy/)

---

## File Structure

```
05-production/03-benchmark/
|-- README.md              <-- You are here
|-- starter/               <-- Starter code (with TODOs to fill in)
|   |-- go.mod
|   +-- main.go
+-- complete/              <-- Complete solution
    |-- go.mod
    |-- main.go
    +-- benchmark_test.go  <-- Go benchmark file
```

- **[starter/](./starter/)** -- Skeleton code that compiles but has `// TODO:` markers for you to complete
- **[complete/](./complete/)** -- Full working solution with benchmark tests

---

## Why Benchmark?

| Question | Benchmark That Answers It |
|---|---|
| How fast are our handlers? | `BenchmarkCreateItem` -> ns/op |
| How much memory per request? | `-benchmem` -> B/op |
| How many heap allocations per request? | `-benchmem` -> allocs/op |
| How well does it handle concurrent load? | `BenchmarkCreateItem_Parallel` |
| Is JSON serialisation a bottleneck? | Compare handler vs raw JSON benchmarks |

> Go has a built-in benchmarking tool -- just write `func Benchmark*(b *testing.B)` in a `*_test.go` file and run `go test -bench=.`

---

## Step-by-Step Guide

### Step 1: Open the starter project

```bash
cd 05-production/03-benchmark/starter
```

Open `main.go` -- you will see the skeleton with `// TODO:` comments for handlers and app setup.

### Step 2: Implement the Application Factory

The `NewApp` function creates a configured Kruda app that both `main()` and benchmark tests can use:

```go
func NewApp(store *ItemStore) *kruda.App {
    app := kruda.New()

    kruda.Get[struct{}, MessageResponse](app, "/health", func(c *kruda.C[struct{}]) (*MessageResponse, error) {
        return &MessageResponse{Message: "ok"}, nil
    })

    kruda.Get[struct{}, []ItemResponse](app, "/items", func(c *kruda.C[struct{}]) (*[]ItemResponse, error) {
        items := store.All()
        return &items, nil
    })

    kruda.Post[CreateItemInput, ItemResponse](app, "/items", func(c *kruda.C[CreateItemInput]) (*ItemResponse, error) {
        item := store.Create(c.In.Name, c.In.Price)
        return &item, nil
    })

    app.Compile()
    return app
}
```

> `app.Compile()` is required for `kruda.NewTestClient` to work in benchmarks

### Step 3: Create the Benchmark Test File

Create `benchmark_test.go` in the same directory as `main.go`:

```go
package main

import (
    "encoding/json"
    "testing"

    "github.com/go-kruda/kruda"
)
```

### Step 4: Write Benchmarks Using TestClient

The `kruda.NewTestClient` lets you exercise the full HTTP handler stack without starting a real server:

```go
func BenchmarkListItems(b *testing.B) {
    store := NewItemStore()
    store.Create("test", 100)
    app := NewApp(store)
    client := kruda.NewTestClient(app)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        client.Get("/items")
    }
}

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
```

> `b.ResetTimer()` resets the timer after setup is complete -- ensuring results measure only handler execution, not setup time

### Step 5: Write JSON Serialisation Benchmarks

```go
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
```

> Compare ns/op of this benchmark with `BenchmarkCreateItem` -- the difference is the overhead of Kruda's handler system

### Step 6: Write Parallel Benchmarks

```go
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
```

> `b.RunParallel()` runs the benchmark with multiple goroutines simultaneously -- simulating concurrent requests like in production

### Step 7: Run Benchmarks

```bash
go test -bench=. -benchmem
```

You will see output like this:

```
goos: linux
goarch: amd64
pkg: github.com/go-kruda/tutorial/05-production/03-benchmark/complete
BenchmarkListItems-8                   200000     6120 ns/op   4096 B/op   12 allocs/op
BenchmarkCreateItem-8                  500000     2340 ns/op    512 B/op    8 allocs/op
BenchmarkHealthCheck-8                2000000      580 ns/op    128 B/op    3 allocs/op
BenchmarkJSONMarshalItemResponse-8    3000000      420 ns/op     64 B/op    1 allocs/op
BenchmarkJSONMarshalItemSlice-8        100000    15200 ns/op   8192 B/op    2 allocs/op
BenchmarkJSONUnmarshalCreateInput-8   2000000      890 ns/op    256 B/op    4 allocs/op
BenchmarkListItems_Parallel-8          500000     3100 ns/op   4096 B/op   12 allocs/op
BenchmarkCreateItem_Parallel-8        1000000     1200 ns/op    512 B/op    8 allocs/op
PASS
ok      github.com/go-kruda/tutorial/05-production/03-benchmark/complete    12.345s
```

---

## How to Read Benchmark Results

### Key Metrics

| Metric | Meaning | Lower is better? |
|---|---|---|
| **ns/op** | Average time per operation (nanoseconds) | Yes |
| **B/op** | Bytes allocated per operation | Yes |
| **allocs/op** | Heap allocations per operation | Yes |
| **N** (iterations) | Number of iterations Go ran | More = more accurate |

### Example Analysis

```
BenchmarkCreateItem-8   500000   2340 ns/op   512 B/op   8 allocs/op
```

- **2340 ns/op** = ~427,000 operations/second on a single core
- **512 B/op** = each request uses 512 bytes of memory
- **8 allocs/op** = 8 heap allocations per request

### Comparing Sequential vs Parallel

```
BenchmarkCreateItem-8            500000   2340 ns/op
BenchmarkCreateItem_Parallel-8  1000000   1200 ns/op
```

> The parallel benchmark is faster because it uses multiple CPU cores. However, if the parallel ns/op is much higher than sequential, that indicates lock contention.

---

## Advanced Techniques

### Run Benchmarks Multiple Times for Accuracy

```bash
go test -bench=. -benchmem -count=5
```

### Use benchstat to Compare Results

```bash
# Install benchstat
go install golang.org/x/perf/cmd/benchstat@latest

# Run benchmark before changes
go test -bench=. -benchmem -count=5 > old.txt

# Make changes...

# Run benchmark after changes
go test -bench=. -benchmem -count=5 > new.txt

# Compare
benchstat old.txt new.txt
```

### Run Specific Benchmarks

```bash
# Run only benchmarks with "Create" in the name
go test -bench=Create -benchmem

# Run only parallel benchmarks
go test -bench=Parallel -benchmem
```

### Set Minimum Benchmark Duration

```bash
# Run each benchmark for at least 3 seconds
go test -bench=. -benchmem -benchtime=3s
```

---

## Compare with complete/

If you get stuck, check the solution in **[complete/](./complete/)** and compare:

```bash
diff starter/main.go complete/main.go
```

---

## Key Concepts Summary

| Concept | Description |
|---|---|
| `func Benchmark*(b *testing.B)` | Benchmark functions must start with `Benchmark` |
| `b.N` | Number of iterations Go determines automatically |
| `b.ResetTimer()` | Reset timer after setup to measure only the code of interest |
| `b.RunParallel()` | Run benchmark with multiple goroutines simultaneously |
| `kruda.NewTestClient(app)` | Create a test client to benchmark the full HTTP handler stack |
| `client.Get(path)` | Send a GET request through the test client |
| `client.Post(path, body)` | Send a POST request through the test client |
| `app.Compile()` | Required before using NewTestClient |
| `-bench=.` | Run all benchmarks (regex pattern) |
| `-benchmem` | Show memory allocation statistics |
| `-count=5` | Run benchmark 5 times for accuracy |
| `benchstat` | Tool for comparing benchmark results |

---

## Congratulations! You have completed the Tutorial!

You have learned everything in the Kruda Tutorial from beginner to production!

What you learned across the tutorial:

- **Beginner** -- REST API + Typed Handlers
- **Auto CRUD** -- Generate CRUD endpoints automatically
- **Intermediate** -- Database, Config, Error Handling
- **Advanced** -- DI Container, Auth, OpenAPI, SSE, MCP, WebSocket, Testing, Architecture
- **Production** -- Monitoring, Docker Deploy, Benchmark

Happy building with Kruda!

--> [Back to main page](../../)
