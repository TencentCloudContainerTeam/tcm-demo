const fs = require('fs');
const util = require('util');
const request = require('request');
const requestPromise = util.promisify(request);
const router = require('koa-router')();

// const userURL = 'http://150.109.18.177:30988/users';
// const discountURL = 'http://150.109.18.177:31474/discount?ids=1,2,3';
// const recommendURL = 'http://150.109.18.177:31139/recommend?ids=1,2,3';

const userURL = 'http://users.base.svc.cluster.local:7000/users';
const discountURL = 'http://discount.base.svc.cluster.local:7000/discount';
const recommendURL = 'http://recommend.base.svc.cluster.local:7000/recommend';

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

router.get('/api/mall', async (ctx, next) => {
  let headers = getForwardHeaders(ctx.request);

  let userName = getCookieKey('user', headers.cookie);
  let userUri = userURL;
  if (userName) {
    userUri += '?name=' + userName;
  }
  console.log("getting mall data")
  let userOptions = {
    uri: userUri,
    method: 'GET',
    headers: headers,
    json: true
  };

  let discountOptions = {
    uri: discountURL,
    method: 'GET',
    headers: headers,
    json: true
  };

  let recommendOptions = {
    uri: recommendURL,
    method: 'GET',
    headers: headers,
    json: true
  };

  let [user, discount, recommend] = await Promise.all([
    requestPromise(userOptions),
    requestPromise(discountOptions),
    requestPromise(recommendOptions)
  ]);

  ctx.body = {
    user: user.body,
    discount: discount.body,
    recommend: recommend.body
  };
});

router.get('/api/users', async (ctx, next) => {
  let headers = getForwardHeaders(ctx.request);
  let userName = getCookieKey('user', headers.cookie);
  let userUri = userURL;
  if (userName) {
    userUri += '?name=' + userName;
  }
  let userOptions = {
    uri: userUri,
    method: 'GET',
    headers: headers,
    json: true
  };
  let user = await requestPromise(userOptions);

  ctx.body = user.body;
});

router.post('/api/users', async (ctx, next) => {
  let headers = getForwardHeaders(ctx.request);

  headers = Object.assign({}, { 'Content-Type': 'application/json' }, headers);

  let userOptions = {
    uri: userURL,
    method: 'POST',
    json: ctx.request.body,
    headers: headers
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
    uri: discountURL,
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
    uri: recommendURL,
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

function getCookieKey(key, cookie) {
  if (cookie && cookie.length > 0) {
    c_start = cookie.indexOf(key + '=');
    if (c_start != -1) {
      c_start = c_start + key.length + 1;
      c_end = cookie.indexOf(';', c_start);
      if (c_end == -1) {
        c_end = cookie.length;
      }
      return unescape(cookie.substring(c_start, c_end));
    }
  }
  return '';
}

module.exports = router;
