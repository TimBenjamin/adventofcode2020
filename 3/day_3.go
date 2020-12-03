package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func get_treemap() (treemap [][]string) {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		row := strings.Split(line, "")
		if len(row) > 0 {
			treemap = append(treemap, row)
		}
	}
	return
}

func get_num_trees(x, y int) (num_trees int) {
	treemap := get_treemap()
	cur_x := 0       // our position in the current row
	cur_y := 0       // our current row position in treemap
	var row []string // our current row
	// we need to step x across for each y down, looping the row as if x continues infinitely
	for {
		cur_y += y
		if cur_y >= len(treemap) {
			return
		}
		row = treemap[cur_y]
		cur_x += x
		if cur_x >= len(row) {
			cur_x = cur_x - len(row)
		}
		if row[cur_x] == `#` {
			num_trees++
		}
	}
}

func part_1() (num_trees int) {
	return get_num_trees(3, 1)
}

func part_2() (product int) {
	moves := [][]int{
		{1, 1},
		{3, 1},
		{5, 1},
		{7, 1},
		{1, 2},
	}
	product = 1
	for _, xy := range moves {
		product *= get_num_trees(xy[0], xy[1])
	}
	return
}

func main() {
	fmt.Println("PART 1:")
	num_trees := part_1()
	fmt.Printf("I encountered %v trees\n", num_trees)

	fmt.Println("PART 2:")
	product := part_2()
	fmt.Printf("The product is %v\n", product)
}
