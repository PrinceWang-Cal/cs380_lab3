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
	@echo "  make build           - Build the BST program"
	@echo "  make test-seq        - Test sequential solution on all input files"
	@echo "  make test-seq-simple - Test sequential on simple.txt"
	@echo "  make test-seq-coarse - Test sequential on coarse.txt"
	@echo "  make test-seq-fine   - Test sequential on fine.txt"
	@echo "  make test-hash       - Test Step 1 (hash computation)"
	@echo "  make test-groups     - Test Step 2 (hash groups)"
	@echo "  make test-all        - Run all tests"
	@echo "  make race            - Check for race conditions"
	@echo "  make clean           - Remove build artifacts"
	@echo "  make help            - Show this help message"

# Default target
.DEFAULT_GOAL := help

