package GroupieTracker

import (
	"log"
	"net/http"
	"text/template"
)

func EndGame(w http.ResponseWriter, r *http.Request) {

	rows := 5

	data := Data{Rows: rows}
	template, err := template.ParseFiles("./HTML/EndGame.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, data)
}
