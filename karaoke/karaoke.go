package karaoke

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	fluidsynthBinary = "fluidsynth"
)

// Song holds information about the artist, title and filename for a song.
type Song struct {
	Artist string
	Title  string
	File   string
	Lyrics string
}

// DownloadSong downloads a song from remotepath to localpath.
func DownloadSong(localpath, remotepath string) error {
	if err := os.MkdirAll(filepath.Dir(localpath), 0755); err != nil {
		return fmt.Errorf("creating directory %s failed: %v", filepath.Dir(localpath), err)
	}
	f, err := os.Create(localpath)
	if err != nil {
		return fmt.Errorf("creating %s failed: %v", localpath, err)
	}
	defer f.Close()

	resp, err := http.Get(remotepath)
	if err != nil {
		return fmt.Errorf("request to %s failed: %v", remotepath, err)
	}
	defer resp.Body.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("downloading %s to %s failed: %v", remotepath, localpath, err)
	}

	return nil
}

// GetSongList returns a list of songs from a manifest url.
func GetSongList(uri string) (songs map[string]Song, err error) {
	resp, err := http.Get(uri)
	if err != nil {
		return songs, fmt.Errorf("request to %s failed: %v", uri, err)
	}
	defer resp.Body.Close()

	// decode the body
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&songs); err != nil {
		return songs, fmt.Errorf("decoding response from %s failed: %v", uri, err)
	}

	return songs, nil
}

// FluidsynthBinaryExists checks if the fluidsynth binary exists.
func FluidsynthBinaryExists() bool {
	_, err := exec.LookPath(fluidsynthBinary)
	if err != nil {
		return false
	}

	// Return true when there is no error.
	return err == nil
}
