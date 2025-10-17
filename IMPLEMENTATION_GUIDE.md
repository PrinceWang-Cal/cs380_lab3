# BST Assignment Implementation Guide

## Quick Start

1. **Compile:**
   ```bash
   go build BST.go
   ```

2. **Run tests:**
   ```bash
   chmod +x test.sh
   ./test.sh
   ```

## Implementation Order

### Phase 1: Basic Infrastructure (Start Here!)
1. **ParseInputFile()** - Read and parse input file
   - Use `ioutil.ReadFile()` to read file
   - Split by newlines using `strings.Split()`
   - For each line, split by spaces and parse integers with `strconv.Atoi()`
   - Create BST and insert values in order

2. **Insert()** - BST insertion
   - If current node is nil, create new TreeNode with value
   - If value < current.Value, insert to left subtree
   - Otherwise, insert to right subtree

3. **InOrderTraversal()** - Traverse BST
   - Recursive: left subtree → current node → right subtree
   - Return slice of integers in sorted order

4. **ComputeHash()** - Hash function (MUST BE EXACT)
   ```go
   hash := 1
   values := bst.Root.InOrderTraversal()
   for _, value := range values {
       newValue := value + 2
       hash = (hash * newValue + newValue) % 1000
   }
   return hash
   ```

5. **AreEqual()** - Compare two BSTs
   - Get in-order traversals of both trees
   - Compare if slices are equal

### Phase 2: Step 1 - Hash Computation
6. **ComputeHashesSequential()**
   - Simple loop through all BSTs, compute hash for each

7. **ComputeHashesParallel()**
   - Create `sync.WaitGroup`
   - Divide BSTs among workers (consider using work distribution)
   - Each goroutine computes hashes for its subset
   - Use mutex or separate maps per worker to avoid race conditions
   - Wait for all goroutines to complete

### Phase 3: Step 2 - Hash Groups
8. **BuildHashGroupsSequential()**
   - Compute hash for each BST
   - Build map[hash][]treeID

9. **BuildHashGroupsChannel()**
   - Create channel for (hash, id) pairs
   - Spawn hash worker goroutines that:
     - Compute hashes for their assigned BSTs
     - Send results to channel
   - Spawn 1 manager goroutine that:
     - Receives from channel
     - Updates map (no mutex needed, only one writer)
   - Close channel when all workers done

10. **BuildHashGroupsMutex()**
    - Create shared map and mutex
    - Spawn worker goroutines that:
      - Compute hashes for assigned BSTs
      - Lock mutex, update map, unlock mutex
    - Use WaitGroup to wait for completion

### Phase 4: Step 3 - Tree Comparison
11. **CompareTreesSequential()**
    - For each hash group with multiple trees
    - Compare all pairs of trees
    - Build equivalence groups (union-find or adjacency matrix)

12. **CompareTreesParallelPool()**
    - Create work channel (buffer of comparison tasks)
    - Create result structure (adjacency matrix or similar)
    - Spawn worker goroutines that:
      - Read (id1, id2) pairs from channel
      - Compare trees[id1] with trees[id2]
      - Update result structure if equal (with mutex protection)
    - Main thread enqueues all comparison tasks
    - Close channel when done, wait for workers
    - Build equivalence groups from results

## Testing Strategy

1. **Start with simple.txt** - small, easy to debug with print statements
2. **Verify correctness** - hash values and groups should match between implementations
3. **Test with coarse.txt** - should show good speedup with parallelism
4. **Test with fine.txt** - many small trees, tests overhead

## Common Pitfalls

1. **Hash function MUST be exact** - autograder checks correctness
2. **Output format MUST match** - check spacing, colons, time format
3. **Race conditions** - use `-race` flag: `go run -race BST.go`
4. **Empty lines in input** - handle gracefully
5. **Time measurement** - exclude tree construction, include only what's specified
6. **Groups with single tree** - do NOT print these

## Debugging Tips

```bash
# Check for race conditions
go run -race BST.go -hash-workers=4 -input=input/simple.txt

# Verbose output (add your own debug prints)
# But remove before submission!

# Compare outputs
./BST -hash-workers=1 -data-workers=1 -input=input/simple.txt > out1.txt
./BST -hash-workers=4 -data-workers=1 -input=input/simple.txt > out2.txt
diff out1.txt out2.txt  # Should be same except timing
```

## Example Expected Output

### Step 1 (hash only):
```
hashTime: 0.000123
```

### Step 2 (hash groups):
```
hashGroupTime: 0.000234
156: 1 5 7
892: 2 3 4
```

### Step 3 (with tree comparison):
```
hashGroupTime: 0.000234
156: 1 5 7
892: 2 3 4
compareTreeTime: 0.000156
group 0: 2 3 4
```

## Useful Go Snippets

### WaitGroup
```go
var wg sync.WaitGroup
for i := 0; i < numWorkers; i++ {
    wg.Add(1)
    go func(workerID int) {
        defer wg.Done()
        // do work
    }(i)
}
wg.Wait()
```

### Channel with goroutine
```go
resultChan := make(chan Result, 100)
go func() {
    for result := range resultChan {
        // process result
    }
}()
// ... send to channel ...
close(resultChan)
```

### Mutex
```go
var mu sync.Mutex
mu.Lock()
// critical section
mu.Unlock()
```

## Submission Checklist

- [ ] All TODOs implemented
- [ ] Tested on all three input files
- [ ] Outputs match expected format exactly
- [ ] No race conditions (`go run -race` passes)
- [ ] Removed all debug print statements
- [ ] File named `BST.go` (or `BST_opt.go` if optional completed)
- [ ] Code + report packaged in tar file

