package lyrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	musicMatchAPIVersion = "1.1"
	musicMatchURI        = "http://api.musixmatch.com/ws/" + musicMatchAPIVersion
)

type Client struct {
	Token string
}

type Header struct {
	StatusCode  int     `json:"status_code"`
	ExecuteTime float64 `json:"execute_time"`
}

type Response struct {
	Message Message `json:"message"`
}

type Message struct {
	Header Header `json:"header"`
	Body   Body   `json:"body"`
}

type Body struct {
	TrackList []map[string]Track `json:"track_list"`
	Lyrics    Lyrics             `json:"lyrics"`
}

type Track struct {
	ID   int64  `json:"track_id"`
	MBID string `json:"track_mbid"`
	Name string `json:"track_name"`
	//HasLyrics    bool      `json:"has_lyrics"`
	//HasSubtitles bool      `json:"has_subtitles"`
	LyricsID    int64  `json:"lyrics_id"`
	SubtitlesID int64  `json:"subtitles_id"`
	AlbumID     int64  `json:"album_id"`
	AlbumName   string `json:"album_name"`
	ArtistID    int64  `json:"artist_id"`
	ArtistMBID  string `json:"artist_mbid"`
	ArtistName  string `json:"artist_name"`
	//UpdatedTime time.Time `json:"updated_time"`
}

type LyricsResponse struct {
	message struct {
		Header Header
		body   struct {
			Lyrics Lyrics `json:"lyrics"`
		}
	}
}

type Lyrics struct {
	ID          int64     `json:"lyrics_id"`
	Body        string    `json:"lyrics_body"`
	Language    string    `json:"lyrics_language"`
	Copyright   string    `json:"lyrics_copyright"`
	UpdatedTime time.Time `json:"updated_time"`
}

func (c *Client) SearchTrack(query string) (track Track, err error) {
	v := url.Values{
		"apikey": []string{c.Token},
		"q":      []string{query},
	}
	uri := fmt.Sprintf("%s/track.search?%s", musicMatchURI, v.Encode())
	resp, err := http.Get(uri)
	if err != nil {
		return track, fmt.Errorf("request to %s failed: %v", uri, err)
	}
	defer resp.Body.Close()

	// decode the body
	var r Response
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&r); err != nil {
		return track, fmt.Errorf("decoding response from %s failed: %v", uri, err)
	}

	if len(r.Message.Body.TrackList) <= 0 {
		return track, fmt.Errorf("could not find any tracks matching %s", query)
	}

	return r.Message.Body.TrackList[0]["track"], nil
}

func (c *Client) GetTrackLyrics(track Track) (string, error) {
	v := url.Values{
		"apikey":     []string{c.Token},
		"track_id":   []string{strconv.FormatInt(track.ID, 10)},
		"track_mbid": []string{track.MBID},
	}
	uri := fmt.Sprintf("%s/track.lyrics.get?%s", musicMatchURI, v.Encode())
	resp, err := http.Get(uri)
	if err != nil {
		return "", fmt.Errorf("request to %s failed: %v", uri, err)
	}
	defer resp.Body.Close()

	// decode the body
	var r Response
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&r); err != nil {
		return "", fmt.Errorf("decoding response from %s failed: %v", uri, err)
	}

	return r.Message.Body.Lyrics.Body, nil
}

/*
#!/usr/bin/env python

import re
from bs4 import BeautifulSoup
from os import path
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
        search_term = ("<!-- Usage of azlyrics.com content by "
                       "any third-party lyrics provider is prohibited "
                       "by our licensing agreement. Sorry about that. -->")
        lyrics = re.search(
            b'%s(?:\r\n)+(.+)(?:\r\n)+</div>' % search_term, html, re.DOTALL)
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
        results = soup.find('td')
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
    # first checking lyrics folder
    home = path.expanduser("~")
    lyrics_file = path.join(home, ".cliaoke/lyrics/" + song_file.replace('.mid','.txt'))
    if path.exists(lyrics_file):
        print 'Lyrics already downloaded! Getting them from the lyrics folder.'
        f = open(lyrics_file, 'r')
        return f.read()

    print 'Lyrics have not been downloaded yet. Downloading lyrics.'
    # clean it
    song_query = clean_input(song_file)

    # search for the lyrics
    link = fetch_search_top_link(song_query)
    if link is not None:
        # storing lyrics for later use
        song_lyrics = fetch_lyrics(link)
        save(lyrics_file, song_lyrics)
        return song_lyrics
    else:
        return None

def save(lyrics_file, song_lyrics):
    print 'Saving lyrics to be accessible later'
    with open(lyrics_file, 'w') as f :
        f.write(song_lyrics)
*/
