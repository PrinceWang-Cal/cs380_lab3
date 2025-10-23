# Step 1 Final Analysis: Per-BST vs Worker Pool

## Assignment Questions & Answers

### Q1: Which implementation is faster (and by how much)?

**Answer: It depends on the workload characteristics!**

#### Coarse.txt (100 Large Trees)
- **Winner: Per-BST** with 1.98x speedup
- Per-BST: 0.044s (spawns ~100 goroutines)
- Worker Pool (16): 0.047s (1.85x speedup)
- **Difference: Per-BST is ~7% faster**

#### Fine.txt (100,000 Small Trees)
- **Winner: Worker Pool (8 workers)** with 1.76x speedup
- Worker Pool (8): 0.034s
- Per-BST: 0.047-0.082s (spawns ~100k goroutines)
- **Difference: Worker Pool is ~40% faster!**

**Conclusion:** Worker pool is more **reliable** and **scalable** across different workloads.

---

### Q2: Can Go manage goroutines well enough that you don't have to worry about how many threads to spawn anymore?

**Answer: NO - But Go makes it much easier than traditional threads!**

#### Evidence Supporting "Go is Good":
✅ **100 goroutines:** Managed excellently
   - Per-BST on coarse.txt: 0.044s (fastest approach)
   - Nearly 2x speedup over sequential
   - Better than manually tuned worker pool

✅ **Better than traditional threads:**
   - Creating 100k pthreads would crash or be extremely slow
   - Go handles it without crashing
   - Still gets 1.26x speedup even with 100k goroutines

#### Evidence Supporting "You Still Need to Think":
❌ **100,000 goroutines:** Significant overhead
   - Per-BST on fine.txt: 0.047-0.082s (highly variable!)
   - Sometimes **slower than sequential** (0.065s vs 0.060s)
   - 40% slower than optimal worker pool (0.034s)

❌ **Over-provisioning hurts:**
   - 128 workers vs 8 workers on fine.txt:
     - 128 workers: 0.039s
     - 8 workers: 0.034s
     - Loss: 15% slower

❌ **Unpredictability at scale:**
   - Per-BST times vary significantly run-to-run
   - Worker pool is consistent and predictable

**Conclusion:** Go is **impressive** but goroutines aren't **magic**. You still need to match parallelism to your hardware and workload for optimal performance.

---

## Detailed Performance Data

### Coarse.txt Results (100 Trees, Large)

| Strategy | Workers/Goroutines | Time (s) | Speedup | Notes |
|----------|-------------------|----------|---------|-------|
| Sequential | 1 | 0.0870 | 1.00x | Baseline |
| Worker Pool | 8 | 0.0508 | 1.71x | Good |
| Worker Pool | 16 | 0.0465 | 1.87x | Better |
| Worker Pool | 32 | 0.0469 | 1.85x | Slight decline |
| Worker Pool | 64 | 0.0567 | 1.53x | Overhead increasing |
| Worker Pool | 128 | 0.0534 | 1.63x | High overhead |
| **Per-BST** | **~100** | **0.0439** | **1.98x** | **⭐ FASTEST** |

**Key Observations:**
- Per-BST wins for coarse-grained parallelism
- Worker pool optimal at 8-16 workers (close to 8 CPU cores)
- Performance degrades beyond 16 workers

### Fine.txt Results (100,000 Trees, Small)

| Strategy | Workers/Goroutines | Time (s) | Speedup | Notes |
|----------|-------------------|----------|---------|-------|
| Sequential | 1 | 0.0595 | 1.00x | Baseline |
| **Worker Pool** | **8** | **0.0338** | **1.76x** | **⭐ FASTEST** |
| Worker Pool | 16 | 0.0364 | 1.64x | Good |
| Worker Pool | 32 | 0.0367 | 1.62x | Slight decline |
| Worker Pool | 64 | 0.0397 | 1.50x | Overhead visible |
| Worker Pool | 128 | 0.0390 | 1.53x | High overhead |
| Per-BST | ~100,000 | 0.0473 | 1.26x | Inconsistent |
| Per-BST | ~100,000 | 0.0651* | 0.91x | ❌ Slower than seq! |
| Per-BST | ~100,000 | 0.0826* | 0.72x | ❌ Much slower! |

