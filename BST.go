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

// bst

type TreeNode struct {
	Value int
	Left  *TreeNode
	Right *TreeNode
}

type BST struct {
	Root *TreeNode
	ID   int
}

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

func (t *TreeNode) InOrderTraversal() []int {
	if t == nil {
		return []int{}
	}
	left := t.Left.InOrderTraversal()
	right := t.Right.InOrderTraversal()
	return append(append(left, t.Value), right...)
}

func (bst *BST) ComputeHash() int {
	hash := 1
	values := bst.Root.InOrderTraversal()
	for _, value := range values {
		newValue := value + 2
		hash = (hash * newValue + newValue) % 1000
	}
	return hash
}

func AreEqual(bst1, bst2 *BST) bool {
	values1 := bst1.Root.InOrderTraversal()
	values2 := bst2.Root.InOrderTraversal()
	return reflect.DeepEqual(values1, values2)
}

// parse inputs
func ParseInputFile(filename string) ([]*BST, error) {
	data, err := ioutil.ReadFile(filename)
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

// core functions
// 1. Computation of Hash of Each Tree (sequential)
func ComputeHashesSequential(bsts []*BST) map[int]int {
	hashes := make(map[int]int)
	for _, bst := range bsts {
		hashes[bst.ID] = bst.ComputeHash()
	}
	return hashes
}

// 1. Computation of Hash of Each Tree (spawn per bst tree)
func ComputeHashesParallelPerBST(bsts []*BST) map[int]int {
	hashes := make(map[int]int)
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	for _, bst := range bsts {
		wg.Add(1)
		go func(b *BST) {
			defer wg.Done()
			hash := b.ComputeHash()
			mu.Lock()
			hashes[b.ID] = hash
			mu.Unlock()
		}(bst)
	}
	
	wg.Wait()
	return hashes
}

// 1. Computation of Hash of Each Tree (worker pool)
func ComputeHashesParallelWorkerPool(bsts []*BST, numWorkers int) map[int]int {
	hashes := make(map[int]int)
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	bstsPerWorker := len(bsts) / numWorkers
	if bstsPerWorker == 0 {
		bstsPerWorker = 1
	}
	
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		
		go func(workerID int) {
			defer wg.Done()
			
			start := workerID * bstsPerWorker
			end := start + bstsPerWorker
			
			if workerID == numWorkers-1 {
				end = len(bsts)
			}
			
			if start >= len(bsts) {
				return
			}
			if end > len(bsts) {
				end = len(bsts)
			}
			
			for j := start; j < end; j++ {
				hash := bsts[j].ComputeHash()
				mu.Lock()
				hashes[bsts[j].ID] = hash
				mu.Unlock()
			}
		}(i)
	}
	
	wg.Wait()
	return hashes
}

func ComputeHashesParallel(bsts []*BST, numWorkers int, perBSTStrategy bool) map[int]int {
	if perBSTStrategy {
		return ComputeHashesParallelPerBST(bsts)
	}
	return ComputeHashesParallelWorkerPool(bsts, numWorkers)
}

type HashGroup struct {
	Hash    int
	TreeIDs []int
}

// 2. Computation of Hash Groups (sequential)
func BuildHashGroupsSequential(bsts []*BST) map[int][]int {
	hashGroups := make(map[int][]int)
	for _, bst := range bsts {
		hash := bst.ComputeHash()
		hashGroups[hash] = append(hashGroups[hash], bst.ID)
	}
	return hashGroups
}

type HashResult struct {
	Hash int
	ID   int
}

