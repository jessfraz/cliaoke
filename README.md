# cli-aoke

Command Line Karaoke

## How it Works

The first time you run `cli-aoke songs` since you don't have the song files, they will be scraped from the web and placed in a `.cli-aoke/` directory in your home directory.
Running `cli-aoke sing <song_file_name>` finds the song file, searches azlyrics.com for the lyrics based off the filename, scrapes the site for the lyrics, then starts a python thread with the `.mid` file playing through `fluidsynth` & the lyrics being displayed line by line in your terminal.

**Caveats**

- You *must* install `fluidsynth` according to the instructions below.
- Sometimes the search for lyrics selects the wrong one, *whomp whomp*.
- The lines being printed are not synced with the song. (maybe there is a way to parse the lyric metadata from a .mid file?)

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
