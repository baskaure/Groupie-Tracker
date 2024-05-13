package GroupieTracker

import (
	"net/http"
)

func Server() {
	var M *DataMusic
	var PBM *PetitBacData
	var L *DataLyrics

	M = new(DataMusic)
	PBM = new(PetitBacData)
	L = new(DataLyrics)

	API(M, L, PBM)

	http.HandleFunc("/Game/2", OptionBlindTest)
	http.HandleFunc("/Game/1", OptionDeafTest)
	http.HandleFunc("/Game/3", OptionPetitBac)
	http.HandleFunc("/Game", GameInfo)

	http.HandleFunc("/2", func(w http.ResponseWriter, r *http.Request) {
		GameBlind(w, r, M)
	})
	http.HandleFunc("/GameBlindTestTimer", func(w http.ResponseWriter, r *http.Request) {
		GameBlindTimer(w, r, M)
	})
	http.HandleFunc("/1", func(w http.ResponseWriter, r *http.Request) {
		GameDeaf(w, r, L)
	})
	http.HandleFunc("/GameDeafTestTimer", func(w http.ResponseWriter, r *http.Request) {
		GameDeafTimer(w, r, L)
	})
	http.HandleFunc("/3", func(w http.ResponseWriter, r *http.Request) {
		PetitBac(w, r, PBM)
	})
	http.HandleFunc("/GamePetitBacCheck", func(w http.ResponseWriter, r *http.Request) {
		GamePetitBacCheck(w, r, PBM)
	})
	http.HandleFunc("/ConfigureRoom", func(w http.ResponseWriter, r *http.Request) {
		HandleConfigureRoom(w, r, M, L, PBM)
	})

	http.HandleFunc(websocketRoute, WebSocketHandler)
	go BroadcastMessages()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, loginRoute, http.StatusSeeOther)
	})
	http.HandleFunc(loginRoute, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			http.ServeFile(w, r, "HTML/login.html")
		} else {
			HandleLogin(w, r)
		}
	})
	http.HandleFunc(signupRoute, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			http.ServeFile(w, r, "HTML/signup.html")
		} else {
			HandleSignup(w, r)
		}
	})

	http.HandleFunc(homeRoute, HomePage)
	http.HandleFunc(logoutRoute, HandleLogout)
	http.HandleFunc(createRoomRoute, HandleCreateRoom)
	http.HandleFunc(chatRoute, ChatPage)
	http.HandleFunc("/JoinRoom", HandleJoinRoom)
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)

}
