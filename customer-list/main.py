from flask import Flask
import psycopg
from waitress import serve
import os
from querier import querier

database_url = os.environ['DATABASE_URL'] \
    if os.environ['DATABASE_URL'] is not None and os.environ['DATABASE_URL'] != '' \
    else 'postgresql://watcher:password@localhost:5432/watcher?sslmode=disable'

app = Flask(__name__)

with psycopg.connect(database_url) as conn:
    conn.read_only = True

    @app.get('/')
    def hello():
        return 'OK'


    @app.get('/customers')
    def customers():
        return querier(conn)

    serve(app, host="0.0.0.0", port=7201)


