const express = require('express');
const bodyParser = require('body-parser');
const _ = require('lodash');
const uuid = require('uuid/v1');
var app = express();
app.use(bodyParser.json())

var users = [];

function newUser(username, age){
    return {id:uuid(),username:username, age:age};
}

users.push(newUser('test',1));

app.get('/users', function (req, res) {
  console.log('get users');
  res.send(users);
});

app.get('/users/:id', function (req, res) {
  console.log('get user by id');
  var user = null;
  users.forEach(function(item){
    if(item.id===req.params.id){
        user = item;
        res.send(user);
    }
  });
  if(!user){
    res.send('User not found');
  }
});

app.post('/users', function(req, res){
  console.log('post users');
  var user = _.pick(req.body,['username','age']);
  user.id = uuid();
  users.push(user);
  res.send(user);
});

app.listen(process.env.PORT || 8888, function () {
  console.log('Users service listening on port '+ (process.env.PORT || 8888));
});


module.exports.app = app;