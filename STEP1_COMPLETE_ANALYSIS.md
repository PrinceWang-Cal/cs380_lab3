# Step 1: Hash Computation - Complete Worker Scaling Analysis
## Sequential vs Per-BST vs Worker Pool (2, 4, 8, 16, 32 workers)

## Complete Performance Data

### Simple.txt (~10 BSTs, Small Trees)

| Strategy | Workers | Time (s) | Speedup | Notes |
|----------|---------|----------|---------|-------|
| **Sequential** | 1 | 0.000010 | 1.00x | ⭐ BEST (baseline) |
| Per-BST | ~10 | 0.000088 | 0.11x | ❌ 8.8× SLOWER! |
| Pool | 2 | 0.000027 | 0.37x | ❌ |
| Pool | 4 | 0.000049 | 0.20x | ❌ |
| Pool | 8 | 0.000058 | 0.17x | ❌ |
| Pool | 16 | 0.000074 | 0.13x | ❌ |
| Pool | 32 | 0.000094 | 0.10x | ❌ |

**Winner: Sequential (1.00x)**
- All parallel approaches are SLOWER
- Overhead dominates for 10 microseconds of work
- Per-BST is worst (8.8× slower!)

---

### Coarse.txt (~100 BSTs, Large Trees)

| Strategy | Workers | Time (s) | Speedup | Efficiency | Notes |
|----------|---------|----------|---------|------------|-------|
| Sequential | 1 | 0.144 | 1.00x | 100% | Baseline |
| Per-BST | ~100 | 0.074 | 1.93x | 2% | Good! |
| Pool | 2 | 0.067 | 2.14x | 107% | ⭐ Super-linear! |
| Pool | 4 | 0.056 | 2.56x | 64% | |
| **Pool** | **8** | **0.052** | **2.75x** | **34%** | ⭐⭐⭐ BEST |
| Pool | 16 | 0.095 | 1.50x | 9% | ❌ Worse than 8 |
| Pool | 32 | 0.099 | 1.45x | 5% | ❌ Continues decline |

**Winner: Pool (8 workers) - 2.75x speedup**
- Per-BST: 1.93x (competitive)
- Pool peaks at 8 workers
- Performance degrades significantly with 16+ workers
- Super-linear speedup at 2 workers (likely cache effects)

---

### Fine.txt (~100k BSTs, Small Trees)

| Strategy | Workers | Time (s) | Speedup | Efficiency | Notes |
|----------|---------|----------|---------|------------|-------|
| Sequential | 1 | 0.049 | 1.00x | 100% | Baseline |
| Per-BST | ~100k | 0.037 | 1.34x | <0.01% | Surprisingly good! |
| Pool | 2 | 0.042 | 1.16x | 58% | |
| **Pool** | **4** | **0.034** | **1.45x** | **36%** | ⭐⭐⭐ BEST |
| Pool | 8 | 0.038 | 1.28x | 16% | Decline |
| Pool | 16 | 0.052 | 0.94x | 6% | ❌ SLOWER than sequential! |
| Pool | 32 | 0.043 | 1.13x | 4% | Partial recovery |

**Winner: Pool (4 workers) - 1.45x speedup**
- Per-BST: 1.34x (competitive despite 100k goroutines!)
- Pool peaks at 4 workers
- Performance degrades at 8+ workers
- 16 workers actually SLOWER than sequential!

---

## Summary Table

| Input | Sequential | Per-BST | Best Pool | Optimal | Best Speedup |
|-------|-----------|---------|-----------|---------|--------------|
| **simple.txt** | 0.000010s | 0.000088s ❌ | 0.000027s @2 ❌ | Sequential | 1.00x |
| **coarse.txt** | 0.144s | 0.074s ⭐ | 0.052s @8 ⭐⭐ | Pool (8) | 2.75x |
| **fine.txt** | 0.049s | 0.037s ⭐ | 0.034s @4 ⭐⭐ | Pool (4) | 1.45x |

