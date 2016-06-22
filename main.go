// Idea: broadcast current best value
// Use canImprove to prune some search paths
package main

import "fmt"
import "os"
import "time"
import "flag"
import "runtime/pprof"

const outputInfo = true
const workCapacity = 10000

type FuncCall struct {
	idx int
	capacity int
	k *Knapsack
	curValue int
}

func (x FuncCall) isTail() bool {return x.idx >= x.k.Len() - x.k.Len()/10}
func isTail(idx int, len int) bool {return idx >= len - len/10}

// Checks if there is the possibility that this function call could theoretically improve the current best result
// Use the fact that items are ordered by density
func (x FuncCall) canImprove(curBestValue int) bool {
	return curBestValue == 0 || 
		   (int(float32(x.capacity) * x.k.Items[x.idx].Density) >= (curBestValue - x.curValue))
}

type ChanCollection struct {
	work chan FuncCall 			// Used to store work-packages
	res chan int 				// Used to send results to the gatherer
	workerfinished chan int 	// Used by the gatherer to terminate workers and to send preliminary results
	newwork chan int 			// Used to announce new work at the gatherer
	ack chan bool 				// Used by the gatherer to acknowledge new work
}


func knapsack_split(c *ChanCollection, call FuncCall) {
	c.newwork <- call.k.Len() - call.idx
	<-c.ack

	// Loop through each element
	// Assume all prior elements zero and search space for this one
	for i := call.k.Len()-1; i >= call.idx; i-- {
		// Start the search ommiting this element == 0 (will be covered by another call)
		if call.capacity - call.k.Get(i).Weight >= 0 {
			if len(c.work) < cap(c.work) && !isTail(i, call.k.Len()) {
				c.work <- FuncCall{i,
						call.capacity - call.k.Get(i).Weight,
						call.k,
						call.k.Get(i).Value + call.curValue}
			} else { // If the work queue is full, do some work right now
				c.res <- knapsack_rec(c,FuncCall{i, 
												call.capacity - call.k.Get(i).Weight,
												call.k,
												call.curValue + call.k.Get(i).Value})
			}
		} else {
			c.res <- call.curValue
		}
	}
}

func knapsack_rec(c *ChanCollection, call FuncCall) int {
	if call.idx >= call.k.Len() {
		return call.curValue
	}

	// If we still have a lot to do, do some dimensions with the _dim algo
	spawn := false
	/*if call.idx < call.k.Len()/3 {
		limit := call.capacity / call.k.Get(call.idx).Weight + 1
		c.newwork <- limit
		<-c.ack
		spawn = true
	}*/

	var v, w, r int
	v, r = call.curValue, call.curValue
	for call.capacity-w >= 0 {
		if !spawn {
			r = Max(r, knapsack_rec(c, FuncCall{call.idx+1, call.capacity-w, call.k, v}))
		} else {
			c.work <- FuncCall{call.idx+1, call.capacity-w, call.k, v}
		}
		v += call.k.Get(call.idx).Value
		w += call.k.Get(call.idx).Weight
	}
	return r
}

func knapsack_dim(c *ChanCollection, call FuncCall) int {
	numDimensions := 70
	until := call.idx + numDimensions
	allDims := false
	if until >= call.k.Len() {
		until = call.k.Len()-1
		numDimensions = until-call.idx
		allDims = true
	}

	// Initialize a galoisvec for the loop
	n := NewGaloisVec(numDimensions)
	for i:=call.idx; i<until; i++ {
		n.Limit[i-call.idx] = call.capacity / call.k.Items[i].Weight + 1
	}

	// Allocate work in heaps
	workFree := 0

	var r = call.curValue
	for {
		value := call.curValue
		weight := 0

		// Loop through items starting from the right
		// If at a point capacity is exceeded already, everything to the left can be ignored and zeroed
		jump, stop := false, false
		for i:= until-1; i>=call.idx; i-- {
			multiplicator := n.N[i-call.idx]
			if multiplicator == 0 {continue;}

			value += multiplicator * call.k.Items[i].Value
			weight += multiplicator * call.k.Items[i].Weight

			if call.capacity-weight < 0 {
				// If it is the first item we are checking, break the for loop
				if i == until-1 {
					stop = true
					break
				}

				// Otherwise set everything to the left to zero and the item on the right +1
				n.N[i-call.idx+1]++
				for ; i>=call.idx; i-- {
					n.N[i-call.idx] = 0
				}

				// That is fucking expensive... better to cut later than to spend ages in this part
				/*if n.Check() != 0 {
					stop = true
					break
				}*/

				jump = true
				break
			}
		}

		if jump {
			continue
		}
		if stop {
			break
		}

		
		if weight <= call.capacity {
			// If there is still work to be done, invoke the recursive algorithm
			if allDims {
				r = Max(r, value)
			} else {
				//r = Max(r, knapsack_rec(c, FuncCall{until, call.capacity - weight, call.k, value}))
				// Split work
				if workFree <= 0 {
					c.newwork <- 50
					<-c.ack
					workFree = 50
				}

				c.work <- FuncCall{until, call.capacity - weight, call.k, value}
				workFree--
			}
		}

		if n.Add(1) != 0 {
			break
		}
	}

	for workFree > 0 {
		c.res <- 0
		workFree--
	}

	return r
}


