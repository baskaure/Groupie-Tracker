package GroupieTracker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func GameDeaf(w http.ResponseWriter, r *http.Request, L *DataLyrics) {
	fmt.Println("Handle : ", r.FormValue("Handle"))
	if L.Handle == 0 {
		L.Handle = CheckHandle(r.FormValue("Handle"))
	}
	Input := Artist(r.FormValue("GameGuessInput"))
	if L.NameMusic == Input {
		L.CountHandle++
		ArtistPlaylistLyrics(L)
	}
	if L.Handle+1 == L.CountHandle {
		EndGame(w, r)
		fmt.Println("-----------FIN-----------")
		return
	}
	template, err := template.ParseFiles("./HTML/deaf_test.html", "./Template/Timer.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, L)
}

func GameDeafTimer(w http.ResponseWriter, r *http.Request, L *DataLyrics) {
	L.CountHandle++
	if L.Handle+1 == L.CountHandle {
		L.CountHandle = 1
		EndGame(w, r)
		fmt.Println("-----------FIN-----------")
		return
	}
	ArtistPlaylistLyrics(L)

	template, err := template.ParseFiles("./HTML/deaf_test.html", "./Template/Timer.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, L)
}

func SearchLyrics(artist, song string, L *DataLyrics) {
	L.Lyrics = ""
	url := fmt.Sprintf("https://api.lyrics.ovh/v1/%s/%s", artist, song)

	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	// if response.StatusCode != http.StatusOK {
	// 	fmt.Println("request failed with status code %d")
	// }

	var data map[string]string
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		fmt.Println(err)
	}

	lyrics, ok := data["lyrics"]

	ErreurOfAPI := fmt.Sprintf("Désolé nous n'avons pas encore les paroles de %s de %s. Si vous les connaissez, vous pouvez nous les envoyer très simplement en remplissant le formulaire ci dessous.", L.NameMusicCheck, L.Artist)
	if !ok || lyrics == "" || lyrics == ErreurOfAPI {
		ArtistPlaylistLyrics(L)
		SearchLyrics(L.Artist, L.NameMusic, L)
		return
	}

	L.Lyrics = lyrics
	L.Lyrics = RemoveFirstLine(lyrics)

}

func OptionDeafTest(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./HTML/OptionDeaf.html", "./Template/Option.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}
