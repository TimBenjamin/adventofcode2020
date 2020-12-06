package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func part_1() (sum int) {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	alpha := `abcdefghijklmnopqrstuvwxyz`
	sum = 0
	cur_sum := 0
	log := make([]int, 26)

	scanner := bufio.NewScanner(f)
	var line string
	for scanner.Scan() {
		line = scanner.Text()

		if len(line) == 0 || line == "" || line == "\n" {
			// sum and add
			for _, val := range log {
				cur_sum += val
			}
			sum += cur_sum
			log = make([]int, 26)
			cur_sum = 0

		} else {
			// add to current log
			responses := strings.Split(line, "")
			for _, r := range responses {
				idx := strings.Index(alpha, r)
				log[idx] = 1
			}
		}
	}
	// don't forget to sum the last set of responses!
	for _, val := range log {
		cur_sum += val
	}
	sum += cur_sum
	log = make([]int, 26)

	return
}

func part_2() (sum int) {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	alpha := `abcdefghijklmnopqrstuvwxyz`
	sum = 0
	cur_sum := 0
	cur_response_count := 0
	log := make([]int, 26)

	scanner := bufio.NewScanner(f)
	var line string
	for scanner.Scan() {
		line = scanner.Text()

		if len(line) == 0 || line == "" || line == "\n" {
			// sum and add
			for _, val := range log {
				if val == cur_response_count {
					cur_sum += 1
				}
			}
			fmt.Printf("Completed a record %v - got sum: %v\n", log, cur_sum)
			sum += cur_sum
			log = make([]int, 26)
			cur_sum = 0
			cur_response_count = 0

		} else {
			// add to current log
			responses := strings.Split(line, "")
			fmt.Println("Analysing response:", responses)
			for _, r := range responses {
				idx := strings.Index(alpha, r)
				fmt.Printf("Found index %v for letter %v\n", idx, r)
				log[idx] += 1
			}
			cur_response_count++
		}
	}
	// don't forget to sum the last set of responses!
	for _, val := range log {
		if val == cur_response_count {
			cur_sum += 1
		}
	}
	fmt.Printf("Completed final record %v - got sum: %v\n", log, cur_sum)
	sum += cur_sum
	log = make([]int, 26)

	return
}

func main() {
	fmt.Println("PART 1:")
	sum := part_1()
	fmt.Printf("The sum of the responses is %v\n", sum)

	fmt.Println("PART 2:")
	sum = part_2()
	fmt.Printf("The sum of the responses is now %v\n", sum)
}