*Run-to-run variation

**Key Observations:**
- Worker pool dominant for fine-grained parallelism
- Per-BST highly variable (0.047-0.082s)
- 100k goroutines create significant scheduling overhead

---

## Performance Visualizations

### Coarse.txt Performance Curve

```
Time (seconds)
0.090 ┤ █ Sequential
      │
0.085 ┤
      │
0.080 ┤
      │
0.075 ┤
      │
0.070 ┤
      │
0.065 ┤
      │
0.060 ┤               █ (64w)
      │         █ (32w) █ (128w)
0.055 ┤
      │   █ (8w)
0.050 ┤
      │      █ (16w)
0.045 ┤         █ Per-BST ⭐
      │
0.040 ┤
      └────┬────┬────┬────┬────┬────┬────┬─────
           1    8   16   32   64  128  100
                    Worker Count / Goroutines

Legend: █ = Data point, ⭐ = Fastest
```

### Fine.txt Performance Curve

```
Time (seconds)
0.085 ┤
      │                      █ Per-BST (worst)
0.080 ┤
      │
0.075 ┤
      │
0.070 ┤
      │
0.065 ┤                █ Per-BST (bad)
      │
0.060 ┤ █ Sequential
      │
0.055 ┤
      │
0.050 ┤           █ Per-BST (best case)
      │
0.045 ┤
      │
0.040 ┤                  █ (64w)
      │                 █ (128w) █ (32w)
0.035 ┤              █ (16w)
      │   █ (8w) ⭐ OPTIMAL
0.030 ┤
      └────┬────┬────┬────┬────┬────┬─────────
           1    8   16   32   64  128  100k
                    Worker Count / Goroutines
```

---

## Analysis: Why These Results?

### Per-BST Wins on Coarse.txt

**Reasons:**
1. **Only 100 goroutines** - manageable overhead
2. **Maximum parallelism** - all trees computed simultaneously
3. **Large trees** - computation time >> goroutine overhead
4. **No contention** - each goroutine independent

**Math:**
```
Per goroutine overhead: ~2 KB stack + scheduling
100 goroutines: ~200 KB + minimal scheduling
Large tree hash time: ~1 ms
Overhead ratio: < 1%
```

### Worker Pool Wins on Fine.txt

**Reasons:**
1. **100k goroutines** - massive overhead
2. **Context switching** - 100k goroutines on 8 cores
3. **Memory overhead** - 200 MB+ for stacks
4. **Scheduler pressure** - Go runtime overwhelmed
5. **Small trees** - computation time ≈ goroutine overhead

**Math:**
```
Per goroutine overhead: ~2 KB stack + scheduling
100k goroutines: ~200 MB + heavy scheduling
Small tree hash time: ~10 μs
Context switch: ~1-2 μs
Overhead ratio: 10-20%!
```

### Optimal Worker Count = 8

**Why 8 workers perform best:**
- System has **8 CPU cores**
- Each worker gets dedicated core
- Minimal context switching
- Full hardware utilization
- Low scheduling overhead

**Why more workers hurt:**
- **16+ workers:** More workers than cores
- Context switching increases
- Cache thrashing
- Scheduler overhead grows
- Diminishing returns

---

## Key Takeaways for Writeup

### 1. Workload Characteristics Matter

| Workload | Per-BST | Worker Pool | Winner |
|----------|---------|-------------|--------|
| Few large tasks | Excellent | Good | Per-BST |
| Many small tasks | Poor | Excellent | Worker Pool |
| Unknown workload | Risky | Safe | Worker Pool |

