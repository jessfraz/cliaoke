#!/usr/bin/env python

import re
from bs4 import BeautifulSoup
import urllib2


def fetch_lyrics(url):
    try:
        response = urllib2.urlopen(url)
    except urllib2.URLError, e:
        if hasattr(e, 'reason'):
            print 'We failed to reach a server.'
            print 'Reason: ', e.reason
        elif hasattr(e, 'code'):
            print 'The server couldn\'t fulfill the request.'
            print 'Error code: ', e.code
    else:
        # go on with your life
        html = response.read()

        lyrics = re.search(
            b'<!-- start of lyrics -->(?:\r\n)+(.+)(?:\r\n)+<!-- end of lyrics -->', html, re.DOTALL)
        if lyrics:
            # Strip html tags from decoded lyrics
            return re.sub(r'<.+>', '', lyrics.group(1).decode('utf8'))
        else:
            return None


def fetch_search_top_link(song_query):
    # get the url
    uri = 'http://search.azlyrics.com/search.php?q=' + song_query

    try:
        response = urllib2.urlopen(uri)
    except urllib2.URLError, e:
        if hasattr(e, 'reason'):
            print 'We failed to reach a server.'
            print 'Reason: ', e.reason
        elif hasattr(e, 'code'):
            print 'The server couldn\'t fulfill the request.'
            print 'Error code: ', e.code
    else:
        # go on with your life
        html = response.read()

        # pass the html to BeautifulSoup
        soup = BeautifulSoup(html)
        results = soup.find(id="inn")
        if results is not None:
            atag = results.find('a')
            if atag is not None:
                return atag.get('href')

    return None


def clean_input(song_file):
    song = song_file.replace('.mid', '').replace('.kar', '')
    song = re.sub('[^a-zA-Z0-9\n\.]', '+', song)
    return song


def get(song_file):
    # clean it
    song_query = clean_input(song_file)

    # search for the lyrics
    link = fetch_search_top_link(song_query)
    if link is not None:
        return fetch_lyrics(link)
    else:
        return None