// 2. Computation of Hash Groups (channel based)
func BuildHashGroupsChannel(bsts []*BST, numHashWorkers int) map[int][]int {
	hashGroups := make(map[int][]int)
	resultChan := make(chan HashResult, numHashWorkers)
	
	var wg sync.WaitGroup
	
	bstsPerWorker := len(bsts) / numHashWorkers
	if bstsPerWorker == 0 {
		bstsPerWorker = 1
	}
	
	for i := 0; i < numHashWorkers; i++ {
		wg.Add(1)
		
		go func(workerID int) {
			defer wg.Done()
			
			start := workerID * bstsPerWorker
			end := start + bstsPerWorker
			if workerID == numHashWorkers-1 {
				end = len(bsts)
			}
			if start >= len(bsts) {
				return
			}
			if end > len(bsts) {
				end = len(bsts)
			}
			for j := start; j < end; j++ {
				hash := bsts[j].ComputeHash()
				resultChan <- HashResult{Hash: hash, ID: bsts[j].ID}
			}
		}(i)
	}
	
	done := make(chan bool)
	go func() {
		for result := range resultChan {
			hashGroups[result.Hash] = append(hashGroups[result.Hash], result.ID)
		}
		done <- true
	}()
	
	wg.Wait()
	
	close(resultChan)
	
	<-done
	
	return hashGroups
}

// 2. Computation of Hash Groups (mutex based)
func BuildHashGroupsMutex(bsts []*BST, numWorkers int) map[int][]int {
	hashGroups := make(map[int][]int)
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	bstsPerWorker := len(bsts) / numWorkers
	if bstsPerWorker == 0 {
		bstsPerWorker = 1
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		
		go func(workerID int) {
			defer wg.Done()
			
			start := workerID * bstsPerWorker
			end := start + bstsPerWorker
			
			if workerID == numWorkers-1 {
				end = len(bsts)
			}
			
			if start >= len(bsts) {
				return
			}
			if end > len(bsts) {
				end = len(bsts)
			}
			
			for j := start; j < end; j++ {
				hash := bsts[j].ComputeHash()
				id := bsts[j].ID
	
				mu.Lock()
				hashGroups[hash] = append(hashGroups[hash], id)
				mu.Unlock()
			}
		}(i)
	}
	
	wg.Wait()
	return hashGroups
}


func BuildEquivalenceGroupsFromMatrix(adjMatrix [][]bool) [][]int {
	n := len(adjMatrix)
	visited := make([]bool, n)
	var equivalenceGroups [][]int
	
	for i := 0; i < n; i++ {
		if !visited[i] {
			group := []int{}
			queue := []int{i}
			visited[i] = true
			
			for len(queue) > 0 {
				current := queue[0]
				queue = queue[1:]
				group = append(group, current)
				
				for j := 0; j < n; j++ {
					if adjMatrix[current][j] && !visited[j] {
						visited[j] = true
						queue = append(queue, j)
					}
				}
			}
			
			if len(group) > 1 {
				equivalenceGroups = append(equivalenceGroups, group)
			}
		}
	}
	
	return equivalenceGroups
}

