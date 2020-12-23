package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Node struct {
	data int
	next *Node
}

func get_cups(input string) (cups []int) {
	s := strings.Split(input, "")
	for _, n := range s {
		c, _ := strconv.Atoi(n)
		cups = append(cups, c)
	}
	return
}

func print_list(n *Node, num_to_print int, current int) {
	for i := 0; i < num_to_print; i++ {
		if n.data == current {
			fmt.Printf("(%v) ", n.data)
		} else {
			fmt.Printf("%v ", n.data)
		}
		n = n.next
		if n == nil {
			break
		}
	}
	fmt.Println()
}

var test_input = "389125467"
var real_input = "653427918"

func main() {

	var highest_cup int // 1000000 // 9 in part 1
	var move_limit int  // 10000000 in part 2 // 100 in part 1
	var cups []int
	var PART int
	input := real_input

	if len(os.Args) > 1 && os.Args[1] == "2" {
		highest_cup = 1000000
		move_limit = 10000000
		PART = 2
		// highest to begin with is 9
		// so add cups to start_circle from 10 to 1000000
		cups = get_cups(input)
		for i := 10; i <= highest_cup; i++ {
			cups = append(cups, i)
		}
	} else {
		PART = 1
		highest_cup = 9
		move_limit = 100
		cups = get_cups(input)
	}

	head := Node{data: cups[0]}
	n := &head
	for i := 1; i < len(cups); i++ {
		n.next = &Node{data: cups[i]}
		n = n.next
	}
	n.next = &head

	current_cup := &head
	current_move := 1

	//fmt.Printf("-- move %v --\n", current_move)
	//fmt.Printf("cups: ")
	//print_list(current_cup, highest_cup, current_cup.data)

	// Before the crab starts, it will designate the first cup in your list as the current cup
	// current_cup is the head of the linked list...
	start := time.Now()
	for {
		// The crab picks up the three cups that are immediately clockwise of the current cup.
		// They are removed from the circle; cup spacing is adjusted as necessary to maintain the circle.
		picked_up_cups := []int{}
		n := current_cup.next
		for {
			picked_up_cups = append(picked_up_cups, n.data)
			if len(picked_up_cups) == 3 {
				break
			}
			n = n.next
		}

		// The crab selects a destination cup: the cup with a label equal to the current cup's label minus one.
		// If this would select one of the cups that was just picked up, the crab will keep subtracting one until
		// it finds a cup that wasn't just picked up.
		// If at any point in this process the value goes below the lowest value on any cup's label,
		// it wraps around to the highest value on any cup's label instead.
		destination_cup := current_cup.data - 1
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

		//fmt.Println("pick up: ", picked_up_cups)
		//fmt.Println("destination: ", destination_cup)
		//fmt.Println()

		// The crab places the cups it just picked up so that they are immediately clockwise of the destination cup.
		// They keep the same order as when they were picked up.
		cup := current_cup
		var old_third_pointer *Node
		for {
			if cup.data == destination_cup {
				// the existing next cup is where we want the 3rd picked up one to point
				new_next := cup.next
				// and
				// destination cup.next should be set to the first of the picked_up ones
				first_pick := cup.next
				for {
					if first_pick.data == picked_up_cups[0] {
						cup.next = first_pick
						second_pick := first_pick.next
						third_pick := second_pick.next
						old_third_pointer = third_pick.next
						third_pick.next = new_next
						break
					}
					first_pick = first_pick.next
				}
				cup.next = first_pick
				break
			}
			cup = cup.next
		}
		// now the cup before first_pick (i.e. current_cup) needs to point to the cup that third_pick was pointing to
		current_cup.next = old_third_pointer

		// The crab selects a new current cup: the cup which is immediately clockwise of the current cup.
		// NB, the new current cup is going to be the next cup after the pick-up

		// and the current_cup needs resetting to the one next to existing current_cup
		current_cup = current_cup.next

		if current_move == move_limit {
			break
		}
		current_move++

		//fmt.Printf("-- move %v --\n", current_move)
		//fmt.Printf("cups: ")
		//print_list(current_cup, highest_cup, current_cup.data)

	}

	if PART == 1 {
		// After the crab is done, what order will the cups be in?
		// Starting after the cup labeled 1,
		// collect the other cups' labels clockwise into a single string with no extra characters;
		// each number except 1 should appear exactly once.
		fmt.Printf("-- final --\n")
		fmt.Printf("cups: ")
		print_list(current_cup, highest_cup, current_cup.data)
		n := current_cup
		answer := ""
		for {
			if n.data == 1 {
				for i := 0; i < 9; i++ {
					answer += strconv.Itoa(n.next.data)
					n = n.next
				}
				break
			}
			n = n.next
		}
		fmt.Println("The answer was:", answer)
	} else {
		elapsed := time.Since(start)
		fmt.Printf("Took %v\n", elapsed)
		// Determine which two cups will end up immediately clockwise of cup 1.
		// What do you get if you multiply their labels together?
		n := current_cup
		var first int
		var second int
		for {
			if n.data == 1 {
				first = n.next.data
				n = n.next
				second = n.next.data
				break
			}
			n = n.next
		}
		fmt.Println("First after 1 was:", first)
		fmt.Println("Second after 1 was:", second)
		fmt.Println("The answer was:", first*second)
	}

}
