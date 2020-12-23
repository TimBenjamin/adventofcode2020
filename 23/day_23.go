package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
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

// int version of in_array
func in_array(val int, array []int) (ok bool) {
	for _, i := range array {
		if ok = i == val; ok {
			return
		}
	}
	return
}

func get_circle(input string) (circle Circle) {
	circle.cups = []int{}
	s := strings.Split(input, "")
	for _, n := range s {
		c, _ := strconv.Atoi(n)
		circle.cups = append(circle.cups, c)
	}
	return
}

func get_lowest_cup(circle Circle) (cup int) {
	for _, c := range circle.cups {
		if cup == 0 || c < cup {
			cup = c
		}
	}
	return
}

func get_highest_cup(circle Circle) (cup int) {
	for _, c := range circle.cups {
		if c > cup {
			cup = c
		}
	}
	return
}

func get_cup_position(cup int, circle Circle) int {
	for i, c := range circle.cups {
		if c == cup {
			return i
		}
	}
	panic(errors.New("Did not find cup in circle in get_cup_position()"))
}

func run(circle Circle) (int, Circle) {
	current_move := 0

	// Before the crab starts, it will designate the first cup in your list as the current cup
	current_cup := circle.cups[0]

	for {
		current_move++
		if current_move%100 == 0 {
			fmt.Println("move: ", current_move)
		}

		// formatted output:
		fmt.Printf("-- move %v --\n", current_move)
		fmt.Printf("cups: ")
		for _, c := range circle.cups {
			if c == current_cup {
				fmt.Printf(" (%v)", c)
			} else {
				fmt.Printf(" %v", c)
			}
		}
		fmt.Printf("\n")

		// The crab picks up the three cups that are immediately clockwise of the current cup.
		// They are removed from the circle; cup spacing is adjusted as necessary to maintain the circle.
		picked_up_cups := []int{}
		remainder := []int{}

		pos := get_cup_position(current_cup, circle)
		for i := pos + 1; i < pos+4; i++ {
			picked_up_cups = append(picked_up_cups, circle.getCupAt(i)) // 3 cups after current cup
		}
		for _, c := range circle.cups {
			if in_array(c, picked_up_cups) {
				continue
			}
			remainder = append(remainder, c)
		}

		// formatted output:
		fmt.Printf("pick up: %v\n", picked_up_cups)

		// The crab selects a destination cup: the cup with a label equal to the current cup's label minus one.
		// If this would select one of the cups that was just picked up, the crab will keep subtracting one until
		// it finds a cup that wasn't just picked up.
		// If at any point in this process the value goes below the lowest value on any cup's label,
		// it wraps around to the highest value on any cup's label instead.
		destination_cup := current_cup
		for {
			destination_cup--
			if destination_cup < 1 {
				destination_cup = highest_cup
			}
			if !in_array(destination_cup, picked_up_cups) {
				break
			}
		}

		// formatted output:
		fmt.Printf("destination: %v\n\n", destination_cup)

		// The crab places the cups it just picked up so that they are immediately clockwise of the destination cup.
		// They keep the same order as when they were picked up.
		circle = Circle{}
		for _, c := range remainder {
			circle.cups = append(circle.cups, c)
			if c == destination_cup {
				for _, p := range picked_up_cups {
					circle.cups = append(circle.cups, p)
				}
			}
		}

		// The crab selects a new current cup: the cup which is immediately clockwise of the current cup.
		current_cup = circle.getCupAt(get_cup_position(current_cup, circle) + 1)

		if current_move == move_limit {
			break
		}
	}
	return current_cup, circle
}

var test_input = "389125467"
var real_input = "653427918"
var highest_cup int // 1000000 // 9 in part 1
var move_limit int  // 10000000 in part 2 // 100 in part 1

// Run the program with the argument "2" to run part 2, or anything else for part 1.
func main() {

	var start_circle Circle

	if len(os.Args) > 1 && os.Args[1] == "2" {
		// Part 2!
		highest_cup = 1000000
		move_limit = 10000000
		input := test_input
		// the remaining cups are just numbered in an increasing fashion starting from
		// the number after the highest number in your list and proceeding one by one until one million is reached.
		start_circle = get_circle(input)
		// highest to begin with is 9
		// so add cups to start_circle from 10 to 1000000
		for i := 10; i <= highest_cup; i++ {
			start_circle.cups = append(start_circle.cups, i)
		}
		current, final_circle := run(start_circle)
		fmt.Println("-- final --")
		fmt.Println("cups: ", final_circle.cups)
		fmt.Println("current: ", current)

		// Part 2 answer:
		// Determine which two cups will end up immediately clockwise of cup 1.
		// What do you get if you multiply their labels together?
		pos_1 := final_circle.getCupAt(1)
		answer := final_circle.cups[pos_1+1] * final_circle.cups[pos_1+2]
		fmt.Println("Part 2 answer:", answer)
	} else {
		// Part 1!
		highest_cup = 9
		move_limit = 10
		input := real_input
		start_circle = get_circle(input)
		current, final_circle := run(start_circle)
		fmt.Println("-- final --")
		fmt.Println("cups: ", final_circle.cups)
		fmt.Println("current: ", current)

		// After the crab is done, what order will the cups be in?
		// Starting after the cup labeled 1,
		// collect the other cups' labels clockwise into a single string with no extra characters;
		// each number except 1 should appear exactly once.
		labels := []string{}
		pos := get_cup_position(1, final_circle)
		for i := 1; i < 9; i++ {
			pos++
			labels = append(labels, strconv.Itoa(final_circle.getCupAt(pos)))
		}
		fmt.Println("\nFinal string:", strings.Join(labels, ""))
	}
}
