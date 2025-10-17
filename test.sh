#!/bin/bash

# Test script for BST assignment
# Make sure to compile first: go build BST.go

echo "==================================="
echo "Testing Step 1: Hash Computation"
echo "==================================="

echo -e "\n--- Simple input, 1 worker ---"
./BST -hash-workers=1 -input=input/simple.txt

echo -e "\n--- Simple input, 2 workers ---"
./BST -hash-workers=2 -input=input/simple.txt

echo -e "\n--- Simple input, 4 workers ---"
./BST -hash-workers=4 -input=input/simple.txt

echo -e "\n==================================="
echo "Testing Step 2: Hash Groups"
echo "==================================="

echo -e "\n--- Sequential (1,1) ---"
./BST -hash-workers=1 -data-workers=1 -input=input/simple.txt

echo -e "\n--- Channel-based (2,1) ---"
./BST -hash-workers=2 -data-workers=1 -input=input/simple.txt

echo -e "\n--- Mutex-based (2,2) ---"
./BST -hash-workers=2 -data-workers=2 -input=input/simple.txt

echo -e "\n==================================="
echo "Testing Step 3: Tree Comparison"
echo "==================================="

echo -e "\n--- Sequential comparison ---"
./BST -hash-workers=1 -data-workers=1 -comp-workers=1 -input=input/simple.txt

echo -e "\n--- Parallel comparison (2 workers) ---"
./BST -hash-workers=2 -data-workers=2 -comp-workers=2 -input=input/simple.txt

echo -e "\n==================================="
echo "Done!"
echo "==================================="

