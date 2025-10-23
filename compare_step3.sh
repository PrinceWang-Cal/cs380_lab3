#!/bin/bash

# Script to compare Step 3 implementations (Tree Comparison)
# Sequential vs Unbounded vs Worker Pool

INPUT_FILE=${1:-coarse}
INPUT_PATH="input/${INPUT_FILE}.txt"

echo "=========================================="
echo "Step 3 Comparison: Tree Comparison"
echo "Input: ${INPUT_FILE}.txt"
echo "=========================================="
echo ""

echo "1. Sequential (Baseline)"
echo "------------------------------------------------"
./BST -hash-workers=8 -data-workers=8 -comp-workers=1 -input=$INPUT_PATH | grep "compareTreeTime"

echo ""
echo "2. Parallel Unbounded (goroutine per comparison)"
echo "------------------------------------------------"
echo "   - Spawns one goroutine for each tree pair"
echo "   - Maximum parallelism"
echo "   - No control over goroutine count"
echo ""
./BST -hash-workers=8 -data-workers=8 -comp-workers=2 -comp-strategy=unbounded -input=$INPUT_PATH | grep "compareTreeTime"

echo ""
echo "3. Parallel Worker Pool (fixed workers)"
echo "------------------------------------------------"
echo "   - Fixed number of worker goroutines"
echo "   - Workers process from bounded buffer (channel)"
echo "   - Controlled parallelism"
echo ""

WORKERS=(2 4 8 16)
for w in "${WORKERS[@]}"; do
    echo -n "  $w workers: "
    ./BST -hash-workers=8 -data-workers=8 -comp-workers=$w -input=$INPUT_PATH | grep "compareTreeTime" | awk '{print $2}'
done

echo ""
echo "=========================================="
echo "Analysis Questions"
echo "=========================================="
echo "1. How do the performance and complexity compare?"
echo "   → Unbounded: Simple, maximum parallelism"
echo "   → Pool: More complex, controlled parallelism"
echo ""
echo "2. How do they scale compared to single thread?"
echo "   → Compare to sequential baseline"
echo ""
echo "3. Is the additional complexity worthwhile?"
echo "   → Trade-off: complexity vs performance"
echo ""

