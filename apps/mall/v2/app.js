const Koa = require('koa');
const path = require('path');
const Router = require('koa-router');
const static = require('koa-static');
const bodyParser = require('koa-bodyparser');

const app = new Koa();
const staticPath = './static/dist';

app.use(bodyParser());

app.use(static(path.join(__dirname, staticPath)));

app.use(require('./router/index').routes());

app.listen(7000);
console.log('start success, port: 7000');
