package main

import "fmt"
import "os"
import "strconv"
import "math"
import "time"

const outputInfo = true

type FuncCall struct {
	idx int
	capacity int
	k *Knapsack
	curValue int
}

func worker(work chan FuncCall, res chan int) {
	// While work channel is open, invoke a knapsack_rec and send the result back
	for inv := range work {
		res <- inv.curValue + Knapsack_rec(inv.idx, inv.capacity, inv.k)
	}
}

func gatherer(res chan int, totalworkchan chan int, finalres chan int) {
	work := 0
	maxRes := 0
	totalwork := math.MaxInt64
	running := true
	tick := time.Tick(1*time.Second)
	for running {
		select {
			case r := <-res: {
				maxRes = Max(maxRes, r)
				work++
				if work>=totalwork {
					running = false
				}
			}
			case tw := <- totalworkchan: {
				totalwork = tw
				if work >= totalwork {
					running = false
				}
			}
			case now := <- tick: {
				if outputInfo {
					if totalwork < math.MaxInt64 {
						fmt.Println("Gatherer state:", work, "items,", totalwork, "to do (", now, ")")
					} else {
						fmt.Println("Gatherer state:", work, "items, ? to do (", now, ")")
					}
				}
			}
		}
	}

	if outputInfo {
		fmt.Println("Gatherer total work:", work)
	}
	finalres <- maxRes
	close(finalres)
}

func main() {
	// Parse number of worker threads
	numWorkerThreads := 1
	if len(os.Args) >= 2 {
		if n, err := strconv.Atoi(os.Args[1]); err != nil {
			fmt.Println(err)
		} else {
			numWorkerThreads = n
		}
	}

	// Read the knapsack from command line
	k := ReadKnapsack()

	if outputInfo {
		fmt.Println("Knapsack size:", k.Len(), "capacity:", k.Capacity)
		fmt.Println("Num threads:", numWorkerThreads)
	}


	// Spawn workers
	// Add some buffer in the channels so we won't end up with send-blocking threads
	send := make(chan FuncCall, numWorkerThreads*100)
	intermediate_res := make(chan int, numWorkerThreads)
	for i := 0; i < numWorkerThreads; i++ {
		go worker(send, intermediate_res)
	}
	// Spawn gatherer
	totalwork := make(chan int)
	finalres := make(chan int)
	go gatherer(intermediate_res, totalwork, finalres)

	workCount := 0
	// Loop through each element
	// Assume all prior elements zero and search space for this one
	for i := 0; i < k.Len(); i++ {
		for n := 1; n <= k.Get(i).MaxCount; n++ {
			workCount++
			send <- FuncCall{i+1, // Let it search from the next index on
					k.Capacity - n * k.Get(i).Weight,
					k,
					n *  k.Get(i).Value}
		}
	}

	// Collect results
	close(send);
	totalwork <- workCount
	maxRes := <-finalres
	fmt.Println(maxRes)
}


/*func main() {
	numWorkerThreads := 1
	numMasterDimensions := 4
	

	// Parse number of worker threads
	if len(os.Args) >= 2 {
		if n, err := strconv.Atoi(os.Args[1]); err != nil {
			fmt.Println(err)
		} else {
			numWorkerThreads = n
		}
	}

	// Parse number of dimensions
	if len(os.Args) >= 3 {
		if n, err := strconv.Atoi(os.Args[2]); err != nil {
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
	// Add some buffer in the channels so we won't end up with send-blocking threads
	send := make(chan FuncCall, numWorkerThreads*10)
	intermediate_res := make(chan int, numWorkerThreads)
	for i := 0; i < numWorkerThreads; i++ {
		go worker(send, intermediate_res)
	}
	// Spawn gatherer
	totalwork := make(chan int)
	finalres := make(chan int)
	go gatherer(intermediate_res, totalwork, finalres)


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
	totalwork <- workCount
	maxRes := <-finalres
	fmt.Println(maxRes)
}*/