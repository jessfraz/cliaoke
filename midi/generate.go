// +build ignore

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudflare/cfssl/log"
	"github.com/jfrazelle/cliaoke/karaoke"
	"github.com/jfrazelle/cliaoke/lyrics"
)

func getSongArtistAndTitle(name string) (artist string, title string) {
	name = strings.TrimSuffix(name, ".mid")
	name = strings.Replace(name, "_", " ", -1)

	parts := strings.SplitN(name, "-", 2)
	if len(parts) < 2 {
		// then the song has no artist, for example "Sonic The Hedgehog"
		title = strings.TrimSpace(parts[0])
	} else {
		artist = strings.TrimSpace(parts[0])
		title = strings.TrimSpace(parts[1])
	}

	// clean up grammar for searching for lyrics
	title = strings.Replace(title, "Dont", "Don't", -1)
	title = strings.Replace(title, "Adams", "Adam's", -1)

	return artist, title
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

			// search for the lyrics for the track
			s.Lyrics, err = lyrics.Search(s.Artist + " " + s.Title)
			if err != nil {
				log.Errorf("[%s]: %v", s.Title, err)
			}

			songs = append(songs, s)
		}
	}

	b, err := json.MarshalIndent(songs, "", "    ")
	if err != nil {
		panic(err)
	}
	out.Write(b)
}