---

## Key Findings

### 1. Different Optimal Worker Counts by Input

**Coarse.txt:** Optimal at 8 workers (2.75x)
- Large trees = expensive hash computation
- Benefits from more workers
- Peaks at 8, degrades at 16+

**Fine.txt:** Optimal at 4 workers (1.45x)
- Small trees = cheap hash computation
- Less work per BST
- Peaks at 4, degrades at 8+

**Pattern:** 2× difference in optimal worker count!

### 2. Per-BST Performance Varies Dramatically

| Input | Per-BST Speedup | Goroutines | Notes |
|-------|----------------|------------|-------|
| simple.txt | 0.11x | ~10 | ❌ Terrible (overhead >> work) |
| coarse.txt | 1.93x | ~100 | ⭐ Competitive! |
| fine.txt | 1.34x | ~100k | ⭐ Surprisingly good! |

**Surprising:** Per-BST handles 100k goroutines reasonably well!

### 3. Worker Pool Scaling Patterns

**Coarse.txt - Continuous then Decline:**
```
3.0x ┤
     │      ⭐ 8w
2.5x ┤         █
     │      █
2.0x ┤   █  Per-BST
     │█
1.0x ┤              █ 16w
     └──┴───┴──┴───┴───┴───
        Seq 2  4  8  16 32
```

**Fine.txt - Peak then Degrade:**
```
1.5x ┤    ⭐ 4w
     │       █  Per-BST
     │    █  █
1.0x ┤█  █        ❌16w
     │         █     █
0.5x ┤
     └──┴───┴──┴───┴───┴───
        Seq 2  4  8  16 32
```

### 4. Performance Degradation with Too Many Workers

**Coarse.txt:**
- 8 workers: 2.75x ⭐
- 16 workers: 1.50x ❌ (45% slower!)
- 32 workers: 1.45x (continues decline)

**Fine.txt:**
- 4 workers: 1.45x ⭐
- 8 workers: 1.28x (decline starts)
- 16 workers: 0.94x ❌ (SLOWER than sequential!)

**Why?** Overhead, contention, and context switching dominate

### 5. Super-linear Speedup at 2 Workers (Coarse.txt)

- 2 workers: 2.14x speedup (107% efficiency!)
- Likely cache effects:
  - Sequential may cause cache misses
  - 2 workers may benefit from better cache utilization
  - Or measurement variance

---

## Detailed Analysis

### Why Different Optimal Points?

#### Coarse.txt (Optimal: 8 workers)

**Characteristics:**
- ~100 BSTs
- Large trees (many nodes)
- Expensive hash computation (~1-2ms per BST)

**Why 8 workers win:**
- Enough work to keep workers busy
- Parallelism benefit > overhead
- Good balance

**Why 16+ workers fail:**
- Not enough BSTs to keep all workers busy
- Some workers idle
- Context switching overhead increases
- Diminishing returns

#### Fine.txt (Optimal: 4 workers)

**Characteristics:**
- ~100k BSTs
- Small trees (few nodes)
- Cheap hash computation (~0.5μs per BST)

**Why 4 workers win:**
- Balance between parallelism and overhead
- Less context switching than 8+
- Good work distribution

**Why 8+ workers fail:**
- Work too granular
- Context switching overhead dominates
- Cache thrashing
- Lock contention (if any)

---

## Per-BST vs Worker Pool Comparison

### Per-BST Strengths:

1. **Simple implementation** (~20 lines)
2. **Maximum parallelism** (every BST processed simultaneously)
3. **No queue management** needed
4. **Competitive on coarse.txt** (1.93x vs 2.75x)
5. **Handles 100k goroutines** reasonably well (fine.txt: 1.34x)

### Per-BST Weaknesses:

1. **Overhead for tiny workloads** (simple.txt: 8.8× slower!)
2. **Unpredictable** (depends on Go runtime)
3. **Not optimal** for any input
4. **Worse than best pool** by 30-50%

