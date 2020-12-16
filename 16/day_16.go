package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Rule struct {
	index   int // we need an order for the rules to solve part 2
	name    string
	min     int // min of the whole range
	max     int // max ditto
	not_min int // exclude these numbers in the middle
	not_max int //
}

// I have just typed this in from my input
// Mistake really as I got one number out by one and that threw my final answer off. D'oh.
var rules = []Rule{
	Rule{0, "departure location", 27, 974, 375, 394},
	Rule{1, "departure station", 40, 953, 288, 294},
	Rule{2, "departure platform", 27, 961, 555, 569},
	Rule{3, "departure track", 40, 958, 605, 617},
	Rule{4, "departure date", 43, 972, 843, 849},
	Rule{5, "departure time", 30, 952, 303, 314},
	Rule{6, "arrival location", 32, 950, 479, 495},
	Rule{7, "arrival station", 48, 969, 734, 754},
	Rule{8, "arrival platform", 37, 954, 261, 275},
	Rule{9, "arrival track", 40, 964, 513, 518},
	Rule{10, "class", 34, 966, 278, 283},
	Rule{11, "duration", 25, 961, 649, 671},
	Rule{12, "price", 28, 956, 685, 704},
	Rule{13, "route", 30, 950, 158, 175},
	Rule{14, "row", 47, 970, 882, 902},
	Rule{15, "seat", 38, 959, 706, 726},
	Rule{16, "train", 40, 961, 196, 216},
	Rule{17, "type", 28, 958, 859, 878},
	Rule{18, "wagon", 31, 967, 544, 553},
	Rule{19, "zone", 49, 953, 791, 815},
}

// test rules
// part 1:
/*
var rules = []Rule{
	Rule{0, "class", 1, 7, 4, 4},
	Rule{1, "row", 6, 44, 12, 32},
	Rule{2, "seat", 13, 50, 41, 44},
}
*/
// part 2:
/*
var rules = []Rule{
	Rule{0, "class", 0, 19, 2, 3},
	Rule{1, "row", 0, 19, 6, 7},
	Rule{2, "seat", 0, 19, 14, 15},
}
*/

//your ticket:
var my_ticket = []int{103, 79, 61, 97, 109, 67, 89, 83, 59, 53, 139, 131, 101, 113, 149, 127, 71, 73, 107, 137}

// I modified the puzzle input such that input.txt contains only "nearby tickets"
var tickets [][]int

func get_tickets() {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var line string
	for scanner.Scan() {
		var set []int
		line = scanner.Text()
		split_line := strings.Split(line, ",")
		for _, s := range split_line {
			n, _ := strconv.Atoi(s)
			set = append(set, n)
		}
		tickets = append(tickets, set)
	}
	return
}

func check(number int) bool {
	// we must find at least one rule that is good for this number
	for _, rule := range rules {
		// assume the rule is good to start with
		rule_is_good := true
		if number < rule.min || number > rule.max { // number is completely outside the range
			rule_is_good = false
		}
		if number >= rule.not_min && number <= rule.not_max { // number is in the forbidden inner range
			rule_is_good = false
		}
		if rule_is_good { // we found one good rule, that's all we need for now
			return true
		}
	}
	return false // did not find any good rule for this number
}

func sanitise_tickets() {
	var good_tickets [][]int
	for _, ticket := range tickets {
		ticket_is_good := true
		for _, number := range ticket {
			if !check(number) {
				ticket_is_good = false
				break
			}
		}
		if ticket_is_good {
			good_tickets = append(good_tickets, ticket)
		}
	}
	tickets = good_tickets
}

func part_1() (sum int) {
	// find all the numbers in the "nearby tickets" data that do not match any rule and add them
	if len(tickets) == 0 {
		get_tickets()
	}
	var bad_numbers []int
	for _, ticket := range tickets {
		for _, number := range ticket {
			if !check(number) {
				bad_numbers = append(bad_numbers, number)
				sum += number
			}
		}
	}
	return
}

