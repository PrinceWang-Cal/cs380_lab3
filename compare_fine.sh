#!/bin/bash

# Specialized comparison script for fine.txt (100k BSTs)
# Tests both strategies with various worker counts

INPUT_FILE="input/fine.txt"

echo "=========================================="
echo "Fine-Grained Comparison: fine.txt (~100k BSTs)"
echo "=========================================="
echo ""
echo "This test will take several minutes..."
echo ""

WORKERS=(1 8 64 128 1024 10000)

echo "Strategy: Worker Pool (Fixed number of workers)"
echo "------------------------------------------------"
for w in "${WORKERS[@]}"; do
    echo -n "Workers=$w: "
    ./BST -hash-workers=$w -hash-strategy=pool -input=$INPUT_FILE | grep hashTime
done

echo ""
echo "Strategy: Per-BST (One goroutine per BST)"
echo "------------------------------------------------"
echo "Note: This spawns ~100,000 goroutines (one per BST)"
echo -n "Creating ~100k goroutines: "
./BST -hash-workers=1 -hash-strategy=perbst -input=$INPUT_FILE | grep hashTime

echo ""
echo "=========================================="
echo "Analysis"
echo "=========================================="
echo "Key Observations:"
echo "1. Worker Pool Performance:"
echo "   - How does performance scale with worker count?"
echo "   - Is there a sweet spot before overhead dominates?"
echo "   - Do 10,000 workers perform worse than 1,024?"
echo ""
echo "2. Per-BST Strategy:"
echo "   - Can Go really handle 100k goroutines efficiently?"
echo "   - How does it compare to the optimal worker pool size?"
echo ""
echo "3. Goroutine Overhead:"
echo "   - Creating/scheduling 100k goroutines has overhead"
echo "   - Worker pool should show you still need to think about parallelism"
echo ""
echo "4. Practical Takeaway:"
echo "   - Go is good at managing goroutines, but they're not free!"
echo "   - Controlled parallelism (worker pool) often wins for large workloads"

