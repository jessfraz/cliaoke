// +build ignore

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/jessfraz/cliaoke/karaoke"
	"github.com/jessfraz/cliaoke/lyrics"
)

func getSongArtistAndTitle(name string) (artist string, title string) {
	name = strings.TrimSuffix(name, ".mid")
	name = strings.Replace(name, "_", " ", -1)

	parts := strings.SplitN(name, " - ", 2)
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
	title = strings.Replace(title, "Ill", "I'll", -1)
	title = strings.Replace(title, "Im", "I'm", -1)
	title = strings.Replace(title, "Gangstas", "Gangsta's", -1)

	artist = strings.Replace(artist, "Destinys", "Destiny's", -1)
	artist = strings.Replace(artist, "Dr ", "Dr. ", -1)

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

	// get all the songs
	remoteSongs, err := karaoke.GetSongList("https://s3.j3ss.co/cliaoke/midi/manifest.json")
	if err != nil {
		logrus.Fatal(err)
	}

	localSongs := map[string]karaoke.Song{}
	for _, f := range fs {
		// get all the mid files
		if strings.HasSuffix(f.Name(), ".mid") {
			s := karaoke.Song{
				File: f.Name(),
			}

			s.Artist, s.Title = getSongArtistAndTitle(f.Name())
			key := strings.Replace(strings.Replace(strings.ToLower(s.Title), " ", "_", -1), "'", "", -1)

			// search for the lyrics for the track if we don't already have it
			if rs, exists := remoteSongs[key]; exists && rs.Lyrics != "" {
				s.Lyrics = rs.Lyrics
			} else {
				s.Lyrics, err = lyrics.Search(s.Artist + " " + s.Title)
				if err != nil {
					logrus.Errorf("[%s]: %v", s.Title, err)
				}
			}

			// make sure the key does not already exist
			if _, exists := localSongs[key]; exists {
				logrus.Errorf("%s already exists in the map, not adding %s", key, s.Title)
				continue
			}

			localSongs[key] = s
		}
	}

	b, err := json.MarshalIndent(localSongs, "", "    ")
	if err != nil {
		panic(err)
	}
	out.Write(b)
}
