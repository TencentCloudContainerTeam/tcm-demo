const fs = require('fs');
const util = require('util');
const request = require('request');
const requestPromise = util.promisify(request);
const router = require('koa-router')();

router.get('/', async (ctx, next) => {
  'use strict';
  ctx.redirect('/index');
});

router.get('/index', async (ctx, next) => {
  'use strict';
  ctx.response.type = 'html';
  ctx.response.body = fs.createReadStream('./static/dist/index.html');
});

router.get('/login', async (ctx, next) => {
  'use strict';
  ctx.response.type = 'html';
  ctx.response.body = fs.createReadStream('./static/dist/login.html');
});

router.get('/api/users', async (ctx, next) => {
  let headers = getForwardHeaders(ctx.request);
  let userName = headers.cookie.user || '';
  let userOptions = {
    uri: 'http://users.base.svc.cluster.local:7000/users?name=' + userName,
    method: 'GET',
    headers: headers,
    json: true
  };
  let user = await requestPromise(userOptions);

  ctx.body = user.body;
});

router.post('/api/users', async (ctx, next) => {
  let headers = getForwardHeaders(ctx.request);

  headers = Object.assign({}, { 'content-type': 'application/json' }, headers);

  let userOptions = {
    uri: 'http://users.base.svc.cluster.local:7000/users',
    method: 'POST',
    data: JSON.stringify(ctx.request.body),
    headers: headers,
    json: true
  };
  let user = await requestPromise(userOptions);

  ctx.body = {
    code: 0,
    data: user.body
  };
});

router.get('/api/products', async (ctx, next) => {
  let headers = getForwardHeaders(ctx.request);
  let discountOptions = {
    uri: 'http://discount.base.svc.cluster.local:7000/discount',
    method: 'GET',
    headers: headers,
    json: true
  };
  let discount = await requestPromise(discountOptions);

  ctx.body = discount.body;
});

router.get('/api/recommends', async (ctx, next) => {
  let headers = getForwardHeaders(ctx.request);

  let recommendOptions = {
    uri: 'http://recommend.base.svc.cluster.local:7000/recommend',
    method: 'GET',
    headers: headers,
    json: true
  };
  let recommend = await requestPromise(recommendOptions);
  ctx.body = recommend.body;
});

router.all('*', async (ctx, next) => {
  'use strict';
  ctx.response.type = 'html';
  ctx.response.body = fs.createReadStream('./static/dist/404.html');
});

function getForwardHeaders(request) {
  headers = {};
  forwardHeaders = [
    'cookie',
    'x-request-id',
    'x-b3-traceid',
    'x-b3-spanid',
    'x-b3-parentspanid',
    'x-b3-sampled',
    'x-b3-flags',
    'x-ot-span-context'
  ];

  for (let i = 0, len = forwardHeaders.length; i < len; i++) {
    let key = forwardHeaders[i];
    if (request.headers[key]) {
      headers[key] = request.headers[key];
    }
  }

  return headers;
}

module.exports = router;
