package hangmanweb

import (
	"math/rand"
	"strings"
	"time"
)

func ValidChars(seq string) bool {
	/*
		Checks if a string is made out of
		lower characters

		:param seq (string): the string target
		:return: true if it's valid else false
	*/
	for _, char := range seq {
		if char < 97 || seq[0] > 122 {
			return false
		}
	}
	return true
}

func ValidFileName(seq string) bool {
	/*
		Checks if a string is a correct
		filename

		:param seq (string): the string to test
		:return: true if it's a correct filename
	*/
	if len(seq) <= 1 {
		return false
	}

	if len(strings.Split(seq, "/")) > 1 {
		return false
	}

	if len(seq) >= 3 && seq[len(seq)-3:] == ".go" {
		return false
	}

	return true
}

func Lower(seq string) string {
	/*
		Lower an entire string

		:param seq (string): the string to lower
		:return: the lowered string
	*/
	result := ""
	for _, char := range seq {
		if rune(char) >= 65 && rune(char) <= 90 {
			result += string(rune(char) + 32)
		} else {
			result += string(rune(char))
		}
	}
	return result
}

func RandInt(max int) int {
	/*
		Returns a random int

		:param max (int): the maximum value (excluded)
		:return: the random int
	*/
	return rand.Intn(max)
}

func RandInit() {
	/*
		Init the random seed
	*/
	rand.Seed(time.Now().UnixNano())
}

func IsIn(l string, seq []string) bool {
	/*
		Checks if an element l is in the sequence seq

		:param l (string): the element
		:param seq ([]string): the iterable
		:return: true if l in seq else false
	*/
	for _, char := range seq {
		if char == l {
			return true
		}
	}
	return false
}

func IsInInt(i int, seq []int) bool {
	/*
		Same as IsIn, but with different argument types

		:param l (int): the element
		:param seq ([]int): the iterable
		:return: true if l in seq else false
	*/
	for _, char := range seq {
		if char == i {
			return true
		}
	}
	return false
}

func IsInString(l string, seq string) bool {
	/*
		Same as IsIn, but with different argument types

		:param l (string): the element
		:param seq (string): the iterable
		:return: true if l in seq else false
	*/
	for _, char := range seq {
		if string(char) == l {
			return true
		}
	}
	return false
}

func RemoveAccents(seq string) string {
	/*
		Function that remplace every accent
		in a given string by a "normal" letter

		:param seq (string): The string to format
	*/
	noAccent := ""
	for _, char := range seq {
		isAccent := false
		for k, v := range Accents {
			if IsInString(string(char), v) {
				noAccent += k
				isAccent = true
				break
			}
		}
		if !isAccent {
			noAccent += string(char)
		}
	}
	return noAccent
}
