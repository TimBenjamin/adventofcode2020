package main

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"regexp"
	"strconv"

	"github.com/llgcode/draw2d/draw2dimg"
)

var directions [][]string

func get_directions() {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		r, _ := regexp.Compile(`^(\w)(\d+)?$`)
		matches := r.FindAllStringSubmatch(line, -1)
		direction := []string{matches[0][1], matches[0][2]}
		directions = append(directions, direction)
	}

	return
}

var cur_x int
var cur_y int

func part_1() (sum int) {
	if len(directions) == 0 {
		get_directions()
	}

	facing := "E"
	cur_x = 0
	cur_y = 0

	shift := float64(1500)
	factor := float64(1)

	dest := image.NewRGBA(image.Rect(0, 0, 2000.0, 2000.0))
	gc := draw2dimg.NewGraphicContext(dest)
	gc.Clear()
	gc.SetStrokeColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
	gc.SetLineWidth(1)
	gc.BeginPath() // Initialize a new path

	// Move to a position to start the new path
	x := (factor * float64(cur_x)) + shift
	y := (factor * float64(cur_y)) + shift
	gc.MoveTo(x, y)

	for _, direction := range directions {
		//fmt.Printf("Instruction: %v\n", direction)
		//fmt.Printf("x: %v, y: %v\n", cur_x, cur_y)
		amt, _ := strconv.Atoi(direction[1])
		if direction[0] == "F" {
			move(facing, amt)
		} else if direction[0] == "L" || direction[0] == "R" {
			facing = rotate(facing, direction[0], amt)
			continue // no drawing
		} else {
			move(direction[0], amt)
		}

		// convert x and y to floats, and move them into positive space
		x = (factor * float64(cur_x)) + shift
		y = (factor * float64(cur_y)) + shift
		fmt.Printf("draw x: %v / y: %v\n", x, y)
		if x < 0 || y < 0 {
			return
		}
		gc.LineTo(x, y)
		gc.Stroke()
		gc.MoveTo(x, y)
	}

	// Save to file
	draw2dimg.SaveToPngFile("part_1_plot.png", dest)

	return
}

func abs(v int) (a int) {
	return int(math.Abs(float64(v)))
}

// This is clunky as anything but it works...
func rotate(cur_facing string, direction string, amt int) (new_facing string) {
	cur_deg := 0 // N
	if cur_facing == "E" {
		cur_deg = 90
	} else if cur_facing == "S" {
		cur_deg = 180
	} else if cur_facing == "W" {
		cur_deg = 270
	}

	if direction == "R" {
		cur_deg += amt
	} else {
		cur_deg -= amt
	}
	cur_deg = cur_deg % 360

	switch cur_deg {
	case 0:
		new_facing = "N"
	case -180:
		new_facing = "S"
	case 90:
		new_facing = "E"
	case -270:
		new_facing = "E"
	case 180:
		new_facing = "S"
	case 270:
		new_facing = "W"
	case -90:
		new_facing = "W"
	}

	if len(new_facing) == 0 {
		fmt.Printf("rotate - cur_facing: %v / direction: %v / amt: %v => new_facing: %v\n", cur_facing, direction, amt, new_facing)
		panic(errors.New("Cannot calculate new facing"))
	}

	//fmt.Printf(" => rotate from %v to: %v\n", cur_facing, new_facing)
	return
}

func move(direction string, amt int) {
	switch direction {
	case "N":
		cur_y += amt
	case "S":
		cur_y -= amt
	case "E":
		cur_x += amt
	case "W":
		cur_x -= amt
	}
	//fmt.Printf(" move %v => cur_x: %v / cur_y: %v\n", direction, cur_x, cur_y)
}

var way_y int
var way_x int

func part_2() (sum int) {
	if len(directions) == 0 {
		get_directions()
	}

	// The waypoint starts 10 units east and 1 unit north
	// waypoint position:
	way_y = 1
	way_x = 10
	// ship position:
	cur_y = 0
	cur_x = 0

	shift := float64(1000)
	factor := float64(0.01)
	fmt.Println("draw part 2")

	dest := image.NewRGBA(image.Rect(0, 0, 1500.0, 1500.0))
	gc := draw2dimg.NewGraphicContext(dest)
	gc.Clear()
	gc.SetStrokeColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
	gc.SetLineWidth(1)
	gc.BeginPath() // Initialize a new path

	// Move to a position to start the new path
	x := (factor * float64(cur_x)) + shift
	y := (factor * float64(cur_y)) + shift
	gc.MoveTo(x, y)

	for _, direction := range directions {
		//fmt.Printf("Instruction: %v\n", direction)
		amt, _ := strconv.Atoi(direction[1])
		if direction[0] == "F" {
			// means move the ship in the direction of the waypoint by the number given
			// F10 means move N 10 * way_y and E 10 * way_x
			move("N", amt*way_y)
			move("E", amt*way_x)
		} else if direction[0] == "L" || direction[0] == "R" {
			// means to rotate the waypoint, not the facing, which is now irrelevant
			way_y, way_x = rotate_wp(way_y, way_x, direction[0], amt)
			continue // nothing to draw
		} else {
			// NESW in this part mean to move the waypoint, not the ship
			move_wp(direction[0], amt)
		}

		// convert x and y to floats, and move them into positive space
		x = (factor * float64(cur_x)) + shift
		y = (factor * float64(cur_y)) + shift
		fmt.Printf("draw x: %v / y: %v\n", x, y)
		if x < 0 || y < 0 {
			return
		}
		if x > 1500 || y > 1500 {
			return
		}
		gc.LineTo(x, y)
		gc.Stroke()
		gc.MoveTo(x, y)
	}

	// Save to file
	draw2dimg.SaveToPngFile("part_2_plot.png", dest)

	// finally work out the Manhattan product
	// "between that location and the ship's starting position"
	fmt.Printf("Final x: %v / final y: %v\n", cur_x, cur_y)
	sum = abs(cur_x) + abs(cur_y)
	fmt.Printf(" => sum: %v\n", sum)
	return
}

func rotate_wp(way_y int, way_x int, direction string, amt int) (new_way_y int, new_way_x int) {
	// R90 makes:
	//  way_y => way_x
	//  way_x => -way_y
	// L90 makes:
	//  way_y => -way_x
	//  way_x => +way_y
	// we can split down any other, 180 / 270 into 90 degree turns

	turns := amt / 90
	for i := 0; i < turns; i++ {
		if direction == "R" {
			new_way_y = 0 - way_x
			new_way_x = way_y
		} else if direction == "L" {
			new_way_y = way_x
			new_way_x = 0 - way_y
		}
		way_y = new_way_y
		way_x = new_way_x
	}

	//fmt.Printf(" => rotate WP to: y:%v/x:%v\n", new_way_y, new_way_x)
	return
}

func move_wp(direction string, amt int) {
	switch direction {
	case "N":
		way_y += amt
	case "S":
		way_y -= amt
	case "E":
		way_x += amt
	case "W":
		way_x -= amt
	}
	//fmt.Printf(" move WP %v => way_x: %v / way_y: %v\n", direction, way_x, way_y)
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
