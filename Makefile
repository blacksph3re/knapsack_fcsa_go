build:
	go build -o knapsack

run:
	cat input/huge.in | ./knapsack -n 128

profile:
	cat input/medhuge.in | ./knapsack -n 1 -cpuprofile=knapsack.prof
	#go tool pprof knapsack knapsack.prof

scalingtest:
	sh ./scalingtest.sh

inputhuge:
	python inputer.py -n 50000000 -m 100000 > input/huge.in

inputmed:
	python inputer.py -n 10000000 -m 100000 > input/medhuge.in

inputcseq:
	python inputer.py -m 300 -m 100 > input/large4.in