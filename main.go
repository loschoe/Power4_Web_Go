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

// =================== Variables globales ===================
var (
	columns          map[int][]int
	rows, cols       int
	moveCount 		 int
	currentPlayer    = 1
	j1Global         string
	j2Global         string
	difficultyGlobal string
	gravityDown		     bool = true
	gravityEnabled    bool = true 
)

// =================== Structures ===================
type GameData struct {
	Grid          [][]int
	Cols          []int
	J1            string
	J2            string
	CurrentPlayer int
	Message 	  string
}

// =================== Initialisation ===================
func init() {
	rand.Seed(time.Now().UnixNano())
}

// Initialiser les colonnes du jeu 
func initColumns() {
	columns = make(map[int][]int)
	for c := 0; c < cols; c++ {
		columns[c] = make([]int, rows)
	}
	currentPlayer = 1
}

// Reset le jeu après une partie 
func resetGame() {
	initColumns()
	moveCount = 0
	gravityDown = true
}

// Placer un bloc plein
func placeBlocks(num int) {
	placed := 0
	for placed < num {
		row := rand.Intn(rows)
		col := rand.Intn(cols)
		if columns[col][row] == 0 {
			columns[col][row] = 3 								 
			placed++
		}
	}
}

// Gravité 
func Gravity() {
	if !gravityEnabled {					// Ne pas toucher à la gravité du mode easy 
		return 
	}
	moveCount++								// Initialisation des 6 coups pour les autres difficultés 
	if moveCount%6 == 0 {
		gravityDown = !gravityDown
	}
}

// =================== Logique du jeu ===================
// Détection d'un gagnant 
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

// Détection d'une égalité (grille pleine sans gagnant)
func isDraw(grid [][]int) bool {
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == 0 {
				return false // encore une case vide → pas égalité
			}
		}
	}
	return true
}

// =================== Handlers (Afficher les différentes pages)===================
func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "game.html")))
	tmpl.Execute(w, nil)
}

func handleInit(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "init_game.html")))
	tmpl.Execute(w, nil)
}

// Config de la difficulté 
func setupGame(j1, j2, difficulty string) GameData {
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

	// Gravité désactivée en mode facile
	if difficulty == "easy" {
		gravityEnabled = false
	} else {
		gravityEnabled = true
	}

	// Le nombre de blocs pleins à placer 
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

	return data
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/init", http.StatusSeeOther)
		return
	}

	j1 := r.FormValue("j1")
	j2 := r.FormValue("j2")
	difficulty := r.FormValue("difficulty")

	// Sauvegarde
	j1Global = j1
	j2Global = j2
	difficultyGlobal = difficulty

	data := setupGame(j1, j2, difficulty)

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
	c, err := strconv.Atoi(colStr)
	if err == nil && c >= 0 && c < cols {
		if gravityDown {
			// Gravité normale : les jetons tombent vers le bas
			for r := 0; r < rows; r++ {
				if columns[c][r] == 0 {
					columns[c][r] = currentPlayer
					currentPlayer = 3 - currentPlayer
					break
				}
			}
		} else {
			// Gravité inversée : les jetons "montent" vers le haut
			for r := rows - 1; r >= 0; r-- {
				if columns[c][r] == 0 {
					columns[c][r] = currentPlayer
					currentPlayer = 3 - currentPlayer
					break
				}
			}
		}
		Gravity()
	}

	// Mise à jour de la grille
	grid := make([][]int, rows)
	for r := 0; r < rows; r++ {
		grid[r] = make([]int, cols)
	}
	for c := 0; c < cols; c++ {
		for r := 0; r < rows; r++ {
			grid[rows-1-r][c] = columns[c][r]
		}
	}

	// Détermination du gagnant ou égalité
	winner := detectWinner(grid)
	winnerName := ""
	if winner == 1 {
		winnerName = j1Global
	} else if winner == 2 {
		winnerName = j2Global
	} else if isDraw(grid) {
		winnerName = "Aucun gagnant"
	}

	if winnerName != "" {
		http.Redirect(w, r, "/winner?name="+winnerName, http.StatusSeeOther)
		return
	}

	// les messages d'alerte du changement de gravité 
	msg := ""
	if gravityEnabled && moveCount%6 == 0 {
		if gravityDown {
			msg = "💡 Gravité réactivée — les jetons retombent !"
		} else {
			msg = "⚠️ Gravité désactivée — les jetons restent en haut !"
		}
	}

	data := GameData{
		Grid:          grid,
		Cols:          make([]int, cols),
		J1:            j1Global,
		J2:            j2Global,
		CurrentPlayer: currentPlayer,
		Message:       msg,
	}
	for i := 0; i < cols; i++ {
		data.Cols[i] = i
	}

	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "start_game.html")))
	tmpl.Execute(w, data)
}

func handleWinner(w http.ResponseWriter, r *http.Request) {
	winner := r.URL.Query().Get("name")

	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "winner.html")))
	tmpl.Execute(w, struct {
		Winner string
		J1     string
		J2     string
	}{
		Winner: winner,
		J1:     j1Global,
		J2:     j2Global,
	})
}

// La revanche
func handleRevanche(w http.ResponseWriter, r *http.Request) {
	data := setupGame(j1Global, j2Global, difficultyGlobal)

	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "start_game.html")))
	tmpl.Execute(w, data)
}

// =================== Main ===================
func main() {
	// Fichiers statiques
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes principales
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/init", handleInit)
	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/play", handlePlay)
	http.HandleFunc("/winner", handleWinner)
	http.HandleFunc("/revanche", handleRevanche)

	log.Println("✅ Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
