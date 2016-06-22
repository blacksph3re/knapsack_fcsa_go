build:
	go build -o knapsack

run:
	cat input/large.in | ./knapsack -n 64

profile:
	cat input/medium.in | ./knapsack -n 1 -cpuprofile=knapsack.prof
	#go tool pprof knapsack knapsack.prof

scalingtest:
	sh ./scalingtest.sh