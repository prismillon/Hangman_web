package hangmanweb

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type WebDisplay struct {
	Total       string
	Current     string
	UsedLetters []string
	Repetition  string
	BadChar     bool
	ErrorLeft   int
	Tries       int
	IsRunning   bool
	MmrDiff     int
}

type UserData struct {
	Game    *GameState
	Display *WebDisplay
	MMR     *int
}

type Players struct {
	Users  map[string]*UserData
	Scores [][]string
}

func DisplayErrorPage(w http.ResponseWriter, data string, templatePages []string) {
	/*
		This helper function allow an easy way to display error pages with custom errors.

		:param w (http.ResponseWritter): the webpage
		:param data (string): the error to display
		:param templatePages ([]string): the html files to render
	*/
	type Data struct {
		Error string
	}
	content := Data{Error: data}
	parsedTemplate, _ := template.ParseFiles(templatePages...)
	parsedTemplate.Execute(w, content)
}

func UpdateDisplay(game *GameState, display *WebDisplay) {
	/*
		Update the display object with the game object
	*/

	display.Total = string(game.Total)
	display.Current = string(game.Current)
	display.ErrorLeft = 10 - game.ErrorCount
	display.UsedLetters = game.UsedLetters
	display.Tries = game.Tries
}

func (players *Players) IndexHandler(w http.ResponseWriter, req *http.Request) {
	/*
		Method of Players
		This function is the handler of the root path of the website.
		It will display different files depending of the context.
		If the user doesnt have any cookie (username) set, it'll display the register page
		If the user have a cookie (username) but Players doesnt contain such username, it'll display the Difficulty choice page
		else, it'll display the current hangman game
	*/
	// Cookie check
	cookie, noCookie := req.Cookie("username")
	if noCookie != nil {
		parsedTemplate, _ := template.ParseFiles(RegisterPage...)
		parsedTemplate.Execute(w, nil)
		return
	}

	// Cookie check
	username := cookie.Value
	if players.Users[username] == nil {
		parsedTemplate, _ := template.ParseFiles(DifficultyPage...)
		parsedTemplate.Execute(w, nil)
		return
	}

	// Game display
	currentPlayer := players.Users[username]
	parsedTemplate, _ := template.ParseFiles(PlayPage...)
	UpdateDisplay(currentPlayer.Game, currentPlayer.Display)
	err := parsedTemplate.Execute(w, currentPlayer.Display)
	if err != nil {
		fmt.Println("[ERROR] - Error while executing template :", err)
		return
	}
}

