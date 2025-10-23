package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"reflect"
	"sort"
	"strings"
	"strconv"
	"sync"
	"time"
)

// ====================
// BST Node Structure
// ====================

type TreeNode struct {
	Value int
	Left  *TreeNode
	Right *TreeNode
}

type BST struct {
	Root *TreeNode
	ID   int
}

// ====================
// BST Operations
// ====================

// Insert adds a value to the BST
func (t *TreeNode) Insert(value int) *TreeNode {
	if t == nil {
		return &TreeNode{Value: value}
	}
	if value < t.Value {
		t.Left = t.Left.Insert(value)
	} else {
		t.Right = t.Right.Insert(value)
	}
	return t
}

// InOrderTraversal returns values in sorted order
func (t *TreeNode) InOrderTraversal() []int {
	// TODO: Implement in-order traversal (left, root, right)
	// Hint: Use recursion or iterative approach with stack
	if t == nil {
		return []int{}
	}
	left := t.Left.InOrderTraversal()
	right := t.Right.InOrderTraversal()
	return append(append(left, t.Value), right...)
}

// ComputeHash computes the hash of a BST using in-order traversal
func (bst *BST) ComputeHash() int {
	// TODO: Implement the exact hash function from submission guide:
	// hash = 1
	// for each value in tree.in_order_traversal() {
	//   new_value = value + 2
	//   hash = (hash * new_value + new_value) % 1000
	// }
	hash := 1
	values := bst.Root.InOrderTraversal()
	for _, value := range values {
		newValue := value + 2
		hash = (hash * newValue + newValue) % 1000
	}
	return hash
}

// AreEqual checks if two BSTs contain the same values in the same order
func AreEqual(bst1, bst2 *BST) bool {
	// TODO: Implement tree equality check
	// Hint: Compare in-order traversals
	values1 := bst1.Root.InOrderTraversal()
	values2 := bst2.Root.InOrderTraversal()
	return reflect.DeepEqual(values1, values2)
}

// ====================
// File Parsing
// ====================

// ParseInputFile reads the input file and constructs BSTs
func ParseInputFile(filename string) ([]*BST, error) {
	// TODO: Read file using ioutil.ReadFile
	// TODO: Split by newlines, parse each line as a BST
	// TODO: For each line, split by spaces and insert values in order
	
	data, err := ioutil.ReadFile(filename)
	// TODO: Remove this line if it is not needed
	if err != nil {
		return nil, err
	}
	
	var bsts []*BST
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	
	for id, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
	
		// Parse space-separated integers from line
		// Create BST and insert values in order
		// Set BST ID to the line index
		bst := &BST{ID: id}
		values := strings.Split(line, " ")
		for _, value := range values {
			intValue, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			bst.Root = bst.Root.Insert(intValue)
		}
		bsts = append(bsts, bst)
	}
	
	return bsts, nil
}

// ====================
// Step 1: Hash Computation
// ====================

// ComputeHashesSequential computes hashes for all BSTs in main thread
func ComputeHashesSequential(bsts []*BST) map[int]int {
	// TODO: Iterate through all BSTs and compute their hashes
	// TODO: Return a map from BST ID to hash value
	hashes := make(map[int]int)
	for _, bst := range bsts {
		hashes[bst.ID] = bst.ComputeHash()
	}
	return hashes
}

// ComputeHashesParallelPerBST spawns one goroutine per BST
func ComputeHashesParallelPerBST(bsts []*BST) map[int]int {
	hashes := make(map[int]int)
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	// Spawn a goroutine for each BST
	for _, bst := range bsts {
		wg.Add(1)
		go func(b *BST) {
			defer wg.Done()
			hash := b.ComputeHash()
			
			// Safely store result
			mu.Lock()
			hashes[b.ID] = hash
			mu.Unlock()
		}(bst)
	}
	
	wg.Wait()
	return hashes
}

