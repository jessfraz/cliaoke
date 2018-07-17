package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/genuinetools/pkg/cli"
	"github.com/jessfraz/cliaoke/karaoke"
	"github.com/jessfraz/cliaoke/version"
	"github.com/sirupsen/logrus"
)

const (
	defaultSongStore = ".cliaoke"
	midiURI          = "https://raw.githubusercontent.com/jessfraz/cliaoke/master/midi"
	manifestURI      = midiURI + "/manifest.json"
)

var (
	debug bool
)

//go:generate go run midi/generate.go

func main() {
	// Create a new cli program.
	p := cli.NewProgram()
	p.Name = "cliaoke"
	p.Description = "Command line karaoke"

	// Set the GitCommit and Version.
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	// Setup the global flags.
	p.FlagSet = flag.NewFlagSet("global", flag.ExitOnError)
	p.FlagSet.BoolVar(&debug, "d", false, "enable debug logging")

	// Set the before function.
	p.Before = func(ctx context.Context) error {
		// Set the log level.
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		return nil
	}

	// Set the main program action.
	p.Action = func(ctx context.Context, args []string) error {
		// On ^C, or SIGTERM handle exit.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, syscall.SIGTERM)
		go func() {
			for sig := range c {
				logrus.Infof("Received %s, exiting.", sig.String())
				os.Exit(0)
			}
		}()

		// get all the songs
		songs, err := karaoke.GetSongList(manifestURI)
		if err != nil {
			logrus.Fatal(err)
		}

		// If they did not request a song then print the list of songs.
		if len(args) < 1 {
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
			return nil
		}

		// play requested song
		fmt.Println("DJ cliaoke on the request line.")

		// Find the song
		// If we only got one arg, assume they sent the command for the song
		songRequested := args[0]
		if len(args) > 1 {
			songRequested = strings.ToLower(strings.Join(args, "_"))
		}
		song, ok := songs[songRequested]
		if !ok {
			logrus.Fatalf("Could not find %s in song list, run with no arguments to see the list of songs available.", songRequested)
		}

		// download the song if we can't find it locally
		fmt.Printf("Bringing up the track %s...\n", song.Title)

		// Make sure we have fluidsynth installed.
		if !karaoke.FluidsynthBinaryExists() {
			return errors.New("fluidsynth is not installed and is required")
		}

		home, err := getHome()
		if err != nil {
			logrus.Fatal(err)
		}

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
		return nil
	}

	// Run our program.
	p.Run()
}

func getHome() (string, error) {
	home := os.Getenv(homeKey)
	if home != "" {
		return home, nil
	}

	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return u.HomeDir, nil
}
