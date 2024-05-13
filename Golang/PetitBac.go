package GroupieTracker

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"text/template"
	"time"
)

func PetitBac(w http.ResponseWriter, r *http.Request, PBM *PetitBacData) {
	PBM.Artiste = (r.FormValue("GamePBArtiste"))
	PBM.Album = r.FormValue("GamePBAlbum")
	PBM.MusicGroup = r.FormValue("GamePBMusicGroup")
	PBM.MusicalInstrument = r.FormValue("GamePBMusicalInstrument")
	PBM.Featuring = r.FormValue("GamePBFeaturing")
	PBM.CountHandle++

	// fmt.Println(r.FormValue("HandleName"))
	if PBM.Handle == 0 {
		PBM.Handle = CheckHandle(r.FormValue("HandleName"))
	}
	if PBM.Handle+1 == PBM.CountHandle {
		PBM.CountHandle = 0
		EndGame(w, r)
		fmt.Println("-----------FIN-----------")
		return
	}
	fmt.Println(PBM.Handle)
	fmt.Println(PBM.CountHandle)
	PBM.RandomLetter = RandomLetter()

	template, err := template.ParseFiles("./HTML/petit_bac.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, PBM)
}

func GamePetitBacCheck(w http.ResponseWriter, r *http.Request, PBM *PetitBacData) {

	template, err := template.ParseFiles("./HTML/PetitBacCheck.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, PBM)
	PBM.Artiste = ""
	PBM.Album = ""
	PBM.MusicGroup = ""
	PBM.MusicalInstrument = ""
	PBM.Featuring = ""
}

func OptionPetitBac(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./HTML/OptionPetitBac.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func RandomLetter() string {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(26)
	randomLetter := 'A' + rune(randomIndex)

	return string(randomLetter)
}
func GenerateCategory() {

}