// ComputeHashesParallelWorkerPool uses fixed number of worker goroutines
func ComputeHashesParallelWorkerPool(bsts []*BST, numWorkers int) map[int]int {
	hashes := make(map[int]int)
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	// Calculate work distribution
	bstsPerWorker := len(bsts) / numWorkers
	if bstsPerWorker == 0 {
		bstsPerWorker = 1
	}
	
	// Spawn numWorkers goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		
		go func(workerID int) {
			defer wg.Done()
			
			// Calculate range for this worker
			start := workerID * bstsPerWorker
			end := start + bstsPerWorker
			
			// Last worker takes any remainder
			if workerID == numWorkers-1 {
				end = len(bsts)
			}
			
			// Make sure we don't go out of bounds
			if start >= len(bsts) {
				return
			}
			if end > len(bsts) {
				end = len(bsts)
			}
			
			// Process assigned BSTs
			for j := start; j < end; j++ {
				hash := bsts[j].ComputeHash()
				
				// Safely store result
				mu.Lock()
				hashes[bsts[j].ID] = hash
				mu.Unlock()
			}
		}(i)
	}
	
	wg.Wait()
	return hashes
}

// ComputeHashesParallel dispatcher function that chooses strategy
func ComputeHashesParallel(bsts []*BST, numWorkers int, perBSTStrategy bool) map[int]int {
	if perBSTStrategy {
		return ComputeHashesParallelPerBST(bsts)
	}
	return ComputeHashesParallelWorkerPool(bsts, numWorkers)
}

// ====================
// Step 2: Hash Groups
// ====================

// HashGroup represents trees with the same hash
type HashGroup struct {
	Hash    int
	TreeIDs []int
}

// BuildHashGroupsSequential builds hash groups in main thread
func BuildHashGroupsSequential(bsts []*BST) map[int][]int {
	hashGroups := make(map[int][]int)
	for _, bst := range bsts {
		hash := bst.ComputeHash()
		hashGroups[hash] = append(hashGroups[hash], bst.ID)
	}
	return hashGroups
}

// HashResult represents a hash computation result
type HashResult struct {
	Hash int
	ID   int
}

// BuildHashGroupsChannel builds hash groups using channel-based coordination
func BuildHashGroupsChannel(bsts []*BST, numHashWorkers int) map[int][]int {
	hashGroups := make(map[int][]int)
	
	// Create channel for hash workers to send results
	resultChan := make(chan HashResult, numHashWorkers)
	
	var wg sync.WaitGroup
	
	// Calculate work distribution
	bstsPerWorker := len(bsts) / numHashWorkers
	if bstsPerWorker == 0 {
		bstsPerWorker = 1
	}
	
	// Spawn hash worker goroutines
	for i := 0; i < numHashWorkers; i++ {
		wg.Add(1)
		
		go func(workerID int) {
			defer wg.Done()
			
			// Calculate range for this worker
			start := workerID * bstsPerWorker
			end := start + bstsPerWorker
			
			// Last worker takes remainder
			if workerID == numHashWorkers-1 {
				end = len(bsts)
			}
			
			// Bounds checking
			if start >= len(bsts) {
				return
			}
			if end > len(bsts) {
				end = len(bsts)
			}
			
			// Compute hashes and send to channel
			for j := start; j < end; j++ {
				hash := bsts[j].ComputeHash()
				resultChan <- HashResult{Hash: hash, ID: bsts[j].ID}
			}
		}(i)
	}
	
	// Spawn central manager goroutine to collect results
	done := make(chan bool)
	go func() {
		for result := range resultChan {
			// Only one goroutine modifies the map - no mutex needed!
			hashGroups[result.Hash] = append(hashGroups[result.Hash], result.ID)
		}
		done <- true
	}()
	
	// Wait for all hash workers to finish
	wg.Wait()
	
	// Close channel to signal manager to stop
	close(resultChan)
	
	// Wait for manager to finish processing all results
	<-done
	
	return hashGroups
}

// BuildHashGroupsMutex builds hash groups using mutex-protected map
func BuildHashGroupsMutex(bsts []*BST, numWorkers int) map[int][]int {
	hashGroups := make(map[int][]int)
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	// Calculate work distribution
	bstsPerWorker := len(bsts) / numWorkers
	if bstsPerWorker == 0 {
		bstsPerWorker = 1
	}
	
	// Spawn worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		
		go func(workerID int) {
			defer wg.Done()
			
			// Calculate range for this worker
			start := workerID * bstsPerWorker
			end := start + bstsPerWorker
			
			// Last worker takes remainder
			if workerID == numWorkers-1 {
				end = len(bsts)
			}
			
			// Bounds checking
			if start >= len(bsts) {
				return
			}
			if end > len(bsts) {
				end = len(bsts)
			}
			
			// Process assigned BSTs
			for j := start; j < end; j++ {
				hash := bsts[j].ComputeHash()
				id := bsts[j].ID
				
				// Acquire mutex, update map, release mutex
				mu.Lock()
				hashGroups[hash] = append(hashGroups[hash], id)
				mu.Unlock()
			}
		}(i)
	}
	
	wg.Wait()
	return hashGroups
}

