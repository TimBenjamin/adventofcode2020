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
var test_input_file = "./test_input.txt"
var input_file = real_input_file

var PART = 1

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

func solve_multiplied_terms(line string) string {
	re := regexp.MustCompile(`\([\s\d\*]+\)`)
	for {
		loc := re.FindStringSubmatchIndex(line)
		if loc == nil {
			break
		}

		eq := line[loc[0]+1 : loc[1]-1]
		// each of these is a mini sum that does not have brackets
		sum := strings.Split(eq, " ")
		start, _ := strconv.Atoi(sum[0])
		replacement := solve(sum[1:], start)
		line = line[0:loc[0]] + strconv.Itoa(replacement) + line[loc[1]:]
	}
	//fmt.Println("line after solving multiply clauses:", line)
	return line
}

func solve_added_terms(line string) string {
	re := regexp.MustCompile(`\([\s\d\+]+\)`)
	for {
		loc := re.FindStringSubmatchIndex(line)
		if loc == nil {
			break
		}

		eq := line[loc[0]+1 : loc[1]-1]
		// each of these is a mini sum that does not have brackets
		sum := strings.Split(eq, " ")
		start, _ := strconv.Atoi(sum[0])
		replacement := solve(sum[1:], start)
		line = line[0:loc[0]] + strconv.Itoa(replacement) + line[loc[1]:]
	}
	//fmt.Println("line after solving add clauses:", line)
	return line
}

func bracket_add_terms(line string) string {
	re := regexp.MustCompile(`\b[\d\s\+]+\d`)
	locations := re.FindAllStringSubmatchIndex(line, -1)
	var adds int
	for _, loc := range locations {
		loc[0] += adds
		loc[1] += adds
		if loc[0] > 0 && line[loc[0]-1] == '(' && line[loc[1]] == ')' {
			// the term is already in brackets
			continue
		}
		clause := "(" + line[loc[0]:loc[1]] + ")"
		// the clause must actually contain a +!
		// could just be >1 digits, my regex isn't good
		plus_matched, _ := regexp.MatchString(`\+`, clause)
		if !plus_matched {
			continue
		}
		new_line := ""
		if loc[0] == 0 {
			new_line = clause + line[loc[1]:]
		} else {
			new_line = line[0:loc[0]] + clause + line[loc[1]:]
		}
		adds += len(new_line) - len(line) // offset subsequent locs thanks to the two new brackets
		line = new_line
	}
	//fmt.Println("line after bracketing add terms", line)
	return line
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
		fmt.Println("Input line:", line)
		// NB the input and test input are all single digit numbers, woo
		// The divide and subtract operators are also never used! Just +, *
		// So let's get evaluate the bracket terms straight away

		// NB part 2! addition clauses go first
		// so let's put them in brackets and let the part 1 stuff do its work
		// but there are cases like
		// 2 * 3 + (4 * 5)
		// so let's deal with multiplied terms in brackets first
		if PART == 2 {
			// do these 3 steps to simplify the sum until it doesn't change
			for {
				new_line := line
				new_line = solve_multiplied_terms(new_line)
				new_line = bracket_add_terms(new_line)
				new_line = solve_added_terms(new_line)
				new_line = bracket_add_terms(new_line)
				if new_line == line {
					break
				}
				line = new_line
			}
		}

		// full solve
		// NB! This should only be done if not already reduced to a simple int
		space_matched, _ := regexp.MatchString(`\s`, line)
		if space_matched {
			re := regexp.MustCompile(`\([\s\d\+\*]+\)`)
			for {
				loc := re.FindStringSubmatchIndex(line)
				if loc == nil {
					break
				}

				eq := line[loc[0]+1 : loc[1]-1]
				// each of these is a mini sum that does not have brackets
				sum := strings.Split(eq, " ")
				start, _ := strconv.Atoi(sum[0])
				replacement := solve(sum[1:], start)
				line = line[0:loc[0]] + strconv.Itoa(replacement) + line[loc[1]:]
			}
		}
		_sum := strings.Split(line, " ")
		sum := []string{}
		for _, s := range _sum {
			if s != " " {
				sum = append(sum, s)
			}
		}
		fmt.Println(" After simplification =====>", sum)
		sums = append(sums, sum)
	}

	return
}

func solve(sum []string, solution int) int {
	index := 0
	for {
		number, _ := strconv.Atoi(sum[index])
		switch sum[index] {
		case "+":
			number, _ = strconv.Atoi(sum[index+1])
			solution += number
		case "*":
			number, _ = strconv.Atoi(sum[index+1])
			solution *= number
		}
		index += 2
		if index >= len(sum) {
			break
		}
	}
	return solution
}

func do_it() (total int) {
	if len(sums) == 0 {
		get_sums()
	}
	for _, sum := range sums {
		var solution int
		if len(sum) == 1 {
			// it's already been simplified to the solution
			solution, _ = strconv.Atoi(sum[0])
		} else {
			solution, _ = strconv.Atoi(sum[0])
			solution = solve(sum[1:], solution)
		}
		total += solution
		fmt.Println("--------------------------")
		fmt.Println(sum)
		fmt.Println("solution:", solution)
		fmt.Println("==========================")
	}
	return
}

// Run the program with the argument "2" to run part 2, or anything else for part 1.
func main() {
	if len(os.Args) > 1 && os.Args[1] == "2" {
		fmt.Println("PART 2:")
		PART = 2
		answer := do_it()
		fmt.Printf("The grand total is %v\n", answer)
	} else {
		fmt.Println("PART 1:")
		answer := do_it()
		fmt.Printf("The grand total is %v\n", answer)
	}
}
