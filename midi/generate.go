// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jfrazelle/cliaoke/karaoke"
	"github.com/jfrazelle/cliaoke/lyrics"
)

func getSongArtistAndTitle(name string) (string, string) {
	name = strings.TrimSuffix(name, ".mid")
	name = strings.Replace(name, "_", " ", -1)

	parts := strings.SplitN(name, "-", 2)
	if len(parts) < 2 {
		// then the song has no artist, for example "Sonic The Hedgehog"
		return "", strings.TrimSpace(parts[0])
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

// Reads all .mid files in the current folder and creates a manifest.json with
// the song information.
func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	path := filepath.Join(wd, "midi")
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	out, err := os.Create(filepath.Join(path, "manifest.json"))
	if err != nil {
		panic(err)
	}
	defer out.Close()

	var songs []karaoke.Song
	for _, f := range fs {
		// get all the mid files
		if strings.HasSuffix(f.Name(), ".mid") {
			s := karaoke.Song{
				File: f.Name(),
			}

			s.Artist, s.Title = getSongArtistAndTitle(f.Name())

			// initialize the lyrics client
			mmAPIKey := os.Getenv("MUSIXMATCH_APIKEY")
			if mmAPIKey == "" {
				panic("MUSIXMATCH_APIKEY cannot be blank.")
			}
			c := &lyrics.Client{Token: mmAPIKey}

			// search for the track
			track, err := c.SearchTrack(s.Artist + " " + s.Title)
			if err != nil {
				panic(err)
			}

			fmt.Printf("Got track for %s - %s:\n%#v\n\n", s.Artist, s.Title, track)

			// get the lyrics
			s.Lyrics, err = c.GetTrackLyrics(track)
			if err != nil {
				panic(err)
			}

			fmt.Printf("got song:\n%#v\n\n", s)

			songs = append(songs, s)
		}
	}

	b, err := json.MarshalIndent(songs, "", "    ")
	if err != nil {
		panic(err)
	}
	out.Write(b)
}
