from shared.db import db

# https://likumi.lv/ta/id/239759-geodeziskas-atskaites-sistemas-un-topografisko-karsu-sistemas-noteikumi
class TKS93MapTile(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(20), unique=True) # tile nomenclature
    ulc = db.Column(db.Float) # upper-left corner coordinate in EPSG:4326
    urc = db.Column(db.Float) # upper-right corner coordinate in EPSG:4326
    blc = db.Column(db.Float) # bottom-left corner coordinates in EPSG:4326
    brc = db.Column(db.Float) # bottom-right corner coordinates in EPSG:4326
    tfwURL = db.Column(db.String(200)) # URL for downloading tiff world file
    rgbURL = db.Column(db.String(200)) # URL for downloading visible light images in tiff format
    cirURL = db.Column(db.String(200)) # URL for downloading near infrared images in tiff format
    rgb = db.Column(db.LargeBinary) # for storing visible light images in jpeg
    cir = db.Column(db.LargeBinary) # for storing near infrared images in jpeg
    ndvi = db.Column(db.LargeBinary) # for storing ndvi images in jpeg

    def __init__(self, name, **kwargs):
        self.name = name
        self.ulc = kwargs.get('ulc', None)
        self.urc = kwargs.get('urc', None)
        self.blc = kwargs.get('blc', None)
        self.brc = kwargs.get('brc', None)
        self.tfwURL = kwargs.get('tfwURL', None)
        self.rgbURL = kwargs.get('rgbURL', None)
        self.cirURL = kwargs.get('cirURL', None)
        self.rgb = kwargs.get('rgb', None)
        self.cir = kwargs.get('cir', None)
        self.ndvi = kwargs.get('ndvi', None)