### 2. Goroutine Economics

**Cheap but not Free:**
- Creating goroutine: ~2 KB + scheduler work
- 100 goroutines: negligible
- 100k goroutines: significant

**Rule of Thumb:**
- < 1,000 goroutines: Go handles easily
- 1,000-10,000: Starts to matter
- > 10,000: Consider worker pool

### 3. Predictability vs Peak Performance

**Per-BST:**
- ✅ Simple to implement
- ✅ Maximum parallelism
- ✅ Best case can be fastest
- ❌ Unpredictable at scale
- ❌ Risk of poor performance

**Worker Pool:**
- ✅ Consistent performance
- ✅ Tunable for workload
- ✅ Scalable
- ✅ Predictable behavior
- ⚠️ Requires tuning

### 4. Go's Concurrency Philosophy

**"Share memory by communicating" doesn't mean:**
- Spawn infinite goroutines
- Ignore performance characteristics
- Parallelism is always free

**It means:**
- Goroutines are cheaper than threads
- Channels simplify coordination
- Concurrency patterns are easier
- **But still need to design carefully!**

---

## Recommendations

### For This Assignment:
Use **worker pool with 8 workers** going forward:
- Consistent performance on all workloads
- Matches CPU core count
- Good speedup (1.7-1.8x)
- Predictable and reliable

### For Real-World Go:

**Use Per-Goroutine Pattern When:**
- Task count < 1,000
- Tasks are I/O-bound
- Need fire-and-forget semantics
- Simplicity > performance

**Use Worker Pool Pattern When:**
- Task count > 1,000
- Tasks are CPU-bound
- Need predictable performance
- Want tunable parallelism

---

## Comparison Table: At a Glance

| Aspect | Per-BST | Worker Pool |
|--------|---------|-------------|
| **Simplicity** | ✅ Very simple | ⚠️ More code |
| **Peak Performance** | ⭐ Can be fastest | ✅ Consistently good |
| **Predictability** | ❌ Variable | ⭐ Stable |
| **Scalability** | ❌ Poor at scale | ⭐ Excellent |
| **Memory** | ❌ N × 2KB | ✅ Workers × 2KB |
| **Overhead** | ⚠️ Depends on N | ✅ Fixed |
| **Tunability** | ❌ None | ⭐ Fully tunable |
| **Use Case** | Few large tasks | Many small tasks |

---

## Conclusion

**The Big Picture:**

Go's goroutines are **impressively lightweight** compared to traditional threads:
- ✅ Can create 100k goroutines (vs ~few thousand threads)
- ✅ Each goroutine uses ~2KB (vs ~1MB for threads)
- ✅ Fast creation and destruction
- ✅ Efficient scheduling

**But they're not magic:**
- ❌ Still have overhead (memory + scheduling)
- ❌ Performance degrades at scale (100k goroutines)
- ❌ Need to match parallelism to hardware
- ❌ Worker pool pattern still valuable

**Answer to "Can you stop worrying about thread count?":**

**NO** - but Go makes it much easier:
- You CAN create thousands of goroutines safely
- You CAN'T create infinite goroutines without cost
- You SHOULD still design for your workload
- You WILL get better performance with thought

**The sweet spot:** Match goroutine count to your workload and hardware (typically 1-2× CPU cores for CPU-bound tasks).

---

## Final Recommendation for Your Writeup

Discuss:
1. ✅ Per-BST won on coarse (1.98x) due to maximum parallelism
2. ✅ Worker pool won on fine (1.76x) due to controlled overhead
3. ✅ Go handles goroutines well BUT overhead exists
4. ✅ Optimal worker count ≈ CPU cores (8 in your case)
5. ✅ Worker pool more practical for unknown workloads

This demonstrates understanding of:
- Go's concurrency model
- Trade-offs between approaches
- Importance of workload characteristics
- Performance tuning principles

