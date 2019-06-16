from http.server import BaseHTTPRequestHandler, HTTPServer
from urllib.parse import urlparse
import json

class Server(BaseHTTPRequestHandler):
    def _set_headers(self):
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()

    def do_GET(self):
        parsed_path = urlparse(self.path)
        if parsed_path.path != '/scores':
            self.send_error(404, "not found")
            return

        try:
            params = dict([p.split('=') for p in parsed_path[4].split('&')])
        except:
            params = {}

        if 'ids' not in params:
            self.send_error(400, "parameter ids is required")
            return

        print('getting scores of ids:', params["ids"])
        ids = params["ids"].split(",")
        scores = {}
        for id in ids:
            try:
                id_num = int(id)
            except ValueError:
                continue # ignore
            score = get_score(id_num)
            if score != '':
                scores[id_num] = score

        self._set_headers()
        self.wfile.write(json.dumps(scores).encode())

def get_score(id):
    mock_DB = {
        1: 9.8,
        2: 10,
        3: 9.7,
        4: 9.5,
        5: 9.9,
        6: 10,
        7: 9.8,
        8: 9.3,
        9: 9.4,
        10: 10,
        11: 8.9,
        12: 9.2,
        13: 9.6,
        14: 8.2,
    }

    if id in mock_DB:
        #return mock_DB[id] # v1
        return mock_DB[id] * 10 #v2
    return ''

def run(server_class=HTTPServer, handler_class=Server, port=7000):
    server_address = ('', port)
    httpd = server_class(server_address, handler_class)

    print('Starting scores service on port %d...' % port)
    httpd.serve_forever()

if __name__ == "__main__":
    run()