func (players *Players) HangmanHandler(w http.ResponseWriter, req *http.Request) {
	/*
		Method of Players
		This function is the handler of the "hangman" path of the website.
		We are checking if the method is POST (if not, an error is diplayed) and that
		the player exist (by looking at the cookie value).
		Then, the game logic is executed for the player
	*/
	// Method check
	if req.Method != "POST" {
		errorText := fmt.Sprintf("%v - %v", http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		http.Error(w, errorText, http.StatusMethodNotAllowed)
		return
	}

	// Cookie check
	cookie, noCookie := req.Cookie("username")
	if noCookie != nil || cookie == nil {
		errorText := fmt.Sprintf("%v - %v", http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		http.Error(w, errorText, http.StatusUnauthorized)
		return
	}

	username := cookie.Value
	currentPlayer := players.Users[username]

	if currentPlayer == nil {
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	currentPlayer.Display.Repetition = ""
	currentPlayer.Display.BadChar = false
	currentPlayer.Display.IsRunning = true

	if currentPlayer.Game.IsFinish() {
		currentPlayer.Display.IsRunning = false
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	req.ParseForm()
	guess := req.Form.Get("entry")
	guess = RemoveAccents(guess)
	if !ValidChars(guess) || guess == "" {
		currentPlayer.Display.BadChar = true
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	if IsIn(guess, currentPlayer.Game.UsedLetters) {
		currentPlayer.Display.Repetition = guess
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	currentPlayer.Game.Tries++

	if len(guess) > 1 {
		if guess == RemoveAccents(string(currentPlayer.Game.Total)) {
			currentPlayer.Game.CompleteWord()
		} else {
			currentPlayer.Game.ErrorCount += 2
			currentPlayer.Game.UsedLetters = append(currentPlayer.Game.UsedLetters, guess)
			if currentPlayer.Game.ErrorCount > 10 {
				currentPlayer.Game.ErrorCount = 10
			}
		}
	} else if !currentPlayer.Game.AddLetter(guess) {
		currentPlayer.Game.ErrorCount++
		currentPlayer.Game.UsedLetters = append(currentPlayer.Game.UsedLetters, guess)
	} else {
		currentPlayer.Game.UsedLetters = append(currentPlayer.Game.UsedLetters, guess)
	}

	if currentPlayer.Game.IsFinish() {
		var facteur float64
		currentPlayer.Display.IsRunning = false
		if currentPlayer.Game.ErrorCount < 10 {
			switch currentPlayer.Game.Difficulty {
			case "easy":
				facteur = 1
			case "medium":
				facteur = 2
			case "hard":
				facteur = 3
			}
			triescc := float64(currentPlayer.Game.Tries + currentPlayer.Game.ErrorCount)
			triescc = triescc / (18 + float64(len(currentPlayer.Game.Total))) * float64(len(currentPlayer.Game.Total))
			addpoint := (((float64(len(currentPlayer.Game.Total)) - triescc) / float64(len(currentPlayer.Game.Total))) * 10 * facteur)
			floatconvert := float64(*currentPlayer.MMR) + addpoint
			*currentPlayer.MMR = int(floatconvert)
			currentPlayer.Display.MmrDiff = int(addpoint)
		} else {
			switch currentPlayer.Game.Difficulty {
			case "easy":
				facteur = 3
			case "medium":
				facteur = 2
			case "hard":
				facteur = 1
			}
			var delpoint float64
			for _, char := range currentPlayer.Game.Current {
				if char == '_' {
					delpoint++
				}
			}
			delpoint = (delpoint / float64(len(currentPlayer.Display.Total)-(len(currentPlayer.Game.Total)/2-1))) * 10 * facteur
			floatconvert := float64(*currentPlayer.MMR) - delpoint
			*currentPlayer.MMR = int(floatconvert)
			currentPlayer.Display.MmrDiff = int(delpoint)
		}

		for index, player := range players.Scores {
			if player[0] == username {
				players.Scores = append(players.Scores[:index], players.Scores[index+1:]...)
				break
			}
		}

		for index, value := range players.Scores {
			leaderboardScore, _ := strconv.Atoi(value[1])
			if leaderboardScore < *currentPlayer.MMR {
				players.Scores = append(players.Scores[:index+1], players.Scores[index:]...)
				players.Scores[index] = []string{username, strconv.Itoa(*currentPlayer.MMR)}
				http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
				return
			}
		}
		players.Scores = append(players.Scores, []string{username, strconv.Itoa(*currentPlayer.MMR)})
		Save(players.Scores)
	}

	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
}

func (players *Players) ResetHandler(w http.ResponseWriter, req *http.Request) {
	/*
		Method of Players
		This function is the handler of the "reset" path of the website.
		We are checking if the method is POST (if not, an error is diplayed) and that
		the player exist (by looking at the cookie value).
		This path expects a form containing a "difficulty" variable that should contain the name of a
		.txt file in ../data/
		It will either reset or create a Game object for the current player
	*/
	// Method check
	if req.Method != "POST" {
		errorText := fmt.Sprintf("%v - %v", http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		http.Error(w, errorText, http.StatusMethodNotAllowed)
		return
	}

	// Cookie check
	cookie, noCookie := req.Cookie("username")
	if noCookie != nil || cookie == nil {
		errorText := fmt.Sprintf("%v - %v", http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		http.Error(w, errorText, http.StatusUnauthorized)
		return
	}

	req.ParseForm()
	difficulty := req.Form.Get("difficulty")

	if !IsIn(difficulty, []string{"easy", "medium", "hard"}) {
		parsedTemplate, _ := template.ParseFiles(DifficultyPage...)
		parsedTemplate.Execute(w, nil)
		return
	}

	currentPlayer := players.Users[cookie.Value]
	if currentPlayer == nil {
		game := NewGame(GetRandomWord("./data/" + difficulty + ".txt"))
		display := &WebDisplay{ErrorLeft: 10, Tries: 0, IsRunning: true}
		SavedMMR := 1000
		for _, name := range players.Scores {
			if cookie.Value == name[0] {
				SavedMMR, _ = strconv.Atoi(name[1])
				break
			}
		}
		players.Users[cookie.Value] = &UserData{game, display, &SavedMMR}

	} else {
		*currentPlayer.Game = *NewGame(GetRandomWord("./data/" + difficulty + ".txt"))
		currentPlayer.Display.IsRunning = true
		currentPlayer.Display.Tries = 0
		currentPlayer.Display.ErrorLeft = 10
	}

	players.Users[cookie.Value].Game.Difficulty = difficulty
	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
}

func (players *Players) Register(w http.ResponseWriter, req *http.Request) {
	/*
		Method of Players
		This function is the handler of the "register" path of the website.
		We are checking if the method is POST (if not, an error is diplayed)
		This simply creates a simple register form with error gestion (username already existing...).
	*/
	// Method check
	if req.Method != "POST" {
		errorText := fmt.Sprintf("%v - %v", http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		http.Error(w, errorText, http.StatusMethodNotAllowed)
		return
	}

	req.ParseForm()
	username := req.Form.Get("username")
	if username == "" {
		DisplayErrorPage(w, "Please enter an username.", RegisterPage)
		return
	}

	if len(username) < 3 {
		DisplayErrorPage(w, "Username is too short.", RegisterPage)
		return
	}

	for key := range players.Users {
		if username == key {
			DisplayErrorPage(w, "Username already exists.", RegisterPage)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{Name: "username", Value: username})
	fmt.Println("[INFO] - New player : ", username)
	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
}

func (players *Players) LeaderBoardHandler(w http.ResponseWriter, req *http.Request) {
	/*
		Method of Players
		simply send the leaderboard page.
	*/
	parsedTemplate, _ := template.ParseFiles(LeaderBoardPage...)
	parsedTemplate.Execute(w, players)
}
