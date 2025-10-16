package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

// --- VARIABLES GLOBALES ---
var columns = map[int][]int{
	0: {}, 1: {}, 2: {}, 3: {}, 4: {}, 5: {}, 6: {},
}
var currentPlayer = 1
var joueur1Name string
var joueur2Name string

// --- STRUCTURES ---
type GameData struct {
	Grid          [6][7]int
	Cols          []int
	J1            string
	J2            string
	CurrentPlayer int
}

// --- FONCTIONS UTILES ---
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

// --- HANDLERS ---
func handleInit(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join("templates", "init_game.html"))
	if err != nil {
		http.Error(w, "Erreur chargement personnalisation : "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		joueur1Name = r.FormValue("j1")
		joueur2Name = r.FormValue("j2")
		resetGame()

		data := GameData{
			Cols:          []int{0, 1, 2, 3, 4, 5, 6},
			J1:            joueur1Name,
			J2:            joueur2Name,
			CurrentPlayer: currentPlayer,
		}

		tmpl, err := template.ParseFiles(filepath.Join("templates", "start_game.html"))
		if err != nil {
			http.Error(w, "Erreur template start_game : "+err.Error(), 500)
			return
		}
		tmpl.Execute(w, data)
		return
	}
	http.Redirect(w, r, "/init", http.StatusSeeOther)
}

func handlePlay(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		colStr := r.FormValue("col")
		c, err := strconv.Atoi(colStr)
		if err == nil && c >= 0 && c <= 6 {
			if len(columns[c]) < 6 {
				columns[c] = append(columns[c], currentPlayer)
				currentPlayer = 3 - currentPlayer
			}
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
		http.Redirect(w, r, "/winner1", http.StatusSeeOther)
		return
	} else if winner == 2 {
		http.Redirect(w, r, "/winner2", http.StatusSeeOther)
		return
	}

	data := GameData{
		Grid:          grid,
		Cols:          []int{0, 1, 2, 3, 4, 5, 6},
		J1:            joueur1Name,
		J2:            joueur2Name,
		CurrentPlayer: currentPlayer,
	}

	tmpl, err := template.ParseFiles(filepath.Join("templates", "start_game.html"))
	if err != nil {
		http.Error(w, "Impossible de charger la page : "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func handleWinner1(w http.ResponseWriter, r *http.Request) {
	resetGame()
	tmpl, err := template.ParseFiles(filepath.Join("templates", "winner1.html"))
	if err != nil {
		http.Error(w, "Erreur chargement page victoire J1 : "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func handleWinner2(w http.ResponseWriter, r *http.Request) {
	resetGame()
	tmpl, err := template.ParseFiles(filepath.Join("templates", "winner2.html"))
	if err != nil {
		http.Error(w, "Erreur chargement page victoire J2 : "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func handleAccueil(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/init", http.StatusSeeOther)
}

// --- MAIN ---
func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handleAccueil)
	http.HandleFunc("/init", handleInit)
	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/play", handlePlay)
	http.HandleFunc("/winner1", handleWinner1)
	http.HandleFunc("/winner2", handleWinner2)

	log.Println("✅ Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}