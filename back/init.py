import os
import sys
import re
import requests
import config
import shutil
import io
import json
from models.tile import TKS93MapTile
from time import sleep
from tqdm import tqdm
from pyproj import Proj, transform

LKS92 = Proj('epsg:3059')
WGS84 = Proj('epsg:3857')

def getCoordinates(tfwURL):
    tfw = requests.get(tfwURL, headers={'User-Agent': 'PostmanRuntime/7.28.4', 'Host':'s3.storage.pub.lvdc.gov.lv', 'Postman-Token': '34dfa909-0656-45e1-b11c-717cc67960ec'}).content.decode("utf-8")
    buf = io.StringIO(tfw)
    width = float(buf.readline().strip())*10000
    buf.readline() # empty rotation
    buf.readline() # empty rotation
    height = float(buf.readline().strip())*10000
    y = float(buf.readline().strip())
    x = float(buf.readline().strip())
    ulcx, ulcy = x, y
    ulcx, ulcy = transform(LKS92,WGS84, ulcx, ulcy)
    urcx, urcy = x+width, y
    urcx, urcy = transform(LKS92,WGS84, urcx, urcy)
    blcx, blcy = x, y+height
    blcx, blcy = transform(LKS92,WGS84, blcx, blcy)
    brcx, brcy = x+width, y+height
    brcx, brcy = transform(LKS92,WGS84, brcx, brcy)
    return ulcx, ulcy, urcx, urcy, brcx, brcy, blcx, blcy

def defaultfill(db):
    # download rgb aka visible light map urls file
    rgbURLs = requests.get(config.rgbURL).content.decode("utf-8")

    print('Creating tiles from rgbURLsFile.')

    tile_count = 10775
    
    with io.StringIO(rgbURLs) as file:
        for i in tqdm(range(tile_count)):
            tfwURL = file.readline().strip()
            rgbURL = file.readline().strip()
            if not tfwURL or not rgbURL:
                break

            # skip weird links
            if re.search('.xml', tfwURL) is not None:
                tfwURL = rgbURL
                rgbURL = file.readline()

            # last links are broken
            if re.search('.tif', tfwURL) is not None:
                break

            temp = re.search('v6/(.+?).tfw', tfwURL)
            if temp is None:
                print(tfwURL)
                
            part1 = re.search('/(.+?)-', temp.group(1)).group(1)
            part2 = re.search('-(.+?)$', temp.group(1)).group(1)

            tile_name = f'{part1}-{part2}'

            # if it exists and is fully completed, skip it
            tile = TKS93MapTile.query.filter_by(name=tile_name).first()
            if tile is not None:
                if (tile.brcy is not None and
                    tile.tfwURL is not None and
                    tile.rgbURL is not None):
                    continue

            #ulc, urc, blc, brc = getCoordinates(tfwURL)
            
            ulcx, ulcy, urcx, urcy, brcx, brcy, blcx, blcy = getCoordinates(tfwURL)

            # create a new tile
            newTile = TKS93MapTile(
                name=tile_name, tfwURL=tfwURL, rgbURL=rgbURL,
                ulcx=ulcx,ulcy=ulcy,urcx=urcx,urcy=urcy,
                blcx=blcx,blcy=blcy,brcx=brcx,brcy=brcy)

            # if it exists, update it
            tile = TKS93MapTile.query.filter_by(name=tile_name).first()
            if tile is not None:
                tile.name = newTile.name
                tile.tfwURL = newTile.tfwURL
                tile.rgbURL = newTile.rgbURL
                tile.ulcx = newTile.ulcx
                tile.ulcy = newTile.ulcy
                tile.urcx = newTile.urcx
                tile.urcy = newTile.urcy
                tile.blcx = newTile.blcx
                tile.blcy = newTile.blcy
                tile.brcx = newTile.brcx
                tile.brcy = newTile.brcy
            else:
                db.session.add(newTile)

    db.session.commit()

    # download rgb aka visible light map urls file
    cirURLs = requests.get(config.cirURL).content.decode("utf-8")

    print('Updating tiles from cirURLsFile.')

    with io.StringIO(cirURLs) as file:
        for i in tqdm(range(tile_count)):
            tfwURL = file.readline().strip()
            cirURL = file.readline().strip()
            if not tfwURL or not cirURL:
                break

            # skip weird links
            if re.search('.xml', tfwURL) is not None:
                tfwURL = cirURL
                cirURL = file.readline()

            # last links are broken
            if re.search('.tif', tfwURL) is not None:
                break

            temp = re.search('v6/(.+?).tfw', tfwURL)
            if temp is None:
                print(tfwURL)
                
            part1 = re.search('/(.+?)-', temp.group(1)).group(1)
            part2 = re.search('-(.+?)$', temp.group(1)).group(1)

            tile_name = f'{part1}-{part2}'

            # if it exists, update it
            tile = TKS93MapTile.query.filter_by(name=tile_name).first()
            tile.cirURL = cirURL

    db.session.commit()