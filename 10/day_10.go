package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
)

func prepend_int(p int, s []int) (n []int) {
	return append([]int{p}, s...)
}

func reverseAny(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

var adapters []int

func get_adapters() {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	adapters = []int{0}
	scanner := bufio.NewScanner(f)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		a, _ := strconv.Atoi(line)
		adapters = append(adapters, a)
	}
	sort.Ints(adapters)

	// there is an extra final adapter which is +3 to the biggest in "adapters"
	adapters = append(adapters, adapters[len(adapters)-1]+3)
	return
}

func part_1() (product int) {
	if len(adapters) == 0 {
		get_adapters()
	}
	diffs := get_diffs(adapters)
	count_diff_1 := 0
	count_diff_3 := 0
	for _, diff := range diffs {
		if diff == 1 {
			count_diff_1++
		} else if diff == 3 {
			count_diff_3++
		}
	}
	fmt.Printf("There are %v x 1-jolt diffs and %v x 3-volt diffs\n", count_diff_1, count_diff_3)
	product = count_diff_1 * count_diff_3
	return
}

func get_diffs(batch []int) (diffs []int) {
	for i, a := range batch {
		diff := 0
		if i > 0 {
			diff = a - batch[i-1]
			diffs = append(diffs, diff)
		}
	}

	return
}

func part_2() (branches int) {
	if len(adapters) == 0 {
		get_adapters()
	}

	// each time there are two or more diffs of 1, that is a viable "branch"
	// the answer is the number of possible endings once all branches are considered
	// the "magic" numbers 2, 4, 7 I found by manually working out the number of paths through these 1's
	// there are no other paths through a set of diffs like 1, 3, 1
	// only when there are 2+ 1's together
	// e.g. 1, 1, 3 ....
	// the input only contains chains of 2, 3, and 4 ones
	// but the next number for 5 ones would be 11, it's a "tribonacci" sequence
	branches = 0

	diffs := get_diffs(adapters)
	branches = 1
	for i := 0; i < len(diffs); i++ {
		d := diffs[i]
		// how many 1's together in a chain?
		if d == 1 {
			num_ones := 1
			for k := i + 1; k < len(diffs); k++ {
				if diffs[k] == 1 {
					num_ones++
				} else {
					break
				}
			}

			if num_ones > 1 {
				if num_ones == 2 {
					branches = branches * 2
				} else if num_ones == 3 {
					branches = branches * 4
				} else if num_ones == 4 {
					branches = branches * 7
				}
				i += num_ones - 1 // skip past this sequence now
				continue
			}
		}
	}

	return
}

// Run the program with the argument "2" to run part 2, or anything else for part 1.
func main() {
	if len(os.Args) > 1 && os.Args[1] == "2" {
		fmt.Println("PART 2:")
		answer := part_2()
		fmt.Printf("The answer is %v\n", answer)
	} else {
		fmt.Println("PART 1:")
		answer := part_1()
		fmt.Printf("The answer is %v\n", answer)
	}
}
