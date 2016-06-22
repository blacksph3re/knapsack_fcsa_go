###
### Example usage:
###   $ python inputer.py -n 100 -m 500
###
### This produces knapsack with capacity 500 and
### 100 items with random value and weight, sending to stdout.
###

from __future__ import print_function
import argparse
import random
import sys

parser = argparse.ArgumentParser(description='Random input generator for knapsack problem')
parser.add_argument('-n', '--numitems',
                    required=True,
                    help='Number of items')
parser.add_argument('-m', '--capacity',
                    required=True,
                    help='Capacity of knapsack')
parser.add_argument('-s', '--seed',
                    default="0",
                    help='Seed for random generator')

def main():
    args = parser.parse_args()

    N = int(args.numitems)
    M = int(args.capacity)
    random.seed(int(args.seed))

    print(N, M)
    for i in range(N):
        value  = random.randint(1, M)
        weight = random.randint(1, M)
        print(value, weight)

if __name__ == "__main__":
    main()