### Worker Pool Strengths:

1. **Consistently faster** than per-BST (when tuned)
2. **Tunable performance** (can optimize for workload)
3. **Bounded resources** (predictable goroutine count)
4. **Best absolute performance** on all meaningful inputs

### Worker Pool Weaknesses:

1. **More complex** (~40 lines)
2. **Requires tuning** (need to find optimal worker count)
3. **Can degrade badly** if poorly tuned (fine.txt @ 16: 0.94x!)
4. **Not one-size-fits-all**

---

## Scaling Efficiency Analysis

### Strong Scaling (Fixed Problem Size)

**Coarse.txt:**
```
Workers:    1     2     4     8    16    32
Speedup:   1.00  2.14  2.56  2.75  1.50  1.45
Ideal:     1.00  2.00  4.00  8.00 16.00 32.00
Actual%:   100%  107%   64%   34%    9%    5%
```

**Fine.txt:**
```
Workers:    1     2     4     8    16    32
Speedup:   1.00  1.16  1.45  1.28  0.94  1.13
Ideal:     1.00  2.00  4.00  8.00 16.00 32.00
Actual%:   100%   58%   36%   16%   (-)    4%
```

**Observations:**
- Coarse: Super-linear at 2, then sub-linear
- Fine: Sub-linear throughout, degrades at 16

---

## Why Pool Degrades with More Workers

### Coarse.txt (8 → 16 workers)

**Performance drop:** 2.75x → 1.50x (45% slower!)

**Likely causes:**

1. **Insufficient work:**
   - Only ~100 BSTs
   - 16 workers = ~6 BSTs per worker
   - Some workers idle

2. **Context switching:**
   - 16 goroutines on ~8 cores
   - OS scheduling overhead
   - CPU time wasted

3. **Channel contention:**
   - Workers compete for work channel
   - More blocking/wakeup cycles
   - Coordination overhead

### Fine.txt (4 → 16 workers)

**Performance drop:** 1.45x → 0.94x (slower than sequential!)

**Likely causes:**

1. **Tiny work units:**
   - Each hash computation: ~0.5μs
   - Context switch: ~1-2μs
   - Overhead > work!

2. **Cache thrashing:**
   - 16 workers across cores
   - Cache lines invalidated
   - Memory bandwidth saturated

3. **Synchronization overhead:**
   - Channel operations dominate
   - More workers = more coordination
   - Overhead kills performance

---

## Recommendations

### By Input Type

**Simple.txt (or any tiny workload):**
- ✅ **Use Sequential**
- All parallel approaches slower
- 10 microseconds too fast

**Coarse.txt (large trees, ~100 BSTs):**
- ✅ **Use Pool with 8 workers** (2.75x)
- Per-BST competitive (1.93x)
- Avoid 16+ workers (degrades)

**Fine.txt (small trees, ~100k BSTs):**
- ✅ **Use Pool with 4 workers** (1.45x)
- Per-BST competitive (1.34x)
- Avoid 8+ workers (degrades)

### General Guidelines

```
if BSTs < 100:
    use sequential (overhead not worth it)
    
else if large_trees:
    use pool with 8 workers
    (expensive computation, benefits from parallelism)
    
else if small_trees:
    use pool with 4 workers
    (cheap computation, need fewer workers)
    
else:
    profile and tune!
```

---

## Surprising Discoveries

### 1. 🔍 Per-BST handles 100k goroutines well!

- Expected: Crash or terrible performance
- Reality: 1.34x speedup on fine.txt
- Go's runtime is impressive!

### 2. 🔍 Super-linear speedup at 2 workers (coarse.txt)

- 2 workers: 2.14x (107% efficiency!)
- Likely cache effects
- Or measurement variance

### 3. 🔍 16 workers SLOWER than sequential (fine.txt)

- Expected: Some speedup
- Reality: 0.94x (slowdown!)
- Overhead completely dominates

