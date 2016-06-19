build:
	go build -o knapsack

run:
	cat input/large2.in | ./knapsack 2

scalingtest:
	sh ./scalingtest.sh