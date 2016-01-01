// +build linux

package main

import (
	"fmt"
	"os/exec"
	"strings"
)

const soundfontPath = "/usr/share/sounds/sf2/FluidR3_GM.sf2"

func play(localmid string) error {
	cmd := exec.Command("fluidsynth", "-a", "alsa", "-m", "alsa_seq", "-l", "-i", soundfontPath, localmid)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("running `%s %s` failed: %s, %v", cmd.Path, strings.Join(cmd.Args, " "), out, err)
	}

	return nil
}