func part_2() (answer int) {
	if len(tickets) == 0 {
		get_tickets()
	}
	// first replace tickets with a good set containing no invalid tickets
	sanitise_tickets()

	// I need to find which number relates to which rule.
	// I'll keep a tally for each field, as to which rules they match
	// It's a matrix with fields in the rows, and rules in the columns
	/*
	   			rule
	   field	0	1	2	3
	     0		T	T	T	T
	     1		T	T	T	T
	     2		T	T	T	T
	     3		T	T	T	T
	*/
	// assume to start with that every field matches every rule
	//   field:0 {true, true, true} // if there are just 3 rules
	// when we find a field in a ticket that does not match a rule, make it false
	num_of_fields := len(tickets[0]) // same as the number of rules of course
	field_validation := [][]bool{}
	for i := 0; i < num_of_fields; i++ {
		rule_row := make([]bool, num_of_fields)
		for j := 0; j < num_of_fields; j++ {
			rule_row[j] = true
		}
		field_validation = append(field_validation, rule_row)
	}

	for _, ticket := range tickets {
		for field, number := range ticket {
			for _, rule := range rules {
				if (number < rule.min || number > rule.max) || (number >= rule.not_min && number <= rule.not_max) {
					// this rule, for this field, is bad
					field_validation[field][rule.index] = false
				}
			}
		}
	}

	// each index of this will be one field
	// each row is the list of rules that match that field
	// e.g. Field: 1 matches rules:[5 7 12 13 14 16 19]
	var field_evaluation [][]int

	for _, rules := range field_validation {
		matches := []int{}
		for i := 0; i < len(rules); i++ {
			if rules[i] {
				matches = append(matches, i)
			}
		}
		field_evaluation = append(field_evaluation, matches)
	}

	// final evaluation
	// if a field matches only one rule, it must be made false for all the other fields
	for i := 0; i < len(field_evaluation); i++ {
		if len(field_evaluation[i]) == 1 {
			// all other fields must remove this rule
			fmt.Printf("only field %v can have rule %v\n", i, field_evaluation[i][0])
			rule_index := field_evaluation[i][0]
			field_evaluation = evaluation_remove_rule(field_evaluation, rule_index, i)
			print_state(field_evaluation)
		}
	}

	// we are left with several fields that match 2 rules, so we have to exhaustively pick through them.
	for {
		found_pair := false
		for i := 0; i < len(field_evaluation); i++ {
			if len(field_evaluation[i]) == 0 {
				panic(errors.New("ERROR: zero length row found"))
			}
			if len(field_evaluation[i]) == 2 {
				found_pair = true
				rule_index := field_evaluation[i][0] // this "0" was a piece of luck... just so happens that the first of each pair is the right number. Huh.
				fmt.Printf("=> remove %v from all except field %v\n", rule_index, i)
				field_evaluation = evaluation_remove_rule(field_evaluation, rule_index, i)
				print_state(field_evaluation)
				break
			}
		}
		if !found_pair {
			fmt.Println("No more pairs available!")
			break
		}
	}

	// hopefully now we have a one to one correspondence between fields and rules
	// test it against my ticket
	fmt.Println("Validate solution against my ticket...")
	for field, number := range my_ticket {
		rule_index := field_evaluation[field][0]
		for _, rule := range rules {
			if rule.index == rule_index {
				if (number < rule.min || number > rule.max) || (number >= rule.not_min && number <= rule.not_max) {
					fmt.Printf("FAIL! number %v fails rule %v!\n", number, rule.name)
				} else {
					fmt.Printf("pass - number %v passes rule %v\n", number, rule.name)
				}
				break
			}
		}
	}

	// print out my ticket just to be sure:
	fmt.Println("\nMY TICKET WITH FIELDS:")
	var value int
	for rule_index, rule := range rules {
		for i := 0; i < len(field_evaluation); i++ {
			if field_evaluation[i][0] == rule_index {
				value = my_ticket[i]
				break
			}
		}
		fmt.Printf(" %v : %v\n", rule.name, value)
	}
	fmt.Println("================")

	// "look for the six fields on your ticket that start with the word departure.
	// What do you get if you multiply those six values together?"
	departure_rules := []int{}
	for _, rule := range rules {
		match, _ := regexp.MatchString(`^departure`, rule.name)
		if match {
			departure_rules = append(departure_rules, rule.index)
		}
	}
	fmt.Println("Departure rules are indexes:", departure_rules)
	fields := []int{}
	magic_numbers := []int{} // the numbers to multiply together
	for i := 0; i < len(field_evaluation); i++ {
		for _, rule_index := range departure_rules {
			if field_evaluation[i][0] == rule_index {
				fields = append(fields, i)
				break
			}
		}
	}
	fmt.Println("Got relevant fields:", fields)

	for _, field := range fields {
		magic_numbers = append(magic_numbers, my_ticket[field])
	}
	fmt.Println("Got magic numbers:", magic_numbers)

	answer = 1 // don't want to multiply by 0...
	for _, magic := range magic_numbers {
		answer = answer * magic
	}

	return
}

// remove rule_index from all rows of field_evaluation, except for removing it from row "except"
// in row "except", remove everything else to leave only rule_index
func evaluation_remove_rule(field_evaluation [][]int, rule_index int, except int) (new_evaluation [][]int) {
	for i := 0; i < len(field_evaluation); i++ {
		old_matches := field_evaluation[i]
		new_matches := []int{}
		if i == except {
			new_matches = []int{rule_index}
		} else {
			for m := 0; m < len(old_matches); m++ {
				if old_matches[m] != rule_index {
					new_matches = append(new_matches, old_matches[m])
				}
			}
		}
		new_evaluation = append(new_evaluation, new_matches)
	}
	return
}

func print_state(field_validation [][]int) {
	fmt.Println("State of the validation:")
	for field, matches := range field_validation {
		fmt.Printf("Field: %v matches rules:%v\n", field, matches)
	}
}

// Run the program with the argument "2" to run part 2, or anything else for part 1.
func main() {
	if len(os.Args) > 1 && os.Args[1] == "2" {
		fmt.Println("PART 2:")
		answer := part_2()
		fmt.Printf("The answer is %v\n", answer)
	} else {
		fmt.Println("PART 1:")
		answer := part_1()
		fmt.Printf("The answer is %v\n", answer)
	}
}
