package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

const preamble = 25

var numbers []int

func get_numbers() error {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		num, _ := strconv.Atoi(line)
		numbers = append(numbers, num)
	}
	return nil
}

func part_1() (bad int) {
	// I am looking for two numbers in the previous 'preamble'-length slice that add to current number
	// if we don't find one, success! that is the number we want to return
outer:
	for position := preamble; position < len(numbers); position++ {
		preamble_slice := numbers[position-preamble : position]
		current_num := numbers[position]
		found := false
		for i := 0; i < len(preamble_slice); i++ {
			for j := 0; j < len(preamble_slice); j++ {
				// the two numbers must be different
				if j == i {
					continue
				}
				if preamble_slice[i]+preamble_slice[j] == current_num {
					found = true
					continue outer
				}
			}
		}
		if !found {
			bad = current_num
			return
		}
	}
	return
}

// find the set of contiguous numbers that add up to bad_num
// but we have to return the sum of the smallest and largest in this range!
func part_2() int {
	bad_num := part_1()
outer:
	for i := 0; i < len(numbers); i++ {
		running_total := numbers[i]
		// pointless cases...
		if numbers[i] > bad_num || numbers[i] == bad_num {
			continue
		}
		for j := i + 1; j < len(numbers); j++ {
			running_total += numbers[j]
			if running_total == bad_num {
				// success!
				return get_part_2_answer(numbers, i, j)
			}
			if running_total > bad_num {
				// this will happen most of the time
				continue outer
			}
		}
	}
	return 0
}

func get_part_2_answer(numbers []int, i int, j int) (sum int) {
	number_slice := numbers[i : j+1]
	sort.Ints(number_slice)
	return number_slice[0] + number_slice[len(number_slice)-1]
}

// Run the program with the argument "2" to run part 2, or anything else for part 1.
func main() {
	get_numbers()
	if len(os.Args) > 1 && os.Args[1] == "2" {
		fmt.Println("PART 2:")
		answer := part_2()
		fmt.Printf("The answer is %v\n", answer)
	} else {
		fmt.Println("PART 1:")
		answer := part_1()
		fmt.Printf("The bad number is %v\n", answer)
	}
}
