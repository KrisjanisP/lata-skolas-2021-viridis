from shared.db import db

# https://likumi.lv/ta/id/239759-geodeziskas-atskaites-sistemas-un-topografisko-karsu-sistemas-noteikumi
class TKS93MapTile(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(20), unique=True) # tile nomenclature
    ulcx = db.Column(db.Float) # upper-left corner coordinate x in EPSG:4326
    ulcy = db.Column(db.Float) # upper-left corner coordinate y in EPSG:4326
    urcx = db.Column(db.Float) # upper-right corner coordinate x in EPSG:4326
    urcy = db.Column(db.Float) # upper-right corner coordinate y in EPSG:4326
    blcx = db.Column(db.Float) # bottom-left corner coordinates x in EPSG:4326
    blcy = db.Column(db.Float) # bottom-left corner coordinates y in EPSG:4326
    brcx = db.Column(db.Float) # bottom-right corner coordinates x in EPSG:4326
    brcy = db.Column(db.Float) # bottom-right corner coordinates y in EPSG:4326
    tfwURL = db.Column(db.String(200)) # URL for downloading tiff world file
    rgbURL = db.Column(db.String(200)) # URL for downloading visible light images in tiff format
    cirURL = db.Column(db.String(200)) # URL for downloading near infrared images in tiff format
    rgb = db.Column(db.LargeBinary) # for storing visible light images in jpeg
    cir = db.Column(db.LargeBinary) # for storing near infrared images in jpeg
    ndvi = db.Column(db.LargeBinary) # for storing ndvi images in jpeg

    def __init__(self, name, **kwargs):
        self.name = name
        self.ulcx = kwargs.get('ulcx', None)
        self.ulcy = kwargs.get('ulcy', None)
        self.urcx = kwargs.get('urcx', None)
        self.urcy = kwargs.get('urcy', None)
        self.blcx = kwargs.get('blcx', None)
        self.blcy = kwargs.get('blcy', None)
        self.brcx = kwargs.get('brcx', None)
        self.brcy = kwargs.get('brcy', None)
        self.tfwURL = kwargs.get('tfwURL', None)
        self.rgbURL = kwargs.get('rgbURL', None)
        self.cirURL = kwargs.get('cirURL', None)
        self.rgb = kwargs.get('rgb', None)
        self.cir = kwargs.get('cir', None)
        self.ndvi = kwargs.get('ndvi', None)