// 3. Tree Comparison (sequential)
func CompareTreesSequential(bsts []*BST, hashGroups map[int][]int) [][]bool {
	n := len(bsts)
	
	adjMatrix := make([][]bool, n)
	for i := 0; i < n; i++ {
		adjMatrix[i] = make([]bool, n)
		adjMatrix[i][i] = true
	}
	
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

// 3. Tree Comparison (parallel unbounded)
func CompareTreesParallelUnbounded(bsts []*BST, hashGroups map[int][]int) [][]bool {
	n := len(bsts)
	
	adjMatrix := make([][]bool, n)
	for i := 0; i < n; i++ {
		adjMatrix[i] = make([]bool, n)
		adjMatrix[i][i] = true
	}
	
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	for _, hashGroup := range hashGroups {
		if len(hashGroup) > 1 {
			for i := 0; i < len(hashGroup); i++ {
				for j := i + 1; j < len(hashGroup); j++ {
					wg.Add(1)
					
					go func(id1, id2 int) {
						defer wg.Done()
						
						if AreEqual(bsts[id1], bsts[id2]) {
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
	wg.Wait()
	return adjMatrix
}

type ComparisonWork struct {
	ID1 int
	ID2 int
}

// 3. Tree Comparison (worker pool)
func CompareTreesParallelPool(bsts []*BST, hashGroups map[int][]int, numWorkers int) [][]bool {
	n := len(bsts)
	
	adjMatrix := make([][]bool, n)
	for i := 0; i < n; i++ {
		adjMatrix[i] = make([]bool, n)
		adjMatrix[i][i] = true
	}
	
	workChan := make(chan ComparisonWork, numWorkers)
	
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			
			for work := range workChan {
				if AreEqual(bsts[work.ID1], bsts[work.ID2]) {
					mu.Lock()
					adjMatrix[work.ID1][work.ID2] = true
					adjMatrix[work.ID2][work.ID1] = true
					mu.Unlock()
				}
			}
		}(i)
	}
	
	for _, hashGroup := range hashGroups {
		if len(hashGroup) > 1 {
			for i := 0; i < len(hashGroup); i++ {
				for j := i + 1; j < len(hashGroup); j++ {
					workChan <- ComparisonWork{
						ID1: hashGroup[i],
						ID2: hashGroup[j],
					}
				}
			}
		}
	}
	
	close(workChan)
	wg.Wait()
	
	return adjMatrix
}

// Output Functions
func PrintHashTime(elapsed time.Duration) {
	fmt.Printf("hashTime: %.6f\n", elapsed.Seconds())
}
func PrintHashGroups(hashGroups map[int][]int) {
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

func PrintHashGroupTime(elapsed time.Duration) {
	fmt.Printf("hashGroupTime: %.6f\n", elapsed.Seconds())
}
func PrintTreeGroups(equivalenceGroups [][]int) {
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

func PrintCompareTreeTime(elapsed time.Duration) {
	fmt.Printf("compareTreeTime: %.6f\n", elapsed.Seconds())
}

func main() {
	// Parse flags
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
	
	bsts, err := ParseInputFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}
	
	// step 1
	if dataWorkers == 0 && compWorkers == 0 {
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
	
	// step 2
	if dataWorkers > 0 && compWorkers == 0 {
		start := time.Now()
		var hashGroups map[int][]int
		
		if hashWorkers == 1 && dataWorkers == 1 {
			// sequential
			hashGroups = BuildHashGroupsSequential(bsts)
		} else if hashWorkers > 1 && dataWorkers == 1 {
			// channel-based: hash workers send to central manager
			hashGroups = BuildHashGroupsChannel(bsts, hashWorkers)
		} else if hashWorkers > 1 && dataWorkers == hashWorkers {
			// mutex-based: each worker updates map with mutex
			hashGroups = BuildHashGroupsMutex(bsts, hashWorkers)
		} else {
			fmt.Println("Optional implementation not completed")
			return
		}
		
		elapsed := time.Since(start)
		PrintHashGroupTime(elapsed)
		PrintHashGroups(hashGroups)
		return
	}
	
	// step 3
	if dataWorkers > 0 && compWorkers > 0 {
		hashGroupStart := time.Now()
		var hashGroups map[int][]int
		
		if hashWorkers == 1 && dataWorkers == 1 {
			// sequential
			hashGroups = BuildHashGroupsSequential(bsts)
		} else if hashWorkers > 1 && dataWorkers == 1 {
			// channel-based
			hashGroups = BuildHashGroupsChannel(bsts, hashWorkers)
		} else if hashWorkers > 1 && dataWorkers == hashWorkers {
			// mutex-based
			hashGroups = BuildHashGroupsMutex(bsts, hashWorkers)
		} else {
			fmt.Println("Optional implementation not completed")
			return
		}
		
		hashGroupElapsed := time.Since(hashGroupStart)
		
		// compare trees and populate adjacency matrix (timed)
		compareStart := time.Now()
		var adjMatrix [][]bool
		
		if compWorkers == 1 {
			// sequential
			adjMatrix = CompareTreesSequential(bsts, hashGroups)
		} else if compStrategy == "unbounded" {
			// parallel unbounded
			adjMatrix = CompareTreesParallelUnbounded(bsts, hashGroups)
		} else {
			// parallel worker pool
			adjMatrix = CompareTreesParallelPool(bsts, hashGroups, compWorkers)
		}
		
		compareElapsed := time.Since(compareStart)
		
		equivalenceGroups := BuildEquivalenceGroupsFromMatrix(adjMatrix)
		
		PrintHashGroupTime(hashGroupElapsed)
		PrintHashGroups(hashGroups)
		PrintCompareTreeTime(compareElapsed)
		PrintTreeGroups(equivalenceGroups)
		return
	}
	
	fmt.Println("Invalid flag combination")
}

