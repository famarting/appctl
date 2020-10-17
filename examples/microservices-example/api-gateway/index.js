const express = require('express')  
const httpProxy = require('express-http-proxy')  
const gatewayConfig = require('./api-gateway-config.json');
const app = express()

gatewayConfig.services.forEach(function(service) {
  console.log("mapping service "+service.name);
  const proxy = httpProxy(service.url);
  app.all(service.mapping,(req,res,next)=>{
    proxy(req,res,next);
  })
});

app.listen(gatewayConfig.port || process.env.PORT || 10000, function () {
  console.log('Gateway listening on port '+ (gatewayConfig.port || process.env.PORT || 10000));
});