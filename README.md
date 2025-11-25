# 🎮 Projet Power4 - Go + Web 

## 🚀 Présentation   :
Bienvenue dans le dépôt GitHub du **Projet Power4**, un jeu sur le navigateur qui reprend le célèbre Puissance 4 avec des ajouts personnels (blocks, gravity)...
Ce jeu est développé dans le cadre d'un module à **STRASBOURG Ynov Campus**.

## 📄 Fonctionnalités :
- Plusieurs grilles sont disponibles selon la difficulté : 6x7 ; 6x9 ; 7x8
- Certains difficultés ont de la *gravité* tout les 6 coups et ont des *blocs pleins* 
- Deux joueurs jouent à tour de rôle sur la même machine 
- Détection automatique de victoire ou d'égalité
- Interface simple pour jouer directement depuis le navigateur

## 🛠️ Installation et exécution :
### 1. Cloner le dépôt
```bash
git clone https://github.com/loschoe/Power4_Web_Go.git
```
### 2. Installer les dépendances Go
```bash
go mod tidy
```
### 3. Lancer le serveur
```bash
go run main.go
```
### 4. Jouer 
Ouvrez votre navigateur et allez sur ```http://localhost:8080```.

## 💡 Langages & tech utilisés :
- Backend : Golang
- Frontend : HTML / CSS
- Serveur HTTP : net/http de Go
