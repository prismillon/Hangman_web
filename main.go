package main

import (
	"fmt"
	hangmanweb "hangmanweb/controller"
	"net/http"
)

func main() {
	hangmanweb.RandInit()
	// Creating our main object that stores all of our players and all of the scores
	players := hangmanweb.Players{Users: map[string]*hangmanweb.UserData{}, Scores: [][]string{}}
	players.Scores = hangmanweb.Load()

	fsCss := http.FileServer(http.Dir("./view/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fsCss))

	http.HandleFunc("/", players.IndexHandler)
	http.HandleFunc("/register", players.Register)
	http.HandleFunc("/hangman", players.HangmanHandler)
	http.HandleFunc("/reset", players.ResetHandler)
	http.HandleFunc("/leaderboard", players.LeaderBoardHandler)

	fmt.Println("[INFO] - Starting the server...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("[ERROR] - Server could not start properly.\n ", err)
	}
}
