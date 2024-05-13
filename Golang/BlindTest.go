package GroupieTracker

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func GameBlind(w http.ResponseWriter, r *http.Request, M *DataMusic) {

	fmt.Println("MMM", M.CountHandle)

	fmt.Println("Handle : ", r.FormValue("Handle"))
	if M.Handle == 0 {
		M.Handle = CheckHandle(r.FormValue("Handle"))
	}
	fmt.Println(M.Handle)

	// fmt.Println("timerSong", r.FormValue("Song"))
	// fmt.Println("timerAnswer", r.FormValue("Answer"))

	Input := Artist(r.FormValue("GameBlindInput"))
	fmt.Println("reponce")
	// fmt.Println(Input)
	fmt.Println(M.NameMusic)
	if M.NameMusic == Input {
		M.CountHandle++
		fmt.Println("Juste")
		ArtistPlaylistMusic(M)
	}
	if M.Handle+1 == M.CountHandle {
		EndGame(w, r)
		fmt.Println("-----------FIN-----------")
		return
	}
	fmt.Println("MMM", M.CountHandle)

	template, err := template.ParseFiles("./HTML/blind_test.html", "./Template/Timer.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, M)

}
func GameBlindTimer(w http.ResponseWriter, r *http.Request, M *DataMusic) {
	fmt.Println()
	fmt.Println("Timer")
	fmt.Println()
	M.CountHandle++
	if M.Handle+1 == M.CountHandle {
		M.CountHandle = 1
		EndGame(w, r)
		fmt.Println("-----------FIN-----------")
		return
	}
	ArtistPlaylistMusic(M)
	template, err := template.ParseFiles("./HTML/Blind.html", "./Template/Timer.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, M)

}
func OptionBlindTest(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./HTML/OptionBlind.html", "./Template/Option.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}