// ====================
// Step 3: Tree Comparison
// ====================

// BuildEquivalenceGroupsFromMatrix builds equivalence groups from adjacency matrix using BFS
func BuildEquivalenceGroupsFromMatrix(adjMatrix [][]bool) [][]int {
	n := len(adjMatrix)
	visited := make([]bool, n)
	var equivalenceGroups [][]int
	
	for i := 0; i < n; i++ {
		if !visited[i] {
			// Start a new group with BFS
			group := []int{}
			queue := []int{i}
			visited[i] = true
			
			for len(queue) > 0 {
				current := queue[0]
				queue = queue[1:]
				group = append(group, current)
				
				// Find all trees equivalent to current
				for j := 0; j < n; j++ {
					if adjMatrix[current][j] && !visited[j] {
						visited[j] = true
						queue = append(queue, j)
					}
				}
			}
			
			// Only add groups with more than 1 tree
			if len(group) > 1 {
				equivalenceGroups = append(equivalenceGroups, group)
			}
		}
	}
	
	return equivalenceGroups
}

// CompareTreesSequential compares trees with matching hashes sequentially
// Returns adjacency matrix where adjMatrix[i][j] = true means BST i equals BST j
func CompareTreesSequential(bsts []*BST, hashGroups map[int][]int) [][]bool {
	n := len(bsts)
	
	// Create adjacency matrix to track equivalence
	adjMatrix := make([][]bool, n)
	for i := 0; i < n; i++ {
		adjMatrix[i] = make([]bool, n)
		adjMatrix[i][i] = true // Tree is equivalent to itself
	}
	
	// Compare all pairs of trees with the same hash
	for _, hashGroup := range hashGroups {
		if len(hashGroup) > 1 {
			for i := 0; i < len(hashGroup); i++ {
				for j := i + 1; j < len(hashGroup); j++ {
					id1 := hashGroup[i]
					id2 := hashGroup[j]
					if AreEqual(bsts[id1], bsts[id2]) {
						// Mark as equivalent (symmetrically)
						adjMatrix[id1][id2] = true
						adjMatrix[id2][id1] = true
					}
				}
			}
		}
	}
	
	return adjMatrix
}

// CompareTreesParallelUnbounded spawns a goroutine for each comparison
// Returns adjacency matrix where adjMatrix[i][j] = true means BST i equals BST j
func CompareTreesParallelUnbounded(bsts []*BST, hashGroups map[int][]int) [][]bool {
	n := len(bsts)
	
	// Create adjacency matrix to track equivalence
	adjMatrix := make([][]bool, n)
	for i := 0; i < n; i++ {
		adjMatrix[i] = make([]bool, n)
		adjMatrix[i][i] = true // Tree is equivalent to itself
	}
	
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	// Spawn a goroutine for each pair comparison
	for _, hashGroup := range hashGroups {
		if len(hashGroup) > 1 {
			// Compare all pairs within this hash group
			for i := 0; i < len(hashGroup); i++ {
				for j := i + 1; j < len(hashGroup); j++ {
					wg.Add(1)
					
					// Spawn goroutine for this comparison
					go func(id1, id2 int) {
						defer wg.Done()
						
						if AreEqual(bsts[id1], bsts[id2]) {
							// Mark as equivalent (symmetrically)
							mu.Lock()
							adjMatrix[id1][id2] = true
							adjMatrix[id2][id1] = true
							mu.Unlock()
						}
					}(hashGroup[i], hashGroup[j])
				}
			}
		}
	}
	
	// Wait for all comparisons to complete
	wg.Wait()
	
	return adjMatrix
}

// ComparisonWork represents a pair of BSTs to compare
type ComparisonWork struct {
	ID1 int
	ID2 int
}

