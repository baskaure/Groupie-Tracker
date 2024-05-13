package GroupieTracker

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan map[string]string)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erreur lors de la connexion WebSocket :", err)
		return
	}
	defer func() {
		log.Println("WebSocket fermé pour", ws.RemoteAddr())
		ws.Close()
	}()

	clients[ws] = true

	for {
		var msg map[string]string
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Erreur lors de la lecture du message :", err)
			delete(clients, ws)
			break
		}
		log.Println("Message reçu de", msg["pseudo"], ":", msg["message"])
		broadcast <- msg
	}
}

func BroadcastMessages() {
	for msg := range broadcast {
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("Erreur lors de l'envoi du message :", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
