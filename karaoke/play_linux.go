// +build linux

package karaoke

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/cemmanouilidis/go.platform"
)

var soundfontPaths = map[string]string{
	"arch":   "/usr/share/soundfonts/FluidR3_GM2-2.sf2",
	"debian": "/usr/share/sounds/sf2/FluidR3_GM.sf2",
}

// Play takes a given .midi file and plays it.
func Play(localmid string) error {
	// detect linux distribution
	dist, _, _, err := platform.LinuxDistribution()
	if err != nil {
		return fmt.Errorf("detecting linux distribution failed: ", err)
	}

	// set soundfontPath depending on linux distribution
	soundfontPath := soundfontPaths["debian"] //default
	if _, ok := soundfontPaths[dist]; ok {
		soundfontPath = soundfontPaths[dist]
	}

	cmd := exec.Command("fluidsynth", "-a", "alsa", "-m", "alsa_seq", "-l", "-i", soundfontPath, localmid)

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("running `%s %s` failed: %s, %v", cmd.Path, strings.Join(cmd.Args, " "), out, err)
	}

	return nil
}
