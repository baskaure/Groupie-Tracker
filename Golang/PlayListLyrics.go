package GroupieTracker

import (
	"math/rand"
	"time"
)

func ArtistPlaylistLyrics(L *DataLyrics) {
	RandomLyrics(L)
	// M.RandomIndex = 30
	randomTrack := L.Playlist.Tracks.Tracks[L.RandomIndex]

	L.NameMusic = ""

	L.Artist = ""

	L.NameMusicCheck = randomTrack.Track.Name

	L.NameMusic = Artist(randomTrack.Track.Name)

	if len(randomTrack.Track.Artists) > 1 {
		for _, artist := range randomTrack.Track.Artists {
			L.Artist += artist.Name + ", "
		}
	} else {
		for _, artist := range randomTrack.Track.Artists {
			L.Artist += artist.Name
		}
	}
	SearchLyrics(L.Artist, L.NameMusic, L)

}

func RandomLyrics(L *DataLyrics) {
	L.RandomIndex = 0
	rand.Seed(time.Now().UnixNano())
	L.RandomIndex = rand.Intn(len(L.Playlist.Tracks.Tracks))

	CheckIndexLurics(L)
}
