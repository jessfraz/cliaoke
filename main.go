package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/homedir"
	"github.com/docker/docker/pkg/term"
	"github.com/jfrazelle/cliaoke/karaoke"
)

const (
	// BANNER is what is printed for help/info output
	BANNER = `      _ _             _
  ___| (_) __ _  ___ | | _____
 / __| | |/ _` + "`" + ` |/ _ \| |/ / _ \
| (__| | | (_| | (_) |   <  __/
 \___|_|_|\__,_|\___/|_|\_\___|

 Command Line Karaoke
 Version: %s

`
	// VERSION is the binary version.
	VERSION = "v0.1.0"

	defaultSongStore = ".cliaoke"
	midiURI          = "https://s3.j3ss.co/cliaoke/midi"
)

var (
	songRequested string
	debug         bool
	version       bool
)

func init() {
	// parse flags
	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&version, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&debug, "d", false, "run in debug mode")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, VERSION))
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() >= 1 {
		songRequested = strings.Join(flag.Args(), " ")
	}

	if songRequested == "help" {
		usageAndExit("", 0)
	}

	if version || songRequested == "version" {
		fmt.Printf("%s", VERSION)
		os.Exit(0)
	}

	// set log level
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

//go:generate go run midi/generate.go

func main() {
	// get all the songs
	songs, err := getSongList()
	if err != nil {
		logrus.Fatal(err)
	}

	if songRequested == "" {
		// print songs table
		_, stdout, _ := term.StdStreams()
		w := tabwriter.NewWriter(stdout, 20, 1, 3, ' ', 0)

		// print header
		fmt.Fprintln(w, "COMMAND\tTITLE\tARTIST")

		// print the keys alphabetically
		printSorted := func(m map[string]karaoke.Song) {
			mk := make([]string, len(m))
			i := 0
			for k := range m {
				mk[i] = k
				i++
			}
			sort.Strings(mk)

			for _, key := range mk {
				fmt.Fprintf(w, "%s\t%s\t%s\n", key, m[key].Title, m[key].Artist)
			}
		}

		printSorted(songs)
		w.Flush()
		return
	}

	// play requested song
	fmt.Println("DJ cliaoke on the request line.")

	// find the song
	song, ok := songs[songRequested]
	if !ok {
		logrus.Fatalf("Could not find %s in song list, run with no arguments to see the list of songs available.", songRequested)
	}

	// download the song if we can't find it locally
	fmt.Printf("Bringing up the track %s...\n", song.Title)
	home := homedir.Get()
	localmid := filepath.Join(home, defaultSongStore, song.File)
	if _, err := os.Stat(localmid); os.IsNotExist(err) {
		if err := downloadSong(localmid, song.File); err != nil {
			logrus.Fatal(err)
		}
	}

	// play the song
	if err := play(localmid); err != nil {
		logrus.Fatal(err)
	}
}

func downloadSong(localpath, remotepath string) error {
	if err := os.MkdirAll(filepath.Dir(localpath), 0755); err != nil {
		return fmt.Errorf("creating directory %s failed: %v", filepath.Dir(localpath), err)
	}
	f, err := os.Create(localpath)
	if err != nil {
		return fmt.Errorf("creating %s failed: %v", localpath, err)
	}
	defer f.Close()

	uri := midiURI + "/" + remotepath
	resp, err := http.Get(uri)
	if err != nil {
		return fmt.Errorf("request to %s failed: %v", uri, err)
	}
	defer resp.Body.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("downloading %s to %s failed: %v", uri, localpath, err)
	}

	return nil
}

func getSongList() (songs map[string]karaoke.Song, err error) {
	uri := midiURI + "/manifest.json"
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

func usageAndExit(message string, exitCode int) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(exitCode)
}
