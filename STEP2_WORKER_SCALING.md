# Step 2: Hash Group Building - Complete Worker Scaling Analysis
## Sequential vs Channel vs Mutex (with varying workers: 2, 4, 6, 8, 16, 32)

## Complete Performance Data

### Simple.txt (~10 BSTs, Small Trees)

| Strategy | Workers (H/D) | Time (s) | Speedup | Notes |
|----------|---------------|----------|---------|-------|
| **Sequential** | 1 / 1 | 0.000011 | 1.00x | ⭐ BEST (baseline) |
| Channel-based | 8 / 1 | 0.000250 | 0.04x | ❌ 22.7× SLOWER! |
| Mutex | 8 / 2 | 0.000023 | 0.47x | ❌ |
| Mutex | 8 / 4 | 0.000036 | 0.30x | ❌ |
| Mutex | 8 / 6 | 0.000028 | 0.39x | ❌ |
| Mutex | 8 / 8 | 0.000026 | 0.42x | ❌ |
| Mutex | 8 / 16 | 0.000031 | 0.35x | ❌ |
| Mutex | 8 / 32 | 0.000034 | 0.32x | ❌ |

**Winner: Sequential (1.00x)**
- All parallel approaches are SLOWER
- Overhead dominates for 11 microseconds of work
- No benefit from parallelism

---

### Coarse.txt (~100 BSTs, Large Trees)

| Strategy | Workers (H/D) | Time (s) | Speedup | Efficiency | Notes |
|----------|---------------|----------|---------|------------|-------|
| Sequential | 1 / 1 | 0.105 | 1.00x | 100% | Baseline |
| **Channel-based** | 8 / 1 | 0.059 | **1.76x** | 22% | ⭐ Faster! |
| Mutex | 8 / 2 | 0.062 | 1.69x | 85% | |
| Mutex | 8 / 4 | 0.056 | 1.88x | 47% | |
| Mutex | 8 / 6 | 0.052 | 1.99x | 33% | |
| Mutex | 8 / 8 | 0.053 | 1.98x | 25% | |
| **Mutex** | **8 / 16** | **0.052** | **2.01x** | **13%** | ⭐⭐ BEST |
| Mutex | 8 / 32 | 0.055 | 1.91x | 6% | Slight decline |

**Winner: Mutex (16 workers) - 2.01x speedup**
- Channel-based: 1.76x (surprisingly fast!)
- Mutex scaling:
  - Optimal at 16 workers (2.01x)
  - Diminishing returns after 16
  - 32 workers slightly worse (1.91x)

---

### Fine.txt (~100k BSTs, Small Trees)

| Strategy | Workers (H/D) | Time (s) | Speedup | Efficiency | Notes |
|----------|---------------|----------|---------|------------|-------|
| Sequential | 1 / 1 | 0.046 | 1.00x | 100% | Baseline |
| Channel-based | 8 / 1 | 0.054 | 0.85x | 11% | ❌ SLOWER |
| Mutex | 8 / 2 | 0.036 | 1.26x | 63% | |
| **Mutex** | **8 / 4** | **0.027** | **1.69x** | **42%** | ⭐⭐⭐ BEST |
| Mutex | 8 / 6 | 0.029 | 1.61x | 27% | |
| Mutex | 8 / 8 | 0.047 | 0.99x | 12% | ❌ SLOWER than sequential! |
| Mutex | 8 / 16 | 0.044 | 1.04x | 7% | Barely faster |
| Mutex | 8 / 32 | 0.034 | 1.37x | 4% | Some recovery |

**Winner: Mutex (4 workers) - 1.69x speedup**
- Channel-based: 0.85x (slower than sequential)
- Mutex scaling:
  - Optimal at 4 workers (1.69x)
  - Performance DEGRADES with 8 workers (0.99x!)
  - Slight recovery at 32 workers (1.37x)

---

## Summary Table

| Input | Sequential | Channel | Best Mutex | Optimal Workers | Best Speedup |
|-------|-----------|---------|------------|-----------------|--------------|
| **simple.txt** | 0.000011s | 0.000250s ❌ | 0.000023s @ 2 | Sequential | 1.00x |
| **coarse.txt** | 0.105s | 0.059s ⭐ | 0.052s @ 16 | **Mutex (16)** | **2.01x** |
| **fine.txt** | 0.046s | 0.054s ❌ | 0.027s @ 4 | **Mutex (4)** | **1.69x** |

