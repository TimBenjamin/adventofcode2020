package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Return the number of valid passwords according to the rules
func part_1() (num_valid int) {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// example: 1-3 a: abcde
	r, _ := regexp.Compile(`^(\d+)\-(\d+)\s([a-zA-Z])\:\s(.+)$`)
	scanner := bufio.NewScanner(f)
	num_valid = 0
	for scanner.Scan() {
		line := scanner.Text()
		matches := r.FindAllStringSubmatch(line, -1)
		min, _ := strconv.Atoi(matches[0][1])
		max, _ := strconv.Atoi(matches[0][2])
		letter := matches[0][3]
		password := matches[0][4]

		m, _ := regexp.Compile(letter)
		count := len(m.FindAllString(password, -1))

		if count >= min && count <= max {
			num_valid = num_valid + 1
		} else {
			//fmt.Printf("Rule %v failed - checking password %v: found %v matches of %v (min: %v / max: %v)\n", line, password, count, letter, min, max)
		}
	}

	return num_valid
}

func part_2() (num_valid int) {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// example: 1-3 a: abcde
	// this means position 1 and position 3
	// ONLY one of these positions must contain the letter "a"
	// the positions start at 1 not 0!
	r, _ := regexp.Compile(`^(\d+)\-(\d+)\s([a-zA-Z])\:\s(.+)$`)
	scanner := bufio.NewScanner(f)
	num_valid = 0
	for scanner.Scan() {
		line := scanner.Text()
		matches := r.FindAllStringSubmatch(line, -1)
		pos_1, _ := strconv.Atoi(matches[0][1])
		pos_2, _ := strconv.Atoi(matches[0][2])
		letter := matches[0][3]
		password := matches[0][4]

		p := strings.Split(password, "")
		match_pos_1 := false
		match_pos_2 := false

		// test position 1:
		if p[pos_1-1] == letter {
			//fmt.Printf("Password %v is OK, contains %v at pos_1 %v\n", password, letter, pos_1)
			match_pos_1 = true
		}

		// test position 2:
		if len(p) >= pos_2 {
			if p[pos_2-1] == letter {
				//fmt.Printf("Password %v is OK, contains %v at pos_2 %v\n", password, letter, pos_1)
				match_pos_2 = true
			}
		} else {
			fmt.Printf("Password %v is too short (%v chars) for the second position %v\n", password, len(p), pos_2)
		}

		if match_pos_1 && match_pos_2 {
			// not allowed
			//fmt.Printf("Password %v is not OK, contains %v at both pos_1 %v and pos_2 %v\n", password, letter, pos_1, pos_2)
		} else if match_pos_1 || match_pos_2 {
			num_valid = num_valid + 1
		}
	}

	return num_valid
}

func main() {
	fmt.Println("PART 1:")
	num_valid := part_1()
	fmt.Printf("I found %v valid passwords\n", num_valid)

	fmt.Println("PART 2:")
	num_valid = part_2()
	fmt.Printf("I found %v valid passwords\n", num_valid)
}
