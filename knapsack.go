package main

import "fmt"
import "sort"

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
	for idx,_ := range retval.Items {
		var value, weight int
		fmt.Scanln(&value, &weight)
		retval.Items[idx] = NewItem(value, weight, retval.Capacity)
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

func Knapsack_rec(idx, capacity int, k *Knapsack) int {
	if idx >= k.Len() {
		return 0
	}
	var v, w, r int

	for capacity-w >= 0 {
		r = Max(r, v + Knapsack_rec(idx+1, capacity-w, k))
		v += k.Get(idx).Value
		w += k.Get(idx).Weight
	}
	return r
}