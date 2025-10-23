# Step 3: Parallel Tree Comparison (Unbounded)

## Implementation Complete ✅

### What Was Implemented

**`CompareTreesParallelUnbounded`** - Spawns one goroutine per tree comparison

```go
// For each pair of trees with same hash:
for i := 0; i < len(hashGroup); i++ {
    for j := i + 1; j < len(hashGroup); j++ {
        wg.Add(1)
        go func(id1, id2 int) {
            defer wg.Done()
            if AreEqual(bsts[id1], bsts[id2]) {
                mu.Lock()
                adjMatrix[id1][id2] = true
                adjMatrix[id2][id1] = true
                mu.Unlock()
            }
        }(hashGroup[i], hashGroup[j])
    }
}
```

### Key Features

1. **Adjacency Matrix**: Tracks equivalence relationships
   - `adjMatrix[i][j] = true` means BST i equals BST j
   - Initialized with diagonal = true (tree equals itself)
   - Updated symmetrically when match found

2. **Maximum Parallelism**: Spawns goroutine for each comparison
   - Simple implementation
   - No work distribution needed
   - All comparisons run simultaneously

3. **Thread Safety**: Mutex protects adjacency matrix
   - Only locks when writing to matrix
   - Comparison happens outside lock (parallel)
   - No race conditions

4. **Connected Components**: Uses BFS to build equivalence groups
   - Same algorithm as sequential version
   - Finds all transitively equivalent trees
   - Filters groups with only 1 tree

## Performance Results

### Simple.txt (10 Trees)
```
Sequential:           0.000056s
Parallel Unbounded:   0.000137s (2.4x SLOWER!)
```
**Analysis:** Overhead dominates for tiny workload

### Coarse.txt (100 Large Trees)
```
Sequential:           3.716667s
Parallel Unbounded:   1.639896s (2.27x FASTER ⭐)
```
**Analysis:** Significant speedup on large trees

### Expected Performance on Fine.txt
For fine.txt with many small trees, parallel unbounded would spawn **many goroutines**:
- Depends on hash collision patterns
- Could be thousands of goroutines
- Overhead may be significant

## Usage

### Test Unbounded Strategy
```bash
# Sequential (baseline)
./BST -hash-workers=8 -data-workers=8 -comp-workers=1 -input=input/coarse.txt

# Parallel unbounded
./BST -hash-workers=8 -data-workers=8 -comp-workers=2 -comp-strategy=unbounded -input=input/coarse.txt
```

Note: `-comp-workers` value is ignored for unbounded strategy (spawns one goroutine per comparison)

### Race Detection
```bash
go run -race BST.go -hash-workers=8 -data-workers=8 -comp-workers=2 -comp-strategy=unbounded -input=input/simple.txt
```
✅ Passes with no race conditions

## Implementation Details

### 1. Adjacency Matrix Creation
```go
n := len(bsts)
adjMatrix := make([][]bool, n)
for i := 0; i < n; i++ {
    adjMatrix[i] = make([]bool, n)
    adjMatrix[i][i] = true  // Self-equivalence
}
```

### 2. Parallel Comparison
```go
var mu sync.Mutex
var wg sync.WaitGroup

for _, hashGroup := range hashGroups {
    if len(hashGroup) > 1 {
        for i := 0; i < len(hashGroup); i++ {
            for j := i + 1; j < len(hashGroup); j++ {
                wg.Add(1)
                go func(id1, id2 int) {
                    defer wg.Done()
                    if AreEqual(bsts[id1], bsts[id2]) {
                        mu.Lock()
                        adjMatrix[id1][id2] = true
                        adjMatrix[id2][id1] = true
                        mu.Unlock()
                    }
                }(hashGroup[i], hashGroup[j])
            }
        }
    }
}
wg.Wait()
```

### 3. Build Equivalence Groups (BFS)
```go
visited := make([]bool, n)
var equivalenceGroups [][]int

for i := 0; i < n; i++ {
    if !visited[i] {
        group := []int{}
        queue := []int{i}
        visited[i] = true
        
        for len(queue) > 0 {
            current := queue[0]
            queue = queue[1:]
            group = append(group, current)
            
            for j := 0; j < n; j++ {
                if adjMatrix[current][j] && !visited[j] {
                    visited[j] = true
                    queue = append(queue, j)
                }
            }
        }
        
        if len(group) > 1 {
            equivalenceGroups = append(equivalenceGroups, group)
        }
    }
}
```

## Advantages vs Sequential

### ✅ Pros
- **Simple**: No work distribution logic
- **Maximum parallelism**: All comparisons simultaneous
- **Good speedup**: 2.27x on coarse.txt
- **Scales naturally**: More comparisons = more parallelism

### ⚠️ Cons
- **High overhead**: Small workloads slower
- **Many goroutines**: Could spawn thousands
- **No control**: Can't tune worker count
- **Memory**: Each goroutine uses ~2KB stack

## Comparison: Unbounded vs Sequential

| Aspect | Sequential | Unbounded |
|--------|-----------|-----------|
| **Simplicity** | ✅ Simple | ✅ Simple |
| **Performance (coarse)** | 3.72s | 1.64s (2.27x) |
| **Performance (simple)** | 0.000056s | 0.000137s (2.4x slower) |
| **Goroutines** | 0 | N comparisons |
| **Overhead** | None | Goroutine creation |
| **Tunability** | N/A | ❌ None |
| **Scalability** | Poor | Good for large tasks |

## When to Use Unbounded

**Good for:**
- ✅ Large tree comparisons (like coarse.txt)
- ✅ Moderate number of comparisons (<10,000)
- ✅ When simplicity matters
- ✅ I/O-bound or long-running comparisons

**Not good for:**
- ❌ Many small comparisons
- ❌ Tens of thousands of comparisons
- ❌ Need predictable performance
- ❌ Limited memory

## Next Step: Worker Pool

The assignment mentions implementing a **worker pool** version next:
- Fixed number of workers (e.g., 4 or 8)
- Workers process comparisons from channel
- More controlled overhead
- Better for fine-grained tasks

This will allow comparison of:
1. Sequential (baseline)
2. Unbounded (maximum parallelism)
3. Worker pool (controlled parallelism)

## Correctness Verification

✅ **Same output as sequential:**
- All three hash groups found correctly
- All equivalence groups match
- No missing or extra trees

✅ **Thread-safe:**
- Passes race detector
- Mutex protects shared state
- No data races

✅ **Transitive closure:**
- BFS correctly finds connected components
- If A=B and B=C, then A=C in same group

## Summary

**Parallel unbounded implementation is complete and working!**

- ✅ 2.27x speedup on coarse.txt
- ✅ Correct results
- ✅ Thread-safe
- ✅ Simple to implement
- ⚠️ Overhead for small workloads
- 📝 Ready for comparison with worker pool approach

The next step is to implement `CompareTreesParallelPool` for more controlled parallelism, which will likely perform better on workloads with many small comparisons.

