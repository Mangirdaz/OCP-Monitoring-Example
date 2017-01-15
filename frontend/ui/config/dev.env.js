var merge = require('webpack-merge')
var prodEnv = require('./prod.env')

module.exports = merge(prodEnv, {
  NODE_ENV: '"development"',
  API_URL: '"http://api-svc-myproject.192.168.1.108.xip.io/api/v1/"',
  IMG_URL: '"http://img-svc-myproject.192.168.1.108.xip.io/api/v1/img"'
})
