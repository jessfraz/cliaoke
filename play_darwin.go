// +build darwin

package main

import (
	"fmt"
	"os/exec"
	"strings"
)

const soundfontPath = "/usr/local/share/fluidsynth/generaluser.v.1.44.sf2"

func play(localmid string) error {
	cmd := exec.Command("fluidsynth", "-i", soundfontPath, localmid)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("running `%s %s` failed: %s, %v", cmd.Path, strings.Join(cmd.Args, " "), out, err)
	}

	return nil
}