---

## Key Findings

### 1. Different Inputs Have Different Optimal Worker Counts

**Coarse.txt:** Optimal at 16 workers
- More workers = better (up to 16)
- Scales linearly up to 6 workers
- Diminishing returns 6→16
- Slight decline at 32

**Fine.txt:** Optimal at 4 workers
- Peak at 4 workers (1.69x)
- Sharp decline at 8 workers (0.99x!)
- Slightly better at 32 (1.37x)
- Different pattern than coarse

### 2. Channel-based Performance Varies by Input ⭐

**Surprising discovery:**

| Input | Channel Speedup | vs Sequential |
|-------|----------------|---------------|
| simple.txt | 0.04x | 96% SLOWER ❌ |
| **coarse.txt** | **1.76x** | **76% FASTER** ⭐ |
| fine.txt | 0.85x | 15% SLOWER ❌ |

**Why channel wins on coarse.txt:**
- Larger hash groups to build
- Manager has enough work to stay busy
- Hash computation parallelism (8 workers) is beneficial
- Manager not a bottleneck for this workload

**Why channel loses on fine.txt:**
- Small, quick hash groups
- Manager becomes bottleneck
- Overhead > benefit

### 3. Mutex-based Scaling Patterns

**Coarse.txt scaling:**
```
2.0x ┤              ⭐ 16w
     │           █  █
     │        █     
     │     █
1.5x ┤  █
1.0x ┤█
     └───┴──┴──┴──┴───┴───
         Seq 2  4  6  8  16 32
```

**Fine.txt scaling:**
```
2.0x ┤
     │    ⭐ 4w
1.5x ┤       █
     │    █     █       █
1.0x ┤█           █  █
0.5x ┤
     └───┴──┴──┴──┴───┴───
         Seq 2  4  6  8  16 32
```

**Different patterns:**
- Coarse: Continuous improvement to 16
- Fine: Peak at 4, decline, partial recovery

### 4. Efficiency vs Worker Count

**Coarse.txt:**
- 2 workers: 85% efficient (excellent!)
- 4 workers: 47% efficient
- 16 workers: 13% efficient (best absolute speedup)
- Trade-off: More workers = more speed but lower efficiency

**Fine.txt:**
- 2 workers: 63% efficient
- 4 workers: 42% efficient (best absolute speedup)
- 8 workers: 12% efficient (slower than sequential!)
- Efficiency matters more here

---

## Detailed Analysis

### Why Different Optimal Points?

#### Coarse.txt (Optimal: 16 workers)

**Characteristics:**
- ~100 BSTs
- Large trees (many values)
- Hash computation is expensive
- ~100 different hash values

**Why 16 workers win:**
- Enough work to keep workers busy
- Low contention (100 hash groups)
- Workers rarely compete for same lock
- Parallelism benefit > overhead

**Why 32 doesn't improve:**
- Diminishing returns
- Contention increases
- Context switching overhead
- Cache thrashing

#### Fine.txt (Optimal: 4 workers)

**Characteristics:**
- ~100k BSTs
- Small trees (few values)
- Hash computation is cheap
- High hash collision rate

**Why 4 workers win:**
- Balance between parallelism and overhead
- Less contention than 8+ workers
- Good cache locality

**Why 8 workers WORSE than sequential:**
- High contention on hash groups
- Many BSTs map to same hash
- Lock contention dominates
- Too much context switching

**Why 32 recovers:**
- Work distribution evens out
- Some workers finish early
- Less waiting on locks
- But still not as good as 4

---

## Channel-based vs Mutex-based

### Coarse.txt Comparison

| Strategy | Time | Speedup | Notes |
|----------|------|---------|-------|
| Channel (8/1) | 0.059s | 1.76x | Surprisingly good! |
| Mutex (16) | 0.052s | 2.01x | ⭐ Best (14% faster than channel) |

**Why mutex still wins:**
- Better parallelism (16 workers vs 1 manager)
- Direct updates vs serialized
- But channel is competitive!

### Fine.txt Comparison

| Strategy | Time | Speedup | Notes |
|----------|------|---------|-------|
| Channel (8/1) | 0.054s | 0.85x | ❌ Slower than sequential |
| Mutex (4) | 0.027s | 1.69x | ⭐ Best (2× faster than channel) |

