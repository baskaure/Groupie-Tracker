package GroupieTracker

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const (
	loginRoute  = "/Login"
	signupRoute = "/Signup"
	homeRoute   = "/Home"
	logoutRoute = "/Logout"
	staticRoute = "/static/"

	// deafTestRoute  = "/Game/1"
	// blindTestRoute = "/Game/2"
	// petitBacRoute        = "/Game/3"
	// JoinRoom             = "/JoinRoom"
	// ConfigureRoomRoute   = "/ConfigureRoom"
	createRoomRoute      = "/CreateRoom"
	chatRoute            = "/chat"
	websocketRoute       = "/ws"
	sessionCookieName    = "session-name"
	databasePath         = "./BDD.db"
	dbConnectionError    = "Connexion à la base de données échouée"
	dbQueryError         = "Erreur de requête de base de données"
	dbQueryLogError      = "Erreur requête DB:"
	dbConnectionLogError = "Erreur connexion DB:"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func init() {
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
	}

	rand.Seed(time.Now().UnixNano())
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateRoomCode(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func isValidPassword(password string) bool {
	if len(password) < 12 {
		return false
	}
	var (
		uppercase   = regexp.MustCompile(`[A-Z]`)
		lowercase   = regexp.MustCompile(`[a-z]`)
		number      = regexp.MustCompile(`[0-9]`)
		specialChar = regexp.MustCompile(`[#$&+,:;=?@^_|~]`)
	)
	return uppercase.MatchString(password) && lowercase.MatchString(password) &&
		number.MatchString(password) && specialChar.MatchString(password)
}

func HandleSignup(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleSignup: Début du traitement")

	if r.Method != "POST" {
		log.Println("HandleSignup: Méthode non autorisée")
		http.Redirect(w, r, signupRoute, http.StatusSeeOther)
		return
	}
	r.ParseForm()
	username := r.FormValue("pseudo")
	email := r.FormValue("email")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("password_confirm")
	log.Println("HandleSignup: pseudo =", username, ", email =", email)

	if password != passwordConfirm {
		errorMessage := "Les mots de passe ne correspondent pas"
		data := map[string]interface{}{
			"ErrorMessage": errorMessage,
		}
		tmpl, err := template.ParseFiles("HTML/signup.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)
		return
	}
	if !isValidPassword(password) {
		errorMessage := "Le mot de passe doit contenir au moins une majuscule, un chiffre, un caractère spécial et avoir au moins 12 caractères de long."
		data := map[string]interface{}{
			"ErrorMessage": errorMessage,
		}
		tmpl, err := template.ParseFiles("HTML/signup.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)
		return
	}

	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Println(dbConnectionLogError, err)
		http.Error(w, dbConnectionError, http.StatusInternalServerError)
		return
	}
	defer db.Close()
	log.Println("HandleSignup: Connexion à la base de données réussie")

	var userCount int
	err = db.QueryRow("SELECT COUNT(*) FROM USER WHERE pseudo = ? OR email = ?", username, email).Scan(&userCount)
	if err != nil {
		log.Println(dbQueryLogError, err)
		http.Error(w, dbQueryError, http.StatusInternalServerError)
		return
	}
	if userCount > 0 {
		errorMessage := "Nom d'utilisateur ou email déjà enregistré"
		data := map[string]interface{}{
			"ErrorMessage": errorMessage,
		}
		tmpl, err := template.ParseFiles("HTML/signup.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("HandleSignup: Erreur lors du hashage du mot de passe:", err)
		http.Error(w, "Erreur lors du traitement du mot de passe", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO USER (pseudo, email, password) VALUES (?, ?, ?)", username, email, hashedPassword)
	if err != nil {
		log.Println("HandleSignup: Erreur lors de l'insertion de l'utilisateur:", err)
		http.Error(w, "Échec de l'inscription de l'utilisateur", http.StatusInternalServerError)
		return
	}
	log.Println("HandleSignup: Utilisateur inscrit avec succès")

	http.Redirect(w, r, loginRoute, http.StatusSeeOther)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleLogin: Début du traitement")

	if r.Method != "POST" {
		log.Println("HandleLogin: Méthode non autorisée")
		http.Redirect(w, r, loginRoute, http.StatusSeeOther)
		return
	}
	r.ParseForm()
	login := r.FormValue("login")
	password := r.FormValue("password")
	log.Println("HandleLogin: login =", login)

	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Println(dbConnectionLogError, err)
		http.Error(w, dbConnectionError, http.StatusInternalServerError)
		return
	}
	defer db.Close()
	log.Println("HandleLogin: Connexion à la base de données réussie")

	var storedPassword string
	var username string
	err = db.QueryRow("SELECT pseudo, password FROM USER WHERE pseudo = ? OR email = ?", login, login).Scan(&username, &storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			errorMessage := "Identifiant ou mot de passe incorrect"
			data := map[string]interface{}{
				"ErrorMessage": errorMessage,
			}
			tmpl, err := template.ParseFiles("HTML/login.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, data)
			return
		}
		log.Println(dbQueryLogError, err)
		http.Error(w, dbQueryError, http.StatusInternalServerError)
		return
	}
	log.Println("HandleLogin: Utilisateur trouvé:", username)

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		errorMessage := "Identifiant ou mot de passe incorrect"
		data := map[string]interface{}{
			"ErrorMessage": errorMessage,
		}
		tmpl, err := template.ParseFiles("HTML/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)
		return
	}
	log.Println("HandleLogin: Mot de passe correct")

	session, err := store.Get(r, sessionCookieName)
	if err != nil {
		log.Println("Erreur lors de l'obtention de la session:", err)
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	session.Values["pseudo"] = username
	if err := session.Save(r, w); err != nil {
		log.Println("Erreur lors de l'enregistrement de la session:", err)
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	log.Println("HandleLogin: Session mise à jour avec succès")

	http.Redirect(w, r, homeRoute, http.StatusSeeOther)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionCookieName)
	pseudo, ok := session.Values["pseudo"]

	if !ok {
		http.Redirect(w, r, loginRoute, http.StatusSeeOther)
		return
	}

	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Println(dbConnectionLogError, err)
		http.Error(w, dbConnectionError, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM GAMES")
	if err != nil {
		log.Println(dbQueryLogError, err)
		http.Error(w, dbQueryError, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	games := []struct {
		ID   int
		Name string
	}{}

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Println("Erreur lecture DB:", err)
			continue
		}
		games = append(games, struct {
			ID   int
			Name string
		}{ID: id, Name: name})
	}

	data := map[string]interface{}{
		"Pseudo": pseudo,
		"Games":  games,
	}

	tmpl, err := template.ParseFiles("HTML/homeTest.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func DeafTestPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./HTML/deaf_test.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, r)
}

func BlindTestPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./HTML/blind_test.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, r)
}

func PetitBacPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./HTML/petit_bac.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func getSessionPseudo(r *http.Request) (string, error) {
	session, err := store.Get(r, sessionCookieName)
	if err != nil {
		return "", err
	}

	if pseudo, ok := session.Values["pseudo"].(string); ok && pseudo != "" {
		return pseudo, nil
	}

	return "", fmt.Errorf("session invalide ou expirée")
}

func redirectUserToRoom(w http.ResponseWriter, r *http.Request, roomCode string) {
	http.Redirect(w, r, fmt.Sprintf("/Room/%s", roomCode), http.StatusSeeOther)
}

func ChatPage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionCookieName)
	pseudo, ok := session.Values["pseudo"].(string)
	if !ok {
		http.Redirect(w, r, loginRoute, http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"Pseudo": pseudo,
	}

	tmpl, err := template.ParseFiles("HTML/chat.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionCookieName)
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, loginRoute, http.StatusSeeOther)
}
func HandleConfigureRoom(w http.ResponseWriter, r *http.Request, M *DataMusic, L *DataLyrics, PBM *PetitBacData) {

	if M.Handle == 0 {
		M.Handle = CheckHandle(r.FormValue("Handle"))
	}
	if L.Handle == 0 {
		L.Handle = CheckHandle(r.FormValue("Handle"))
	}
	if PBM.Handle == 0 {
		PBM.Handle = CheckHandle(r.FormValue("Handle"))
	}

	session, err := store.Get(r, sessionCookieName)
	if err != nil {
		log.Println("Erreur lors de l'accès à la session:", err)
		http.Error(w, "Erreur de session", http.StatusInternalServerError)
		return
	}

	pseudo, ok := session.Values["pseudo"]
	if !ok {
		log.Println("Session invalide ou expirée")
		http.Redirect(w, r, loginRoute, http.StatusSeeOther)
		return
	}

	log.Println("Session active pour l'utilisateur:", pseudo)
	fmt.Println("gameid", r.FormValue("game_id"))
	gameID := r.FormValue("game_id") //r.URL.Query().Get("game_id")
	fmt.Println(gameID)
	if gameID == "" {
		http.Error(w, "Jeu non spécifié", http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"Pseudo": pseudo,
		"GameID": gameID,
	}

	tmpl, err := template.ParseFiles("HTML/configure_room.html")
	if err != nil {
		log.Println("Erreur lors du chargement de la page de configuration:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}
func GameInfo(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("HTML/game.html")
	if err != nil {
		log.Println("Erreur lors du chargement de la page de configuration:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, r)
}
