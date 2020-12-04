package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var validating bool // whether we are going to validate or not (i.e. part 2 or not)

func check(entry string) (valid bool) {
	var byr, iyr, eyr, hgt, hcl, ecl, pid bool // we don't care about cid
	records := strings.Split(entry, " ")
	for _, r := range records {
		record := strings.Split(r, ":")
		if len(record) != 2 {
			// badly formated field:value
			return
		}
		switch record[0] {
		case "byr":
			byr = validate("byr", record[1])
		case "iyr":
			iyr = validate("iyr", record[1])
		case "eyr":
			eyr = validate("eyr", record[1])
		case "hgt":
			hgt = validate("hgt", record[1])
		case "hcl":
			hcl = validate("hcl", record[1])
		case "ecl":
			ecl = validate("ecl", record[1])
		case "pid":
			pid = validate("pid", record[1])
		}
	}
	if byr && iyr && eyr && hgt && hcl && ecl && pid {
		valid = true
	} else {
		valid = false
	}
	byr, iyr, eyr, hgt, hcl, ecl, pid = false, false, false, false, false, false, false
	return
}

func validate(field string, value string) (valid bool) {
	if validating {
		/*
			byr (Birth Year) - four digits; at least 1920 and at most 2002.
			iyr (Issue Year) - four digits; at least 2010 and at most 2020.
			eyr (Expiration Year) - four digits; at least 2020 and at most 2030.
			hgt (Height) - a number followed by either cm or in:

				If cm, the number must be at least 150 and at most 193.
				If in, the number must be at least 59 and at most 76.

			hcl (Hair Color) - a # followed by exactly six characters 0-9 or a-f.
			ecl (Eye Color) - exactly one of: amb blu brn gry grn hzl oth.
			pid (Passport ID) - a nine-digit number, including leading zeroes.
		*/
		value = strings.TrimSpace(value)
		switch field {
		case "byr":
			matched, _ := regexp.MatchString(`^\d{4}$`, value)
			if matched {
				date, _ := strconv.Atoi(value)
				if date >= 1920 && date <= 2002 {
					valid = true
				}
			}
		case "iyr":
			matched, _ := regexp.MatchString(`^\d{4}$`, value)
			if matched {
				date, _ := strconv.Atoi(value)
				if date >= 2010 && date <= 2020 {
					valid = true
				}
			}
		case "eyr":
			matched, _ := regexp.MatchString(`^\d{4}$`, value)
			if matched {
				date, _ := strconv.Atoi(value)
				if date >= 2020 && date <= 2030 {
					valid = true
				}
			}
		case "hgt":
			r, _ := regexp.Compile(`^(\d+)(\w+)$`)
			matches := r.FindAllStringSubmatch(value, -1)
			height, _ := strconv.Atoi(matches[0][1])
			if matches[0][2] == "cm" {
				if height >= 150 && height <= 193 {
					valid = true
				}
			} else if matches[0][2] == "in" {
				if height >= 59 && height <= 76 {
					valid = true
				}
			}
		case "hcl":
			matched, _ := regexp.MatchString(`^#[0-9a-f]{6}$`, value)
			if matched {
				valid = true
			}
		case "ecl":
			matched, _ := regexp.MatchString(`^(amb|blu|brn|gry|grn|hzl|oth)$`, value)
			if matched {
				valid = true
			}
		case "pid":
			matched, _ := regexp.MatchString(`^\d{9}$`, value)
			if matched {
				valid = true
			}
		}
	} else {
		if len(value) > 0 {
			valid = true
		}
	}
	return valid
}

func analyse() (num_valid int) {
	f, err := os.Open("./input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var line string
	var entry string
	for scanner.Scan() {
		// records might run over several lines
		line = scanner.Text()
		entry = strings.TrimSpace(entry)
		if len(entry) > 0 && (len(line) == 0 || line == "" || line == "\n") {
			// evaluate and go to next
			if check(entry) {
				num_valid++
			}
			entry = ""
		} else {
			entry += " " + line
		}
	}
	// don't forget the last entry!
	if len(entry) > 0 {
		if check(entry) {
			num_valid++
		}
	}
	return
}

func main() {
	fmt.Println("PART 1:")
	validating = false
	num_valid := analyse()
	fmt.Printf("There are %v valid passports\n", num_valid)

	fmt.Println("PART 2:")
	validating = true
	num_valid = analyse()
	fmt.Printf("There are %v valid passports after validation\n", num_valid)
}
