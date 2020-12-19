package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var real_input_file = "./input.txt"
var test_input_file = "./test_input_2.txt"
var input_file = test_input_file

const PART = 2

// the input parsed into characters with spaces removed
var sums [][]string

func splitter(line string) (split []string) {
	_sum := strings.Split(line, "")
	sum := []string{}
	for _, s := range _sum {
		if s != " " {
			sum = append(sum, s)
		}
	}
	// this is fine except that 2+ -digit numbers are now in consecutive places
	for _, k := range sum {
		if k == "+" || k == "-" || k == "*" || k == "(" || k == ")" {
			split = append(split, k)
		} else {
			if len(split) > 1 {
				matched, _ := regexp.MatchString(`\d+`, split[len(split)-2])
				if matched {
					split[len(split)-2] += k
					continue
				}
			}
			split = append(split, k)
		}
	}
	return split
}

func get_sums() {
	f, err := os.Open(input_file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		fmt.Println("line:", line)
		// NB the input and test input are all single digit numbers, woo
		// The divide operator is also never used! Just +, -, *
		// let's get rid of the bracket terms straight away

		// NB part 2 - seems  that we should add some brackets on added terms first!
		if PART == 2 {
			re := regexp.MustCompile(`[^(]*\d+\s\+\s\d+[^)]*`) // doesn't work
			for {
				loc := re.FindStringSubmatchIndex(line)
				if loc == nil {
					break
				}
				fmt.Println("found a plus term", loc)
				break
			}
		}

		re := regexp.MustCompile(`\([\s\d\+\-\*]+\)`)
		for {
			loc := re.FindStringSubmatchIndex(line)
			if loc == nil {
				break
			}

			eq := line[loc[0]+1 : loc[1]-1]
			// each of these is a mini sum that does not have brackets
			//fmt.Println("eq:", eq)
			sum := strings.Split(eq, " ")
			start, _ := strconv.Atoi(sum[0])
			//fmt.Println("attempt to solve:", sum)
			replacement := solve(sum[1:], start)
			//fmt.Println("replace:", replacement)
			line = line[0:loc[0]] + strconv.Itoa(replacement) + line[loc[1]:]
			//fmt.Println("made:", line)
		}
		_sum := strings.Split(line, " ")
		sum := []string{}
		for _, s := range _sum {
			if s != " " {
				sum = append(sum, s)
			}
		}
		fmt.Println("=====>", sum)
		sums = append(sums, sum)
	}

	return
}

func solve(sum []string, solution int) int {
	index := 0
	fmt.Println("start with:", solution)
	//fmt.Println("sum:", sum)
	for {
		number, _ := strconv.Atoi(sum[index])
		switch sum[index] {
		case "+":
			//fmt.Printf("%v %v\n", sum[index], sum[index+1])
			number, _ = strconv.Atoi(sum[index+1])
			solution += number
		case "-":
			//fmt.Printf("%v %v\n", sum[index], sum[index+1])
			number, _ = strconv.Atoi(sum[index+1])
			solution -= number
		case "*":
			//fmt.Printf("%v %v\n", sum[index], sum[index+1])
			number, _ = strconv.Atoi(sum[index+1])
			solution *= number
		}
		fmt.Println("solution is now:", solution)
		index += 2
		if index >= len(sum) {
			break
		}
	}
	return solution
}

func part_1() (total int) {
	if len(sums) == 0 {
		get_sums()
	}
	for _, sum := range sums {
		solution, _ := strconv.Atoi(sum[0])
		solution = solve(sum[1:], solution)
		total += solution
		fmt.Println("==========================")
	}
	return
}

// Run the program with the argument "2" to run part 2, or anything else for part 1.
func main() {
	/*if len(os.Args) > 1 && os.Args[1] == "2" {
		fmt.Println("PART 2:")
		answer := part_2()
		fmt.Printf("The answer is %v\n", answer)
	} else {
		fmt.Println("PART 1:")
		answer := part_1()
		fmt.Printf("The answer is %v\n", answer)
	}*/
	fmt.Println("PART 1:")
	answer := part_1()
	fmt.Printf("The answer is %v\n", answer)
}