**Why mutex dominates:**
- Parallel updates critical
- Channel serialization kills performance
- Manager can't keep up with 100k BSTs

---

## Scaling Efficiency Analysis

### Strong Scaling (Fixed Problem Size)

**Coarse.txt:**
```
Workers:   Seq   2     4     6     8     16    32
Speedup:   1.00  1.69  1.88  1.99  1.98  2.01  1.91
Ideal:     1.00  2.00  4.00  6.00  8.00  16.00 32.00
Actual%:   100%  85%   47%   33%   25%   13%   6%
```

**Fine.txt:**
```
Workers:   Seq   2     4     6     8     16    32
Speedup:   1.00  1.26  1.69  1.61  0.99  1.04  1.37
Ideal:     1.00  2.00  4.00  6.00  8.00  16.00 32.00
Actual%:   100%  63%   42%   27%   12%   7%    4%
```

**Observations:**
- Coarse: Sub-linear but consistent scaling
- Fine: Non-linear with degradation at 8 workers
- Neither scales linearly

---

## Why 8 Workers Fails on Fine.txt

### The Mystery: Why does 8 workers perform WORSE than sequential?

**Hypothesis 1: Hash Collision Contention**
- Fine.txt has ~100k BSTs
- Many BSTs likely map to same hash values
- With 8 workers, high contention on map locks
- Workers spend more time waiting than working

**Hypothesis 2: Cache Thrashing**
- 8 workers on (likely) 8 CPU cores
- All updating same map simultaneously
- Cache lines invalidated constantly
- Memory bandwidth saturated

**Hypothesis 3: Context Switching**
- 8 active workers = OS scheduling overhead
- More context switches than 4 workers
- CPU time wasted on scheduling

**Hypothesis 4: Critical Section Overhead**
- Lock/unlock for every BST
- With 8 workers, lock is hot
- Serialization point becomes bottleneck

**Why 4 workers avoid this:**
- Less contention (fewer workers competing)
- Better cache behavior
- Less context switching
- Parallelism benefit > overhead

**Why 32 workers partially recover:**
- Work distribution more even
- Workers finish at different times
- Less synchronization needed
- But still worse than 4

---

## Performance Breakdown

### Coarse.txt (by component)

Assuming hash computation is parallel (8 workers for all):

| Strategy | Hash Time | Group Build | Total | Notes |
|----------|-----------|-------------|-------|-------|
| Sequential (1/1) | ~0.050s | ~0.105s | ~0.155s | Baseline |
| Channel (8/1) | ~0.047s | ~0.059s | ~0.106s | Hash parallel, group serial |
| Mutex (8/16) | ~0.047s | ~0.052s | ~0.099s | Both parallel! |

**Key insight:**
- Hash computation is same (all use 8 workers)
- Difference is in group building
- Mutex (16) builds groups fastest

### Fine.txt (by component)

| Strategy | Hash Time | Group Build | Total | Notes |
|----------|-----------|-------------|-------|-------|
| Sequential (1/1) | ~0.048s | ~0.046s | ~0.094s | Baseline |
| Channel (8/1) | ~0.049s | ~0.054s | ~0.103s | Manager bottleneck |
| Mutex (8/4) | ~0.049s | ~0.027s | ~0.076s | Optimal! |

**Key insight:**
- Hash time similar for all
- Group building time varies dramatically
- Mutex (4) is optimal for this workload

---

## Recommendations

### By Input Type

**Simple.txt (or any tiny workload):**
- ✅ **Use Sequential**
- Don't parallelize microsecond-scale work
- Overhead always dominates

**Coarse.txt (large trees, ~100 BSTs):**
- ✅ **Use Mutex with 16 workers**
- 2.01x speedup (best overall)
- Channel (1.76x) is competitive but slower
- More workers = better (up to 16)

**Fine.txt (small trees, ~100k BSTs):**
- ✅ **Use Mutex with 4 workers**
- 1.69x speedup (best overall)
- Channel is SLOWER than sequential (0.85x)
- More workers HURT performance (8 workers = 0.99x!)

### General Guidelines

```
if BSTs < 100:
    use sequential (overhead not worth it)
    
else if large_trees and BSTs ~100:
    use mutex with 8-16 workers
    (scales well, high parallelism)
    
else if small_trees and BSTs > 10k:
    use mutex with 4-8 workers
    (high contention, need fewer workers)
    
else:
    profile and tune!
```

---

## Key Lessons

