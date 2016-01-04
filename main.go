package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/homedir"
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
	manifestURI      = midiURI + "/manifest.json"
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
		fmt.Printf("%s\n", VERSION)
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
	songs, err := karaoke.GetSongList(manifestURI)
	if err != nil {
		logrus.Fatal(err)
	}

	if songRequested == "" {
		// print songs table
		w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)

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
		remotemid := midiURI + "/" + song.File
		if err := karaoke.DownloadSong(localmid, remotemid); err != nil {
			logrus.Fatal(err)
		}
	}

	var wg sync.WaitGroup

	// play the song
	wg.Add(1)
	go func(s string) {
		defer wg.Done()
		if err := karaoke.Play(s); err != nil {
			logrus.Error(err)
		}
	}(localmid)

	// show the lyrics
	wg.Add(1)
	go func(l string) {
		defer wg.Done()

		lines := strings.Split(l, "\n")
		for _, line := range lines {
			fmt.Println(line)
			time.Sleep(2 * time.Second)
		}
	}(song.Lyrics)

	// wait
	wg.Wait()
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
