var express = require('express');
var bodyParser = require('body-parser')
var app = express();
app.use(bodyParser.json())

var fruits = [];

function addFruit(name, color){
    return {name:name, color:color};
}

fruits.push(addFruit('orange','orange'));
fruits.push(addFruit('pear','green'));

app.get('/fruits', function (req, res) {
  console.log('get fruits')
  res.send(fruits);
});

app.listen(process.env.PORT || 9999, function () {
  console.log('fruits service listening on port '+ (process.env.PORT || 9999));
});

