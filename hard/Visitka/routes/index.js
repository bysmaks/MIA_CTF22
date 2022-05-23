var express = require('express');
var router = express.Router();

const flag = "CTF{bEst_lay0ut}"

router.get('/', function(req, res, next) {
 	res.render('index')
});

router.post('/', function(req, res, next) {
	var profile = req.body.profile
 	res.render('index', profile)
});

module.exports = router;
