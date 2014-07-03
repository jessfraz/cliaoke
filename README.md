# cli-aoke

Command Line Karaoke

## Requirements

This assumes you have setup `fluidsynth` in the following way:

```bash
$ brew install fluidsynth
$ wget http://www.schristiancollins.com/soundfonts/GeneralUser_GS_1.44-FluidSynth.zip
$ unzip GeneralUser_GS_1.44-FluidSynth.zip
$ mkdir -p /usr/local/share/fluidsynth
$ mv GeneralUser\ GS\ 1.44\ FluidSynth/GeneralUser\ GS\ FluidSynth\ v1.44.sf2 /usr/local/share/fluidsynth/generaluser.v.1.44.sf2
```

## Installation

```bash
pip install cli-aoke
```

## Usage

**View song choices**

```bash
$ cli-aoke songs
```

**Sing a song**

```bash
$ cli-aoke sing <song_file_name>
```