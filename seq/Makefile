CC    = gcc
FLAGS = -O3 -Wall -Wno-unused-result

SRC   = knapsack.c
EXEC  = bin/knapsack

all: make_dirs $(EXEC)

clean:
	rm -f -r bin

make_dirs:
	mkdir -p bin

$(EXEC): $(SRC)
	$(CC) $(FLAGS) $^ -o $@

run:
	./$(EXEC) < ../input/large4.in

