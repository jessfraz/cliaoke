import threading
import time
from . import songs


class StoppableThread(threading.Thread):

    """Thread class with a stop() method. The thread itself has to check
    regularly for the stopped() condition."""

    def __init__(self, target):
        super(StoppableThread, self).__init__(target=target)
        self._stop = threading.Event()

    def stop(self):
        self._stop.set()

    def stopped(self):
        return self._stop.isSet()


class Karaoke(object):
    # start threads

    def __init__(self, file, lyrics):
        self.file = file
        self.lyrics = lyrics

        self.song_thread = StoppableThread(self.PlaySong)
        self.lyrics_thread = StoppableThread(self.ServeLyrics)

        self.song_thread.start()
        self.lyrics_thread.start()

    # setup song thread function
    def PlaySong(self):
        songs.play(self.file)

    # setup lyric thread function
    def ServeLyrics(self):
        lines = self.lyrics.splitlines()
        while (self.lyrics_thread.stopped() is False and len(lines) > 0):
            print lines[0]
            del lines[0]
            time.sleep(1)
