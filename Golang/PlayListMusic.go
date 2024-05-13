package GroupieTracker

import (
	"math/rand"
	"time"

	"github.com/zmb3/spotify"
)

func ArtistPlaylistMusic(M *DataMusic) {
	RandomMusic(M)
	// M.RandomIndex = 30
	randomTrack := M.Playlist.Tracks.Tracks[M.RandomIndex]

	M.NameMusic = ""
	M.PreviewURL = ""

	M.PreviewURL = randomTrack.Track.PreviewURL
	M.NameMusicCheck = randomTrack.Track.Name

	M.NameMusic = Artist(randomTrack.Track.Name)

	CheckPreview(M)

}

func StartPlayback(client *spotify.Client) error {
	err := client.Play()
	if err != nil {
		return err
	}
	return nil
}

func RandomMusic(M *DataMusic) {
	M.RandomIndex = 0

	rand.Seed(time.Now().UnixNano())
	M.RandomIndex = rand.Intn(len(M.Playlist.Tracks.Tracks))

	CheckIndex(M)
}
