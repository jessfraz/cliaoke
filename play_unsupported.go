// +build !linux,!darwin

package main

import "fmt"

const soundfontPath = ""

func play(localmid string) error {
	return fmt.Errorf("Sorry, cliaoke is only supported on linux and darwin, consider making a PR for what I can only imagine is openbsd or windows. Clippy says DO IT!")
}
