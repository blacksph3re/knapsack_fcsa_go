#!/bin/sh

for i in 1 2 4 8 16 32 64 128 256
do
	echo "-------------"
	echo "Running strong scalingtest for N=$i"
	time cat input/huge.in | ./knapsack -n $i
done
