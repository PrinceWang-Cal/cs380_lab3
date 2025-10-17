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

// ComputeHashesParallel computes hashes using multiple goroutines
func ComputeHashesParallel(bsts []*BST, numWorkers int) map[int]int {
	// TODO: Implement parallel hash computation
	// TODO: Create numWorkers goroutines
	// TODO: Divide BSTs among workers
	// TODO: Use sync.WaitGroup to wait for all goroutines
	// TODO: Store results in a thread-safe manner
	
	hashes := make(map[int]int)
	
	// Hint: You might want to use a mutex or channels here
	
	return hashes
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
	// TODO: Compute hashes and group BSTs by hash value
	// TODO: Return map from hash to list of BST IDs
	hashGroups := make(map[int][]int)
	for _, bst := range bsts {
		hashGroups[bst.ComputeHash()] = append(hashGroups[bst.ComputeHash()], bst.ID)
	}
	return hashGroups
}

// BuildHashGroupsChannel builds hash groups using channel-based coordination
func BuildHashGroupsChannel(bsts []*BST, numHashWorkers int) map[int][]int {
	// TODO: Spawn numHashWorkers goroutines to compute hashes
	// TODO: Each worker sends (hash, BST ID) pairs to a channel
	// TODO: Spawn 1 central manager goroutine to receive from channel and update map
	// TODO: Use sync.WaitGroup for synchronization
	
	hashGroups := make(map[int][]int)
	
	// Hint: Create a channel for (hash, id) pairs
	// type HashResult struct {
	// 	Hash int
	// 	ID   int
	// }
	
	return hashGroups
}

// BuildHashGroupsMutex builds hash groups using mutex-protected map
func BuildHashGroupsMutex(bsts []*BST, numWorkers int) map[int][]int {
	// TODO: Spawn numWorkers goroutines to compute hashes
	// TODO: Each worker updates the shared map after acquiring mutex
	// TODO: Use sync.WaitGroup for synchronization
	
	hashGroups := make(map[int][]int)
	var mu sync.Mutex
	_ = mu // Use the mutex
	
	return hashGroups
}

// ====================
// Step 3: Tree Comparison
// ====================

// CompareTreesSequential compares trees with matching hashes sequentially
func CompareTreesSequential(bsts []*BST, hashGroups map[int][]int) [][]int {
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
	
	// Build connected components from adjacency matrix
	visited := make([]bool, n)
	var equivalenceGroups [][]int
	
	for i := 0; i < n; i++ {
		if !visited[i] {
			// Start a new group with BFS/DFS
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

// CompareTreesParallelUnbounded spawns a goroutine for each comparison
func CompareTreesParallelUnbounded(bsts []*BST, hashGroups map[int][]int) [][]int {
	// TODO: For each pair of trees with same hash, spawn a goroutine to compare
	// TODO: Use a 2D adjacency matrix or other structure to track equivalence
	// TODO: Build equivalence groups from comparison results
	// TODO: Use sync.WaitGroup to wait for all comparisons
	
	var equivalenceGroups [][]int
	return equivalenceGroups
}

// CompareTreesParallelPool uses a fixed pool of worker goroutines
func CompareTreesParallelPool(bsts []*BST, hashGroups map[int][]int, numWorkers int) [][]int {
	// TODO: Create a concurrent buffer (channel) for work items
	// TODO: Work item = (BST ID 1, BST ID 2) pair to compare
	// TODO: Spawn numWorkers goroutines to process work from channel
	// TODO: Main thread adds work items to channel
	// TODO: Workers update adjacency matrix when trees match
	// TODO: Build equivalence groups from adjacency matrix
	
	var equivalenceGroups [][]int
	return equivalenceGroups
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
	
	flag.Parse()
	
	hashWorkers := *hashWorkersPtr
	dataWorkers := *dataWorkersPtr
	compWorkers := *compWorkersPtr
	inputFile := *inputFilePtr
	
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
			_ = ComputeHashesParallel(bsts, hashWorkers)
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
		
		// Then, compare trees
		compareStart := time.Now()
		var equivalenceGroups [][]int
		
		if compWorkers == 1 {
			// Sequential tree comparison
			equivalenceGroups = CompareTreesSequential(bsts, hashGroups)
		} else {
			// Parallel tree comparison with worker pool
			equivalenceGroups = CompareTreesParallelPool(bsts, hashGroups, compWorkers)
		}
		
		compareElapsed := time.Since(compareStart)
		
		// Output results
		PrintHashGroupTime(hashGroupElapsed)
		PrintHashGroups(hashGroups)
		PrintCompareTreeTime(compareElapsed)
		PrintTreeGroups(equivalenceGroups)
		return
	}
	
	fmt.Println("Invalid flag combination")
}

