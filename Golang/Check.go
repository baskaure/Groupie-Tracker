package GroupieTracker

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func CheckPreview(M *DataMusic) {
	if M.PreviewURL == "" {
		ArtistPlaylistMusic(M)
	}
}

func CheckHandle(Check string) int {
	value, err := strconv.Atoi(Check)
	if err != nil {
		fmt.Println("Erreur de conversion en entier:", err)
		return 5
	}
	return value
}

func CheckIndex(M *DataMusic) {
	for _, value := range M.ArtistIndex {
		if len(M.ArtistIndex) == len(M.Playlist.Tracks.Tracks) {

			fmt.Println("vous avez trouver tout la liste")
			return
		} else if M.RandomIndex == value {
			RandomMusic(M)
			return
		}
	}
	M.ArtistIndex = append(M.ArtistIndex, M.RandomIndex)
}
func CheckIndexLurics(L *DataLyrics) {
	for _, value := range L.ArtistIndex {
		if len(L.ArtistIndex) == len(L.Playlist.Tracks.Tracks) {

			fmt.Println("vous avez trouver tout la liste")
			return
		} else if L.RandomIndex == value {
			RandomLyrics(L)
			return
		}
	}
	L.ArtistIndex = append(L.ArtistIndex, L.RandomIndex)
}

func Artist(Input string) string {
	Word := ReplaceAndWithAmpersand(Input)
	output := strings.Map(func(r rune) rune {
		switch r {
		case ' ', '\'', '"', ',', ':', '.', '?', '+', '!', '€', '$':
			return -1 // Supprime le caractère
		default:
			return r
		}
	}, RemoveAccents(strings.ToLower(RemoveAfterDash(Word))))
	re := regexp.MustCompile(`\([^)]*\)`)
	return re.ReplaceAllString(output, "")
}

func RemoveAccents(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r)
	}), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

func RemoveAfterDash(s string) string {
	parts := strings.SplitN(s, " - ", 2)
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}
	return s
}

func RemoveFirstLine(s string) string {
	lines := strings.Split(s, "\n")
	rand.Seed(time.Now().UnixNano())
	Index := 0
	Index = rand.Intn(len(lines)-12) + 1

	result := ""
	if len(lines) > 1 {
		for i := 0; i < 12; i++ {
			result += lines[Index] + "<br>"
			Index++
		}
		return result
	}
	return ""
}
func ReplaceAndWithAmpersand(s string) string {
	s = strings.ReplaceAll(s, " et ", " & ")
	s = strings.ReplaceAll(s, " and ", " & ")
	s = strings.ReplaceAll(s, " plus ", " + ")

	return s
}
