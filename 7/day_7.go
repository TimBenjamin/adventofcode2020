package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func get_rules() map[string][]string {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rules := make(map[string][]string)

	scanner := bufio.NewScanner(f)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		// plaid aqua bags contain 5 wavy silver bags, 5 faded silver bags.
		// mirrored crimson bags contain no other bags.
		line = strings.TrimSuffix(line, `.`)
		segments := strings.Split(line, ` bags contain `)
		colour := segments[0]
		rule := []string{}
		if segments[1] != "no other bags" {
			bags := strings.Split(segments[1], ", ")
			for _, v := range bags {
				// 5 wavy silver bags
				r, _ := regexp.Compile(`^([\w\s]+)\sbags?$`)
				matches := r.FindStringSubmatch(v)
				rule = append(rule, matches[0])
			}
		}
		rules[colour] = rule
	}
	return rules
}

func string_in_slice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Returns a slice containing bags that can contain my_bag
func get_can_contain(my_bag string, rules map[string][]string) (can_contain []string) {
	for b, bags := range rules {
		for _, r := range bags {
			if strings.Index(r, my_bag) != -1 {
				// b can contain my bag
				can_contain = append(can_contain, b)
			}
		}
	}
	return
}

// Function to recursively find bags that can contain a bag
func get_cc(cc0 []string, rules map[string][]string) {
	for _, bag := range cc0 {
		if !string_in_slice(bag, master_can_contain) {
			master_can_contain = append(master_can_contain, bag)
			cc1 := get_can_contain(bag, rules)
			if len(cc1) > 0 {
				get_cc(cc1, rules)
			}
		}
	}
}

var master_can_contain []string

func part_1() (count int) {
	rules := get_rules()
	my_bag := `shiny gold`
	count = 0
	cc0 := get_can_contain(my_bag, rules)
	if len(cc0) > 0 {
		get_cc(cc0, rules)
	}
	return len(master_can_contain)
}

func get_bag(my_bag string, quantity int, rules map[string][]string) (total int) {
	//fmt.Println("==== Bag:", my_bag)
	cc0 := rules[my_bag]
	//fmt.Printf(" The bag %v can contain: %v\n", my_bag, cc0)
	if len(cc0) > 0 {
		for _, bag_rule := range cc0 {
			r, _ := regexp.Compile(`^(\d+)\s([\w\s]+)\sbags?$`)
			matches := r.FindAllStringSubmatch(bag_rule, -1)
			num, _ := strconv.Atoi(matches[0][1])
			bag := matches[0][2]
			// we need to add num x (total for bag)
			//fmt.Printf("  Got to look at bag %v...\n", bag)
			num_in_bag := get_bag(bag, num, rules)
			//fmt.Printf("  The total in bag %v was %v\n", bag, num_in_bag)
			if num_in_bag > 0 {
				//fmt.Printf("  Adding %v x %v = %v to the total\n", num_in_bag, num, (num_in_bag * num))
				total += num_in_bag * num
			}
			//fmt.Printf("  Adding %v to the total\n", num)
			total += num
		}
	} else {
		//fmt.Printf("Bag %v contains nothing, adding 0\n", my_bag)
		total = 0
	}
	//fmt.Printf("Returning total %v for bag %v\n", total, my_bag)
	return total
}

var part_2_grand_total int

func part_2() (sum int) {
	rules := get_rules()
	my_bag := `shiny gold`

	sum = get_bag(my_bag, 1, rules)
	part_2_grand_total = 1 * sum

	return
}

func main() {
	fmt.Println("PART 1:")
	answer := part_1()
	fmt.Printf("The answer is %v\n", answer)

	fmt.Println("PART 2:")
	answer = part_2()
	fmt.Printf("The answer is %v\n", answer)

}
