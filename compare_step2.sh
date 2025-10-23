#!/bin/bash

# Script to compare Step 2 implementations (Hash Groups)
# Channel-based vs Mutex-based synchronization

INPUT_FILE=${1:-simple}
INPUT_PATH="input/${INPUT_FILE}.txt"

echo "=========================================="
echo "Step 2 Comparison: Hash Groups"
echo "Input: ${INPUT_FILE}.txt"
echo "=========================================="
echo ""

WORKERS=(1 2 4 8 16)

echo "1. Sequential (Baseline)"
echo "------------------------------------------------"
./BST -hash-workers=1 -data-workers=1 -input=$INPUT_PATH | grep hashGroupTime

echo ""
echo "2. Channel-based (hash-workers, 1 manager)"
echo "------------------------------------------------"
echo "   - Hash workers compute in parallel"
echo "   - Send results to channel"
echo "   - 1 central manager updates map (no mutex needed)"
echo ""
for w in "${WORKERS[@]}"; do
    if [ $w -eq 1 ]; then
        continue  # Skip, already tested as sequential
    fi
    echo -n "  $w workers: "
    ./BST -hash-workers=$w -data-workers=1 -input=$INPUT_PATH | grep hashGroupTime
done

echo ""
echo "3. Mutex-based (hash-workers, direct map access)"
echo "------------------------------------------------"
echo "   - Hash workers compute in parallel"
echo "   - Each worker locks mutex to update map"
echo "   - Workers 'fight' for the lock"
echo ""
for w in "${WORKERS[@]}"; do
    if [ $w -eq 1 ]; then
        continue  # Skip, already tested as sequential
    fi
    echo -n "  $w workers: "
    ./BST -hash-workers=$w -data-workers=$w -input=$INPUT_PATH | grep hashGroupTime
done

echo ""
echo "=========================================="
echo "Analysis Questions"
echo "=========================================="
echo "1. Which approach has more overhead?"
echo "   → Compare times for same worker count"
echo ""
echo "2. How much faster are they vs single thread?"
echo "   → Compare to sequential baseline"
echo ""
echo "3. Which approach is simpler?"
echo "   → Channel: Separate concerns, no mutex"
echo "   → Mutex: Direct access, explicit locking"
echo ""