func worker(c *ChanCollection) {
	// While work channel is not empty, invoke a knapsack_rec and send the result back
	running := true
	bestRes := 0
	for running {
		// We need a select here as last item might be stolen
		select {
			case inv := <- c.work : {
				start := time.Now()
				if inv.canImprove(bestRes) {
					c.res <- knapsack_rec(c, inv)
				} else {
					c.res <- 0
				}
				dif := time.Since(start)
				if dif > 5*time.Second {
					fmt.Println("Long invocation: %v, total time: %v", inv, dif)
				}
			}
			case res := <- c.workerfinished: {
				// Negative number is termination command
				if res < 0 {
					running = false
				} else {
					bestRes = res
				}
			}
		}
	}
}

func gatherer(c *ChanCollection, finalres chan int, numWorkerThreads int) {
	work := 0
	maxRes := 0
	running := true
	totalWork := 0
	tick := time.Tick(1*time.Second)
	sendPreliminaryResults := time.Tick(100*time.Millisecond)
	for running {
		select {
			case r := <-c.res: {
				// Send first result asap
				if maxRes == 0 && r != 0 {
					for i:=0; i<numWorkerThreads; i++ {
						c.workerfinished <- maxRes
					}				
				}
				maxRes = Max(maxRes, r)
				work++
				if work >= totalWork {
					for i:=0; i<numWorkerThreads; i++ {
						c.workerfinished <- -1
					}
					running = false
				}
			}
			case newWork := <- c.newwork: {
				totalWork += newWork
				c.ack <- true
			}
			case <- tick: {
				if outputInfo {
						fmt.Printf("Gatherer state: %v items of %v (diff: %v), best: %v\n", work, totalWork, totalWork-work, maxRes)
				}
			}

			// The gatherer sends preliminary results from time to time to prune branches
			case <- sendPreliminaryResults: {
				for i:=0; i<numWorkerThreads; i++ {
					c.workerfinished <- maxRes
				}
			}
		}
	}

	finalres <- maxRes
}



func main() {
	var numWorkerThreads int =  1
	flag.IntVar(&numWorkerThreads, "n", 1, "Number of threads");
	var cpuprofile = flag.String("cpuprofile", "", "Enable CPU profiling, write to passed file")


	flag.Parse()

	if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            fmt.Println("Fatal error", err)
            return
        }
        pprof.StartCPUProfile(f)
        fmt.Println("profiling")
        defer pprof.StopCPUProfile()
    }

	// Read the knapsack from command line
	k := ReadKnapsack()

	if outputInfo {
		fmt.Println("Knapsack size:", k.Len(), "capacity:", k.Capacity)
		fmt.Println("Num threads:", numWorkerThreads)
	}

	// Create a fuckload of channels
	// Add some buffer in the channels so we won't end up with send-blocking threads
	c := ChanCollection{
		work: make(chan FuncCall, workCapacity + k.Len()),
		res: make(chan int, numWorkerThreads*2),
		workerfinished: make(chan int, numWorkerThreads+2),
		newwork: make(chan int, numWorkerThreads+2),
		ack: make(chan bool, numWorkerThreads+2)}
	finalres := make(chan int)
	
	// Spawn gatherer first, so he can acknowledge new work already
	go gatherer(&c, finalres, numWorkerThreads)

	for i := 0; i < numWorkerThreads; i++ {
		go worker(&c)
	}

	knapsack_split(&c, FuncCall{0, k.Capacity, k, 0})
	

	// Collect results
	maxRes := <-finalres
	fmt.Println(maxRes)
}