package hangmanweb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func Save(table [][]string) {
	/*
		Save the score into the database.
	*/
	encoded, err := json.Marshal(table)
	if err == nil {
		ioutil.WriteFile("Basededonnée.json", encoded, 0777)
	} else {
		fmt.Println(err)
	}
}

func Load() [][]string {
	/*
		Load the scores from the database.
	*/
	var Table [][]string
	data, _ := ioutil.ReadFile("Basededonnée.json")
	dele := json.Unmarshal(data, &Table)
	if dele != nil {
		fmt.Println(dele)
	}
	return Table
}
