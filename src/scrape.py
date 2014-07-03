#!/usr/bin/env python

import urllib2
from bs4 import BeautifulSoup
from os import path, makedirs


def download_file(cliaoke_dir, url):
    file_name = url.split('/')[-1]
    u = urllib2.urlopen(url)
    f = open(path.join(cliaoke_dir, file_name), 'wb')
    meta = u.info()
    file_size = int(meta.getheaders("Content-Length")[0])
    print "Downloading: %s Bytes: %s" % (file_name, file_size)

    file_size_dl = 0
    block_sz = 8192
    while True:
        buffer = u.read(block_sz)
        if not buffer:
            break

        file_size_dl += len(buffer)
        f.write(buffer)
        status = r"%10d  [%3.2f%%]" % (
            file_size_dl, file_size_dl * 100. / file_size)
        status = status + chr(8) * (len(status) + 1)
        print status,

    f.close()


def do(cliaoke_dir):

    # get the url
    base_uri = 'http://www.albinoblacksheep.com/audio/midi/'
    response = urllib2.urlopen(base_uri)
    html = response.read()

    # pass the html to BeautifulSoup
    soup = BeautifulSoup(html)

    # find all the option tags
    for option in soup.find_all('option'):
        url_slug = option.get('value')
        if url_slug is not None:
            # get the slug's html
            response = urllib2.urlopen(base_uri + url_slug)
            html = response.read()
            slug_soup = BeautifulSoup(html)
            # find the embed
            embed = slug_soup.embed
            if embed is not None:
                # get the file url
                file_url = embed.get('src')
                if file_url is not None:
                    download_file(cliaoke_dir, file_url)
