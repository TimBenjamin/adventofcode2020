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

var rule_a string
var rule_b string

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
		if line_split[1] == "\"a\"" {
			rule_a = line_split[0]
		} else if line_split[1] == "\"b\"" {
			rule_b = line_split[0]
		}
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
	fmt.Println(rule_index)
	//fmt.Printf("handle_rule with index %v [%v]\n", rule_index, rules[rule_index])
	// is it a numeric one or a letter one?
	rule := rules[rule_index]
	matched, _ := regexp.MatchString(`\d`, rule)
	if matched {
		// handle a numeric rule
		//fmt.Println("handle numeric rule:", rule)
		options := strings.Split(rule, "|")
		// there might be only 1 option, or 2
		// save the data in case one option fails and we have to try the other
		data_original := data
	outer:
		for option_index, option := range options {
			option = strings.TrimSpace(option)
			rule_sequence := strings.Split(option, " ")
			fmt.Printf(" try option %v which is %v\n", option_index, rule_sequence)
			for i, r := range rule_sequence {
				//fmt.Println("Apply rule:", r)
				data, result = handle_rule(r, data) // this disappears down a chain of rules and comes back when it hits an a/b rule
				//fmt.Printf("-> result: %v - data: %v\n", result, data)
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
			fmt.Println("   removed a, data is now:", data)
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
			fmt.Println("   removed b, data is now:", data)
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

func chop(chain string, letter byte) (new_chain string, result bool) {
	//fmt.Println("chop:", string(letter))
	if len(chain) == 0 {
		//fmt.Println("The chain is empty!")
		return chain, false
	} else if chain[0] == letter {
		if len(chain) > 1 {
			chain = chain[1:]
		} else {
			chain = ""
		}
		return chain, true
	} else {
		//fmt.Printf("The chain does not match '%v'!\n", string(letter))
		return chain, false
	}
}

func apply(rule string, chain string) (new_chain string, result bool) {
	if len(chain) == 0 {
		fmt.Println("Zero length chain, but rules remain - rule is:", rule)
		return chain, false
	}
	fmt.Println(" options:", rule)
	original_chain := chain
	options := strings.Split(rule, " | ")
	option_good := false
	for _, option := range options {
		option_good = true // assume it worked until it doesn't
		sequence := strings.Split(option, " ")
		fmt.Println("sequence: ", sequence)
	inner:
		for _, s := range sequence {
			fmt.Println(" s: ", s)
			if s == rule_a {
				fmt.Println("a")
				chain, result = chop(chain, 'a')
				if !result {
					fmt.Println(" -> Restore chain to:", original_chain)
					option_good = false
					break inner
				}
				fmt.Println(" chain is now:", chain)
			} else if s == rule_b {
				fmt.Println("b")
				chain, result = chop(chain, 'b')
				if !result {
					fmt.Println(" -> Restore chain to:", original_chain)
					option_good = false
					break inner
				}
				fmt.Println(" chain is now:", chain)
			} else {
				chain, result = apply(rules[s], chain)
				if result {
					// do the next "s" in sequence, as we had a good result
					continue inner
				} else {
					fmt.Printf("Sequence %v has failed\n", sequence)
					option_good = false
					break inner
				}
			}
		}
		fmt.Println("Reached the end of sequence:", sequence)
		if option_good {
			// don't want to do the next option because the first one worked
			return chain, true
		} else {
			chain = original_chain
		}
	}
	if !option_good {
		// no good options found
		fmt.Println("no good options found")
		chain = original_chain
		return chain, false
	}
	// we can get here if there were no option, so option_good never got set to false
	fmt.Println("Reached the end of the rule:", rule)
	return chain, true
}

// yet another version ... this time with a rule like 1 2 | 3 4
// - pass [1 2] back to the function rather than having the inner loop
func do(rule string, chain string) (new_chain string, result bool) {
	fmt.Println("do:", rule)
	if len(chain) == 0 {
		fmt.Printf("chain is zero length, rule failed: %v\n", rule)
		return "", false
	}
	// first, is this a fork?
	options := strings.Split(rule, " | ")
	if len(options) > 1 {
		for _, option := range options {
			new_chain, result = do(option, chain)
			if result {
				chain = new_chain
				fmt.Printf("option %v true\n", option)
				return chain, true // bust out of the options loop as this one worked
			} else {
				fmt.Printf("> option %v failed, try next\n", option)
				fmt.Println("  - length of resulting chain was:", len(new_chain))
			}
		}
		// neither option worked
		fmt.Println("> both options failed")
		return chain, false
	} else {
		// not a fork, handle the sequence
		original_chain := chain
		sequence := strings.Split(rule, " ")
		for _, s := range sequence {
			if s == "11" {
				fmt.Println(">>>>>>>Rule 11, sequence is:", sequence)
			}
			if s == rule_a {
				chain, result = chop(chain, 'a')
				if !result {
					fmt.Printf("> chop a failed, rule %v in sequence %v\n", s, rule)
					return original_chain, false
				}
			} else if s == rule_b {
				chain, result = chop(chain, 'b')
				if !result {
					fmt.Printf("> chop b failed, rule %v in sequence %v\n", s, rule)
					return original_chain, false
				}
			} else {
				chain, result = do(rules[s], chain)
				if result {
					// do the next "s" in sequence, as we had a good result
					continue
				} else {
					fmt.Println("> seq failed on rule", s)
					return original_chain, false
				}
			}
		}
		// we made it through that sequence
		fmt.Printf("seq %v true\n", rule)
		return chain, true
	}
}

// Run the program with the argument "2" to run part 2, or anything else for part 1.
func main() {
	get_rules()
	fmt.Println("rule_a is:", rule_a)
	fmt.Println("rule_b is:", rule_b)
	if len(os.Args) > 1 && os.Args[1] == "2" {
		fmt.Println("part 2")
		// some rules to change for part 2
		rules["8"] = "42 | 42 8"
		rules["11"] = "42 31 | 42 11 31"
	}
	get_input()
	num_valid := 0
	valid := []string{}

	// single case testing:
	input = []string{"aaaaabbaabaaaaababaa"} // invalid in part 1, valid in part 2
	//input = []string{"bbabbbbaabaabba"} // valid in both parts

	for _, chain := range input {
		fmt.Println("Chain:", chain)
		final_chain, result := do(rules["0"], chain)
		if result {
			if len(final_chain) == 0 {
				fmt.Println("valid!")
				num_valid++
				valid = append(valid, chain)
			} else {
				fmt.Println("valid, but data remains, so invalid")
			}
		} else {
			fmt.Println("failed")
		}
	}
	fmt.Println("Number valid:", num_valid)
	fmt.Println(valid)
}
