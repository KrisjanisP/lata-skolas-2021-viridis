import os
from flask import Flask, jsonify
from flask_marshmallow import Marshmallow
from init import defaultfill
from models.tile import TKS93MapTile
from geojson import Point, Feature, FeatureCollection, Polygon
from flask_cors import CORS, cross_origin
from flask_headers import headers
# init app
app = Flask(__name__)

cors = CORS(app)
app.config['CORS_HEADERS'] = 'Content-Type'

basedir = os.path.abspath(os.path.dirname(__file__))

# database config
app.config['SQLALCHEMY_DATABASE_URI'] = 'sqlite:///'+os.path.join(basedir, 'db.sqlite')
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False

# init db
from shared.db import db
db.init_app(app)
with app.app_context():
    db.create_all()
    defaultfill(db)

# init ma
ma = Marshmallow(app)

# tile schema
class TilesSchema(ma.Schema):
    class Meta:
        fields = ('name', 'ulcx', 'ulcy', 'urcx', 'urcy', 'blcx', 'blcy', 'brcx', 'brcy')

# init schema
tile_schema = TilesSchema()
tiles_schema = TilesSchema(many=True)

# Get all tiles
@app.route('/tks93tiles', methods=['GET'])
@cross_origin()
@headers({'content-type':'application/geo+json'})
def get_tiles():
    tiles = TKS93MapTile.query.all()
    features = []
    for tile in tiles:
        if tile.ulcx is None:
            continue
        point = Point((tile.ulcx, tile.ulcy))
        polygon = Polygon([[(tile.ulcx, tile.ulcy),(tile.urcx, tile.urcy),(tile.brcx, tile.brcy),(tile.blcx, tile.blcy),(tile.ulcx, tile.ulcy)]])
        feature = Feature(geometry=polygon, id=tile.name)
        features.append(feature)
    return FeatureCollection(features)

app.run()