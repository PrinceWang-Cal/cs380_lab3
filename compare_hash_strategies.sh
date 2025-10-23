#!/bin/bash

# Script to compare hash computation strategies
# Usage: ./compare_hash_strategies.sh [simple|coarse|fine]

INPUT_FILE=${1:-simple}
INPUT_PATH="input/${INPUT_FILE}.txt"

echo "=========================================="
echo "Comparing Hash Strategies: ${INPUT_FILE}.txt"
echo "=========================================="
echo ""

WORKERS=(1 8 16 32 64 128)

echo "Strategy: Worker Pool (Fixed number of workers)"
echo "------------------------------------------------"
for w in "${WORKERS[@]}"; do
    echo -n "Workers=$w: "
    ./BST -hash-workers=$w -hash-strategy=pool -input=$INPUT_PATH | grep hashTime
done

echo ""
echo "Strategy: Per-BST (One goroutine per BST)"
echo "------------------------------------------------"
for w in "${WORKERS[@]}"; do
    echo -n "Workers=$w (ignored): "
    ./BST -hash-workers=$w -hash-strategy=perbst -input=$INPUT_PATH | grep hashTime
done

echo ""
echo "=========================================="
echo "Analysis"
echo "=========================================="
echo "Note: Per-BST strategy ignores -hash-workers flag and spawns"
echo "      one goroutine per BST (${INPUT_FILE}.txt has ~10 BSTs for simple,"
echo "      ~100 for coarse, ~100000 for fine)"

