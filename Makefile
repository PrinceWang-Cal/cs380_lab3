.PHONY: build clean test-seq test-seq-simple test-seq-coarse test-seq-fine test-all help

# Binary name
BINARY=BST

# Build the program
build:
	@echo "Building $(BINARY)..."
	go build -o $(BINARY) $(BINARY).go
	@echo "Build complete!"

# Test sequential solution (all flags = 1) on all input files
test-seq: build test-seq-simple test-seq-coarse test-seq-fine

# Test sequential on simple.txt
test-seq-simple: build
	@echo "=========================================="
	@echo "Testing Sequential Solution: simple.txt"
	@echo "=========================================="
	./$(BINARY) -hash-workers=1 -data-workers=1 -comp-workers=1 -input=input/simple.txt
	@echo ""

# Test sequential on coarse.txt
test-seq-coarse: build
	@echo "=========================================="
	@echo "Testing Sequential Solution: coarse.txt"
	@echo "=========================================="
	./$(BINARY) -hash-workers=1 -data-workers=1 -comp-workers=1 -input=input/coarse.txt
	@echo ""

# Test sequential on fine.txt
test-seq-fine: build
	@echo "=========================================="
	@echo "Testing Sequential Solution: fine.txt"
	@echo "=========================================="
	./$(BINARY) -hash-workers=1 -data-workers=1 -comp-workers=1 -input=input/fine.txt
	@echo ""

# Test only hash computation (Step 1)
test-hash: build
	@echo "=========================================="
	@echo "Testing Step 1: Hash Computation Only"
	@echo "=========================================="
	@echo "--- 1 worker ---"
	./$(BINARY) -hash-workers=1 -input=input/simple.txt
	@echo ""
	@echo "--- 2 workers ---"
	./$(BINARY) -hash-workers=2 -input=input/simple.txt
	@echo ""
	@echo "--- 4 workers ---"
	./$(BINARY) -hash-workers=4 -input=input/simple.txt
	@echo ""

# Compare hash strategies
compare-strategies: build
	@echo "=========================================="
	@echo "Comparing Hash Strategies"
	@echo "=========================================="
	./compare_hash_strategies.sh simple

compare-strategies-coarse: build
	./compare_hash_strategies.sh coarse

compare-strategies-fine: build
	@echo "WARNING: This will take a while with fine.txt..."
	./compare_hash_strategies.sh fine

# Detailed fine.txt comparison (1, 8, 64, 128, 1024, 10000 workers)
compare-fine: build
	@echo "Running comprehensive fine.txt analysis..."
	./compare_fine.sh

# Complete Step 1 analysis: Sequential vs Per-BST vs Pool (2,4,8,16,32)
compare-step1-complete: build
	@echo "Running complete Step 1 analysis..."
	./analyze_step1_complete.sh

# Test hash groups (Step 2)
test-groups: build
	@echo "=========================================="
	@echo "Testing Step 2: Hash Groups"
	@echo "=========================================="
	@echo "--- Sequential (1,1) ---"
	./$(BINARY) -hash-workers=1 -data-workers=1 -input=input/simple.txt
	@echo ""
	@echo "--- Channel-based (2,1) ---"
	./$(BINARY) -hash-workers=2 -data-workers=1 -input=input/simple.txt
	@echo ""
	@echo "--- Mutex-based (2,2) ---"
	./$(BINARY) -hash-workers=2 -data-workers=2 -input=input/simple.txt
	@echo ""

# Compare Step 2 implementations (old individual tests)
compare-step2-old: build
	./compare_step2.sh simple

compare-step2-coarse-old: build
	./compare_step2.sh coarse

compare-step2-fine-old: build
	./compare_step2.sh fine

# Compare Step 2 implementations (comprehensive - all inputs)
compare-step2: build
	@echo "Running comprehensive Step 2 analysis..."
	./analyze_step2_performance.sh

# Compare Step 2 with worker scaling (2, 4, 6, 8, 16, 32 workers)
compare-step2-scaling: build
	@echo "Running Step 2 worker scaling analysis..."
	./analyze_step2_workers.sh

# Compare Step 3 implementations (Sequential vs Unbounded vs Pool)
compare-step3: build
	@echo "Running comprehensive Step 3 analysis..."
	./analyze_step3_performance.sh

# Run all tests
test-all: test-hash test-groups test-seq-simple

# Check for race conditions
race: build
	@echo "=========================================="
	@echo "Testing for Race Conditions"
	@echo "=========================================="
	go run -race $(BINARY).go -hash-workers=4 -data-workers=4 -comp-workers=4 -input=input/simple.txt

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY)
	@echo "Clean complete!"

# Show help
help:
	@echo "Available commands:"
	@echo ""
	@echo "Build:"
	@echo "  make build                      - Build the BST program"
	@echo ""
	@echo "Sequential Testing:"
	@echo "  make test-seq                   - Test sequential solution on all input files"
	@echo "  make test-seq-simple            - Test sequential on simple.txt"
	@echo "  make test-seq-coarse            - Test sequential on coarse.txt"
	@echo "  make test-seq-fine              - Test sequential on fine.txt"
	@echo ""
	@echo "Step 1 (Hash Computation):"
	@echo "  make test-hash                  - Test Step 1 basic"
	@echo "  make compare-strategies         - Compare hash strategies (simple.txt)"
	@echo "  make compare-strategies-coarse  - Compare hash strategies (coarse.txt)"
	@echo "  make compare-fine               - Detailed fine.txt analysis (1,8,64,128,1k,10k)"
	@echo "  make compare-step1-complete     - Complete: seq/per-BST/pool (2,4,8,16,32)"
	@echo ""
	@echo "Step 2 (Hash Groups):"
	@echo "  make test-groups                - Test Step 2 basic"
	@echo "  make compare-step2              - Compare seq/channel/mutex on all inputs"
	@echo "  make compare-step2-scaling      - Worker scaling analysis (2,4,6,8,16,32)"
	@echo ""
	@echo "Step 3 (Tree Comparison):"
	@echo "  make compare-step3              - Compare seq/unbounded/pool on all inputs"
	@echo ""
	@echo "Other:"
	@echo "  make test-all                   - Run all tests"
	@echo "  make race                       - Check for race conditions"
	@echo "  make clean                      - Remove build artifacts"
	@echo "  make help                       - Show this help message"

# Default target
.DEFAULT_GOAL := help

