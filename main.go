package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/Sirupsen/logrus"
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

	midiURI = "https://s3.j3ss.co/cliaoke/midi"
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
	if songRequested == "" {
		// list all songs
		songs, err := getSongList()
		if err != nil {
			logrus.Fatal(err)
		}

		// print songs table
		_, stdout, _ := term.StdStreams()
		w := tabwriter.NewWriter(stdout, 20, 1, 3, ' ', 0)

		// print header
		fmt.Fprintln(w, "TITLE\tARTIST")
		for _, song := range songs {
			fmt.Fprintf(w, "%s\t%s\n", song.Title, song.Artist)
		}
		w.Flush()
		return
	}

	// play requested song
	fmt.Printf("song requested was %s", songRequested)
}

func getSongList() (songs []karaoke.Song, err error) {
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
