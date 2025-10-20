package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

var columns map[int][]int
var rows, cols int
var currentPlayer = 1

type GameData struct {
	Grid          [][]int
	Cols          []int
	J1            string
	J2            string
	CurrentPlayer int
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func initColumns() {
	columns = make(map[int][]int)
	for c := 0; c < cols; c++ {
		columns[c] = make([]int, rows)
	}
	currentPlayer = 1
}

func resetGame() {
	initColumns()
}

func placeBlocks(num int) {
	placed := 0
	for placed < num {
		row := rand.Intn(rows)
		col := rand.Intn(cols)
		if columns[col][row] == 0 {
			columns[col][row] = 3 // 3 = bloc
			placed++
		}
	}
}

func detectWinner(grid [][]int) int {
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			token := grid[row][col]
			if token == 0 || token == 3 {
				continue
			}
			// horizontal
			if col+3 < cols &&
				grid[row][col+1] == token &&
				grid[row][col+2] == token &&
				grid[row][col+3] == token {
				return token
			}
			// vertical
			if row+3 < rows &&
				grid[row+1][col] == token &&
				grid[row+2][col] == token &&
				grid[row+3][col] == token {
				return token
			}
			// diagonale descendante
			if row+3 < rows && col+3 < cols &&
				grid[row+1][col+1] == token &&
				grid[row+2][col+2] == token &&
				grid[row+3][col+3] == token {
				return token
			}
			// diagonale montante
			if row-3 >= 0 && col+3 < cols &&
				grid[row-1][col+1] == token &&
				grid[row-2][col+2] == token &&
				grid[row-3][col+3] == token {
				return token
			}
		}
	}
	return 0
}

// =================== Handlers ===================
func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "game.html")))
	tmpl.Execute(w, nil)
}

func handleInit(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "init_game.html")))
	tmpl.Execute(w, nil)
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/init", http.StatusSeeOther)
		return
	}

	j1 := r.FormValue("j1")
	j2 := r.FormValue("j2")
	difficulty := r.FormValue("difficulty")

	switch difficulty {
	case "easy", "medium":
		rows, cols = 6, 7
	case "hard":
		rows, cols = 6, 9
	case "extrem":
		rows, cols = 7, 8
	default:
		rows, cols = 6, 7
	}

	resetGame()

	switch difficulty {
	case "medium":
		placeBlocks(3)
	case "hard":
		placeBlocks(5)
	case "extrem":
		placeBlocks(7)
	}

	data := GameData{
		Grid:          make([][]int, rows),
		Cols:          make([]int, cols),
		J1:            j1,
		J2:            j2,
		CurrentPlayer: currentPlayer,
	}

	for r := 0; r < rows; r++ {
		data.Grid[r] = make([]int, cols)
	}

	for i := 0; i < cols; i++ {
		data.Cols[i] = i
	}

	for c := 0; c < cols; c++ {
		for r := 0; r < rows; r++ {
			data.Grid[rows-1-r][c] = columns[c][r]
		}
	}

	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "start_game.html")))
	tmpl.Execute(w, data)
}

func handlePlay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/start", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	colStr := r.FormValue("col")
	j1 := r.FormValue("j1")
	j2 := r.FormValue("j2")

	c, err := strconv.Atoi(colStr)
	if err == nil && c >= 0 && c < cols {
		for r := 0; r < rows; r++ {
			if columns[c][r] == 0 {
				columns[c][r] = currentPlayer
				currentPlayer = 3 - currentPlayer
				break
			}
		}
	}

	grid := make([][]int, rows)
	for r := 0; r < rows; r++ {
		grid[r] = make([]int, cols)
	}
	for c := 0; c < cols; c++ {
		for r := 0; r < rows; r++ {
			grid[rows-1-r][c] = columns[c][r]
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
		Cols:          make([]int, cols),
		J1:            j1,
		J2:            j2,
		CurrentPlayer: currentPlayer,
	}
	for i := 0; i < cols; i++ {
		data.Cols[i] = i
	}

	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "start_game.html")))
	tmpl.Execute(w, data)
}

func handleWinner(w http.ResponseWriter, r *http.Request) {
	resetGame()
	winner := r.URL.Query().Get("name")
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "winner.html")))
	tmpl.Execute(w, struct{ Winner string }{Winner: winner})
}

// =================== Main ===================
func main() {
	// Fichiers statiques
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes principales
	http.HandleFunc("/", handleHome)        // Page d'accueil
	http.HandleFunc("/init", handleInit)    // Formulaire de saisie des noms
	http.HandleFunc("/start", handleStart)  // Lancement du jeu
	http.HandleFunc("/play", handlePlay)    // Jouer un coup
	http.HandleFunc("/winner", handleWinner)// Page gagnant

	log.Println("✅ Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