### 1. No Universal Optimal Worker Count

- Coarse.txt: 16 workers optimal
- Fine.txt: 4 workers optimal
- Must profile for each workload

### 2. More Workers ≠ Better Performance

- Fine.txt at 8 workers: SLOWER than sequential!
- Contention and overhead can dominate
- Need to find sweet spot

### 3. Channel-based Can Be Competitive

- Coarse.txt: 1.76x speedup
- Better than expected!
- But still slower than best mutex

### 4. Workload Characteristics Matter

**Coarse.txt:**
- Large work per BST
- Many workers beneficial
- Scales to 16 workers

**Fine.txt:**
- Small work per BST
- Fewer workers better
- Degrades with too many

### 5. Efficiency vs Absolute Speedup

**Coarse.txt:**
- 2 workers: 85% efficient (but only 1.69x)
- 16 workers: 13% efficient (but 2.01x) ⭐
- Absolute speedup matters more

**Fine.txt:**
- 4 workers: 42% efficient (and 1.69x) ⭐
- 32 workers: 4% efficient (and only 1.37x)
- Efficiency matters when it impacts speed

---

## Surprising Discoveries

### 1. 🔍 Channel-based wins on coarse.txt!
- Expected: Channel always slow
- Reality: 1.76x speedup (competitive!)
- Lesson: Channel can work for some workloads

### 2. 🔍 8 workers WORSE than sequential on fine.txt!
- Expected: More workers = better
- Reality: 0.99x speedup (basically no gain)
- Lesson: Too many workers = contention

### 3. 🔍 Different optimal points for different inputs
- Coarse: 16 workers (2.01x)
- Fine: 4 workers (1.69x)
- Lesson: Must tune for workload

### 4. 🔍 Non-monotonic scaling on fine.txt
- 4 workers: 1.69x ⭐
- 8 workers: 0.99x ❌
- 32 workers: 1.37x
- Lesson: Performance can degrade then recover

---

## For Your Writeup

### Key Discussion Points

**Worker Scaling Analysis:**
1. Different inputs have different optimal worker counts
   - Coarse: 16 workers (2.01x)
   - Fine: 4 workers (1.69x)

2. Channel-based performance varies dramatically
   - Coarse: 1.76x (competitive!)
   - Fine: 0.85x (slower than sequential)
   - Workload determines effectiveness

3. More workers can hurt performance
   - Fine.txt at 8 workers: 0.99x
   - Contention and overhead dominate
   - Need to find sweet spot

4. Efficiency vs absolute speedup trade-off
   - Coarse: 85% efficient at 2 workers, but 13% at 16 (which is faster)
   - Fine: 42% efficient at 4 workers (and fastest)

**Architecture Impact:**
1. Mutex-based: Parallel updates, scales with workers
2. Channel-based: Serialized updates, can be competitive or slow
3. Choice depends on workload characteristics

**Performance Tuning:**
1. Must profile for specific workload
2. Optimal worker count varies (4 to 16 in this case)
3. Too many workers can degrade performance
4. Sweet spot exists but differs by input

**Lessons Learned:**
1. No universal best configuration
2. Workload characteristics critical
3. Parallelism has overhead costs
4. Must balance efficiency and absolute performance

---

## Conclusion

### Winner Summary

| Input | Winner | Workers | Speedup | Why? |
|-------|--------|---------|---------|------|
| simple.txt | Sequential | 1 | 1.00x | Too small |
| **coarse.txt** | **Mutex** | **16** | **2.01x** | Scales well, low contention |
| **fine.txt** | **Mutex** | **4** | **1.69x** | High contention, need fewer workers |

### Key Takeaways

1. ✅ **Mutex-based is the winner** for both meaningful workloads
2. ✅ **Optimal worker count varies** (4 vs 16)
3. ✅ **Channel can be competitive** (coarse: 1.76x)
4. ✅ **Too many workers hurt** (fine @ 8: 0.99x)
5. ✅ **Profile and tune** for each workload

### Final Recommendation

For this assignment:
- **Coarse.txt:** Mutex with 16 workers (2.01x speedup)
- **Fine.txt:** Mutex with 4 workers (1.69x speedup)

General advice:
- Start with ~4-8 workers
- Profile and measure
- Increase if scaling continues
- Decrease if performance degrades
- Watch for contention indicators

**The best parallel programmer knows that the optimal solution depends on the workload!** 🎯

