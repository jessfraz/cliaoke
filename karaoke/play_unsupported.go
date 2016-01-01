// +build !linux,!darwin

package karaoke

import "fmt"

const soundfontPath = ""

// Play takes a given .midi file and plays it.
func Play(localmid string) error {
	return fmt.Errorf("Sorry, cliaoke is only supported on linux and darwin, consider making a PR for what I can only imagine is openbsd or windows. Clippy says DO IT!")
}