// CompareTreesParallelPool uses a fixed pool of worker goroutines
// Returns adjacency matrix where adjMatrix[i][j] = true means BST i equals BST j
func CompareTreesParallelPool(bsts []*BST, hashGroups map[int][]int, numWorkers int) [][]bool {
	n := len(bsts)
	
	// Create adjacency matrix to track equivalence
	adjMatrix := make([][]bool, n)
	for i := 0; i < n; i++ {
		adjMatrix[i] = make([]bool, n)
		adjMatrix[i][i] = true // Tree is equivalent to itself
	}
	
	// Create buffered channel for work items (bounded buffer)
	// Buffer size = numWorkers to limit concurrent work
	workChan := make(chan ComparisonWork, numWorkers)
	
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	// Spawn worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			
			// Process work items from channel
			for work := range workChan {
				// Compare the two trees
				if AreEqual(bsts[work.ID1], bsts[work.ID2]) {
					// Mark as equivalent (symmetrically)
					mu.Lock()
					adjMatrix[work.ID1][work.ID2] = true
					adjMatrix[work.ID2][work.ID1] = true
					mu.Unlock()
				}
			}
		}(i)
	}
	
	// Main thread: enqueue all comparison work
	for _, hashGroup := range hashGroups {
		if len(hashGroup) > 1 {
			// Generate all pairs within this hash group
			for i := 0; i < len(hashGroup); i++ {
				for j := i + 1; j < len(hashGroup); j++ {
					// Send work to channel (blocks if buffer full)
					workChan <- ComparisonWork{
						ID1: hashGroup[i],
						ID2: hashGroup[j],
					}
				}
			}
		}
	}
	
	// Close channel to signal no more work
	close(workChan)
	
	// Wait for all workers to finish
	wg.Wait()
	
	return adjMatrix
}

// ====================
// Output Functions
// ====================

// PrintHashTime prints the hash computation time
func PrintHashTime(elapsed time.Duration) {
	fmt.Printf("hashTime: %.6f\n", elapsed.Seconds())
}

// PrintHashGroups prints hash groups (only groups with multiple trees)
func PrintHashGroups(hashGroups map[int][]int) {
	// TODO: Sort hashes for consistent output
	// TODO: Print only hash groups with more than 1 tree
	// Format: "hash: id0 id1 id2 ..."
	
	var hashes []int
	for hash := range hashGroups {
		if len(hashGroups[hash]) > 1 {
			hashes = append(hashes, hash)
		}
	}
	sort.Ints(hashes)
	
	for _, hash := range hashes {
		ids := hashGroups[hash]
		sort.Ints(ids)
		fmt.Printf("%d:", hash)
		for _, id := range ids {
			fmt.Printf(" %d", id)
		}
		fmt.Println()
	}
}

// PrintHashGroupTime prints the hash group computation time
func PrintHashGroupTime(elapsed time.Duration) {
	fmt.Printf("hashGroupTime: %.6f\n", elapsed.Seconds())
}

// PrintTreeGroups prints equivalence groups (only groups with multiple trees)
func PrintTreeGroups(equivalenceGroups [][]int) {
	// TODO: Filter out single-tree groups
	// TODO: Print in format: "group i: id0 id1 id2 ..."
	
	groupNum := 0
	for _, group := range equivalenceGroups {
		if len(group) > 1 {
			sort.Ints(group)
			fmt.Printf("group %d:", groupNum)
			for _, id := range group {
				fmt.Printf(" %d", id)
			}
			fmt.Println()
			groupNum++
		}
	}
}

// PrintCompareTreeTime prints the tree comparison time
func PrintCompareTreeTime(elapsed time.Duration) {
	fmt.Printf("compareTreeTime: %.6f\n", elapsed.Seconds())
}

// ====================
// Main Function
// ====================