### 4. 🔍 Different optimal points (4 vs 8)

- Coarse: 8 workers
- Fine: 4 workers
- 2× difference based on work granularity

---

## Per-BST Deep Dive

### How Does Per-BST Handle 100k Goroutines?

**fine.txt results:**
- 100k goroutines spawned
- 1.34x speedup
- Better than pool @ 16 workers (0.94x!)

**Why it works:**

1. **Go's efficient scheduler:**
   - Goroutines are lightweight (~2KB stack)
   - Multiplexed onto OS threads
   - Efficient context switching

2. **Work is uniform:**
   - All BSTs processed similarly
   - No complex coordination
   - Simple spawn-compute-return pattern

3. **No synchronization:**
   - Each goroutine independent
   - No locks, no channels (for hash only)
   - Pure parallel computation

**Why it's not optimal:**

1. **Overhead still exists:**
   - 100k goroutine creations
   - Scheduling overhead
   - Memory overhead

2. **Pool is faster:**
   - 4 workers: 1.45x vs 1.34x
   - 8% faster despite controlling resources
   - Better work distribution

---

## Efficiency vs Absolute Speedup

### Coarse.txt

| Workers | Speedup | Efficiency | Note |
|---------|---------|------------|------|
| 2 | 2.14x | 107% | Super-linear! |
| 4 | 2.56x | 64% | Good |
| 8 | 2.75x | 34% | ⭐ Best absolute |
| 16 | 1.50x | 9% | Poor |

**Takeaway:** Low efficiency (34%) still gives best absolute performance!

### Fine.txt

| Workers | Speedup | Efficiency | Note |
|---------|---------|------------|------|
| 2 | 1.16x | 58% | Decent |
| 4 | 1.45x | 36% | ⭐ Best absolute |
| 8 | 1.28x | 16% | Declining |
| 16 | 0.94x | -6% | Negative! |

**Takeaway:** Efficiency helps predict optimal point (4 workers)

---

## For Your Writeup

### Key Discussion Points

**1. Worker count matters:**
- Coarse: 8 workers optimal (2.75x)
- Fine: 4 workers optimal (1.45x)
- Must tune for workload

**2. More workers can hurt:**
- Coarse @ 16: 45% slower than optimal
- Fine @ 16: Slower than sequential!
- Overhead and contention dominate

**3. Per-BST surprisingly robust:**
- Handles 100k goroutines (1.34x)
- Simple implementation
- But not optimal

**4. Different scaling patterns:**
- Coarse: Continuous to 8, then decline
- Fine: Peak at 4, degrade, partial recovery
- Workload characteristics critical

**5. Super-linear speedup:**
- Coarse @ 2 workers: 2.14x (107%!)
- Likely cache effects
- Shows complexity of performance analysis

---

## Conclusion

### Winner Summary

| Input | Winner | Workers | Speedup | Why? |
|-------|--------|---------|---------|------|
| simple.txt | Sequential | 1 | 1.00x | Too small |
| **coarse.txt** | **Pool** | **8** | **2.75x** | Expensive hashing, benefits from parallelism |
| **fine.txt** | **Pool** | **4** | **1.45x** | Cheap hashing, need fewer workers |

### Key Takeaways

1. ✅ **Worker pool wins when tuned correctly**
2. ✅ **Optimal worker count varies** (4 vs 8)
3. ✅ **Per-BST is competitive but not optimal**
4. ✅ **Too many workers hurt performance**
5. ✅ **Profile and tune for each workload**

### Final Recommendation

**For this assignment:**
- **Coarse.txt:** Pool with 8 workers (2.75x speedup)
- **Fine.txt:** Pool with 4 workers (1.45x speedup)

**General advice:**
- Start with workers = CPU cores
- Profile and measure
- Increase if scaling continues
- Decrease if performance degrades
- Watch for degradation zones

**The best parallel programmer knows that context matters!** 🎯

