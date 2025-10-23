#!/bin/bash

# Comprehensive Step 1 Performance Analysis
# Compares: Sequential, Per-BST, and Worker Pool with varying workers (2,4,8,16,32)

echo "=========================================="
echo "Step 1: Hash Computation - Complete Analysis"
echo "Sequential vs Per-BST vs Worker Pool (2,4,8,16,32)"
echo "=========================================="
echo ""

POOL_WORKERS=(2 4 8 16 32)
INPUTS=(simple coarse fine)

for input in "${INPUTS[@]}"; do
    echo "=========================================="
    echo "Testing: ${input}.txt"
    echo "=========================================="
    
    input_path="input/${input}.txt"
    
    # Get baseline (sequential)
    echo "Getting baseline (sequential)..."
    baseline=$(./BST -hash-workers=1 -input=$input_path 2>/dev/null | grep "hashTime" | awk '{print $2}')
    
    echo ""
    echo "Strategy              Workers    Time (s)      Speedup"
    echo "-------------------   -------    ----------    -------"
    
    # Sequential baseline
    echo "Sequential            1          $(printf '%10.6f' $baseline)    1.00x  (baseline)"
    
    # Per-BST (one goroutine per BST)
    echo -n "Per-BST               N/A        "
    perbst_time=$(./BST -hash-workers=2 -hash-strategy=per-bst -input=$input_path 2>/dev/null | grep "hashTime" | awk '{print $2}')
    if [ -n "$perbst_time" ] && [ -n "$baseline" ]; then
        speedup=$(echo "scale=2; $baseline / $perbst_time" | bc)
        echo "$(printf '%10.6f' $perbst_time)    ${speedup}x"
    else
        echo "ERROR"
    fi
    
    echo ""
    echo "--- Worker Pool with varying workers ---"
    
    # Worker pool with different worker counts
    for w in "${POOL_WORKERS[@]}"; do
        echo -n "Pool                  $(printf '%2d' $w)         "
        pool_time=$(./BST -hash-workers=$w -hash-strategy=pool -input=$input_path 2>/dev/null | grep "hashTime" | awk '{print $2}')
        
        if [ -n "$pool_time" ] && [ -n "$baseline" ]; then
            speedup=$(echo "scale=2; $baseline / $pool_time" | bc)
            echo "$(printf '%10.6f' $pool_time)    ${speedup}x"
        else
            echo "ERROR"
        fi
    done
    
    echo ""
done

echo "=========================================="
echo "Analysis Complete"
echo "=========================================="
echo ""
echo "Legend:"
echo "  Sequential - hash-workers=1 (single-threaded baseline)"
echo "  Per-BST    - One goroutine per BST (unbounded parallelism)"
echo "  Pool       - Fixed worker pool with N workers"

