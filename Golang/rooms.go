package GroupieTracker

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

func HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionCookieName)
	if err != nil {
		http.Error(w, "Erreur de session", http.StatusInternalServerError)
		return
	}

	pseudo, ok := session.Values["pseudo"].(string)
	if !ok || pseudo == "" {
		http.Redirect(w, r, loginRoute, http.StatusSeeOther)
		return
	}

	r.ParseForm()
	gameID := r.FormValue("game_id")
	// fmt.Println(gameID)
	maxPlayers := r.FormValue("max_players")
	maxPlayersInt, err := strconv.Atoi(maxPlayers)
	if err != nil {
		http.Error(w, "Nombre maximal de joueurs invalide", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Connexion à la base de données échouée", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	code := generateRoomCode(6)
	_, err = db.Exec("INSERT INTO ROOMS (created_by, max_player, name, id_game, code) VALUES (?, ?, ?, ?, ?)",
		pseudo, maxPlayersInt, fmt.Sprintf("Room %s", code), gameID, code)
	if err != nil {
		log.Println("Erreur lors de la création de la salle:", err)
		http.Error(w, "Échec de la création de la salle", http.StatusInternalServerError)
		return
	}

	Join := struct {
		Code   string
		GameID string
	}{Code: code, GameID: gameID}

	tmpl, err := template.ParseFiles("./HTML/Room.html")
	if err != nil {
		log.Println("Erreur lors du chargement de la page de configuration:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, Join)
}

func HandleJoinRoom(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionCookieName)
	if err != nil {
		http.Error(w, "Erreur de session", http.StatusInternalServerError)
		return
	}

	pseudo, ok := session.Values["pseudo"].(string)
	if !ok || pseudo == "" {
		http.Redirect(w, r, loginRoute, http.StatusSeeOther)
		return
	}

	r.ParseForm()
	roomCode := r.FormValue("Room")
	log.Println("Code de la salle saisi:", roomCode)

	if roomCode == "" {
		http.Error(w, "Le code de la salle est requis", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Connexion à la base de données échouée", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var roomID int
	err = db.QueryRow("SELECT id FROM ROOMS WHERE code = ?", roomCode).Scan(&roomID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Code de la salle invalide:", roomCode)
			http.Error(w, "Code de la salle invalide", http.StatusBadRequest)
			return
		}
		log.Println("Erreur lors de la vérification du code de la salle:", err)
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	log.Println("Salle trouvée, ID:", roomID)

	tmpl, err := template.ParseFiles("./HTML/Join.html")
	if err != nil {
		log.Println("Erreur lors du chargement de la page de configuration:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, struct{ Code string }{Code: roomCode})
}
