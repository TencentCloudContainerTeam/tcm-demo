from http.server import BaseHTTPRequestHandler, HTTPServer
from urllib.parse import urlparse
import json
import cgi

class Server(BaseHTTPRequestHandler):
    def _set_headers(self):
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()

    # GET sends back a Hello world message
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
        1: 98,
        2: 100,
        3: 97,
        4: 95,
        5: 99,
        6: 100,
        7: 98,
        8: 93,
        9: 94,
        10: 100,
        11: 89,
        12: 92,
        13: 96,
        14: 99,
        15: 98,
        16: 100,
        17: 91,
        18: 99,
        19: 98,
        20: 100
    }

    if id in mock_DB:
        return '%d%%' % mock_DB[id]
    return ''

def run(server_class=HTTPServer, handler_class=Server, port=9000):
    server_address = ('', port)
    httpd = server_class(server_address, handler_class)

    print('Starting scores service on port %d...' % port)
    httpd.serve_forever()

if __name__ == "__main__":
    from sys import argv

    if len(argv) == 2:
        run(port=int(argv[1]))
    else:
        run()
