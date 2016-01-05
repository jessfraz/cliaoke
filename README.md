# cliaoke

[![Circle CI](https://circleci.com/gh/jfrazelle/cliaoke.svg?style=svg)](https://circleci.com/gh/jfrazelle/cliaoke)

Command Line Karaoke

```console
$ cliaoke -h
      _ _             _
  ___| (_) __ _  ___ | | _____
 / __| | |/ _` |/ _ \| |/ / _ \
| (__| | | (_| | (_) |   <  __/
 \___|_|_|\__,_|\___/|_|\_\___|

 Command Line Karaoke
 Version: v0.1.0

  -d    run in debug mode
  -v    print version and exit (shorthand)
  -version
        print version and exit
```

## Usage

**List all songs**

**NOTE:** This does not mean you have all these files locally, when you choose
a song (if you have no already downloaded it from my s3 bucket) it will be
downloaded for you.

```console
$ cliaoke
COMMAND                             TITLE                               ARTIST
1979                                1979                                Smashing Pumpkins
99_ways_to_die                      99 Ways To Die                      Megadeth
...
hard_knock_life                     Hard Knock Life                     Jay-Z
...
missing_you                         Missing You                         Puff Daddy
mo_money_mo_problems                Mo Money Mo Problems                Notorious BIG
...
```

**Play a song**

```console
$ cliaoke mo_money_mo_problems
DJ cliaoke on the request line.
Bringing up the track Mo Money Mo Problems...
```

## Requirements

### Linux

- Download `fluidsynth` and soundfonts on debian this was `fluid-soundfont-gm`.

**OR use my docker image**

```
$ docker run --rm -it \
    --device /dev/snd \
    jess/cliaoke hard_knock_life
DJ cliaoke on the request line.
Bringing up the track Hard Knock Life...
```


### OSX

This assumes you have setup `fluidsynth` in the following way:

```console
$ brew install fluidsynth
$ wget http://www.schristiancollins.com/soundfonts/GeneralUser_GS_1.44-FluidSynth.zip
$ unzip GeneralUser_GS_1.44-FluidSynth.zip
$ mkdir -p /usr/local/share/fluidsynth
$ mv GeneralUser\ GS\ 1.44\ FluidSynth/GeneralUser\ GS\ FluidSynth\ v1.44.sf2 /usr/local/share/fluidsynth/generaluser.v.1.44.sf2
```

Running `cliaoke` with no arguments will list all the available songs. Once downloaded the songs are saved in a `~/.cliaoke/` directory.

**Caveats for Mac Users**

- You *must* install `fluidsynth` according to the instructions below.
- Sometimes the search for lyrics selects the wrong one, *whomp whomp*.
- The lines being printed are not synced with the song. (maybe there is a way to parse the lyric metadata from a .mid file?)


[![Analytics](https://ga-beacon.appspot.com/UA-29404280-16/cliaoke/README.md)](https://github.com/jfrazelle/cliaoke)
