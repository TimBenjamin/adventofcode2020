package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var real_input_file = "./input.txt"
var real_rules_file = "./rules.txt"
var test_input_file = "./test_input_2.txt"
var test_rules_file = "./test_rules_2.txt"

var rules_file = test_rules_file
var input_file = test_input_file

var rules = map[string]string{}
var input []string

func get_rules() {
	f, err := os.Open(rules_file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		line_split := strings.Split(line, ": ")
		rules[line_split[0]] = line_split[1]
	}
}

func get_input() {
	f, err := os.Open(input_file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		input = append(input, line)
	}
}

var path = []string{}

// recursively handle rules with data, ultimately returning true/false if data is length 0?
func handle_rule(rule_index string, data string) (new_data string, result bool) {
	path = append(path, rule_index)
	fmt.Printf("==> handle_rule with index %v [%v] and data %v\n", rule_index, rules[rule_index], data)
	// is it a numeric one or a letter one?
	rule := rules[rule_index]
	matched, _ := regexp.MatchString(`\d`, rule)
	if matched {
		// handle a numeric rule
		//fmt.Println("handle numeric rule:", rule)
		options := strings.Split(rule, "|")
		// there might be only 1 option, or 2
		data_original := data
	outer:
		for option_index, option := range options {
			option = strings.TrimSpace(option)
			fmt.Printf("try option %v which is >>%v<<\n", option_index, option)
			rule_sequence := strings.Split(option, " ")
			for i, r := range rule_sequence {
				fmt.Println("We have to apply rule:", r)
				data, result = handle_rule(r, data)
				if result {
					// if there are more rules to come, and we have no data left, this should be a dead-end path.
					if i < len(rule_sequence)-1 && len(data) == 0 {
						fmt.Println("No more data, but there are more rules! try another path")
						return data, false
					}
					continue // apply the next r
				} else {
					if option_index == 0 {
						// the option failed, try the other option with the original data
						fmt.Printf("  - this option (%v) failed, try the other...\n", option_index)
						fmt.Println("  - reset data to:", data_original)
						data = data_original
						continue outer
					} else {
						// no other options
						return data, false
					}
				}
			}
			// reached the end of this chain of rules
			return data, result
		}
		// if we are here, no options worked.
		fmt.Println("All options failed!!")
		return data, false
	} else if rule == "\"a\"" {
		if len(data) == 0 {
			fmt.Printf("Rule %v addresses non-existent data %v\n", rule, data)
			return data, false
		}
		if data[0] == 'a' {
			if len(data) > 1 {
				data = data[1:]
			} else {
				data = ""
			}
			fmt.Println("Removed a, data is now:", data)
			return data, true
		} else {
			fmt.Printf("Rule %v fails to match 'a' with data %v\n", rule, data)
			return data, false
		}
	} else if rule == "\"b\"" {
		if len(data) == 0 {
			fmt.Printf("Rule %v addresses non-existent data %v\n", rule, data)
			return data, false
		}
		if data[0] == 'b' {
			if len(data) > 1 {
				data = data[1:]
			} else {
				data = ""
			}
			fmt.Println("Removed b, data is now:", data)
			return data, true
		} else {
			fmt.Printf("Rule %v fails to match 'b' with data %v\n", rule, data)
			return data, false
		}
	} else {
		fmt.Println("WARNING, bad rule:", rule)
		panic(errors.New("Encountered an unhandled type of rule!"))
	}
}

func run() int {
	matching_inputs := []string{}
	for _, data := range input {
		path = []string{}
		end_data, result := handle_rule("0", data)
		fmt.Printf("Result from handling rule 0 for data %v is %v with end data >>%v<<\n", data, result, end_data)
		// if it completely matched, i.e. there is nothing left of the input data, add to total
		if result && len(end_data) == 0 {
			fmt.Printf("Data %v fully matched rule 0!\n", data)
			matching_inputs = append(matching_inputs, data)
			fmt.Println(path)
			fmt.Println("=========================")
		} else {
			fmt.Printf("Data %v DID NOT fully match rule 0\n", data)
			fmt.Println(path)
			fmt.Println("=========================")
		}
	}
	fmt.Println("The following inputs matched:")
	for _, m := range matching_inputs {
		fmt.Println(m)
	}
	return len(matching_inputs)
}

// Run the program with the argument "2" to run part 2, or anything else for part 1.
func main() {
	get_rules()
	if len(os.Args) > 1 && os.Args[1] == "2" {
		fmt.Println("part 2")
		// some rules to change for part 2
		rules["8"] = "42 | 42 8"
		rules["11"] = "42 31 | 42 11 31"
	}
	get_input()
	answer := run()
	fmt.Printf("The answer is %v\n", answer)
}
