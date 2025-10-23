# Hash Computation Strategies

## Two Implementations

### 1. **Per-BST Strategy** (`-hash-strategy=perbst`)
**Approach:** Spawns **one goroutine per BST**

```go
func ComputeHashesParallelPerBST(bsts []*BST) map[int]int {
    for _, bst := range bsts {
        wg.Add(1)
        go func(b *BST) {
            defer wg.Done()
            hash := b.ComputeHash()
            // Store result
        }(bst)
    }
    wg.Wait()
}
```

**Characteristics:**
- ✅ Simple: No work distribution logic needed
- ✅ Maximum parallelism: All BSTs computed simultaneously
- ⚠️ High overhead: Creates many goroutines (especially for fine.txt with 100k BSTs)
- ⚠️ Ignores `-hash-workers` flag

**Best for:** Coarse-grained parallelism (few large trees)

### 2. **Worker Pool Strategy** (`-hash-strategy=pool`) [DEFAULT]
**Approach:** Spawns **fixed number of worker goroutines** that share the work

```go
func ComputeHashesParallelWorkerPool(bsts []*BST, numWorkers int) map[int]int {
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            // Process assigned range of BSTs
            for j := start; j < end; j++ {
                hash := bsts[j].ComputeHash()
                // Store result
            }
        }(i)
    }
    wg.Wait()
}
```

**Characteristics:**
- ✅ Controlled overhead: Creates exactly `numWorkers` goroutines
- ✅ Scalable: Works well with any number of BSTs
- ✅ Tunable: Can experiment with `-hash-workers` values
- ⚠️ Requires work distribution logic

**Best for:** Both coarse and fine-grained parallelism

## Usage

### Test Per-BST Strategy
```bash
./BST -hash-workers=4 -hash-strategy=perbst -input=input/coarse.txt
```

### Test Worker Pool Strategy (Default)
```bash
./BST -hash-workers=4 -hash-strategy=pool -input=input/coarse.txt
# Or simply:
./BST -hash-workers=4 -input=input/coarse.txt
```

### Compare Both Strategies
```bash
./compare_hash_strategies.sh simple   # Small test
./compare_hash_strategies.sh coarse   # Large trees
./compare_hash_strategies.sh fine     # Many small trees (WARNING: slow!)
```

## Performance Observations (coarse.txt)

### Worker Pool Strategy:
```
Workers=1:  0.143592s  (sequential baseline)
Workers=2:  0.063846s  (2.25x speedup)
Workers=4:  0.051180s  (2.81x speedup)
Workers=8:  0.048603s  (2.95x speedup)
Workers=16: 0.052693s  (2.73x speedup - overhead increases)
```

### Per-BST Strategy:
```
All runs:   ~0.047-0.092s
Note: Creates ~100 goroutines (one per BST in coarse.txt)
```

## Key Insights

1. **Per-BST is simpler but less flexible**
   - No tuning needed, but can't control goroutine count
   - Good when number of BSTs matches CPU cores

2. **Worker Pool is more practical**
   - Can tune for optimal performance
   - Better for varying workloads
   - Prevents goroutine explosion on fine.txt

3. **Go manages goroutines well, but overhead exists**
   - Creating thousands of goroutines isn't free
   - Worker pool shows you still need to think about parallelism
   - Sweet spot for coarse.txt: 4-8 workers

4. **Diminishing returns after ~8 workers**
   - Beyond CPU core count, overhead dominates
   - Context switching and synchronization costs increase

## Assignment Questions to Consider

1. **Which implementation is faster?**
   - Depends on workload! Per-BST is competitive for coarse.txt
   - Worker pool wins for fine.txt due to controlled overhead

2. **By how much?**
   - See comparison script results above
   - Speedup varies with worker count

3. **Can Go manage goroutines well enough that you don't have to worry?**
   - **Answer: Not quite!** While Go is good at managing goroutines:
     - Creating 100k goroutines (fine.txt with per-BST) has significant overhead
     - Worker pool approach still shows better scalability
     - You still need to think about work distribution
     - But Go makes it much easier than traditional threads!

## Implementation Details

### Thread Safety
Both implementations use `sync.Mutex` to protect the shared hash map:

```go
mu.Lock()
hashes[b.ID] = hash
mu.Unlock()
```

**Alternative:** Could use channels for result collection (implemented in Step 2).

### Work Distribution (Worker Pool)
```go
bstsPerWorker := len(bsts) / numWorkers
start := workerID * bstsPerWorker
end := start + bstsPerWorker

// Last worker takes remainder
if workerID == numWorkers-1 {
    end = len(bsts)
}
```

### Goroutine Pattern
```go
go func(workerID int) {
    defer wg.Done()  // Always decrement counter
    // Do work
}(i)  // Pass by value to avoid closure issues
```

