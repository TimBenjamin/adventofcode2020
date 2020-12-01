package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const YEAR = 2020

// Read in the lines of data
// Discard anything that doesn't convert to an int
// Also discard any int values that are outside the bounds we're interested in
func get_data() []int {
	f, err := os.Open("./day_1_data.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var data []int
	var line string
	const UPPER_LIMIT = YEAR
	const LOWER_LIMIT = 1
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line = scanner.Text()
		num, err := strconv.Atoi(line)

		// get rid of non-ints and numbers we don't care about
		if err == nil && num <= UPPER_LIMIT && num >= LOWER_LIMIT {
			data = append(data, num)
		} else {
			fmt.Println("...excluded line: ", line)
		}
	}
	return data
}

// Given an input list of numbers, find the two numbers that sum to YEAR.
// Return the product of the two numbers.
func part_1() (found bool, product int) {
	data := get_data()
	for i, num_1 := range data {
		// there is only one value of num_2 we care about:
		target := YEAR - num_1
		for j, num_2 := range data {
			// we don't want the same number twice:
			if i == j {
				continue
			}
			if num_2 == target {
				fmt.Println("Found two numbers that sum to 2020: ", num_1, num_2)
				product = num_1 * num_2
				found = true
				return
			}
		}
	}
	return
}

// Same as part 1, except we need 3 numbers that sum to YEAR
// We return the product of these.
func part_2() (found bool, product int) {
	data := get_data()
	for i, num_1 := range data {
		for j, num_2 := range data {
			// we don't want the same number twice:
			if i == j {
				continue
			}
			pair_sum := num_1 + num_2
			// no point looking for num_3 if our sum is already too big
			if pair_sum > YEAR {
				continue
			}
			// there is only one value of num_3 we care about...
			target := YEAR - pair_sum
			for k, num_3 := range data {
				// we don't want either number already being used:
				if i == k || j == k {
					continue
				}
				if num_3 == target {
					fmt.Println("Found three numbers that sum to 2020: ", YEAR, num_1, num_2, num_3)
					product = num_1 * num_2 * num_3
					found = true
					return
				}
			}
		}
	}
	return
}

func main() {
	fmt.Println("PART 1:")

	found, product := part_1()
	if found && product != 0 {
		fmt.Println("The solution for part 1 is: ", product)
	} else {
		fmt.Println("Could not find the numbers for part 1!")
	}

	fmt.Println("PART 2:")

	found, product = part_2()
	if found && product != 0 {
		fmt.Println("The solution for part 2 is: ", product)
	} else {
		fmt.Println("Could not find the numbers for part 2!")
	}
}
