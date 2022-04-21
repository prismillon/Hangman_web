package hangmanweb

var Accents = map[string]string{
	"e": "eéèêëÉÈÊË",
	"a": "aàâäÀÁÂÄÆ",
	"u": "uùüûÚÙÛÜ",
	"i": "iïîÍÌÎÏ",
	"o": "oôöœÓÒÔÖØ",
	"c": "cçÇ",
}

type GameState struct {
	Total       []rune
	Current     []rune
	NewLetters  []int
	ErrorCount  int
	UsedLetters []string
	Tries       int
	Difficulty  string
}

func NewGame(seq string) *GameState {
	/*
		A constructor of GameState, based on a word provided

		:param seq (string): the word the game will be based on
		:return: an instance of GameState
	*/
	total := []rune{}
	for _, char := range seq {
		if char == 13 || char == 10 {
			continue
		}
		total = append(total, char)
	}
	res := GameState{Total: []rune(total)}

	// Create res.Current
	for range res.Total {
		res.Current = append(res.Current, '_')
	}
	for j := 0; j < len(res.Total)/2-1; j++ {
		index := RandInt(len(total) - (j + 1))
		for i, char := range res.Total {
			if index == i && i < len(res.Current) {
				if res.Current[i] != '_' {
					index++
				} else {
					res.Current[i] = char
				}
			}
		}
	}
	res.ErrorCount = 0
	res.Tries = 0
	return &res
}

func (g *GameState) IsFinish() bool {
	/*
		Method of GameState
		Check if the game is won by the user by searching
		for '_' chars left in g.Current
	*/
	if g.ErrorCount >= 10 {
		return true
	}
	for _, char := range g.Current {
		if char == '_' {
			return false
		}
	}
	return true
}

func (g *GameState) CompleteWord() {
	/*
		Method of GameState
		Complete g.Current automatically
	*/
	for index, el := range g.Total {
		g.Current[index] = el
	}
}

func (g *GameState) AddLetter(seq string) bool {
	/*
		Method of GameState
		Tries to add a letter provided to g.Current
		it also handles special chars like é, ç, à...

		:param seq (string): the letter to add
		:return: true if it succeed, false if it failed
			(-> letter is not in g.Total)
	*/
	valid := false
	for i, char := range g.Total {
		_, exists := Accents[seq]
		if exists {
			for _, accent := range Accents[seq] {
				if char == accent {
					if g.Current[i] == '_' {
						g.NewLetters = append(g.NewLetters, i)
						g.Current[i] = char
						valid = true
					}
				}
			}
		} else {
			if string(char) == seq && g.Current[i] == '_' {
				g.NewLetters = append(g.NewLetters, i)
				g.Current[i] = rune(seq[0])
				valid = true
			}
		}
	}
	return valid
}
