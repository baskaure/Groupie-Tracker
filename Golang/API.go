package GroupieTracker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

func API(M *DataMusic, L *DataLyrics, PBM *PetitBacData) {
	authConfig := &clientcredentials.Config{
		ClientID:     "04c11e40f0574b2ca4494b3920af738e",
		ClientSecret: "c580de3dac594966922342bf68431b8d",
		TokenURL:     spotify.TokenURL,
	}

	accessToken, err := authConfig.Token(context.Background())
	if err != nil {
		log.Fatalf("error retrieve access token: %v", err)

	}

	client := spotify.Authenticator{}.NewClient(accessToken)
	playlistID := spotify.ID("0dsI111HmPsDWHHWSAJQSW")
	M.Playlist, err = client.GetPlaylist(playlistID)
	L.Playlist, err = client.GetPlaylist(playlistID)

	if err != nil {
		log.Fatalf("error retrieve playlist data: %v", err)
	}
	// CheckLyricsMusic(L)
	ArtistPlaylistLyrics(L)
	ArtistPlaylistMusic(M)
	StartPlayback(&client)
	M.CountHandle = 1
	L.CountHandle = 1
	PBM.CountHandle = 0
}

func CheckLyricsMusic(L *DataLyrics) {

	// for i, Music := range L.Playlist.Tracks.Tracks {
	for i := 0; i < len(L.Playlist.Tracks.Tracks); i++ {
		randomTrack := L.Playlist.Tracks.Tracks[i]

		song := randomTrack.Track.Name
		artistss := ""
		if len(randomTrack.Track.Artists) > 1 {
			for _, artist := range randomTrack.Track.Artists {
				artistss += artist.Name + ", "
			}
		} else {
			for _, artist := range randomTrack.Track.Artists {
				artistss += artist.Name
			}
		}
		// fmt.Println(artistss, song)
		L.Lyrics = ""
		url := fmt.Sprintf("https://api.lyrics.ovh/v1/%s/%s", artistss, song)

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
			fmt.Println(i + 1)
		}

	}

}
