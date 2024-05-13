package GroupieTracker

import "github.com/zmb3/spotify"

type DataLyrics struct {
	NameMusic      string
	Lyrics         string
	Artist         string
	NameMusicCheck string
	ArtistIndex    []int
	RandomIndex    int
	Playlist       *spotify.FullPlaylist
	Handle         int
	CountHandle    int
}
type DataMusic struct {
	NameMusic      string
	PreviewURL     string
	ValueInput     string
	Playlist       *spotify.FullPlaylist
	ArtistIndex    []int
	RandomIndex    int
	NameMusicCheck string
	Handle         int
	CountHandle    int
}

type PetitBacData struct {
	RandomLetter      string
	Artiste           string
	Album             string
	MusicGroup        string
	MusicalInstrument string
	Featuring         string
	Handle            int
	CountHandle       int
}

type Data struct {
	Rows int
}
