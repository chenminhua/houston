const WebSocket = require('ws');

var ws_url = 'ws://localhost:8087?auth='
for (var i = 1; i < 10; i++) {
  var ws = new WebSocket(ws_url + i)
  var f = function () {
    var m = i
    return function incoming(data) {
      console.log(data + m);
    }
  }()
  ws.on('message', f);
}
