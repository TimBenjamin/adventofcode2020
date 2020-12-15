package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

//const input = "0,3,6" // solution @ 2020: 436 and 175594 @ 30000000
//const input = "1,3,2" // solution @ 2020: 1 and 2578 @ 30000000
//const input = "2,1,3" // solution @ 2020: 10 and 3544142 @ 30000000
//const input = "1,2,3" // solution @ 2020: 27 and 261214 @ 30000000

// the actual input:
const input = "13,16,0,12,15,1" // solution @ 2020: 319 and 2424 @ 30000000

// game container:
var register []int // number n was last spoken on turn register[n]

// the last number spoken:
var last_number []int

// number of turns taken:
var turns int

func get_start(num_turns int) (next_number int) {
	register = make([]int, num_turns) // the biggest possible number is the max number of turns
	input_split := strings.Split(input, ",")
	for _, i := range input_split {
		next_number, _ := strconv.Atoi(i)
		turns++
		add_to_register(next_number)
	}
	return
}

func add_to_register(number int) (age int) {
	if register[number] > 0 { // i.e. it has been spoken before
		age = turns - register[number]
		register[number] = turns // = last seen on turn "turns"
		return
	}
	// it's never been spoken before
	register[number] = turns
	return 0
}

func run_game(num_turns int) int {
	next_number := get_start(num_turns)
	for {
		turns++
		// we are after this nth number, so:
		if turns == num_turns {
			return next_number
		}
		next_number = add_to_register(next_number)
	}
}

// Run the program with the argument "2" to run part 2, or anything else for part 1.
func main() {
	if len(os.Args) > 1 && os.Args[1] == "2" {
		fmt.Println("PART 2:")
		answer := run_game(30000000)
		fmt.Printf("The answer is %v\n", answer)
	} else {
		fmt.Println("PART 1:")
		answer := run_game(2020)
		fmt.Printf("The answer is %v\n", answer)
	}
}