func main() {
	// Parse command-line flags
	hashWorkersPtr := flag.Int("hash-workers", 1, "number of hash workers")
	dataWorkersPtr := flag.Int("data-workers", 0, "number of data workers")
	compWorkersPtr := flag.Int("comp-workers", 0, "number of comparison workers")
	inputFilePtr := flag.String("input", "", "path to input file")
	hashStrategyPtr := flag.String("hash-strategy", "pool", "hash strategy: 'perbst' (one goroutine per BST) or 'pool' (worker pool)")
	compStrategyPtr := flag.String("comp-strategy", "pool", "comparison strategy: 'unbounded' (goroutine per comparison) or 'pool' (worker pool)")
	
	flag.Parse()
	
	hashWorkers := *hashWorkersPtr
	dataWorkers := *dataWorkersPtr
	compWorkers := *compWorkersPtr
	inputFile := *inputFilePtr
	hashStrategy := *hashStrategyPtr
	compStrategy := *compStrategyPtr
	
	if inputFile == "" {
		fmt.Println("Error: -input flag is required")
		return
	}
	
	// Parse input file and construct BSTs
	bsts, err := ParseInputFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}
	
	// ====================
	// STEP 1: Hash Computation Only
	// ====================
	if dataWorkers == 0 && compWorkers == 0 {
		// Only compute hashes (no hash groups or tree comparison)
		start := time.Now()
		
		if hashWorkers == 1 {
			_ = ComputeHashesSequential(bsts)
		} else {
			perBSTStrategy := (hashStrategy == "perbst")
			_ = ComputeHashesParallel(bsts, hashWorkers, perBSTStrategy)
		}
		
		elapsed := time.Since(start)
		PrintHashTime(elapsed)
		return
	}
	
	// ====================
	// STEP 2: Hash Groups
	// ====================
	if dataWorkers > 0 && compWorkers == 0 {
		start := time.Now()
		var hashGroups map[int][]int
		
		if hashWorkers == 1 && dataWorkers == 1 {
			// Sequential implementation
			hashGroups = BuildHashGroupsSequential(bsts)
		} else if hashWorkers > 1 && dataWorkers == 1 {
			// Channel-based: hash workers send to central manager
			hashGroups = BuildHashGroupsChannel(bsts, hashWorkers)
		} else if hashWorkers > 1 && dataWorkers == hashWorkers {
			// Mutex-based: each worker updates map with mutex
			hashGroups = BuildHashGroupsMutex(bsts, hashWorkers)
		} else {
			// TODO: Optional implementation for hash-workers=i, data-workers=j (i>j>1)
			fmt.Println("Optional implementation not completed")
			return
		}
		
		elapsed := time.Since(start)
		PrintHashGroupTime(elapsed)
		PrintHashGroups(hashGroups)
		return
	}
	
	// ====================
	// STEP 3: Tree Comparison
	// ====================
	if dataWorkers > 0 && compWorkers > 0 {
		// First, build hash groups
		hashGroupStart := time.Now()
		var hashGroups map[int][]int
		
		if hashWorkers == 1 && dataWorkers == 1 {
			hashGroups = BuildHashGroupsSequential(bsts)
		} else if hashWorkers > 1 && dataWorkers == 1 {
			hashGroups = BuildHashGroupsChannel(bsts, hashWorkers)
		} else if hashWorkers > 1 && dataWorkers == hashWorkers {
			hashGroups = BuildHashGroupsMutex(bsts, hashWorkers)
		} else {
			// Optional implementation
			fmt.Println("Optional implementation not completed")
			return
		}
		
		hashGroupElapsed := time.Since(hashGroupStart)
		
		// Then, compare trees and populate adjacency matrix (TIMED)
		compareStart := time.Now()
		var adjMatrix [][]bool
		
		if compWorkers == 1 {
			// Sequential tree comparison
			adjMatrix = CompareTreesSequential(bsts, hashGroups)
		} else if compStrategy == "unbounded" {
			// Parallel tree comparison: goroutine per comparison
			adjMatrix = CompareTreesParallelUnbounded(bsts, hashGroups)
		} else {
			// Parallel tree comparison with worker pool
			adjMatrix = CompareTreesParallelPool(bsts, hashGroups, compWorkers)
		}
		
		compareElapsed := time.Since(compareStart)
		
		// Build equivalence groups from adjacency matrix (NOT TIMED)
		equivalenceGroups := BuildEquivalenceGroupsFromMatrix(adjMatrix)
		
		// Output results
		PrintHashGroupTime(hashGroupElapsed)
		PrintHashGroups(hashGroups)
		PrintCompareTreeTime(compareElapsed)
		PrintTreeGroups(equivalenceGroups)
		return
	}
	
	fmt.Println("Invalid flag combination")
}

