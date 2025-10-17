package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

var columns = map[int][]int{
	0: {}, 1: {}, 2: {}, 3: {}, 4: {}, 5: {}, 6: {},
}
var currentPlayer = 1

type GameData struct {
	Grid          [6][7]int
	Cols          []int
	J1            string
	J2            string
	CurrentPlayer int
}

func resetGame() {
	columns = map[int][]int{
		0: {}, 1: {}, 2: {}, 3: {}, 4: {}, 5: {}, 6: {},
	}
	currentPlayer = 1
}

func detectWinner(grid [6][7]int) int {
	for row := 0; row < 6; row++ {
		for col := 0; col < 7; col++ {
			token := grid[row][col]
			if token == 0 {
				continue
			}
			if col+3 < 7 &&
				grid[row][col+1] == token &&
				grid[row][col+2] == token &&
				grid[row][col+3] == token {
				return token
			}
			if row+3 < 6 &&
				grid[row+1][col] == token &&
				grid[row+2][col] == token &&
				grid[row+3][col] == token {
				return token
			}
			if row+3 < 6 && col+3 < 7 &&
				grid[row+1][col+1] == token &&
				grid[row+2][col+2] == token &&
				grid[row+3][col+3] == token {
				return token
			}
			if row-3 >= 0 && col+3 < 7 &&
				grid[row-1][col+1] == token &&
				grid[row-2][col+2] == token &&
				grid[row-3][col+3] == token {
				return token
			}
		}
	}
	return 0
}

func handleGame(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "game.html")))
	tmpl.Execute(w, nil)
}

func handleInit(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "init_game.html")))
	tmpl.Execute(w, nil)
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		j1 := r.FormValue("j1")
		j2 := r.FormValue("j2")
		resetGame()

		data := GameData{
			Cols:          []int{0, 1, 2, 3, 4, 5, 6},
			J1:            j1,
			J2:            j2,
			CurrentPlayer: currentPlayer,
		}

		tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "start_game.html")))
		tmpl.Execute(w, data)
		return
	}
	http.Redirect(w, r, "/init", http.StatusSeeOther)
}

func handlePlay(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		colStr := r.FormValue("col")
		j1 := r.FormValue("j1")
		j2 := r.FormValue("j2")

		c, err := strconv.Atoi(colStr)
		if err == nil && c >= 0 && c <= 6 {
			if len(columns[c]) < 6 {
				columns[c] = append(columns[c], currentPlayer)
				currentPlayer = 3 - currentPlayer
			}
		}

		var grid [6][7]int
		for c := 0; c < 7; c++ {
			for i, val := range columns[c] {
				grid[5-i][c] = val
			}
		}

		winner := detectWinner(grid)
		if winner == 1 {
			http.Redirect(w, r, "/winner?name="+j1, http.StatusSeeOther)
			return
		} else if winner == 2 {
			http.Redirect(w, r, "/winner?name="+j2, http.StatusSeeOther)
			return
		}

		data := GameData{
			Grid:          grid,
			Cols:          []int{0, 1, 2, 3, 4, 5, 6},
			J1:            j1,
			J2:            j2,
			CurrentPlayer: currentPlayer,
		}

		tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "start_game.html")))
		tmpl.Execute(w, data)
	}
}

func handleWinner(w http.ResponseWriter, r *http.Request) {
	resetGame()
	winner := r.URL.Query().Get("name")
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "winner.html")))
	tmpl.Execute(w, struct{ Winner string }{Winner: winner})
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handleGame)
	http.HandleFunc("/init", handleInit)
	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/play", handlePlay)
	http.HandleFunc("/winner", handleWinner)

	log.Println("✅ Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}