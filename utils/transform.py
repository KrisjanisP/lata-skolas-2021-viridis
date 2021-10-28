import sys
from pyproj import Transformer

x = float(sys.argv[1])
y = float(sys.argv[2])
width = float(sys.argv[3])
height = float(sys.argv[4])

transformer = Transformer.from_crs('epsg:3059','epsg:3857')
ulcx, ulcy = transformer.transform(x, y)
urcx, urcy = transformer.transform(x-width, y)
blcx, blcy = transformer.transform(x, y-height)
brcx, brcy = transformer.transform(x-width, y-height)

print(ulcx, ulcy, urcx, urcy, brcx, brcy, blcx, blcy, sep="\n")