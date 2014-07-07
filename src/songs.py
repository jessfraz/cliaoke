#!/usr/bin/env python

from os import walk, path, makedirs
from subprocess import Popen, PIPE
from . import scrape


def exists(song_file):
    home = path.expanduser("~")
    return path.isfile(path.join(home, ".cliaoke", song_file))


def get():
    home = path.expanduser("~")
    cliaoke_dir = path.join(home, ".cliaoke")

    # create the dir if it doesn't exist
    if not path.exists(cliaoke_dir):
        makedirs(cliaoke_dir)
        # scrapre for the songs
        print "Getting .mid files"
        scrape.do(cliaoke_dir)

    # print the songs
    for (dirpath, dirnames, filenames) in walk(cliaoke_dir):
        print "\n".join(filenames)


def play(song_file):
    sound_font_file = '/usr/local/share/fluidsynth/generaluser.v.1.44.sf2'
    home = path.expanduser("~")
    cliaoke_dir = path.join(home, ".cliaoke")
    song_path = path.join(cliaoke_dir, song_file)

    if path.isfile(sound_font_file) is False:
        print "You have not installed fluidsynth correctly"
        print "Please refer to https://github.com/jfrazelle/cli-aoke"
        return None
    elif path.isfile(song_path) is False:
        print "%s does not exist" % song_file
        return None
    else:
        fs = Popen(["fluidsynth", "-i", sound_font_file, song_path],
                   stdout=PIPE, stderr=PIPE)
        stderr = fs.communicate()[1]
        if stderr is not None:
            print stderr
