package main

import "fmt"
import "os"
import "strconv"



type FuncCall struct {
	idx int
	capacity int
	k *Knapsack
	curValue int
}

func worker(work chan FuncCall, res chan int) {
	// While work channel is open, invoke a knapsack_rec and send the result back
	var tmp FuncCall
	for inv := range work {
		tmp = inv
		res <- inv.curValue + Knapsack_rec(inv.idx, inv.capacity, inv.k)
	}
	fmt.Println("worker finished, last invocation", tmp)
}

func main() {
	numWorkerThreads := 1
	numMasterDimensions := 4
	outputInfo := true

	// Parse number of worker threads
	if len(os.Args) >= 2 {
		if n, err := strconv.Atoi(os.Args[1]); err != nil {
			fmt.Println(err)
		} else {
			numWorkerThreads = n
		}
	}

	// Parse number of worker threads
	if len(os.Args) >= 3 {
		if n, err := strconv.Atoi(os.Args[1]); err != nil {
			fmt.Println(err)
		} else {
			numMasterDimensions = n
		}
	}

	// Read the knapsack from command line
	k := ReadKnapsack()

	if k.Len() < numMasterDimensions {
		numMasterDimensions = k.Len()
	}

	// Initialize a galoisvec for the loop
	n := NewGaloisVec(numMasterDimensions)
	workload := 1
	for i:=0; i<numMasterDimensions; i++ {
		n.Limit[i] = k.Capacity / k.Items[i].Weight + 1
		workload *= n.Limit[i]
	}

	if outputInfo {
		fmt.Println("Knapsack size:", k.Len(), "capacity:", k.Capacity)
		fmt.Println("Workload:", workload, "num threads:", numWorkerThreads)
	}

	// Spawn workers
	// Make sure that the total capacity of both channels 
	// exceeds the number of possible assignments in the first numMasterDimensions dimensions
	// Otherwise: deadlock
	send := make(chan FuncCall, numWorkerThreads + workload)
	recv := make(chan int, numWorkerThreads + workload)
	for i := 0; i < numWorkerThreads; i++ {
		go worker(send, recv)
	}

	// Loop through the first dimensions and send work to the workers
	workCount := 0
	for {
		value := 0
		weight := 0
		for idx, multiplicator := range n.N {
			value += multiplicator * k.Items[idx].Value
			weight += multiplicator * k.Items[idx].Weight
		}

		if weight <= k.Capacity {
			workCount++
			send <- FuncCall{numMasterDimensions, k.Capacity-weight, k, value}
		}

		if n.Add(1) != 0 {
			break
		}
	}

	// Collect results
	close(send);
	maxRes := 0
	for ;workCount > 0; workCount-- {
		res := <-recv
		maxRes = Max(res, maxRes)
	}
	fmt.Println(maxRes)
}