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
$ pip install cli-aoke
```

## Usage

**View song choices**

**NOTE**: You must run first  `cli-aoke songs` once before running `sing` because it scrapes a site for the `.mid` files to initialize your songs directory

```bash
$ cli-aoke songs
```

**Sing a song**

```bash
$ cli-aoke sing <song_file_name>
```

## Example

```bash
$ cli-aoke songs
$ cli-aoke sing Jay-Z_-_Hard_Knock_Life.mid
$ cli-aoke sing Blackstreet_-_No_Diggity.mid
$ cli-aoke sing 2Pac_-_California.mid
```
