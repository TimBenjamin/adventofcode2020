package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

var real_input_file = "./input.txt"
var test_input_file = "./test_input.txt"
var input_file = real_input_file

var foods = []Food{}

type Food struct {
	id          int
	ingredients []string
	allergens   []string
}

func in_array(val string, array []string) (ok bool) {
	for _, i := range array {
		if ok = i == val; ok {
			return
		}
	}
	return
}

func remove(val string, array []string) []string {
	result := []string{}
	for _, i := range array {
		if i == val {
			continue
		}
		result = append(result, i)
	}
	return result
}

func get_input() {
	f, err := os.Open(input_file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var line string
	id := 0
	for scanner.Scan() {
		line = scanner.Text()
		line_split := strings.Split(line, " (")
		ingredients := strings.Split(line_split[0], " ")
		// not every food might have allergens
		allergens := []string{}
		if len(line_split) > 1 {
			// "contains dairy, fish)"
			a := line_split[1]
			b := a[9 : len(a)-1]
			allergens = strings.Split(b, ", ")
		}
		food := Food{id, ingredients, allergens}
		foods = append(foods, food)
		id++
	}
	return
}

func run() {
	// we are trying to find clean ingredients that can't contain any allergen
	// start by assuming that all ingredients contain all allergens:
	master_allergens := []string{}
	master_ingredients := map[string][]string{}
	for _, food := range foods {
		for _, a := range food.allergens {
			if !in_array(a, master_allergens) {
				master_allergens = append(master_allergens, a)
			}
		}
		for _, i := range food.ingredients {
			master_ingredients[i] = []string{}
		}
	}
	for i, _ := range master_ingredients {
		list := []string{}
		for _, a := range master_allergens {
			list = append(list, a)
		}
		master_ingredients[i] = list
	}

	// now go through the foods trying to find logical fallacies
	for _, food := range foods {
		for _, i := range food.ingredients {
			if len(food.allergens) == 0 {
				// leave aside for now
				continue
			}
			for _, a := range master_allergens {
				// suppose i contains a
				// then every time a is mentioned, i must be in the ingredients
				// so test that assumption and find when it is false
				for _, f := range foods {
					if food.id == f.id {
						continue
					}
					if in_array(a, f.allergens) {
						if !in_array(i, f.ingredients) {
							// therefore i does not contain a
							list := master_ingredients[i]
							list = remove(a, list)
							master_ingredients[i] = list
						}
					}
				}
			}
		}
	}

	// make a list of the "clean" ingredients that therefore contain no allergens
	clean_ingredients := []string{}
	for i, a := range master_ingredients {
		if len(a) == 0 {
			clean_ingredients = append(clean_ingredients, i)
		}
	}
	fmt.Println("The clean ingredients are:", clean_ingredients)

	// Part 1:
	// How many times do any of those [clean] ingredients appear? [in all the foods' ingredients list]
	count := 0
	for _, f := range foods {
		for _, i := range clean_ingredients {
			if in_array(i, f.ingredients) {
				count++
			}
		}
	}
	fmt.Println("Part 1 answer:", count)

	// Part 2:
	// find which "dirty" ingredient contains which allergen
	// ignore the clean ingredients now
	// hopefully there is one ingredient that contains only one allergen
	// then we can strike that allergen from the other ingredients
	// and so on
	// until there is a 1:1 correspondence!
outer:
	for {
		for i, a := range master_ingredients {
			if len(a) == 1 {
				allergen := a[0]
				// remove a from all master_ingredients, except for i
				for mi, ma := range master_ingredients {
					if mi == i {
						continue
					}
					ma = remove(allergen, ma)
					master_ingredients[mi] = ma
				}
			}
		}
		// still got pairs? if not, break
		for _, a := range master_ingredients {
			if len(a) > 1 {
				continue outer
			}
		}
		break
	}

	// what have we got then?
	fmt.Println("\nPart 2 summary:\n------------------")
	for i, a := range master_ingredients {
		if len(a) > 0 {
			fmt.Printf("%v contains %v\n", i, a[0])
		}
	}
	// to get the solution:
	/*
		Arrange the ingredients alphabetically by their allergen and separate them by commas
		to produce your canonical dangerous ingredient list.
		(There should not be any spaces in your canonical dangerous ingredient list.)
	*/
	dirty_ingredients := []string{}
	sort.Strings(master_allergens)
	for _, allergen := range master_allergens {
		for i, a := range master_ingredients {
			if len(a) > 0 && a[0] == allergen {
				dirty_ingredients = append(dirty_ingredients, i)
			}
		}
	}
	answer := strings.Join(dirty_ingredients, ",")
	fmt.Println("\nSorted dirty ingredients / part 2 answer:\n", answer)
}

func main() {
	get_input()
	run()
}
