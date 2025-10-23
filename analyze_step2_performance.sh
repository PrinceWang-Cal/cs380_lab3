#!/bin/bash

# Comprehensive Step 2 Performance Analysis
# Compares: Sequential, Channel-based, and Mutex-based hash group building

echo "=========================================="
echo "Step 2: Hash Group Building Comparison"
echo "Comparing: Sequential, Channel, and Mutex"
echo "=========================================="
echo ""

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
    
    # Mutex-based (hash-workers=8, data-workers=8)
    echo -n "Mutex-based           8 / 8           "
    mutex_time=$(./BST -hash-workers=8 -data-workers=8 -input=$input_path 2>/dev/null | grep "hashGroupTime" | awk '{print $2}')
    if [ -n "$mutex_time" ] && [ -n "$baseline" ]; then
        speedup=$(echo "scale=2; $baseline / $mutex_time" | bc)
        echo "$(printf '%10.6f' $mutex_time)    ${speedup}x"
    else
        echo "ERROR"
    fi
    
    echo ""
    
    # Show winner
    if [ -n "$channel_time" ] && [ -n "$mutex_time" ]; then
        if (( $(echo "$mutex_time < $channel_time" | bc -l) )); then
            improvement=$(echo "scale=1; ($channel_time - $mutex_time) / $channel_time * 100" | bc)
            echo "→ Winner: Mutex-based (${improvement}% faster than channel)"
        else
            improvement=$(echo "scale=1; ($mutex_time - $channel_time) / $mutex_time * 100" | bc)
            echo "→ Winner: Channel-based (${improvement}% faster than mutex)"
        fi
    fi
    echo ""
done

echo "=========================================="
echo "Analysis Complete"
echo "=========================================="
echo ""
echo "Legend:"
echo "  Sequential    - hash-workers=1, data-workers=1"
echo "  Channel-based - hash-workers=8, data-workers=1 (central manager)"
echo "  Mutex-based   - hash-workers=8, data-workers=8 (direct updates)"

