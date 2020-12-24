package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Circle struct {
	cups []int
}

func (c *Circle) getCupAt(i int) (cup int) {
	// if i is greater than the length of cups
	// start looping from the start again
	n := i % len(c.cups)
	return c.cups[n]
}

func get_cups(input string) (cups []int) {
	s := strings.Split(input, "")
	for _, n := range s {
		c, _ := strconv.Atoi(n)
		cups = append(cups, c)
	}
	return
}

func run(cups []int) []int {
	current_move := 0

	// Before the crab starts, it will designate the first cup in your list as the current cup
	current_cup := cups[0]

	// make these slices at the start then re-use them in the main loop
	remainder := make([]int, len(cups)-3)

	for {
		current_move++

		// The crab picks up the three cups that are immediately clockwise of the current cup.
		// They are removed from the circle; cup spacing is adjusted as necessary to maintain the circle.

		// I want to re-use "cups" at the end so I have to try not to make references to it as I go along

		// always put current_cup at 0 like it is to begin with
		//picked_up_cups := cups[1:4]
		// oddly, it seems faster to make this one fresh every time. huh
		picked_up_cups := []int{}
		for i := 1; i < 4; i++ {
			picked_up_cups = append(picked_up_cups, cups[i])
		}

		// remainder is pos 4 to the end, then the previous current cup ([0])
		//remainder := cups[4:]
		// this is a potential bottleneck
		//remainder := make([]int, len(cups)-3)
		// definitely a lot faster to make this one once, maybe because it's so big
		r := 0
		for i := 4; i < len(cups); i++ {
			remainder[r] = cups[i]
			r++
		}
		remainder[r] = current_cup

		// The crab selects a destination cup: the cup with a label equal to the current cup's label minus one.
		// If this would select one of the cups that was just picked up, the crab will keep subtracting one until
		// it finds a cup that wasn't just picked up.
		// If at any point in this process the value goes below the lowest value on any cup's label,
		// it wraps around to the highest value on any cup's label instead.
		destination_cup := current_cup - 1
		for {
			if destination_cup < 1 {
				destination_cup = highest_cup
			}
			picked_up := false
			for _, p := range picked_up_cups {
				if p == destination_cup {
					picked_up = true
				}
			}
			if !picked_up {
				break
			}
			destination_cup--
		}

		// formatted output:
		/*
			fmt.Printf("-- move %v --\n", current_move)
			fmt.Printf("cups: ")
			for _, c := range cups {
				if c == current_cup {
					fmt.Printf(" (%v)", c)
				} else {
					fmt.Printf(" %v", c)
				}
			}
			fmt.Printf("\n")
		*/
		//fmt.Printf("pick up: %v\n", picked_up_cups)
		//fmt.Printf("destination: %v\n\n", destination_cup)

		// The crab places the cups it just picked up so that they are immediately clockwise of the destination cup.
		// They keep the same order as when they were picked up.
		// new_current_cup, remainder up to destination, destination, pickup, then remainder after destination
		// should be faster re-using the already allocated "cups"
		for i := 0; i < len(remainder); i++ {
			cups[i] = remainder[i]
			//cups = append(cups, remainder[i]) // new current cup and remainder up to destination
			if remainder[i] == destination_cup {
				// destination cup was just added
				c := i + 1
				for _, p := range picked_up_cups {
					cups[c] = p
					c++
				}
				//cups = append(cups, picked_up_cups...)
				for j := i + 1; j < len(remainder); j++ {
					cups[c] = remainder[j]
					c++
				}
				//cups = append(cups, remainder[i+1:]...)
				break
			}
		}

		// The crab selects a new current cup: the cup which is immediately clockwise of the current cup.
		// NB, the new current cup is going to be the next cup after the pick-up
		// i.e. the first of the remainder
		current_cup = remainder[0]

		if current_move == move_limit {
			break
		}
	}
	return cups
}

var test_input = "389125467"
var real_input = "653427918"
var highest_cup int // 1000000 // 9 in part 1
var move_limit int  // 10000000 in part 2 // 100 in part 1

// Run the program with the argument "2" to run part 2, or anything else for part 1.
func main() {

	if len(os.Args) > 1 && os.Args[1] == "2" {
		// Part 2!
		highest_cup = 1000000
		move_limit = 10000000
		input := test_input
		// the remaining cups are just numbered in an increasing fashion starting from
		// the number after the highest number in your list and proceeding one by one until one million is reached.
		cups := get_cups(input)
		// highest to begin with is 9
		// so add cups to start_circle from 10 to 1000000
		for i := 10; i <= highest_cup; i++ {
			cups = append(cups, i)
		}
		start := time.Now()
		final_cups := run(cups)
		elapsed := time.Since(start)
		fmt.Printf("Took %v seconds\n", elapsed)

		// Part 2 answer:
		// Determine which two cups will end up immediately clockwise of cup 1.
		// What do you get if you multiply their labels together?
		final_circle := Circle{final_cups}
		pos_1 := final_circle.getCupAt(1)
		answer := final_circle.cups[pos_1+1] * final_circle.cups[pos_1+2]
		fmt.Println("Part 2 answer:", answer)
	} else {
		// Part 1!
		highest_cup = 9
		move_limit = 100
		input := real_input
		cups := get_cups(input)
		final_cups := run(cups)
		fmt.Println("-- final --")
		fmt.Println("cups: ", final_cups)

		// After the crab is done, what order will the cups be in?
		// Starting after the cup labeled 1,
		// collect the other cups' labels clockwise into a single string with no extra characters;
		// each number except 1 should appear exactly once.
		final_circle := Circle{final_cups}
		labels := []string{}
		var pos int
		for i, c := range final_circle.cups {
			if c == 1 {
				pos = i
			}
		}
		for i := 1; i < 9; i++ {
			pos++
			labels = append(labels, strconv.Itoa(final_circle.getCupAt(pos)))
		}
		fmt.Println("\nFinal string:", strings.Join(labels, ""))
		// answer should be 76952348 for my real input
	}
}
