#!/bin/bash

# Comprehensive Step 2 Performance Analysis with Worker Scaling
# Compares: Sequential, Channel-based, and Mutex-based with varying workers

echo "=========================================="
echo "Step 2: Hash Group Building - Worker Scaling Analysis"
echo "Sequential vs Channel vs Mutex (2,4,6,8,16,32 workers)"
echo "=========================================="
echo ""

MUTEX_WORKERS=(2 4 6 8 16 32)
INPUTS=(simple coarse fine)

for input in "${INPUTS[@]}"; do
    echo "=========================================="
    echo "Testing: ${input}.txt"
    echo "=========================================="
    
    input_path="input/${input}.txt"
    
    # Get baseline (sequential - both hash-workers=1 and data-workers=1)
    echo "Getting baseline (sequential)..."
    baseline=$(./BST -hash-workers=1 -data-workers=1 -input=$input_path 2>/dev/null | grep "hashGroupTime" | awk '{print $2}')
    
    echo ""
    echo "Strategy              Workers (H/D)   Time (s)      Speedup"
    echo "-------------------   -------------   ----------    -------"
    
    # Sequential baseline
    echo "Sequential            1 / 1           $(printf '%10.6f' $baseline)    1.00x  (baseline)"
    
    # Channel-based (hash-workers=8, data-workers=1)
    echo -n "Channel-based         8 / 1           "
    channel_time=$(./BST -hash-workers=8 -data-workers=1 -input=$input_path 2>/dev/null | grep "hashGroupTime" | awk '{print $2}')
    if [ -n "$channel_time" ] && [ -n "$baseline" ]; then
        speedup=$(echo "scale=2; $baseline / $channel_time" | bc)
        echo "$(printf '%10.6f' $channel_time)    ${speedup}x"
    else
        echo "ERROR"
    fi
    
    echo ""
    echo "--- Mutex-based with varying workers ---"
    
    # Mutex-based with different worker counts
    for w in "${MUTEX_WORKERS[@]}"; do
        echo -n "Mutex                 8 / $(printf '%2d' $w)          "
        mutex_time=$(./BST -hash-workers=8 -data-workers=$w -input=$input_path 2>/dev/null | grep "hashGroupTime" | awk '{print $2}')
        
        if [ -n "$mutex_time" ] && [ -n "$baseline" ]; then
            speedup=$(echo "scale=2; $baseline / $mutex_time" | bc)
            echo "$(printf '%10.6f' $mutex_time)    ${speedup}x"
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
echo "  Sequential    - hash-workers=1, data-workers=1"
echo "  Channel-based - hash-workers=8, data-workers=1 (central manager)"
echo "  Mutex         - hash-workers=8, data-workers=N (direct updates)"

