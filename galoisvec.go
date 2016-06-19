package main

type GaloisVec struct {
	Limit []int
	N []int
}

func NewGaloisVec(size int) *GaloisVec {
	newVec := new(GaloisVec)
	newVec.Limit = make([]int, size)
	newVec.N = make([]int, size)
	return newVec
}

// Appends another value to the vector and initializes it with 0
func (x *GaloisVec) Append(limit int) {
	x.Limit = append(x.Limit, limit)
	x.N = append(x.N, 0)
}

// Increments the GaloisVec by step and returns the carry over the last "bit"
// If carry is != 0, there was an overflow
func (x *GaloisVec) Add(step int) int {
	carry := step
	for idx,_ := range x.N {
		if carry <= 0 {break;}

		x.N[idx] += carry
		carry = x.N[idx]/x.Limit[idx]
		x.N[idx] = x.N[idx]%x.Limit[idx]
	}
	return carry
}
