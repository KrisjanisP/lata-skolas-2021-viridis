import re
import io
import requests
from pyproj import Proj, transform
from tqdm import tqdm
import config
from models.tile import TKS93MapTile

LKS92 = Proj('epsg:3059')
WGS84 = Proj('epsg:3857')

def get_req_headers():
    user_agent_header = 'PostmanRuntime/7.28.4'
    host_header = 's3.storage.pub.lvdc.gov.lv'
    return {'User-Agent': user_agent_header, 'Host':host_header}

def get_coordinates(tfw_url):
    """
    Download .tfw file.\n
    Read x, y coordinates from the file.\n
    Transform these coordinates into epsg:3857.
    """

    tfw = requests.get(tfw_url, headers=get_req_headers()).content.decode("utf-8")

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

def process_rgb_urls(db):
    tile_count = config.TILE_COUNT
    # download rgb aka visible light map urls file
    rgb_urls = requests.get(config.RGB_URL).content.decode("utf-8")

    print('Creating tiles from rgbURLsFile.')
    with io.StringIO(rgb_urls) as file:
        for i in tqdm(range(tile_count)):
            if completed:
                continue
            tfw_url = file.readline().strip()
            rgb_url = file.readline().strip()
            if not tfw_url or not rgb_url:
                break

            # skip weird links
            if re.search('.xml', tfw_url) is not None:
                tfw_url = rgb_url
                rgb_url = file.readline()

            # last links are broken
            if re.search('.tif', tfw_url) is not None:
                break

            temp = re.search('v6/(.+?).tfw', tfw_url)
            if temp is None:
                print(tfw_url)

            part1 = re.search('/(.+?)-', temp.group(1)).group(1)
            part2 = re.search('-(.+?)$', temp.group(1)).group(1)

            tile_name = f'{part1}-{part2}'

            # if it exists and is fully completed, skip it
            tile = TKS93MapTile.query.filter_by(name=tile_name).first()
            if tile is not None:
                if (tile.brcy is not None and
                    tile.tfwURL is not None and
                    tile.rgbURL is not None):
                    completed = True
                    continue

            #ulc, urc, blc, brc = getCoordinates(tfwURL)

            ulcx, ulcy, urcx, urcy, brcx, brcy, blcx, blcy = get_coordinates(tfw_url)

            # create a new tile
            new_tile = TKS93MapTile(
                name=tile_name, tfwURL=tfw_url, rgbURL=rgb_url,
                ulcx=ulcx,ulcy=ulcy,urcx=urcx,urcy=urcy,
                blcx=blcx,blcy=blcy,brcx=brcx,brcy=brcy)

            # if it exists, update it
            tile = TKS93MapTile.query.filter_by(name=tile_name).first()
            if tile is not None:
                tile.name = new_tile.name
                tile.tfwURL = new_tile.tfwURL
                tile.rgbURL = new_tile.rgbURL
                tile.ulcx = new_tile.ulcx
                tile.ulcy = new_tile.ulcy
                tile.urcx = new_tile.urcx
                tile.urcy = new_tile.urcy
                tile.blcx = new_tile.blcx
                tile.blcy = new_tile.blcy
                tile.brcx = new_tile.brcx
                tile.brcy = new_tile.brcy
            else:
                db.session.add(new_tile)

    db.session.commit()

def process_cir_urls(db):
    tile_count = config.TILE_COUNT
    cir_urls = requests.get(config.CIR_URL).content.decode("utf-8")

    print('Updating tiles from cirURLsFile.')
    completed = False
    with io.StringIO(cir_urls) as file:
        for i in tqdm(range(tile_count)):
            if completed:
                continue
            tfw_url = file.readline().strip()
            cir_url = file.readline().strip()
            if not tfw_url or not cir_url:
                break

            # skip weird links
            if re.search('.xml', tfw_url) is not None:
                tfw_url = cir_url
                cir_url = file.readline()

            # last links are broken
            if re.search('.tif', tfw_url) is not None:
                break

            temp = re.search('v6/(.+?).tfw', tfw_url)
            if temp is None:
                print(tfw_url)

            part1 = re.search('/(.+?)-', temp.group(1)).group(1)
            part2 = re.search('-(.+?)$', temp.group(1)).group(1)

            tile_name = f'{part1}-{part2}'

            # if it exists, update it
            tile = TKS93MapTile.query.filter_by(name=tile_name).first()
            tile.cirURL = cir_url

    db.session.commit()

def defaultfill(db):
    process_rgb_urls(db)
    process_cir_urls(db)
