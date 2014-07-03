#!/usr/bin/env python

from os import walk, path


def get():
    f = []
    current_path = path.abspath(path.dirname(__file__))
    song_files_path = path.join(current_path, "..", "mid_files")
    for (dirpath, dirnames, filenames) in walk(song_files_path):
        print f.extend(filenames)
        break
