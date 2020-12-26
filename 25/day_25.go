package main

import (
	"fmt"
	"os"
)

var real_card_public_key = 9093927
var real_door_public_key = 11001876
var test_card_public_key = 5764801
var test_door_public_key = 17807724

var card_public_key = real_card_public_key
var door_public_key = real_door_public_key

var card_loop_size int
var door_loop_size int

var divisor = 20201227

/*
The handshake used by the card and the door involves an operation that transforms a subject number.

To transform a subject number, start with the value 1.
Then, a number of times called the loop size, perform the following steps:
- Set the value to itself multiplied by the subject number.
- Set the value to the remainder after dividing the value by 20201227.

*/
func transform(subject_number int, loop_size int) int {
	value := 1
	for i := 0; i < loop_size; i++ {
		value *= subject_number
		value = value % divisor
	}
	return value
}

func get_encryption_key(subject_number int, loop_size int) int {
	return subject_number * loop_size
}

// Run the program with the argument "2" to run part 2, or anything else for part 1.
func main() {
	if len(os.Args) > 1 && os.Args[1] == "2" {
		// Part 2
	} else {

		// Transforming the subject number (the door's public key) with a loop size (the card's loop size)
		// produces the encryption key.
		// Same should happen when transforming the card's public key with the door's loop size.
		// We need this encryption key for part 1.

		// So we need to find out the two loop_size values.

		subject_number := 7
		test_loop_size := 1
		found_card := false
		found_door := false
		value := 1
		for {
			/*if test_loop_size%10000 == 0 {
				fmt.Println("Testing loop size:", test_loop_size)
			}*/
			value *= subject_number
			value = value % divisor
			if value == card_public_key {
				card_loop_size = test_loop_size
				found_card = true
				if found_card && found_door {
					break
				}
			}
			if value == door_public_key {
				door_loop_size = test_loop_size
				found_door = true
				if found_card && found_door {
					break
				}
			}
			test_loop_size++
		}
		fmt.Println("the card loop size is:", card_loop_size)
		fmt.Println("the door loop size is:", door_loop_size)
		// Solution:
		/*
			the card loop size is: 4535884
			the door loop size is: 14984027
		*/
		/*
			The card transforms the subject number of the door's public key according to the card's loop size. The result is the encryption key.
		*/
		enc_1 := transform(door_public_key, card_loop_size)
		fmt.Println("generated enc key from door public key:", enc_1)
		enc_2 := transform(card_public_key, door_loop_size)
		fmt.Println("generated enc key from card public key:", enc_2)
	}
}
