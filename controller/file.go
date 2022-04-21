package hangmanweb

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func GetRandomWord(filename string) string {
	/*
		Get a random word (line) from a file
		word is tested to check if it's a valid word (has to be
			ascii chars only)
		If 5 words are wrong in a row, an error message is displayed
		and the program will exit


		:param filename (string): the name of the targeted file
		:return (string): a line from the file
	*/
	lines := GetFile(filename, "\n")
	cc := 0
	test := false
	for {
		pick := lines[RandInt(len(lines))]
		if pick == "" {
			cc++
			continue
		}
		for _, char := range pick {
			if char < 65 || char > 90 && char < 97 || char > 122 && char < 224 {
				test = true
				cc++
				if cc > 4 {
					fmt.Println("La liste contient trop d'invalide selection du mot impossible ^^'")
					os.Exit(0)
				}
				fmt.Print("Le mot choisit est invalide selection d'un nouveau mot - Tentative numéro " + string(rune(cc)+48))
			}
		}
		if !test {
			return pick
		}
	}
}

func GetFile(filename, sep string) []string {
	/*
		Get the content of the file in a list format
		If the file provided doesn't exist, it displays an
		error message then exit the program, same if the file
		is empty

		:param filename (string): the name of the targeted file
		:param sep (string): the separator wanted
		:return: the entire file splitten by the separator
	*/
	var testfile []string
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	if string(content) == "" {
		fmt.Printf("Le fichier fournit %v est vide, action impossible\n", filename)
		os.Exit(0)
	} else if sep == "\n" {
		testfile = strings.Split(string(content), "\r\n")
		if len(testfile) <= 1 {
			testfile = strings.Split(string(content), "\n")
		}
	} else {
		testfile = strings.Split(string(content), "\r\n\r\n")
		if len(testfile) <= 1 {
			testfile = strings.Split(string(content), "\n\n")
		}
	}
	return testfile
}
