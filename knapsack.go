package main

import "fmt"
import "sort"
import "bufio"
import "strconv"
import "os"

type Item struct {
	Value int
	Weight int
	Density float32
	MaxCount int
}

func NewItem(value, weight, capacity int) Item {
	var i Item
	i.Value = value
	i.Weight = weight
	i.Density = float32(value)/float32(weight)
	i.MaxCount = weight / capacity +1
	return i
}

type Knapsack struct {
	Items []Item
	Capacity int
}

// Implement sort interface for knapsack
func (a *Knapsack) Len() int {return len(a.Items)}
func (a *Knapsack) Swap(i, j int) {a.Items[i], a.Items[j] = a.Items[j], a.Items[i]}
func (a *Knapsack) Less(i, j int) bool {return a.Items[i].Density < a.Items[j].Density}
func (a *Knapsack) Get(i int) *Item {return &a.Items[i];}

func ReadKnapsack() *Knapsack {
	retval := new(Knapsack)
	var count int
	fmt.Scanln(&count, &retval.Capacity)
	retval.Items = make([]Item, count)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	for idx,_ := range retval.Items {

		if !scanner.Scan() {
			fmt.Println("Error parsing Stdin")
			break
		}
		value, err1 := strconv.Atoi(scanner.Text())
		if !scanner.Scan() || err1!=nil {
			fmt.Println("Error parsing Stdin")
			break
		}
		weight, err2 := strconv.Atoi(scanner.Text())
		if err2!=nil {
			fmt.Println("Error parsing Stdin")
			break
		}

		retval.Items[idx] = NewItem(int(value), int(weight), retval.Capacity)
	}
	sort.Sort(retval)
	return retval
}

func Max(a, b int) int {
	if a<b {
		return b
	}
	return a
}
