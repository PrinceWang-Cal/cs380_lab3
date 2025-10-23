#!/bin/bash

# Comprehensive Step 3 Performance Analysis
# Tests: Sequential, Unbounded, and Worker Pool (2, 4, 8, 16, 32) on simple, coarse, and fine

echo "=========================================="
echo "Step 3: Comprehensive Performance Analysis"
echo "Comparing: Sequential, Unbounded, and Worker Pool"
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
    baseline=$(./BST -hash-workers=8 -data-workers=8 -comp-workers=1 -input=$input_path 2>/dev/null | grep "compareTreeTime" | awk '{print $2}')
    
    echo ""
    echo "Strategy              Workers   Time (s)      Speedup"
    echo "-------------------   -------   ----------    -------"
    
    # Sequential baseline
    echo "Sequential               1      $(printf '%10.6f' $baseline)    1.00x  (baseline)"
    
    # Unbounded (spawns goroutine per comparison)
    echo -n "Unbounded                N/A    "
    unbounded_time=$(./BST -hash-workers=8 -data-workers=8 -comp-workers=2 -comp-strategy=unbounded -input=$input_path 2>/dev/null | grep "compareTreeTime" | awk '{print $2}')
    if [ -n "$unbounded_time" ] && [ -n "$baseline" ]; then
        speedup=$(echo "scale=2; $baseline / $unbounded_time" | bc)
        echo "$(printf '%10.6f' $unbounded_time)    ${speedup}x"
    else
        echo "ERROR"
    fi
    
    # Worker pool with different worker counts
    for w in "${POOL_WORKERS[@]}"; do
        echo -n "Pool                  $(printf '%4d' $w)    "
        time_val=$(./BST -hash-workers=8 -data-workers=8 -comp-workers=$w -input=$input_path 2>/dev/null | grep "compareTreeTime" | awk '{print $2}')
        
        if [ -n "$time_val" ] && [ -n "$baseline" ]; then
            speedup=$(echo "scale=2; $baseline / $time_val" | bc)
            echo "$(printf '%10.6f' $time_val)    ${speedup}x"
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
echo "  Sequential - Single-threaded baseline"
echo "  Unbounded  - One goroutine per comparison (comp-workers ignored)"
echo "  Pool       - Fixed worker pool with bounded buffer"

