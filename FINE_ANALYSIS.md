# Fine-Grained Analysis: 100k BSTs Performance

## Test Results (fine.txt)

### Worker Pool Strategy
```
Workers=1:     0.054227s  (sequential baseline)
Workers=8:     0.031282s  ⭐ FASTEST! (1.73x speedup)
Workers=64:    0.041937s  (1.29x speedup)
Workers=128:   0.039738s  (1.36x speedup)
Workers=1024:  0.044120s  (1.23x speedup)
Workers=10000: 0.055374s  (0.98x - SLOWER than sequential!)
```

### Per-BST Strategy (100k goroutines)
```
~100k goroutines: 0.047206s  (1.15x speedup)
```

## Key Insights

### 1. Sweet Spot: 8 Workers (Best Performance)
- **0.031282s** - Nearly 2x faster than sequential
- Close to CPU core count (optimal utilization)
- Minimal overhead, maximum throughput

### 2. Performance Degradation After 8 Workers
```
Workers    Time        Why?
--------   --------    ----------------------------------
8          0.031282s   ✅ Optimal: matches CPU cores
64         0.041937s   ⚠️  Context switching overhead starts
128        0.039738s   ⚠️  More overhead
1024       0.044120s   ⚠️  Significant scheduling overhead
10000      0.055374s   ❌ Worse than sequential!
```

**Why does performance degrade?**
- **Context switching**: OS must juggle 10k goroutines across ~8 CPU cores
- **Memory overhead**: Each goroutine has stack space (even if small)
- **Scheduler overhead**: Go runtime spends more time scheduling than computing
- **Cache thrashing**: Too many goroutines destroy CPU cache locality

### 3. Per-BST Strategy: Surprisingly Good!
```
100k goroutines: 0.047206s (middle of the pack)
```

**Observations:**
- ✅ Go handles 100k goroutines better than expected!
- ⚠️  Still 34% slower than optimal worker pool (8 workers)
- ⚠️  But faster than over-provisioned worker pools (1024, 10000)

**Why it's competitive:**
- Go's goroutines are lightweight (~2KB stack)
- Goroutine scheduling is efficient
- Small BSTs = short-lived goroutines

### 4. The "Too Many Workers" Problem
```
10000 workers: 0.055374s
Sequential:    0.054227s
Overhead:      ~2% SLOWER!
```

**This proves:** You **CANNOT** just spawn infinite goroutines!
- Go manages them well, but overhead exists
- Need to match workload to resources

## Performance Visualization

```
Time (seconds)
0.055 |                                              ▓ (10k workers)
      |                                  ▓ (1 worker)
0.050 |                           ▓ (100k goroutines)
      |                     ▓ (1024)
0.045 |               ▓ (64)
      |           ▓ (128)
0.040 |
      |     ▓ (8 workers) ⭐ OPTIMAL
0.035 |
0.030 |________________________________________________
      1    8    64   128  1024  10k   100k
                  Worker Count
```

## Answers to Assignment Questions

### Q1: Which implementation is faster (and by how much)?

**Answer:** Worker Pool with **8 workers**
- **1.73x faster** than sequential
- **1.51x faster** than per-BST (100k goroutines)
- **1.77x faster** than over-provisioned pool (10k workers)

### Q2: Can Go manage goroutines well enough that you don't have to worry anymore?

**Answer:** **NO - but it's impressive!**

**Evidence:**
1. ✅ Go handles 100k goroutines reasonably (only 34% slower than optimal)
2. ❌ But you still need to think about parallelism:
   - 8 workers: 0.031s ⭐
   - 100k goroutines: 0.047s (50% slower)
   - 10k workers: 0.055s (77% slower!)

3. ✅ Go makes it **easier** than traditional threads (try creating 100k pthreads!)
4. ❌ But **controlled parallelism** (worker pool) still wins

**The Practical Truth:**
- Goroutines are cheap, but not free
- Overhead scales with goroutine count
- You still need to match parallelism to hardware
- Sweet spot ≈ CPU core count for CPU-bound tasks

## Recommendations for Your Writeup

### Performance Graph Data
```
Workers   Time(s)   Speedup   Notes
1         0.054227  1.00x     Sequential baseline
8         0.031282  1.73x     ⭐ Optimal (≈ CPU cores)
64        0.041937  1.29x     Overhead increasing
128       0.039738  1.36x     
1024      0.044120  1.23x     Heavy overhead
10000     0.055374  0.98x     ❌ Worse than seq!
100k*     0.047206  1.15x     *Per-BST strategy
```

### Discussion Points

1. **Worker Pool is Superior**
   - Best performance at 8 workers (likely matches CPU core count)
   - Tunable: can optimize for different hardware
   - Predictable behavior

2. **Per-BST is Acceptable but Suboptimal**
   - Simple to implement
   - Go handles 100k goroutines without crashing
   - But 50% slower than optimal
   - Not tunable - workload-dependent

3. **Goroutine Overhead is Real**
   - 10k workers slower than sequential!
   - Context switching dominates at high counts
   - Cache locality destroyed
   - Go scheduler becomes bottleneck

4. **Lesson: Controlled Parallelism Wins**
   - Can't just "throw goroutines at the problem"
   - Need to match parallelism to hardware
   - Worker pool pattern is more practical
   - But Go makes experimentation easy!

## Hardware Context

Your results likely come from a system with **~8 CPU cores** (or 4 cores with hyperthreading).

**Evidence:**
- Optimal performance at 8 workers
- Performance degrades beyond that
- Diminishing returns after matching core count

Run this to confirm:
```bash
sysctl -n hw.ncpu  # macOS
# or
nproc              # Linux
```

## Conclusion

Go's goroutines are **impressive** but not magic:
- ✅ Can handle 100k goroutines (amazing!)
- ❌ But controlled parallelism is still better
- ✅ Much better than traditional threads
- ❌ Still need to think about performance

**Best Practice:** Use worker pool pattern matching CPU core count for CPU-bound tasks.

