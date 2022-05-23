var express = require('express');
var app = express();

app.get('/', function (req, res) {
  res.send('flag is /tmp/flag.txt<br>Hello ' + eval(req.query.q));
  console.log(req.query.q);
});

app.listen(1337, function () {
  console.log('Example app listening on port 1337!');
});
