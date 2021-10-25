import os
import sys
import re
import requests
import config

# download rgb aka visible light map urls file

rgbURLsFilePath = os.path.join('temp', config.rgbURLsFileName)

with open(rgbURLsFilePath, 'wb') as f:
    f.write(requests.get(config.rgbURL).content)

# download cir aka near-i.r. map urls file

cirURLsFilePath = os.path.join('temp', config.cirURLsFileName)

with open(cirURLsFilePath, 'wb') as f:
    f.write(requests.get(config.cirURL).content)

# extract TKS-93 map division location identifiers

def getLocIds(filePath):
    res = []
    with open(filePath, encoding='utf8') as file:
        while True:
            line = file.readline()
            if not line:
                break
            found = re.search('v6/(.+?).tfw', line)
            if not found:
                continue
            locId = found.group(1)
            part1 = re.search('(.+?)/', locId).group(1)
            part2 = re.search('/(.+?)-', locId).group(1)
            part3 = re.search('-(.+?)$', locId).group(1)
            # skip broken links
            if part1 != part2:
                continue
            res.append(f'{part2}-{part3}')
    return res

rgbLocIds = getLocIds(rgbURLsFilePath)
cirLocIds = getLocIds(cirURLsFilePath)

if rgbLocIds != cirLocIds:
    sys.exit()

locIds = rgbLocIds
