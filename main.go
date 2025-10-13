package main

import (
	"html/template"			// Relier le go à l'HTML
	"log"					// Surveiller les bugs 
	"net/http"				// Le lien avec le web 
	"path/filepath"			
	"strconv"				// Conversions de types 
)

// --- VARIABLES GLOBALES ---
// Le tableau avec chaque colonne du puissance 4 initialisé : 
var columns = map[int][]int{
	0: {}, 1: {}, 2: {}, 3: {}, 4: {}, 5: {}, 6: {},
}

// currentPlayer est soit le joueur 1 soit le joueur 2 
var currentPlayer = 1

// --- FONCTIONS UTILES AU FONCTIONNEMENT ---

// resetGame réinitialise le plateau 
func resetGame() {
	columns = map[int][]int{
		0: {}, 1: {}, 2: {}, 3: {}, 4: {}, 5: {}, 6: {},
	}
	currentPlayer = 1
}

// detectWinner renvoie 1 ou 2 si un joueur a aligné 4 jetons, sinon 0
// Tocken c'est les jetons 
func detectWinner(grid [6][7]int) int {
	for row := 0; row < 6; row++ {
		for col := 0; col < 7; col++ {
			token := grid[row][col]
			if token == 0 {
				continue
			}

			// 4 horizontaux
			if col+3 < 7 &&
				grid[row][col+1] == token &&
				grid[row][col+2] == token &&
				grid[row][col+3] == token {
				return token
			}

			// 4 verticaux
			if row+3 < 6 &&
				grid[row+1][col] == token &&
				grid[row+2][col] == token &&
				grid[row+3][col] == token {
				return token
			}

			// diagonale ↘
			if row+3 < 6 && col+3 < 7 &&
				grid[row+1][col+1] == token &&
				grid[row+2][col+2] == token &&
				grid[row+3][col+3] == token {
				return token
			}

			// diagonale ↗
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

// Page d'accueil du jeu 
func handleGame(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("templates", "game.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Impossible de charger la page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Page de victoire joueur 1
func handleWinner1(w http.ResponseWriter, r *http.Request) {
	resetGame()
	tmpl, err := template.ParseFiles(filepath.Join("templates", "winner1.html"))
	if err != nil {
		http.Error(w, "Erreur chargement page victoire J1", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Page de victoire joueur 2
func handleWinner2(w http.ResponseWriter, r *http.Request) {
	resetGame()
	tmpl, err := template.ParseFiles(filepath.Join("templates", "winner2.html"))
	if err != nil {
		http.Error(w, "Erreur chargement page victoire J2", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Page principale du jeu
func handlePlay(w http.ResponseWriter, r *http.Request) {
	resetGame()
	if r.Method == http.MethodPost {
		r.ParseForm()
		colStr := r.FormValue("col")
		c, err := strconv.Atoi(colStr)
		if err == nil && c >= 0 && c <= 6 {
			// Vérifie que la colonne n’est pas pleine
			if len(columns[c]) < 6 {
				columns[c] = append(columns[c], currentPlayer)
				currentPlayer = 3 - currentPlayer // alterner 1 <-> 2
			}
		}
	}

	// Reconstruire après un tour 
	var grid [6][7]int
	for c := 0; c < 7; c++ {
		for i, val := range columns[c] {
			grid[5-i][c] = val
		}
	}

	// Vérification du gagnant
	winner := detectWinner(grid)
	if winner == 1 {
		http.Redirect(w, r, "/winner1", http.StatusSeeOther)
		return
	} else if winner == 2 {
		http.Redirect(w, r, "/winner2", http.StatusSeeOther)
		return
	}

	// Données envoyées au template pour le mettre à jour
	data := struct {
		Grid [6][7]int
		Cols []int
	}{
		Grid: grid,
		Cols: []int{0, 1, 2, 3, 4, 5, 6},
	}

	// Afficher la page du jeu
	tmplPath := filepath.Join("templates", "start_game.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Impossible de charger la page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// --- MAIN ---

func main() {
	// Servir les fichiers statiques (CSS, images, etc.)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handleGame)
	http.HandleFunc("/play", handlePlay)
	http.HandleFunc("/winner1", handleWinner1)
	http.HandleFunc("/winner2", handleWinner2)

	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
